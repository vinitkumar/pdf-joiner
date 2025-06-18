# PDF Joiner - Agent Guidelines

## Build/Test Commands
- Build: `go build -o pdf-joiner` or `make build`
- Test: `go test -v ./...` or `make test`
- Run single test: `go test -v -run TestFileExists`
- Clean: `make clean`
- Dependencies: `go mod tidy` or `make deps`

## Architecture
- Single-file Go CLI application
- Uses macOS built-in PDF joining utility at `/System/Library/Automator/Combine PDF Pages.action/Contents/MacOS/join`
- No external dependencies beyond Go standard library
- Simple command-line interface with flag parsing

## Code Style
- Standard Go formatting (use `gofmt`)
- Use standard library packages: `flag`, `fmt`, `os`, `os/exec`, `path/filepath`, `time`
- Constants in ALL_CAPS for paths
- Descriptive variable names: `pdfFiles`, `outputPath`, `tempFile`
- Error handling: check errors immediately and exit with os.Exit(1) on failure
- Use defer for cleanup (file closing, temp file removal)
- Table-driven tests with struct slices
- Test function names: `TestFunctionName`
- Comments for exported functions and complex logic only
