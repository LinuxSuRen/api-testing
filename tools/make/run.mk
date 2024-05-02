# Building API Testing local run makefile defined.
#
# All make targets related to local run are defined in this file.

include tools/make/env.mk

ATEST_UI = console/atest-ui

##@ Local runs

.PHONY: run-server
run-server: ## Run the API Testing server
run-server: build-ui
	go run . server --local-storage 'bin/*.yaml' --console-path ${ATEST_UI}/dist

.PHONY: run-console
run-console: ## Run the API Testing console
run-console:
	cd ${ATEST_UI} && ${FRONT_RUNTIMES} run dev

copy:
	sudo cp ${OUTPUT_DIR}/atest /usr/local/bin/

.PHONY: copy-restart
copy-restart: ## Copy the binary to /usr/local/bin and restart the service
copy-restart: build-embed-ui
	atest service stop
	make copy
	atest service restart
