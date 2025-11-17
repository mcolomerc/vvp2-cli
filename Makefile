.PHONY: build clean install test fmt lint help build-all release version

# Variables
BINARY_NAME=vvp2
GO=$(shell which go || echo "/opt/homebrew/bin/go")
GOFLAGS=-v
DIST_DIR=dist

# Version info (dynamic from git or fallback)
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "v0.1.0-dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

# LDFLAGS for version info
LDFLAGS=-w -s -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.buildTime=$(BUILD_TIME)

# Platform targets for cross-compilation
PLATFORMS := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64

# Default target
all: build

## build: Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	$(GO) build $(GOFLAGS) -ldflags="$(LDFLAGS)" -o $(BINARY_NAME) .

## build-all: Build for all platforms
build-all: clean
	@echo "Building for all platforms..."
	@mkdir -p $(DIST_DIR)
	@for platform in $(PLATFORMS); do \
		os=$$(echo $$platform | cut -d'/' -f1); \
		arch=$$(echo $$platform | cut -d'/' -f2); \
		output_name=$(BINARY_NAME); \
		if [ $$os = "windows" ]; then output_name=$(BINARY_NAME).exe; fi; \
		echo "Building for $$os/$$arch..."; \
		mkdir -p $(DIST_DIR)/$(BINARY_NAME)-$$os-$$arch; \
		CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch \
		$(GO) build $(GOFLAGS) -ldflags="$(LDFLAGS)" -o $(DIST_DIR)/$(BINARY_NAME)-$$os-$$arch/$$output_name .; \
		if [ $$? -eq 0 ]; then \
			if [ $$os = "windows" ]; then \
				(cd $(DIST_DIR) && zip -q $(BINARY_NAME)-$$os-$$arch.zip $(BINARY_NAME)-$$os-$$arch/*); \
			else \
				(cd $(DIST_DIR) && tar -czf $(BINARY_NAME)-$$os-$$arch.tar.gz $(BINARY_NAME)-$$os-$$arch/); \
			fi; \
			rm -rf $(DIST_DIR)/$(BINARY_NAME)-$$os-$$arch; \
		else \
			echo "Failed to build for $$os/$$arch"; \
		fi; \
	done
	@echo "✓ Built all platforms in $(DIST_DIR)/"

## clean: Remove build artifacts
clean:
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)
	@rm -rf $(DIST_DIR)
	@$(GO) clean

## install: Install the binary to GOPATH/bin
install:
	@echo "Installing $(BINARY_NAME)..."
	$(GO) install $(GOFLAGS) .

## test: Run tests
test:
	@echo "Running tests..."
	$(GO) test -v ./...

## fmt: Format code
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...

## lint: Run linter
lint:
	@echo "Running linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Run: curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin"; \
	fi

## tidy: Tidy go modules
tidy:
	@echo "Tidying modules..."
	$(GO) mod tidy

## release: Create release artifacts (requires git tag)
release:
	@if [ -z "$$(git tag --points-at HEAD)" ]; then \
		echo "Error: No git tag found at HEAD. Please create a tag first."; \
		echo "Example: git tag v0.1.0 && git push origin v0.1.0"; \
		exit 1; \
	fi
	@echo "Creating release for $(VERSION)..."
	@$(MAKE) clean build-all
	@echo "✓ Release artifacts created in $(DIST_DIR)/"
	@echo "📦 Upload these files to GitHub release:"
	@ls -la $(DIST_DIR)/

## version: Show version information
version:
	@echo "Version: $(VERSION)"
	@echo "Commit: $(COMMIT)"
	@echo "Build Time: $(BUILD_TIME)"

## help: Show this help message
help:
	@echo "vvp2 CLI Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make <target>"
	@echo ""
	@echo "Targets:"
	@echo "  build        Build binary for current platform"
	@echo "  build-all    Build for all supported platforms"
	@echo "  clean        Remove build artifacts"
	@echo "  install      Install binary to GOPATH/bin"
	@echo "  test         Run tests"
	@echo "  fmt          Format code"
	@echo "  lint         Run linter"
	@echo "  tidy         Tidy go modules"
	@echo "  release      Create release artifacts"
	@echo "  version      Show version information"
	@echo "  help         Show this help message"

## run: Run the application with example config
run:
	@echo "Running $(BINARY_NAME)..."
	./$(BINARY_NAME) --help

## deps: Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GO) mod download

