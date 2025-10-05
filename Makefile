# Makefile for UP Namespaces Repository

# Find all namespace directories with go.mod
NAMESPACES := $(shell find . -maxdepth 2 -name 'go.mod' -exec dirname {} \; | sed 's|^\./||' | sort)

.PHONY: help
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: test
test: ## Run tests for all namespaces
	@for ns in $(NAMESPACES); do \
		if [ -f $$ns/Makefile ]; then \
			echo "Testing $$ns..."; \
			$(MAKE) -C $$ns test || exit 1; \
		fi \
	done

.PHONY: build
build: test ## Build all namespaces
	@for ns in $(NAMESPACES); do \
		if [ -f $$ns/Makefile ]; then \
			echo "Building $$ns..."; \
			$(MAKE) -C $$ns build || exit 1; \
		fi \
	done

.PHONY: clean
clean: ## Clean all namespaces
	@for ns in $(NAMESPACES); do \
		if [ -f $$ns/Makefile ]; then \
			echo "Cleaning $$ns..."; \
			$(MAKE) -C $$ns clean || exit 1; \
		fi \
	done

.PHONY: test-ci
test-ci: ## Run CI tests locally using act (requires: brew install act)
	act --container-architecture linux/amd64 -j test

.DEFAULT_GOAL := build
