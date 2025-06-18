# PDF Joiner

[![Go CI](https://github.com/vinitkumar/pdf-joiner/actions/workflows/ci.yml/badge.svg)](https://github.com/vinitkumar/pdf-joiner/actions/workflows/ci.yml)

A simple command-line tool to join multiple PDF files into a single PDF document on **macOS and Linux**.

## Requirements

### macOS
- macOS (the tool uses the built-in macOS PDF joining utility)
- Go 1.24 or higher (for development)

### Linux
- Linux distribution with one of the following PDF tools installed:
  - **pdfunite** (from poppler-utils) - *Recommended, available on most distributions*
  - **ghostscript** (gs) - *Widely available fallback*
  - **qpdf** - *Alternative option*
- Go 1.24 or higher (for development)

### Installation of Linux PDF Tools

**Ubuntu/Debian:**
```bash
sudo apt-get install poppler-utils  # for pdfunite
# or
sudo apt-get install ghostscript    # for gs
# or
sudo apt-get install qpdf           # for qpdf
```

**RHEL/CentOS/Fedora:**
```bash
sudo dnf install poppler-utils      # for pdfunite
# or
sudo dnf install ghostscript        # for gs
# or
sudo dnf install qpdf               # for qpdf
```

**Arch Linux:**
```bash
sudo pacman -S poppler              # for pdfunite
# or
sudo pacman -S ghostscript          # for gs
# or
sudo pacman -S qpdf                 # for qpdf
```

## Installation

### From Source

```bash
git clone https://github.com/vinitkumar/pdf-joiner.git
cd pdf-joiner
go build -o pdf-joiner
```

### From Releases

Download the latest binary from the [Releases page](https://github.com/vinitkumar/pdf-joiner/releases).

## Usage

```bash
# Join PDFs with default output filename (timestamp-based)
./pdf-joiner file1.pdf file2.pdf file3.pdf

# Join PDFs with a custom output filename
./pdf-joiner -o output.pdf file1.pdf file2.pdf file3.pdf
```

## Features

- **Cross-platform support**: Works on macOS and Linux
- **Automatic tool detection**: Uses the best available PDF tool on each platform
- **Multiple backends on Linux**: Supports pdfunite, ghostscript, and qpdf with intelligent fallback
- Join multiple PDF files into a single document
- Specify custom output path
- Automatic output filename generation with timestamp
- Verification of input files

## How It Works

The tool automatically detects your operating system and uses the appropriate PDF joining method:

- **macOS**: Uses the built-in `/System/Library/Automator/Combine PDF Pages.action/Contents/MacOS/join` utility
- **Linux**: Automatically detects and uses the best available tool in this priority order:
  1. `pdfunite` (from poppler-utils) - Fastest and most reliable
  2. `ghostscript` (`gs`) - Widely available, stable fallback
  3. `qpdf` - Alternative option

If none of the Linux tools are available, the program will provide clear instructions on which packages to install.

## Development

### Testing

```bash
go test -v ./...
```

### Building

```bash
go build -o pdf-joiner
```

## License

MIT

## Author

Vinit Kumar  (mail@vinitkumar.me)
