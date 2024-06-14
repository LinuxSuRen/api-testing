# Building docs makefile defined.
#
# All make targets related to docs are defined in this file.

##@ Desktop

.PHONY: desktop-start
desktop-start: ## Start Electron Desktop
desktop-start:
	cd console/atest-desktop && npm run start

desktop-package: ## Package Electron Desktop
desktop-package: build.embed.ui copy-to-desktop
	cd console/atest-desktop && npm i && npm run package

desktop-make: build.embed.ui ## Make an Electron Desktop
	cd console/atest-desktop && npm i && npm run make

desktop-publish: build.embed.ui ## Publish the Electron Desktop
	cd console/atest-desktop && npm i && npm run publish

desktop-test: ## Run unit tests of the Electron Desktop
	cd console/atest-desktop && npm i && npm test
