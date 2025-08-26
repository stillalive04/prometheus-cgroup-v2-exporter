# Makefile for Prometheus cgroup v2 Exporter

# Variables
BINARY_NAME := prometheus-cgroup-v2-exporter
PACKAGE := github.com/stillalive04/prometheus-cgroup-v2-exporter
VERSION := $(shell git describe --tags --always --dirty)
COMMIT := $(shell git rev-parse --short HEAD)
DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GO_VERSION := $(shell go version | cut -d " " -f 3)

# Build flags
LDFLAGS := -X github.com/prometheus/common/version.Version=$(VERSION) \
           -X github.com/prometheus/common/version.Revision=$(COMMIT) \
           -X github.com/prometheus/common/version.Branch=$(shell git rev-parse --abbrev-ref HEAD) \
           -X github.com/prometheus/common/version.BuildDate=$(DATE) \
           -X github.com/prometheus/common/version.GoVersion=$(GO_VERSION)

# Docker variables
DOCKER_IMAGE := stillalive04/prometheus-cgroup-v2-exporter
DOCKER_TAG := $(VERSION)

# Directories
BIN_DIR := bin
DIST_DIR := dist
COVERAGE_DIR := coverage

# Default target
.DEFAULT_GOAL := build

# Help target
.PHONY: help
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Clean up build artifacts
.PHONY: clean
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	@rm -rf $(BIN_DIR) $(DIST_DIR) $(COVERAGE_DIR)
	@go clean -cache -testcache -modcache

# Format code
.PHONY: fmt
fmt: ## Format Go code
	@echo "Formatting code..."
	@go fmt ./...
	@goimports -w .

# Lint code
.PHONY: lint
lint: ## Run linters
	@echo "Running linters..."
	@golangci-lint run ./...

# Vet code
.PHONY: vet
vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...

# Run tests
.PHONY: test
test: ## Run unit tests
	@echo "Running unit tests..."
	@mkdir -p $(COVERAGE_DIR)
	@go test -v -race -coverprofile=$(COVERAGE_DIR)/coverage.out ./...

# Run integration tests
.PHONY: test-integration
test-integration: ## Run integration tests
	@echo "Running integration tests..."
	@go test -v -tags=integration ./tests/integration/...

# Run end-to-end tests
.PHONY: test-e2e
test-e2e: ## Run end-to-end tests
	@echo "Running end-to-end tests..."
	@go test -v -tags=e2e ./tests/e2e/...

# Generate test coverage report
.PHONY: coverage
coverage: test ## Generate coverage report
	@echo "Generating coverage report..."
	@go tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@go tool cover -func=$(COVERAGE_DIR)/coverage.out

# Run benchmarks
.PHONY: benchmark
benchmark: ## Run benchmarks
	@echo "Running benchmarks..."
	@go test -bench=. -benchmem ./...

# Tidy dependencies
.PHONY: tidy
tidy: ## Tidy go modules
	@echo "Tidying dependencies..."
	@go mod tidy

# Download dependencies
.PHONY: deps
deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download

# Build binary
.PHONY: build
build: deps ## Build binary
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BIN_DIR)
	@CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/$(BINARY_NAME) ./cmd/$(BINARY_NAME)

# Build for multiple platforms
.PHONY: build-all
build-all: ## Build for all platforms
	@echo "Building for all platforms..."
	@mkdir -p $(DIST_DIR)
	
	# Linux AMD64
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" \
		-o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/$(BINARY_NAME)
	
	# Linux ARM64
	@GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" \
		-o $(DIST_DIR)/$(BINARY_NAME)-linux-arm64 ./cmd/$(BINARY_NAME)
	
	# Linux ARM
	@GOOS=linux GOARCH=arm CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" \
		-o $(DIST_DIR)/$(BINARY_NAME)-linux-arm ./cmd/$(BINARY_NAME)
	
	# Darwin AMD64
	@GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" \
		-o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/$(BINARY_NAME)
	
	# Darwin ARM64
	@GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" \
		-o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/$(BINARY_NAME)

# Create release archives
.PHONY: package
package: build-all ## Create release packages
	@echo "Creating release packages..."
	@cd $(DIST_DIR) && \
	for binary in $(BINARY_NAME)-*; do \
		if [ -f "$$binary" ]; then \
			tar -czf "$$binary.tar.gz" "$$binary" && \
			rm "$$binary"; \
		fi; \
	done

# Install binary
.PHONY: install
install: build ## Install binary to $GOPATH/bin
	@echo "Installing $(BINARY_NAME)..."
	@go install -ldflags "$(LDFLAGS)" ./cmd/$(BINARY_NAME)

# Run locally
.PHONY: run
run: build ## Run the exporter locally
	@echo "Running $(BINARY_NAME)..."
	@./$(BIN_DIR)/$(BINARY_NAME) --log.level=debug

# Docker build
.PHONY: docker-build
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .
	@docker tag $(DOCKER_IMAGE):$(DOCKER_TAG) $(DOCKER_IMAGE):latest

# Docker push
.PHONY: docker-push
docker-push: docker-build ## Push Docker image
	@echo "Pushing Docker image..."
	@docker push $(DOCKER_IMAGE):$(DOCKER_TAG)
	@docker push $(DOCKER_IMAGE):latest

# Docker run
.PHONY: docker-run
docker-run: docker-build ## Run Docker container
	@echo "Running Docker container..."
	@docker run --rm -p 9753:9753 \
		--pid=host --privileged \
		-v /sys/fs/cgroup:/sys/fs/cgroup:ro \
		-v /proc:/host/proc:ro \
		$(DOCKER_IMAGE):$(DOCKER_TAG)

# Security scan
.PHONY: security
security: ## Run security scans
	@echo "Running security scans..."
	@gosec ./...
	@nancy sleuth

# Generate mocks
.PHONY: mocks
mocks: ## Generate mocks for testing
	@echo "Generating mocks..."
	@mockery --all --output=./internal/mocks

# Check code quality
.PHONY: check
check: fmt vet lint test ## Run all code quality checks

# Pre-commit checks
.PHONY: pre-commit
pre-commit: check security ## Run pre-commit checks

# Release preparation
.PHONY: release
release: clean check package docker-build ## Prepare release

# Development setup
.PHONY: dev-setup
dev-setup: ## Setup development environment
	@echo "Setting up development environment..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install github.com/securecodewarrior/nancy@latest
	@go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	@go install github.com/vektra/mockery/v2@latest

# Show version
.PHONY: version
version: ## Show version information
	@echo "Version: $(VERSION)"
	@echo "Commit: $(COMMIT)"
	@echo "Date: $(DATE)"
	@echo "Go Version: $(GO_VERSION)"

# Show build info
.PHONY: info
info: version ## Show build information
	@echo "Binary: $(BINARY_NAME)"
	@echo "Package: $(PACKAGE)"
	@echo "Docker Image: $(DOCKER_IMAGE):$(DOCKER_TAG)"
