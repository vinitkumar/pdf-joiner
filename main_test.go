package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileExists(t *testing.T) {
	// Create a temporary file for testing
	tempFile, err := os.CreateTemp("", "test-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

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
	defer os.RemoveAll(tempDir)
	
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