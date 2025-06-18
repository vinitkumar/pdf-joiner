.PHONY: build test clean install all release release-snapshot build-all

# Binary name
BINARY_NAME=pdf-joiner

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOCLEAN=$(GOCMD) clean

# Build flags
LDFLAGS=-ldflags "-s -w"

all: test build

build:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) -v

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME).exe
	rm -rf dist/
	rm -rf builds/

run:
	$(GOBUILD) -o $(BINARY_NAME) -v
	./$(BINARY_NAME)

install: build
	mv $(BINARY_NAME) /usr/local/bin/

# Dependencies
deps:
	$(GOMOD) tidy

# Cross compilation with organized folder structure
build-darwin:
	mkdir -p builds/darwin/amd64 builds/darwin/arm64
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o builds/darwin/amd64/$(BINARY_NAME) -v
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o builds/darwin/arm64/$(BINARY_NAME) -v

# Universal binary for macOS (both Intel and Apple Silicon)
build-universal-darwin: build-darwin
	mkdir -p builds/darwin/universal
	lipo -create -output builds/darwin/universal/$(BINARY_NAME) builds/darwin/amd64/$(BINARY_NAME) builds/darwin/arm64/$(BINARY_NAME)

# Linux builds
build-linux:
	mkdir -p builds/linux/amd64 builds/linux/arm64
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o builds/linux/amd64/$(BINARY_NAME) -v
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o builds/linux/arm64/$(BINARY_NAME) -v

# Build for all platforms (Linux and macOS only)
build-all: build-darwin build-linux

# GoReleaser commands
release:
	goreleaser release --clean

release-snapshot:
	goreleaser release --snapshot --clean --skip=publish

release-dry-run:
	goreleaser release --snapshot --clean --skip=publish --skip=validate

# Development builds with GoReleaser
dev-build:
	goreleaser build --snapshot --clean --single-target

# Check GoReleaser configuration
check-config:
	goreleaser check

# Help
help:
	@echo "Available targets:"
	@echo "  build              - Build for current platform"
	@echo "  build-all          - Build for all platforms (Linux, macOS)"
	@echo "  build-darwin       - Build for macOS (Intel + Apple Silicon)"
	@echo "  build-linux        - Build for Linux (AMD64 + ARM64)"
	@echo "  build-universal-darwin - Create universal macOS binary"
	@echo "  test               - Run tests"
	@echo "  clean              - Clean build artifacts"
	@echo "  install            - Install to /usr/local/bin"
	@echo "  deps               - Tidy dependencies"
	@echo "  release            - Create a release with GoReleaser"
	@echo "  release-snapshot   - Create a snapshot release"
	@echo "  release-dry-run    - Dry run of release process"
	@echo "  dev-build          - Build for current platform with GoReleaser"
	@echo "  check-config       - Validate GoReleaser configuration" 