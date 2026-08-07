# PDF Joiner

[![Go CI](https://github.com/vinitkumar/pdf-joiner/actions/workflows/ci.yml/badge.svg)](https://github.com/vinitkumar/pdf-joiner/actions/workflows/ci.yml)
[![golangci-lint](https://github.com/vinitkumar/pdf-joiner/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/vinitkumar/pdf-joiner/actions/workflows/golangci-lint.yml)

A small command-line tool that joins PDF files on macOS and Linux.

## Installation

Install the latest release without `sudo`:

```bash
curl -fsSL https://raw.githubusercontent.com/vinitkumar/pdf-joiner/main/install.sh | sh
```

The installer detects macOS or Linux and amd64 or arm64, verifies the archive checksum, installs `pdf-joiner` in `~/.local/bin`, and configures your shell's `PATH`. To install a specific release or choose another directory:

```bash
curl -fsSL https://raw.githubusercontent.com/vinitkumar/pdf-joiner/main/install.sh | PDF_JOINER_VERSION=v1.3.0 sh
PDF_JOINER_INSTALL_DIR=/custom/bin ./install.sh
```

You can also download an archive directly from the [GitHub Releases page](https://github.com/vinitkumar/pdf-joiner/releases). Every macOS release binary is signed locally with a Developer ID certificate and notarized by Apple.

To build from source, install Go 1.26.5 or newer and run:

```bash
go install github.com/vinitkumar/pdf-joiner@latest
```

## Platform requirements

On macOS, PDF Joiner uses the system PDF joining utility at:

```text
/System/Library/Automator/Combine PDF Pages.action/Contents/MacOS/join
```

On Linux, install at least one supported backend. PDF Joiner checks them in this order:

1. `pdfunite` from Poppler, recommended
2. `gs` from Ghostscript
3. `qpdf`

For example:

```bash
# Debian or Ubuntu
sudo apt-get install poppler-utils

# Fedora
sudo dnf install poppler-utils

# Arch Linux
sudo pacman -S poppler
```

PDF Joiner never installs system packages or invokes `sudo` itself.

## Usage

```bash
# Timestamped output filename
pdf-joiner first.pdf second.pdf third.pdf

# Explicit output filename
pdf-joiner -o combined.pdf first.pdf second.pdf

# Print the installed version
pdf-joiner -version
```

The output path must be different from every input path.

## Development

```bash
make test              # race-enabled tests and coverage
make vet               # Go static analysis
make lint              # golangci-lint
make modernize         # Go 1.26 safe modernizers
make build             # local binary
make release-snapshot  # local GoReleaser build without publishing
make release-local     # sign, notarize, and package a release on macOS
```

Official releases are built on the maintainer's Mac. Darwin binaries are Developer ID-signed and notarized before their archives are uploaded. The private signing key never leaves the login keychain. See [docs/RELEASING.md](docs/RELEASING.md) for the release procedure.

## License

MIT

## Author

Vinit Kumar (mail@vinitkumar.me)
