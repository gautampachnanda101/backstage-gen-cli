.PHONY: help build test lint clean install

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT_SHA ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

BINARY_NAME := backstage-gen-cli
BUILD_DIR := bin
INSTALL_DIR := /usr/local/bin

LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.CommitSHA=$(COMMIT_SHA) -X main.BuildDate=$(BUILD_DATE) -w -s"

help:
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@echo '  build       Build the binary'
	@echo '  test        Run tests'
	@echo '  lint        Run linters'
	@echo '  clean       Clean build artifacts'
	@echo '  install     Install binary to system'
	@echo '  run         Build and run the binary'

build:
	@echo "Building $(BINARY_NAME) $(VERSION)..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

test:
	@echo "Running tests..."
	go test -v -race ./...

lint:
	@echo "Running linters..."
	go vet ./...
	go fmt ./...

clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)

install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_DIR)..."
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/
	@echo "Installation complete!"

run: build
	@echo "Running $(BINARY_NAME)..."
	$(BUILD_DIR)/$(BINARY_NAME)
