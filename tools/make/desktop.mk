# Building docs makefile defined.
#
# All make targets related to docs are defined in this file.

##@ Desktop

.PHONY: desktop-start
desktop-start: ## Start Electron Desktop
desktop-start:
	cd console/atest-desktop && npm run start

desktop-package: ## Package Electron Desktop
desktop-package: build.embed.ui
	cp ${OUTPUT_DIR}/${BINARY} console/atest-desktop
	cd console/atest-desktop && npm run package
