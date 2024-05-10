# Makefile file for setting common variables
#
# All make targets related to common variables are defined in this file.

# ====================================================================================================
# Configure Make:
# ====================================================================================================

# Turn off .INTERMEDIATE file removal by marking all files as
# .SECONDARY.  .INTERMEDIATE file removal is a space-saving hack from
# a time when drives were small; on modern computers with plenty of
# storage, it causes nothing but headaches.
#
# See: https://news.ycombinator.com/item?id=16486331
# Which tells it 2 things:
# a. never apply property #1; never automatically delete intermediate files
# b. always apply property #2; always let us hop over missing elements in the dependency tree
.SECONDARY:

SHELL := /bin/bash

# ====================================================================================================
# ROOT Options:
# ====================================================================================================

# Set the root package
ROOT_PACKAGE = .

# REV is the short git sha of latest commit.
REV=$(shell git rev-parse --short HEAD)

# Set Root Directory Path
ifeq ($(origin ROOT_DIR),undefined)

# Windows use `$(shell pwd -P)`
# Linux /Mac use `$(abspath $(shell pwd -P))`
ROOT_DIR := $(shell pwd -P)
endif

# Set output directory path
ifeq ($(origin OUTPUT_DIR),undefined)
OUTPUT_DIR := $(ROOT_DIR)/bin
endif

# Supported Platforms for building multi arch binaries.
PLATFORMS ?= darwin_amd64 darwin_arm64 linux_amd64 linux_arm64

# Set a specific PLATFORM
ifeq ($(origin PLATFORM), undefined)
	ifeq ($(origin GOOS), undefined)
		GOOS := $(shell go env GOOS)
	endif
	ifeq ($(origin GOARCH), undefined)
		GOARCH := $(shell go env GOARCH)
	endif
	PLATFORM := $(GOOS)_$(GOARCH)
	# Use linux as the default OS when building images
	IMAGE_PLAT := linux_$(GOARCH)
else
	GOOS := $(word 1, $(subst _, ,$(PLATFORM)))
	GOARCH := $(word 2, $(subst _, ,$(PLATFORM)))
	IMAGE_PLAT := $(PLATFORM)
endif

# ====================================================================================================
# Includes:
# ====================================================================================================
include tools/make/tools.mk
include tools/make/golang.mk
include tools/make/image.mk
include tools/make/lint.mk
include tools/make/docs.mk
include tools/make/helm.mk
include tools/make/ui.mk
include tools/make/run.mk
include tools/make/proto.mk
include tools/make/test.mk

# Log the running target
# \033[0;32m -> green
# \033[0m reset color
LOG_TARGET = echo -e "\033[0;32m===========> Running $@ ... \033[0m"
# Log debugging info
define log
echo -e "\033[36m===========>$1\033[0m"
endef

define errorlog
echo -e "\033[0;31m===========>$1\033[0m"
endef

define USAGE_OPTIONS
Options:
  \033[36mIMAGES\033[0m
		 Backend images to make. Default image is api-testing.
		 This option is available when using: make image
		 Example: \033[36mmake image IMAGES="api-testing"\033[0m
  \033[36mPLATFORM\033[0m
		 The specified platform to build.
		 This option is available when using: make build
		 Example: \033[36mmake build PLATFORM="linux_amd64"\033[0m
		 Supported Platforms: linux_amd64 linux_arm64 darwin_amd64 darwin_arm64
  \033[36mPLATFORMS\033[0m
		 The multiple platforms to build.
		 This option is available when using: make build-multiarch
		 Example: \033[36mmake build-multiarch PLATFORMS="linux_amd64 linux_arm64"\033[0m
		 Default is "linux_amd64 linux_arm64 darwin_amd64 darwin_arm64".
endef
export USAGE_OPTIONS

##@ Common

.PHONY: generate
generate: ## Generate go code from templates and tags
generate: grpc-all

## help: Show this help info for API Testing makefiles.
.PHONY: help
help:
	@echo -e "\033[0;34mAPI Testing is an open source for API testing tool.\033[0m\n"
	@echo -e "Usage:\n  make \033[36m<Target>\033[0m \033[36m<Option>\033[0m\n\nTargets:"
	@awk 'BEGIN {FS = ":.*##"; printf ""} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
	@echo -e "\n$$USAGE_OPTIONS"
