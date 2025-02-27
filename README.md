# PDF Joiner

A simple command-line utility for joining multiple PDF files into a single PDF document on macOS.

## Requirements

- macOS (This tool uses the built-in macOS PDF joining utility)
- Go 1.18 or later (for building from source)

## Installation

### From Source

1. Clone this repository:
   ```
   git clone https://github.com/vinitkumar/pdf-joiner.git
   cd pdf-joiner
   ```

2. Build the binary:
   ```
   make build
   ```

3. (Optional) Install the binary to your system:
   ```
   make install
   ```

## Usage

```
pdf-joiner [-o output.pdf] file1.pdf file2.pdf [file3.pdf ...]
```

### Options

- `-o`: Specify the output file path. If not provided, the output will be saved as `joined-pdf-YYYY-MM-DD-HHMMSS.pdf` in the current directory.

### Examples

Join two PDF files:
```
pdf-joiner file1.pdf file2.pdf
```

Join multiple PDF files with a specific output path:
```
pdf-joiner -o merged.pdf file1.pdf file2.pdf file3.pdf
```

Join all PDF files in a directory:
```
pdf-joiner -o merged.pdf /path/to/directory/*.pdf
```

## How It Works

This tool is a wrapper around the macOS built-in PDF joining utility located at:
```
/System/Library/Automator/Combine PDF Pages.action/Contents/MacOS/join
```

## Development

### Running Tests

```
make test
```

### Building for Different Architectures

Build for both Intel and Apple Silicon Macs:
```
make build-universal-darwin
```

### Cleaning Up

```
make clean
```

## License

MIT

## Author

Vinit Kumar 