# A wrapper to ui related.
#
# All make targets related to ui in this file.

include tools/make/env.mk

ATEST_UI = console/atest-ui

BUILD_FLAGS += \
	-w -s \
	-X github.com/linuxsuren/api-testing/pkg/version.version=${TAG}

.PHONY: build.ui
build.ui:
	cd ${ATEST_UI} && ${FRONT_RUNTIMES} i && ${FRONT_RUNTIMES} run build-only

.PHONY: build.embed.ui
build.embed.ui: embed.ui
	make build clean.embed.ui

.PHONY: embed.ui
embed.ui:
	cd ${ATEST_UI} && ${FRONT_RUNTIMES} i && ${FRONT_RUNTIMES} run build-only
	cp ${ATEST_UI}/dist/index.html cmd/data/index.html
	cp ${ATEST_UI}/dist/assets/*.js cmd/data/index.js
	cp ${ATEST_UI}/dist/assets/*.css cmd/data/index.css

clean.embed.ui:
	git checkout cmd/data/index.html
	git checkout cmd/data/index.js
	git checkout cmd/data/index.css

.PHONY: test.ui
test.ui:
	cd ${ATEST_UI} && ${FRONT_RUNTIMES} run test:unit

.PHONY: test.ui.e2e
test.ui.e2e:
	cd ${ATEST_UI} && ${FRONT_RUNTIMES} i && ${FRONT_RUNTIMES} run test:e2e

##@ UI

.PHONY: build-ui
build-ui: ## Build APT Testing UI
build-ui: build.ui

.PHONY: test-ui
test-ui: ## Test APT Testing UI
test-ui: test.ui

.PHONY: test-ui-e2e
test-ui-e2e: ## Test APT Testing UI E2E
test-ui-e2e: test.ui.e2e

.PHONY: build-embed-ui
build-embed-ui:
build-embed-ui: build.embed.ui

.PHONY: embed-ui
embed-ui:
embed-ui: embed.ui
