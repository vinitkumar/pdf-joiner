package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const (
	// Path to the Mac PDF joiner utility
	pdfJoinerPath = "/System/Library/Automator/Combine PDF Pages.action/Contents/MacOS/join"
)

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

	// Check if the Mac PDF joiner utility exists
	if !fileExists(pdfJoinerPath) {
		fmt.Printf("Error: PDF joiner utility not found at '%s'\n", pdfJoinerPath)
		fmt.Println("This tool only works on macOS systems.")
		os.Exit(1)
	}

	// Prepare the command to join PDFs
	args := []string{"-o", *outputPath}
	args = append(args, pdfFiles...)

	// Execute the command
	cmd := exec.Command(pdfJoinerPath, args...)
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		fmt.Printf("Error joining PDFs: %v\n", err)
		fmt.Printf("Command output: %s\n", output)
		os.Exit(1)
	}

	fmt.Printf("Successfully joined %d PDF files into '%s'\n", len(pdfFiles), *outputPath)
}

// fileExists checks if a file exists and is not a directory
func fileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
} 