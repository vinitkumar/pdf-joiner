// Command pdf-joiner joins PDF files on macOS and Linux.
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const macOSPDFJoinerPath = "/System/Library/Automator/Combine PDF Pages.action/Contents/MacOS/join"

var version = "dev"

type pdfJoiner interface {
	Join(inputFiles []string, outputPath string) error
	Name() string
}

type commandRunner interface {
	CombinedOutput(name string, args ...string) ([]byte, error)
}

type execRunner struct{}

func (execRunner) CombinedOutput(name string, args ...string) ([]byte, error) {
	return exec.Command(name, args...).CombinedOutput()
}

type macOSJoiner struct {
	runner commandRunner
}

func (m macOSJoiner) Join(inputFiles []string, outputPath string) error {
	args := make([]string, 0, len(inputFiles)+2)
	args = append(args, "-o", outputPath)
	args = append(args, inputFiles...)

	output, err := m.runner.CombinedOutput(macOSPDFJoinerPath, args...)
	if err != nil {
		return commandError("macOS PDF joiner", output, err)
	}
	return nil
}

func (macOSJoiner) Name() string {
	return "macOS built-in PDF joiner"
}

type linuxJoiner struct {
	backend string
	command string
	runner  commandRunner
}

func (l linuxJoiner) Join(inputFiles []string, outputPath string) error {
	var args []string
	switch l.backend {
	case "pdfunite":
		args = make([]string, 0, len(inputFiles)+1)
		args = append(args, inputFiles...)
		args = append(args, outputPath)
	case "gs":
		args = []string{"-dBATCH", "-dNOPAUSE", "-q", "-sDEVICE=pdfwrite", "-sOutputFile=" + outputPath}
		args = append(args, inputFiles...)
	case "qpdf":
		args = []string{"--empty", "--pages"}
		args = append(args, inputFiles...)
		args = append(args, "--", outputPath)
	default:
		return fmt.Errorf("unsupported Linux PDF backend %q", l.backend)
	}

	output, err := l.runner.CombinedOutput(l.command, args...)
	if err != nil {
		return commandError(l.backend, output, err)
	}
	return nil
}

func (l linuxJoiner) Name() string {
	return "Linux " + l.backend
}

func commandError(tool string, output []byte, err error) error {
	detail := strings.TrimSpace(string(output))
	if detail == "" {
		return fmt.Errorf("%s failed: %w", tool, err)
	}
	return fmt.Errorf("%s failed: %w (output: %s)", tool, err, detail)
}

func newLinuxJoiner(lookPath func(string) (string, error), runner commandRunner) (*linuxJoiner, error) {
	for _, backend := range []string{"pdfunite", "gs", "qpdf"} {
		command, err := lookPath(backend)
		if err == nil {
			return &linuxJoiner{backend: backend, command: command, runner: runner}, nil
		}
	}

	return nil, errors.New("no PDF joining backend found; install pdfunite (recommended), ghostscript, or qpdf")
}

func newPDFJoiner() (pdfJoiner, error) {
	return newPDFJoinerForOS(runtime.GOOS, fileExists, exec.LookPath, execRunner{})
}

func newPDFJoinerForOS(
	goos string,
	pathExists func(string) bool,
	lookPath func(string) (string, error),
	runner commandRunner,
) (pdfJoiner, error) {
	switch goos {
	case "darwin":
		if !pathExists(macOSPDFJoinerPath) {
			return nil, fmt.Errorf("macOS PDF joiner utility not found at %q", macOSPDFJoinerPath)
		}
		return macOSJoiner{runner: runner}, nil
	case "linux":
		return newLinuxJoiner(lookPath, runner)
	default:
		return nil, fmt.Errorf("unsupported operating system %q; supported systems are macOS and Linux", goos)
	}
}

func run(args []string, stdout, stderr io.Writer, now func() time.Time, joinerFactory func() (pdfJoiner, error)) error {
	flags := flag.NewFlagSet("pdf-joiner", flag.ContinueOnError)
	flags.SetOutput(stderr)
	outputPath := flags.String("o", "", "output path for the joined PDF")
	showVersion := flags.Bool("version", false, "print version and exit")
	var usageErr error
	flags.Usage = func() {
		if _, err := fmt.Fprintln(stderr, "Usage: pdf-joiner [-o output.pdf] file1.pdf file2.pdf [file3.pdf ...]"); err != nil {
			usageErr = err
			return
		}
		flags.PrintDefaults()
	}

	if err := flags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return usageErr
		}
		return err
	}
	if *showVersion {
		if _, err := fmt.Fprintln(stdout, version); err != nil {
			return fmt.Errorf("write version: %w", err)
		}
		return nil
	}

	inputFiles := flags.Args()
	if len(inputFiles) < 2 {
		flags.Usage()
		if usageErr != nil {
			return fmt.Errorf("write usage: %w", usageErr)
		}
		return errors.New("at least two PDF files are required")
	}

	for _, path := range inputFiles {
		if err := validateInputFile(path); err != nil {
			return err
		}
		if !strings.EqualFold(filepath.Ext(path), ".pdf") {
			if _, err := fmt.Fprintf(stderr, "Warning: %q may not be a PDF file\n", path); err != nil {
				return fmt.Errorf("write warning: %w", err)
			}
		}
	}

	if *outputPath == "" {
		*outputPath = fmt.Sprintf("joined-pdf-%s.pdf", now().Format("2006-01-02-150405"))
	}
	if err := ensureDistinctOutput(inputFiles, *outputPath); err != nil {
		return err
	}

	outputDir := filepath.Dir(*outputPath)
	if outputDir != "." {
		if err := os.MkdirAll(outputDir, 0o755); err != nil {
			return fmt.Errorf("create output directory %q: %w", outputDir, err)
		}
	}

	joiner, err := joinerFactory()
	if err != nil {
		return err
	}
	if err := joiner.Join(inputFiles, *outputPath); err != nil {
		return err
	}

	if _, err := fmt.Fprintf(stdout, "Successfully joined %d PDF files into %q using %s\n", len(inputFiles), *outputPath, joiner.Name()); err != nil {
		return fmt.Errorf("write result: %w", err)
	}
	return nil
}

func validateInputFile(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("inspect input file %q: %w", path, err)
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("input path %q is not a regular file", path)
	}
	return nil
}

func ensureDistinctOutput(inputFiles []string, outputPath string) error {
	outputAbsolute, err := filepath.Abs(outputPath)
	if err != nil {
		return fmt.Errorf("resolve output path %q: %w", outputPath, err)
	}
	for _, inputPath := range inputFiles {
		inputAbsolute, err := filepath.Abs(inputPath)
		if err != nil {
			return fmt.Errorf("resolve input path %q: %w", inputPath, err)
		}
		if inputAbsolute == outputAbsolute {
			return fmt.Errorf("output path %q must differ from every input file", outputPath)
		}
	}
	return nil
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Mode().IsRegular()
}

func main() {
	if err := run(os.Args[1:], os.Stdout, os.Stderr, time.Now, newPDFJoiner); err != nil {
		fmt.Fprintf(os.Stderr, "pdf-joiner: %v\n", err)
		os.Exit(1)
	}
}
