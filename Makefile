# Makefile for UP Namespaces Repository

# Find all namespace directories with go.mod
NAMESPACES := $(patsubst %/,%,$(dir $(wildcard */go.mod)))

.PHONY: help
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ''
	@echo 'Namespace targets (generated):'
	@printf "  %-15s %s\n" "test" "Run tests for all namespaces"
	@printf "  %-15s %s\n" "lint" "Run linter for all namespaces"
	@printf "  %-15s %s\n" "build" "Build all namespaces"
	@printf "  %-15s %s\n" "clean" "Clean all namespaces"
	@printf "  %-15s %s\n" "fmt" "Format code for all namespaces"

# Template for creating namespace targets
# Usage: $(eval $(call NS_TARGET,target-name,dependencies))
# Creates both 'target' and 'target-<namespace>' phony targets
define NS_TARGET
$(1)_NAMESPACES := $$(patsubst %,$(1)-%,$$(NAMESPACES))

.PHONY: $(1)
$(1): $(2) $$($(1)_NAMESPACES)

.PHONY: $$($(1)_NAMESPACES)
$$($(1)_NAMESPACES):
	$$(MAKE) -C $$(patsubst $(1)-%,%,$$@) $(1)
endef

# Generate targets for each namespace operation
$(eval $(call NS_TARGET,test))
$(eval $(call NS_TARGET,lint))
$(eval $(call NS_TARGET,build,test))
$(eval $(call NS_TARGET,clean))
$(eval $(call NS_TARGET,fmt))

.PHONY: test-ci
test-ci: ## Run CI tests locally using act (requires: brew install act)
	act --container-architecture linux/amd64 -j test

.PHONY: bump-patch-version
bump-patch-version: ## Bump patch version (e.g., v0.0.1 -> v0.0.2) and create tag locally
	@$(MAKE) tag NEW_VERSION=$(shell go tool svu patch)

.PHONY: bump-minor-version
bump-minor-version: ## Bump minor version (e.g., v0.0.1 -> v0.1.0) and create tag locally
	@$(MAKE) tag NEW_VERSION=$(shell go tool svu minor)

.PHONY: tag
tag: NEW_VERSION ?= $(shell go tool svu patch)
tag:
	@echo "Creating tag: $(NEW_VERSION)"
	git tag -s $(NEW_VERSION) -m "Release version $(NEW_VERSION)"
	@echo "Tag created. Push with: git push origin $(NEW_VERSION)"

.DEFAULT_GOAL := build
