# Building docs makefile defined.
#
# All make targets related to docs are defined in this file.

DOCS_OUTPUT_DIR := site/public

LINKINATOR_IGNORE := "github.com githubusercontent.com github.com _print v0.0.1"

##@ Docs

.PHONY: docs
docs: docs.clean 
	@$(LOG_TARGET)
	cd $(ROOT_DIR)/docs/site && npm install
	cd $(ROOT_DIR)/docs/site && npm run build:production

# Docs site, make by hexo.

.PHONY: docs-serve
docs-serve: ## Start API Testing Site Locally.
	@$(LOG_TARGET)
	cd $(ROOT_DIR)/docs/site && npm run serve

.PHONY: docs-clean
docs-clean: ## Remove all files that are created during builds.
docs-clean: docs.clean

.PHONY: docs.clean
docs.clean:
	@$(LOG_TARGET)
	rm -rf $(DOCS_OUTPUT_DIR)
	rm -rf docs/site/node_modules
	rm -rf docs/site/resources
	rm -f docs/site/package-lock.json
	rm -f docs/site/.hugo_build.lock
