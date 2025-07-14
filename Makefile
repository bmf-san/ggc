# Makefile for Go project

APP_NAME=ggc
OUT?=coverage.out

.PHONY: install-tools deps build build-fast run run-fast test lint clean cover test-cover test-and-lint fmt

# Install required tools
install-tools:
	@echo "Installing required tools..."
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.2.1
	@echo "Tools installed successfully"

# Install dependencies and tools
deps: install-tools
	@echo "Installing dependencies..."
	go mod download
	go mod tidy
	@echo "Dependencies installed successfully"

# Build variables
VERSION := $(shell git describe --tags --always --dirty)
COMMIT := $(shell git rev-parse --short HEAD)
DATE := $(shell date -u '+%Y-%m-%d_%H:%M:%S')

# Full build with version info
build:
	go build -ldflags="-X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${DATE}" -o $(APP_NAME)

# Fast build for development (no version info, faster compilation)
build-fast:
	go build -o $(APP_NAME) main.go

run: build
	./$(APP_NAME)

run-fast: build-fast
	./$(APP_NAME)

fmt:
	go fmt ./...

test:
	go test ./...

lint: install-tools
	golangci-lint run --max-issues-per-linter=0 --max-same-issues=0

clean:
	rm -f $(APP_NAME)

cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out

test-cover:
	go test ./... -coverprofile=$(OUT)

test-and-lint: test lint
	@echo "All tests and lint checks passed"
