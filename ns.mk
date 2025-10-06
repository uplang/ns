# Generic Makefile for UP Namespaces
# This file is symlinked as Makefile in each namespace directory
.PHONY: help test lint build examples clean install

# Automatically determine namespace name from directory
NAMESPACE := $(notdir $(CURDIR))

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Namespace: $(NAMESPACE)'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

test: ## Run tests
	go test -v -race ./...

lint: ## Run linter
	GOFLAGS="-buildvcs=false" golangci-lint run --timeout=5m ./...

build: ## Build the namespace binary
	go build -v -o up-ns-$(NAMESPACE) .

examples: ## List examples
	@ls -1 examples/*.up 2>/dev/null || echo "No examples found"

clean: ## Clean build artifacts
	go clean
	rm -f up-ns-$(NAMESPACE)

install: ## Install dependencies
	go mod download
	go mod tidy

fmt: ## Format code
	go fmt ./...

.DEFAULT_GOAL := help
