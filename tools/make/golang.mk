# Building golang binaries makefile defined.
#
# All make targets related to golang are defined in this file.

include tools/make/env.mk

# Set the build flags
GO_FLAGS += \
	-X github.com/linuxsuren/api-testing/pkg/version.version=$(TAG) \
	-X github.com/linuxsuren/api-testing/pkg/version.date=$(shell date +%Y-%m-%d)

# Binary file name
BINARY ?= atest

GOPATH := $(shell go env GOPATH)
ifeq ($(origin GOBIN), undefined)
	GOBIN := $(GOPATH)/bin
endif

GO_VERSION = $(shell grep -oE "^go [[:digit:]]*\.[[:digit:]]*" go.mod | cut -d' ' -f2)

# Build the target binary in target platform.
# The pattern of build.% is `build.{Platform}.{BINARY}`.
# If we want to build API Testing in linux amd64 platform,
# just execute `make go.build.linux_amd64.api-testing`
.PHONY: go.build.%
go.build.%:
	@$(LOG_TARGET)
	$(eval PLATFORM := $(word 1,$(subst ., ,$*)))
	$(eval OS := $(word 1,$(subst _, ,$(PLATFORM))))
	$(eval ARCH := $(word 2,$(subst _, ,$(PLATFORM))))
	@$(call log, "Building binary $(BINARY) for $(OS) $(ARCH).")
	CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) go build ${TOOLEXEC} -o $(OUTPUT_DIR)/$(OS)/$(ARCH)/${BINARY} -ldflags "$(GO_FLAGS)" $(ROOT_PACKAGE)/main.go

# Build the API Testing binaries in the hosted platforms.
.PHONY: go.build
go.build: $(addprefix go.build., $(addprefix $(PLATFORM)., $(BINARY)))

# Build the API Testing binaries in multi platforms
# It will build the linux/amd64, linux/arm64, darwin/amd64, darwin/arm64 binaries out.
.PHONY: go.build.multiarch
go.build.multiarch: $(foreach p,$(PLATFORMS),$(addprefix go.build., $(addprefix $(p)., $(BINARY))))

.PHONY: go.test.unit
go.test.unit: ## Run go unit tests
	go test -race ./...

.PHONY: go.test.lang
go.test.lang: ## Run go limiter long test
	go test pkg/limit/limiter_long_test.go -v

.PHONY: go.test.coverage
go.test.coverage: ## Run go unit tests with coverage
	@$(LOG_TARGET)
	go test ./... -cover -v -coverprofile=coverage.out
	go tool cover -func=coverage.out

.PHONY: go.test.lang
	@$(LOG_TARGET)
	go test pkg/limit/limiter_long_test.go -v

.PHONY: go.clean
go.clean: ## Clean the building output files
	@$(LOG_TARGET)
	rm -rf $(OUTPUT_DIR)

.PHONY: go.mod.lint
lint: go.mod.lint
go.mod.lint:
	@$(LOG_TARGET)
	@go mod tidy -compat=$(GO_VERSION)
	@go fmt ./...
	@if test -n "$$(git status -s -- go.mod go.sum)"; then \
		git diff --exit-code go.mod; \
		git diff --exit-code go.sum; \
		$(call errorlog, "Error: ensure all changes have been committed!"); \
		exit 1; \
	else \
		$(call log, "API Testing go module looks clean!"); \
   	fi

##@ Golang

.PHONY: build
build: ## Build API Testing for host platform. See Option PLATFORM and BINARY.
build: go.build

.PHONY: build-multiarch
build-multiarch: ## Build API Testing for multiple platforms. See Option PLATFORMS and IMAGES.
build-multiarch: go.build.multiarch

.PHONY: test
test: ## Run all Go test of code sources.
test: go.test.coverage

.PHONY: testlang
testlang: ## Run limiter long test.
testlang: go.test.lang

.PHONY: format
format: ## Update and check dependences with go mod tidy.
format: go.mod.lint

.PHONY: clean
clean: ## Remove all files that are created during builds.
clean: go.clean
