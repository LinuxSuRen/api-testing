# A wrapper to do lint checks
#
# All make targets related to lint are defined in this file.

##@ Lint

GITHUB_ACTION ?=

.PHONY: lint
lint: ## Run all linter of code sources, including golint, yamllint.

# lint-deps is run separately in CI to separate the tooling install logs from the actual output logs generated
# by the lint tooling.
.PHONY: lint-deps
lint-deps: ## Everything necessary to lint

GOLANGCI_LINT_FLAGS ?= $(if $(GITHUB_ACTION),--out-format=github-actions)

.PHONY: lint.golint
lint: lint.golint
lint-deps: $(tools/golangci-lint)
lint.golint: $(tools/golangci-lint)
	@$(LOG_TARGET)
	$(tools/golangci-lint) run $(GOLANGCI_LINT_FLAGS) --build-tags=e2e,celvalidation --config=tools/linter/golangci-lint/.golangci.yml

.PHONY: lint.yamllint
lint: lint.yamllint
lint-deps: $(tools/yamllint)
lint.yamllint: $(tools/yamllint)
	@$(LOG_TARGET)
	$(tools/yamllint) --config-file=tools/linter/yamllint/.yamllint $$(git ls-files :*.yml :*.yaml | xargs -L1 dirname | sort -u)
