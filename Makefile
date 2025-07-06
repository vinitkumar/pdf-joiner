.PHONY: build test clean install all release release-snapshot release-dry-run dev-build build-all build-darwin build-linux build-universal-darwin check-config deps help run

# Binary name
BINARY_NAME=pdf-joiner

# Go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean
GOMOD=$(GOCMD) mod

# Build flags
LDFLAGS=-ldflags "-s -w"

# Default target
all: test build

# Local build
build:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) -v .

# Run locally
run:
	$(GOBUILD) -o $(BINARY_NAME) -v .
	./$(BINARY_NAME)

# Run tests
test:
	$(GOTEST) -v ./...

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME) $(BINARY_NAME).exe
	rm -rf dist/ builds/

# Install binary to /usr/local/bin
install: build
	mv $(BINARY_NAME) /usr/local/bin/

# Dependency tidy
deps:
	$(GOMOD) tidy

# macOS cross compilation
build-darwin:
	mkdir -p builds/darwin/amd64 builds/darwin/arm64
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o builds/darwin/amd64/$(BINARY_NAME) -v .
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o builds/darwin/arm64/$(BINARY_NAME) -v .

# Universal macOS binary
build-universal-darwin: build-darwin
	mkdir -p builds/darwin/universal
	lipo -create -output builds/darwin/universal/$(BINARY_NAME) \
		builds/darwin/amd64/$(BINARY_NAME) \
		builds/darwin/arm64/$(BINARY_NAME)

# Linux cross compilation
build-linux:
	mkdir -p builds/linux/amd64 builds/linux/arm64
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o builds/linux/amd64/$(BINARY_NAME) -v .
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o builds/linux/arm64/$(BINARY_NAME) -v .

# Build for all platforms
build-all: build-darwin build-linux

# Build with proper OS/architecture naming
build-named:
	mkdir -p builds/named
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o builds/named/$(BINARY_NAME)-linux-amd64 -v .
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o builds/named/$(BINARY_NAME)-linux-arm64 -v .
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o builds/named/$(BINARY_NAME)-darwin-amd64 -v .
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o builds/named/$(BINARY_NAME)-darwin-arm64 -v .
	@echo "Built binaries with OS/architecture naming:"
	@ls -la builds/named/

# Release with GoReleaser
release:
	goreleaser release --clean

release-snapshot:
	goreleaser release --snapshot --clean --skip=publish

release-dry-run:
	goreleaser release --snapshot --clean --skip=publish --skip=validate

dev-build:
	goreleaser build --snapshot --clean --single-target

check-config:
	goreleaser check

help:
	@echo "Available targets:"
	@echo "  build                 - Build for current platform"
	@echo "  run                   - Build and run locally"
	@echo "  test                  - Run tests"
	@echo "  clean                 - Clean build artifacts"
	@echo "  install               - Install binary to /usr/local/bin"
	@echo "  deps                  - Tidy dependencies"
	@echo "  build-darwin          - Build for macOS (AMD64 + ARM64)"
	@echo "  build-universal-darwin- Create universal binary for macOS"
	@echo "  build-linux           - Build for Linux (AMD64 + ARM64)"
	@echo "  build-all             - Build for all platforms"
	@echo "  build-named           - Build with OS/architecture in binary names"
	@echo "  release               - Create release with GoReleaser"
	@echo "  release-snapshot      - Snapshot release (local)"
	@echo "  release-dry-run       - Dry-run release (no publish)"
	@echo "  dev-build             - Build locally using GoReleaser"
	@echo "  check-config          - Validate GoReleaser configuration"