tools.bindir = tools/bin
tools.srcdir = tools/src


# `go get`-able things
# ====================
tools/golangci-lint = $(tools.bindir)/golangci-lint
tools/protoc-gen-go = $(tools.bindir)/protoc-gen-go
tools/protoc-gen-go-grpc = $(tools.bindir)/protoc-gen-go-grpc
tools/goreleaser = $(tools.bindir)/goreleaser
tools/buf = $(tools.bindir)/buf
tools/skywalking-eyes = $(tools.bindir)/skywalking-eyes

$(tools.bindir)/%: $(tools.srcdir)/%/pin.go $(tools.srcdir)/%/go.mod
	cd $(<D) && GOOS= GOARCH= go build -o $(abspath $@) $$(sed -En 's,^import _ "(.*)".*,\1,p' pin.go)

# `pip install`-able things
# =========================
tools/yamllint = $(tools.bindir)/yamllint
tools/codespell = $(tools.bindir)/codespell

$(tools.bindir)/%.d/venv: $(tools.srcdir)/%/requirements.txt
	mkdir -p $(@D)
	python3 -m venv $@
	$@/bin/pip3 install -r $< || (rm -rf $@; exit 1)
$(tools.bindir)/%: $(tools.bindir)/%.d/venv
	@if [ -e $(tools.srcdir)/$*/$*.sh ]; then \
		ln -sf ../../$(tools.srcdir)/$*/$*.sh $@; \
	else \
		ln -sf $*.d/venv/bin/$* $@; \
	fi

# `npm install`-able things
tools/markdownlint = $(tools.bindir)/markdownlint
tools/linkinator = $(tools.bindir)/linkinator

$(tools.bindir)/%: $(tools.srcdir)/%/package.json
	cd $(<D) && npm install
	ln -sf $(<D)/node_modules/.bin/$* $@

tools.clean: # Remove all tools
	@$(LOG_TARGET)
	rm -rf $(tools.bindir)

.PHONY: clean
clean: ## Remove all files that are created during builds.
clean: tools.clean

.PHONY: install-precheck
install-precheck: ## Install pre-check tools,
install-precheck:
	cp .github/pre-commit .git/hooks/pre-commit

.PHONY: goreleaser
goreleaser: ## Build the project using goreleaser
goreleaser:
	goreleaser build --rm-dist --snapshot
	make clean-embed-ui
