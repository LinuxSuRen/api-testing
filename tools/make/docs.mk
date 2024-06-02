# Building docs makefile defined.
#
# All make targets related to docs are defined in this file.

##@ Docs

.PHONY: clean
clean: ## Remove all files that are created during builds.
clean: docs.clean

.PHONY: check-links
check-links: ## Check for broken links in the docs.
check-links: docs-check-links

# Clean docs dir.
.PHONY: docs.clean
docs.clean:
	@$(LOG_TARGET)
	rm -rf $(DOCS_OUTPUT_DIR)

.PHONY: docs-check-links
docs-check-links:
	@$(LOG_TARGET)
	# Check for broken links
	npm install -g linkinator@6.0.4
	linkinator docs -r --concurrency 25 -s "github.com"
	