# Makefile for Go project

APP_NAME=ggc
OUT?=coverage.out

.PHONY: install-tools deps build run test lint clean cover test-cover test-and-lint fmt docs demos

# Install required tools
install-tools:
	@if ! command -v golangci-lint &> /dev/null; then \
		echo "Installing required tools..."; \
		go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest; \
		echo "Tools installed successfully"; \
	else \
		echo "Tools already installed"; \
	fi

# Install dependencies and tools
deps: install-tools
	@echo "Installing dependencies..."
	go mod download
	go mod tidy
	@echo "Dependencies installed successfully"

VERSION := $(shell git describe --tags --always --dirty)
COMMIT := $(shell git rev-parse --short HEAD)

# Full build with version info
build:
	go build -ldflags="-X main.version=${VERSION} -X main.commit=${COMMIT}" -o $(APP_NAME)

run: build
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
	go test $$(go list ./... | grep -v testutil) -coverprofile=coverage.out
	go tool cover -func=coverage.out

test-cover:
	go test $$(go list ./... | grep -v testutil) -coverprofile=$(OUT)

test-and-lint: test lint
	@echo "All tests and lint checks passed"

# Update documentation and shell completions from registry
.PHONY: docs completions

docs:
	@echo "Updating README.md command table..."
	@go run tools/cmd/gendocs/main.go
	@echo "README.md command table updated from registry"
	@$(MAKE) completions
	@echo "Documentation, completions updated successfully"


DEMO_SCENARIOS := cli-workflow interactive-overview branch-management stash-cycle
DEMO_WORKSPACE_ROOT := $(CURDIR)/docs/demos/workspaces
DEMO_OUTPUT_DIR := docs/demos/generated

# Generate VHS-powered demo assets
.PHONY: demos

demos: build
	@if ! command -v vhs >/dev/null 2>&1; then \
		echo "Error: vhs is required. Install https://github.com/charmbracelet/vhs"; \
		exit 1; \
	fi
	@if ! command -v ttyd >/dev/null 2>&1; then \
		echo "Error: ttyd is required by vhs. Install https://github.com/tsl0922/ttyd"; \
		exit 1; \
	fi
	@if ! command -v ffmpeg >/dev/null 2>&1; then \
		echo "Error: ffmpeg is required by vhs. Install it via your package manager."; \
		exit 1; \
	fi
	@mkdir -p $(DEMO_OUTPUT_DIR)
	@for scenario in $(DEMO_SCENARIOS); do \
		echo "Preparing fixture for $$scenario"; \
		./tools/demos/reset-fixture.sh $$scenario >/dev/null; \
	done
	@echo "Generating demo GIFs with vhs..."
	@GGC_DEMO_CLI_WORKDIR=$(DEMO_WORKSPACE_ROOT)/cli-workflow PATH="$(CURDIR):$$PATH" vhs docs/demos/scripts/cli-workflow.tape
	@GGC_DEMO_INTERACTIVE_WORKDIR=$(DEMO_WORKSPACE_ROOT)/interactive-overview PATH="$(CURDIR):$$PATH" vhs docs/demos/scripts/interactive-overview.tape
	@GGC_DEMO_BRANCH_WORKDIR=$(DEMO_WORKSPACE_ROOT)/branch-management PATH="$(CURDIR):$$PATH" vhs docs/demos/scripts/branch-management.tape
	@echo "Demo assets written to $(DEMO_OUTPUT_DIR)"

completions:
	@echo "Generating shell completions from registry..."
	@go run ./tools/cmd/gencompletions
	@echo "Shell completions updated from registry"
