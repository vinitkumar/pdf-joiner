package main

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestFileExists(t *testing.T) {
	// Create a temporary file for testing
	tempFile, err := os.CreateTemp("", "test-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer func() {
		if err := tempFile.Close(); err != nil {
			t.Logf("Failed to close temporary file: %v", err)
		}
		if err := os.Remove(tempFile.Name()); err != nil {
			t.Logf("Failed to remove temporary file: %v", err)
		}
	}()

	// Test cases
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "Existing file",
			path:     tempFile.Name(),
			expected: true,
		},
		{
			name:     "Non-existing file",
			path:     "non-existing-file.pdf",
			expected: false,
		},
		{
			name:     "Directory",
			path:     ".",
			expected: false, // fileExists returns false for directories
		},
	}

	// Run test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := fileExists(tc.path)
			if result != tc.expected {
				t.Errorf("fileExists(%s) = %v, expected %v", tc.path, result, tc.expected)
			}
		})
	}
}

func TestOutputPathGeneration(t *testing.T) {
	// This is a mock test to demonstrate how we would test the output path generation
	// In a real test, we would need to mock time.Now() or use dependency injection

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "pdf-joiner-test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Failed to remove temporary directory: %v", err)
		}
	}()

	// Test that we can create directories for output
	testOutputPath := filepath.Join(tempDir, "subdir", "output.pdf")
	outputDir := filepath.Dir(testOutputPath)

	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		t.Errorf("Failed to create output directory: %v", err)
	}

	// Verify the directory was created
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		t.Errorf("Output directory was not created: %v", err)
	}
}

func TestNewPDFJoiner(t *testing.T) {
	joiner, err := NewPDFJoiner()

	// The test should work on both macOS and Linux
	switch runtime.GOOS {
	case "darwin":
		if err != nil {
			// It's OK if macOS joiner is not available in test environment
			t.Logf("macOS PDF joiner not available in test environment: %v", err)
			return
		}
		if joiner != nil && joiner.GetName() != "macOS built-in PDF joiner" {
			t.Errorf("Expected macOS joiner, got: %s", joiner.GetName())
		}
	case "linux":
		if err != nil {
			// It's OK if Linux tools are not available in test environment
			t.Logf("Linux PDF tools not available in test environment: %v", err)
			return
		}
		if joiner != nil {
			expectedPrefixes := []string{"Linux pdfunite", "Linux gs", "Linux qpdf"}
			found := false
			for _, prefix := range expectedPrefixes {
				if joiner.GetName() == prefix {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected Linux joiner with one of %v, got: %s", expectedPrefixes, joiner.GetName())
			}
		}
	default:
		if err == nil {
			t.Errorf("Expected error for unsupported OS %s, but got nil", runtime.GOOS)
		}
	}
}

func TestMacOSJoiner(t *testing.T) {
	joiner := &MacOSJoiner{}

	// Test GetName
	if joiner.GetName() != "macOS built-in PDF joiner" {
		t.Errorf("Expected 'macOS built-in PDF joiner', got: %s", joiner.GetName())
	}

	// Test IsAvailable - this may return false in test environment
	available := joiner.IsAvailable()
	t.Logf("macOS joiner available: %v", available)
}

func TestLinuxJoiner(t *testing.T) {
	// Test creating a Linux joiner
	joiner, err := NewLinuxJoiner()
	if err != nil {
		t.Logf("Linux PDF tools not available in test environment: %v", err)
		return
	}

	// Test GetName
	name := joiner.GetName()
	if name == "" {
		t.Error("Linux joiner name should not be empty")
	}

	// Test IsAvailable
	if !joiner.IsAvailable() {
		t.Error("Linux joiner should be available if it was created successfully")
	}

	t.Logf("Linux joiner: %s", name)
}
