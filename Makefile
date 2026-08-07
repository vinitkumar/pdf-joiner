.PHONY: all build build-all build-darwin build-linux build-named build-universal-darwin check check-config clean deps fmt help install lint modernize release release-snapshot run test vet

BINARY_NAME := pdf-joiner
VERSION ?= dev
GO := go
BUILD_FLAGS := -trimpath
LDFLAGS := -s -w -X main.version=$(VERSION)
GOLANGCI_LINT_VERSION := v2.12.2
GOLANGCI_LINT ?= $(GO) run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

all: check build

build:
	$(GO) build $(BUILD_FLAGS) -ldflags "$(LDFLAGS)" -o $(BINARY_NAME) .

run: build
	./$(BINARY_NAME)

test:
	$(GO) test -race -cover ./...

vet:
	$(GO) vet ./...

fmt:
	gofmt -w .

modernize:
	$(GO) fix ./...

lint:
	$(GOLANGCI_LINT) run ./...

check: vet test lint check-config

clean:
	$(GO) clean
	rm -f $(BINARY_NAME) $(BINARY_NAME).exe
	rm -rf dist builds

install: build
	install -m 0755 $(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)

deps:
	$(GO) mod tidy

build-darwin:
	mkdir -p builds/darwin/amd64 builds/darwin/arm64
	GOOS=darwin GOARCH=amd64 $(GO) build $(BUILD_FLAGS) -ldflags "$(LDFLAGS)" -o builds/darwin/amd64/$(BINARY_NAME) .
	GOOS=darwin GOARCH=arm64 $(GO) build $(BUILD_FLAGS) -ldflags "$(LDFLAGS)" -o builds/darwin/arm64/$(BINARY_NAME) .

build-universal-darwin: build-darwin
	mkdir -p builds/darwin/universal
	lipo -create -output builds/darwin/universal/$(BINARY_NAME) builds/darwin/amd64/$(BINARY_NAME) builds/darwin/arm64/$(BINARY_NAME)

build-linux:
	mkdir -p builds/linux/amd64 builds/linux/arm64
	GOOS=linux GOARCH=amd64 $(GO) build $(BUILD_FLAGS) -ldflags "$(LDFLAGS)" -o builds/linux/amd64/$(BINARY_NAME) .
	GOOS=linux GOARCH=arm64 $(GO) build $(BUILD_FLAGS) -ldflags "$(LDFLAGS)" -o builds/linux/arm64/$(BINARY_NAME) .

build-all: build-darwin build-linux

build-named:
	mkdir -p builds/named
	GOOS=linux GOARCH=amd64 $(GO) build $(BUILD_FLAGS) -ldflags "$(LDFLAGS)" -o builds/named/$(BINARY_NAME)-linux-amd64 .
	GOOS=linux GOARCH=arm64 $(GO) build $(BUILD_FLAGS) -ldflags "$(LDFLAGS)" -o builds/named/$(BINARY_NAME)-linux-arm64 .
	GOOS=darwin GOARCH=amd64 $(GO) build $(BUILD_FLAGS) -ldflags "$(LDFLAGS)" -o builds/named/$(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 $(GO) build $(BUILD_FLAGS) -ldflags "$(LDFLAGS)" -o builds/named/$(BINARY_NAME)-darwin-arm64 .

release:
	goreleaser release --clean --timeout 60m

release-snapshot:
	goreleaser release --snapshot --clean

check-config:
	goreleaser check

help:
	@echo "build, test, vet, fmt, modernize, lint, check, build-all, build-named, release-snapshot, release, clean"
