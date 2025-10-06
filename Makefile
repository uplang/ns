# Makefile for UP Namespaces Repository

# Find all namespace directories with go.mod

NAMESPACES := $(patsubst %/,%,$(dir $(wildcard */go.mod)))

.PHONY: help
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

TEST_NAMESPACES := $(patsubst %,test-%,$(NAMESPACES))

.PHONY: test
test: $(TEST_NAMESPACES) ## Run tests for all namespaces

.PHONY: $(TEST_NAMESPACES)
$(TEST_NAMESPACES):
	$(MAKE) -C $(patsubst test-%,%,$@) test

LINT_NAMESPACES := $(patsubst %,lint-%,$(NAMESPACES))

.PHONY: lint
lint: $(LINT_NAMESPACES) ## Run linter for all namespaces

.PHONY: $(LINT_NAMESPACES)
$(LINT_NAMESPACES):
	$(MAKE) -C $(patsubst lint-%,%,$@) lint

BUILD_NAMESPACES := $(patsubst %,build-%,$(NAMESPACES))

.PHONY: build
build: test $(BUILD_NAMESPACES) ## Build all namespaces

.PHONY: $(BUILD_NAMESPACES)
$(BUILD_NAMESPACES):
	$(MAKE) -C $(patsubst build-%,%,$@) build

CLEAN_NAMESPACES := $(patsubst %,clean-%,$(NAMESPACES))

.PHONY: clean
clean: $(CLEAN_NAMESPACES) ## Clean all namespaces

.PHONY: $(CLEAN_NAMESPACES)
$(CLEAN_NAMESPACES):
	$(MAKE) -C $(patsubst clean-%,%,$@) clean

FMT_NAMESPACES := $(patsubst %,fmt-%,$(NAMESPACES))

.PHONY: fmt
fmt: $(FMT_NAMESPACES) ## FMT all namespaces

.PHONY: $(FMT_NAMESPACES)
$(FMT_NAMESPACES):
	$(MAKE) -C $(patsubst fmt-%,%,$@) fmt

.PHONY: test-ci
test-ci: ## Run CI tests locally using act (requires: brew install act)
	act --container-architecture linux/amd64 -j test

.DEFAULT_GOAL := build
