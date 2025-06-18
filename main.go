package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	// Path to the Mac PDF joiner utility
	macOSPDFJoinerPath = "/System/Library/Automator/Combine PDF Pages.action/Contents/MacOS/join"
)

// PDFJoiner interface for cross-platform PDF joining
type PDFJoiner interface {
	Join(inputFiles []string, outputPath string) error
	IsAvailable() bool
	GetName() string
}

// MacOSJoiner implements PDFJoiner for macOS
type MacOSJoiner struct{}

func (m *MacOSJoiner) Join(inputFiles []string, outputPath string) error {
	args := []string{"-o", outputPath}
	args = append(args, inputFiles...)

	cmd := exec.Command(macOSPDFJoinerPath, args...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("error joining PDFs: %v\nCommand output: %s", err, output)
	}

	return nil
}

func (m *MacOSJoiner) IsAvailable() bool {
	return fileExists(macOSPDFJoinerPath)
}

func (m *MacOSJoiner) GetName() string {
	return "macOS built-in PDF joiner"
}

// LinuxJoiner implements PDFJoiner for Linux
type LinuxJoiner struct {
	backend string
	command string
}

func (l *LinuxJoiner) Join(inputFiles []string, outputPath string) error {
	var cmd *exec.Cmd

	switch l.backend {
	case "pdfunite":
		args := append(inputFiles, outputPath)
		cmd = exec.Command("pdfunite", args...)
	case "gs":
		args := []string{"-dBATCH", "-dNOPAUSE", "-q", "-sDEVICE=pdfwrite", "-sOutputFile=" + outputPath}
		args = append(args, inputFiles...)
		cmd = exec.Command("gs", args...)
	case "qpdf":
		args := []string{"--empty", "--pages"}
		args = append(args, inputFiles...)
		args = append(args, "--", outputPath)
		cmd = exec.Command("qpdf", args...)
	default:
		return fmt.Errorf("unsupported Linux PDF backend: %s", l.backend)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error joining PDFs with %s: %v\nCommand output: %s", l.backend, err, output)
	}

	return nil
}

func (l *LinuxJoiner) IsAvailable() bool {
	_, err := exec.LookPath(l.backend)
	return err == nil
}

func (l *LinuxJoiner) GetName() string {
	return fmt.Sprintf("Linux %s", l.backend)
}

// NewLinuxJoiner creates a new LinuxJoiner with the best available backend
func NewLinuxJoiner() (*LinuxJoiner, error) {
	// Priority order of backends to try
	backends := []string{"pdfunite", "gs", "qpdf"}

	for _, backend := range backends {
		if _, err := exec.LookPath(backend); err == nil {
			return &LinuxJoiner{backend: backend, command: backend}, nil
		}
	}

	return nil, fmt.Errorf("no suitable PDF joining tool found on Linux. Please install one of: %s", strings.Join(backends, ", "))
}

// NewPDFJoiner creates the appropriate PDFJoiner for the current OS
func NewPDFJoiner() (PDFJoiner, error) {
	switch runtime.GOOS {
	case "darwin":
		joiner := &MacOSJoiner{}
		if !joiner.IsAvailable() {
			return nil, fmt.Errorf("macOS PDF joiner utility not found at '%s'.\nThis tool only works on macOS systems", macOSPDFJoinerPath)
		}
		return joiner, nil
	case "linux":
		return NewLinuxJoiner()
	default:
		return nil, fmt.Errorf("unsupported operating system: %s. This tool currently supports macOS and Linux only", runtime.GOOS)
	}
}

func main() {
	// Define command-line flags
	outputPath := flag.String("o", "", "Output path for the joined PDF")
	flag.Parse()

	// Get the PDF files to join from the remaining arguments
	pdfFiles := flag.Args()
	if len(pdfFiles) < 2 {
		fmt.Println("Error: At least two PDF files are required for joining")
		fmt.Println("Usage: pdf-joiner [-o output.pdf] file1.pdf file2.pdf [file3.pdf ...]")
		os.Exit(1)
	}

	// Validate that all input files exist and are PDFs
	for _, file := range pdfFiles {
		if !fileExists(file) {
			fmt.Printf("Error: File '%s' does not exist\n", file)
			os.Exit(1)
		}

		if filepath.Ext(file) != ".pdf" {
			fmt.Printf("Warning: File '%s' may not be a PDF file\n", file)
		}
	}

	// If no output path is provided, create a default one with the current date
	if *outputPath == "" {
		currentTime := time.Now().Format("2006-01-02-150405")
		*outputPath = fmt.Sprintf("joined-pdf-%s.pdf", currentTime)
	}

	// Ensure the output directory exists
	outputDir := filepath.Dir(*outputPath)
	if outputDir != "." && outputDir != "" {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			fmt.Printf("Error creating output directory: %v\n", err)
			os.Exit(1)
		}
	}

	// Create the appropriate PDF joiner for the current OS
	joiner, err := NewPDFJoiner()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Join the PDFs using the platform-specific joiner
	if err := joiner.Join(pdfFiles, *outputPath); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully joined %d PDF files into '%s' using %s\n", len(pdfFiles), *outputPath, joiner.GetName())
}

// fileExists checks if a file exists and is not a directory
func fileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
