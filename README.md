# PDF Joiner

[![Go CI](https://github.com/vinitkumar/pdf-joiner/actions/workflows/ci.yml/badge.svg)](https://github.com/vinitkumar/pdf-joiner/actions/workflows/ci.yml)

A simple command-line tool to join multiple PDF files into a single PDF document on macOS.

## Requirements

- macOS (the tool uses the built-in macOS PDF joining utility)
- Go 1.20 or higher (for development)

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

- Join multiple PDF files into a single document
- Specify custom output path
- Automatic output filename generation with timestamp
- Verification of input files

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
