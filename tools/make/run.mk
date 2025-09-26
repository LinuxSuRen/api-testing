# Building API Testing local run makefile defined.
#
# All make targets related to local run are defined in this file.

include tools/make/env.mk

ATEST_UI = console/atest-ui
AI_PLUGIN_DIR := $(or $(AI_PLUGIN_SOURCE),../atest-ext-ai)

##@ Local runs & init env

.PHONY: build-ai-plugin
build-ai-plugin:
	@if [ -n "$(AI_PLUGIN_BINARY_URL)" ]; then \
		echo "📥 Downloading AI plugin binary from $(AI_PLUGIN_BINARY_URL)..."; \
		mkdir -p bin; \
		curl -L "$(AI_PLUGIN_BINARY_URL)" | tar xz -C bin/ --strip-components=1; \
		echo "✅ AI plugin binary downloaded"; \
	elif [ -d "$(AI_PLUGIN_DIR)" ]; then \
		echo "🔨 Building AI plugin from source..."; \
		cd $(AI_PLUGIN_DIR) && make build; \
		echo "✅ AI plugin built from source"; \
	else \
		echo "⚠️  AI plugin directory not found, skipping"; \
	fi

.PHONY: run-server
run-server: ## Run the API Testing server
run-server: build-ui build-ai-plugin run-backend
run-backend:
	go run . server --local-storage 'bin/*.yaml' --console-path ${ATEST_UI}/dist \
		--extension-registry ghcr.io --download-timeout 10m

run-tls-backend:
	go run . server --local-storage 'bin/*.yaml' --console-path ${ATEST_UI}/dist \
		--extension-registry ghcr.io --download-timeout 10m \
		--tls --cert-file e2e/test.pem \
		--key-file e2e/test.key

.PHONY: run-console
run-console: ## Run the API Testing console
run-console:
	cd ${ATEST_UI} && ${FRONT_RUNTIMES} run dev

.PHONY: copy.%
copy.%:
	@$(LOG_TARGET)
	$(eval PLATFORM := $(word 1,$(subst ., ,$*)))
	$(eval OS := $(word 1,$(subst _, ,$(PLATFORM))))
	$(eval ARCH := $(word 2,$(subst _, ,$(PLATFORM))))
	@$(call log, "Building binary $(BINARY) for $(OS)-$(ARCH).")
	sudo cp $(OUTPUT_DIR)/$(OS)/$(ARCH)/${BINARY} /usr/local/bin/

.PHONY: copy
copy: ## Copy the binary to /usr/local/bin
copy: $(addprefix copy., $(addprefix $(PLATFORM)., $(BINARY)))

.PHONY: copy-to-desktop.%
copy-to-desktop.%:
	$(eval PLATFORM := $(word 1,$(subst ., ,$*)))
	$(eval OS := $(word 1,$(subst _, ,$(PLATFORM))))
	$(eval ARCH := $(word 2,$(subst _, ,$(PLATFORM))))
	@$(call log, "Building binary $(BINARY) for $(OS)-$(ARCH).")
	cp $(OUTPUT_DIR)/$(OS)/$(ARCH)/${BINARY} console/atest-desktop


.PHONY: copy-to-desktop
copy-to-desktop: ## Copy the binary to console/atest-desktop
copy-to-desktop: $(addprefix copy-to-desktop., $(addprefix $(PLATFORM)., $(BINARY)))

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
	hd i protoc-gen-grpc-web@1.5.0
	hd i protoc-gen-grpc-gateway@v2.19.1
	hd get protocolbuffers/protobuf@v25.1 -o protobuf.zip
	unzip protobuf.zip bin/protoc
	rm -rf protobuf.zip
	sudo install bin/protoc /usr/local/bin/
	sudo hd get https://github.com/grpc-ecosystem/grpc-gateway/releases/download/v2.18.1/protoc-gen-openapiv2-v2.18.1-linux-x86_64 -o /usr/local/bin/protoc-gen-openapiv2
	sudo chmod +x /usr/local/bin/protoc-gen-openapiv2

.PHONY: init-env
init-env: ## Install the tools to init env [not support windows]
init-env: hd install-tool
	hd i cli/cli
	gh extension install linuxsuren/gh-dev
