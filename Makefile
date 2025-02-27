.PHONY: build test clean install all

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

run:
	$(GOBUILD) -o $(BINARY_NAME) -v
	./$(BINARY_NAME)

install: build
	mv $(BINARY_NAME) /usr/local/bin/

# Dependencies
deps:
	$(GOMOD) tidy

# Cross compilation
build-darwin:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)-darwin-amd64 -v
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)-darwin-arm64 -v

# Universal binary for macOS (both Intel and Apple Silicon)
build-universal-darwin: build-darwin
	lipo -create -output $(BINARY_NAME) $(BINARY_NAME)-darwin-amd64 $(BINARY_NAME)-darwin-arm64
	rm $(BINARY_NAME)-darwin-amd64 $(BINARY_NAME)-darwin-arm64 