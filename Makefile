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

.DEFAULT_GOAL := build
