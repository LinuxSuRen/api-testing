# Building API Testing local run makefile defined.
#
# All make targets related to local run are defined in this file.

include tools/make/env.mk

ATEST_UI = console/atest-ui

##@ Local runs & init env

.PHONY: run-server
run-server: ## Run the API Testing server
run-server: build-ui
	go run . server --local-storage 'bin/*.yaml' --console-path ${ATEST_UI}/dist

.PHONY: run-console
run-console: ## Run the API Testing console
run-console:
	cd ${ATEST_UI} && ${FRONT_RUNTIMES} run dev

.PHONY: copy
copy: ## Copy the binary to /usr/local/bin
copy:
	sudo cp ${OUTPUT_DIR}/$(GOOS)/$(GOARCH)/atest /usr/local/bin/

.PHONY: copy-restart
copy-restart: ## Copy the binary to /usr/local/bin and restart the service
copy-restart: build-embed-ui
	atest service stop
	make copy
	atest service restart

.PHONY: hd
hd:
	curl https://linuxsuren.github.io/tools/install.sh|bash

.PHONY: install-tool
install-tool: ## Install the tools to init env [not support windows]
install-tool: hd
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	hd i protoc-gen-grpc-web
	hd i protoc-gen-grpc-gateway
	hd get protocolbuffers/protobuf@v25.1 -o protobuf.zip
	unzip protobuf.zip bin/protoc
	rm -rf protobuf.zip
	sudo install bin/protoc /usr/local/bin/
	sudo hd get https://github.com/grpc-ecosystem/grpc-gateway/releases/download/v2.18.1/protoc-gen-openapiv2-v2.18.1-linux-x86_64 -o /usr/local/bin/protoc-gen-openapiv2
	sudo chmod +x /usr/local/bin/protoc-gen-openapiv2

.PHONY: init-env
init-env: ## Install the tools to init env [not support windows]
init-env: hd
	hd i cli/cli
	gh extension install linuxsuren/gh-dev
