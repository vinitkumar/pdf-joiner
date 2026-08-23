package main

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"
)

type fakeRunner struct {
	name   string
	args   []string
	output []byte
	err    error
}

func (f *fakeRunner) CombinedOutput(name string, args ...string) ([]byte, error) {
	f.name = name
	f.args = slices.Clone(args)
	return f.output, f.err
}

type fakeJoiner struct {
	inputs []string
	output string
	err    error
}

func (f *fakeJoiner) Join(inputs []string, output string) error {
	f.inputs = slices.Clone(inputs)
	f.output = output
	return f.err
}

func (*fakeJoiner) Name() string { return "fake joiner" }

func TestFileExists(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "input.pdf")
	if err := os.WriteFile(filePath, []byte("pdf"), 0o600); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		path string
		want bool
	}{
		{name: "file", path: filePath, want: true},
		{name: "missing", path: filepath.Join(tempDir, "missing.pdf"), want: false},
		{name: "directory", path: tempDir, want: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := fileExists(test.path); got != test.want {
				t.Fatalf("fileExists(%q) = %v, want %v", test.path, got, test.want)
			}
		})
	}
}

func TestNewPDFJoinerForOS(t *testing.T) {
	runner := &fakeRunner{}
	lookPath := func(command string) (string, error) {
		if command == "gs" {
			return "/usr/local/bin/gs", nil
		}
		return "", errors.New("not found")
	}

	tests := []struct {
		name      string
		goos      string
		exists    bool
		wantName  string
		wantError string
	}{
		{name: "macOS", goos: "darwin", exists: true, wantName: "macOS built-in PDF joiner"},
		{name: "macOS tool missing", goos: "darwin", wantError: "utility not found"},
		{name: "Linux", goos: "linux", wantName: "Linux gs"},
		{name: "unsupported", goos: "windows", wantError: "unsupported operating system"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			joiner, err := newPDFJoinerForOS(test.goos, func(string) bool { return test.exists }, lookPath, runner)
			if test.wantError != "" {
				if err == nil || !strings.Contains(err.Error(), test.wantError) {
					t.Fatalf("newPDFJoinerForOS() error = %v, want containing %q", err, test.wantError)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if got := joiner.Name(); got != test.wantName {
				t.Fatalf("Name() = %q, want %q", got, test.wantName)
			}
		})
	}
}

func TestLinuxJoinerCommandArguments(t *testing.T) {
	tests := []struct {
		backend string
		want    []string
	}{
		{backend: "pdfunite", want: []string{"one.pdf", "two.pdf", "out.pdf"}},
		{backend: "gs", want: []string{"-dBATCH", "-dNOPAUSE", "-q", "-sDEVICE=pdfwrite", "-sOutputFile=out.pdf", "one.pdf", "two.pdf"}},
		{backend: "qpdf", want: []string{"--empty", "--pages", "one.pdf", "two.pdf", "--", "out.pdf"}},
	}

	for _, test := range tests {
		t.Run(test.backend, func(t *testing.T) {
			runner := &fakeRunner{}
			joiner := linuxJoiner{backend: test.backend, command: "/bin/" + test.backend, runner: runner}
			inputs := []string{"one.pdf", "two.pdf"}
			if err := joiner.Join(inputs, "out.pdf"); err != nil {
				t.Fatal(err)
			}
			if runner.name != joiner.command {
				t.Fatalf("command = %q, want %q", runner.name, joiner.command)
			}
			if !slices.Equal(runner.args, test.want) {
				t.Fatalf("args = %#v, want %#v", runner.args, test.want)
			}
			if !slices.Equal(inputs, []string{"one.pdf", "two.pdf"}) {
				t.Fatalf("Join mutated its input slice: %#v", inputs)
			}
		})
	}
}

func TestMacOSJoinerCommandArguments(t *testing.T) {
	runner := &fakeRunner{}
	joiner := macOSJoiner{runner: runner}
	if err := joiner.Join([]string{"one.pdf", "two.pdf"}, "out.pdf"); err != nil {
		t.Fatal(err)
	}

	want := []string{"-o", "out.pdf", "one.pdf", "two.pdf"}
	if runner.name != macOSPDFJoinerPath || !slices.Equal(runner.args, want) {
		t.Fatalf("command = %q %#v, want %q %#v", runner.name, runner.args, macOSPDFJoinerPath, want)
	}
}

func TestJoinErrorWrapsCauseAndOutput(t *testing.T) {
	cause := errors.New("exit status 1")
	runner := &fakeRunner{output: []byte("bad PDF\n"), err: cause}
	joiner := linuxJoiner{backend: "pdfunite", command: "pdfunite", runner: runner}

	err := joiner.Join([]string{"one.pdf", "two.pdf"}, "out.pdf")
	if !errors.Is(err, cause) {
		t.Fatalf("Join() error %v does not wrap %v", err, cause)
	}
	if !strings.Contains(err.Error(), "bad PDF") {
		t.Fatalf("Join() error %q does not contain command output", err)
	}
}

func TestRun(t *testing.T) {
	tempDir := t.TempDir()
	first := filepath.Join(tempDir, "one.pdf")
	second := filepath.Join(tempDir, "TWO.PDF")
	for _, path := range []string{first, second} {
		if err := os.WriteFile(path, []byte("pdf"), 0o600); err != nil {
			t.Fatal(err)
		}
	}

	joiner := &fakeJoiner{}
	var stdout, stderr bytes.Buffer
	now := func() time.Time { return time.Date(2026, time.August, 8, 12, 34, 56, 0, time.UTC) }
	err := run([]string{first, second}, &stdout, &stderr, now, func() (pdfJoiner, error) { return joiner, nil })
	if err != nil {
		t.Fatal(err)
	}
	if joiner.output != "joined-pdf-2026-08-08-123456.pdf" {
		t.Fatalf("output = %q", joiner.output)
	}
	if !strings.Contains(stdout.String(), "Successfully joined 2 PDF files") {
		t.Fatalf("stdout = %q", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q", stderr.String())
	}
}

func TestRunValidation(t *testing.T) {
	tempDir := t.TempDir()
	first := filepath.Join(tempDir, "one.pdf")
	second := filepath.Join(tempDir, "two.pdf")
	for _, path := range []string{first, second} {
		if err := os.WriteFile(path, []byte("pdf"), 0o600); err != nil {
			t.Fatal(err)
		}
	}

	tests := []struct {
		name string
		args []string
		want string
	}{
		{name: "too few inputs", args: []string{first}, want: "at least two"},
		{name: "missing input", args: []string{first, filepath.Join(tempDir, "missing.pdf")}, want: "inspect input file"},
		{name: "directory input", args: []string{first, tempDir}, want: "not a regular file"},
		{name: "output overwrites input", args: []string{"-o", first, first, second}, want: "must differ"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := run(test.args, &bytes.Buffer{}, &bytes.Buffer{}, time.Now, func() (pdfJoiner, error) {
				return nil, errors.New("factory should not be called")
			})
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("run() error = %v, want containing %q", err, test.want)
			}
		})
	}
}

func TestRunVersion(t *testing.T) {
	oldVersion := version
	version = "v1.2.3"
	t.Cleanup(func() { version = oldVersion })

	var stdout bytes.Buffer
	if err := run([]string{"-version"}, &stdout, &bytes.Buffer{}, time.Now, func() (pdfJoiner, error) {
		return nil, errors.New("factory should not be called")
	}); err != nil {
		t.Fatal(err)
	}
	if got := strings.TrimSpace(stdout.String()); got != "v1.2.3" {
		t.Fatalf("version output = %q", got)
	}
}
