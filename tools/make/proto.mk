# Building grpc generate related makefile defined.
#
# All make targets related grpc generate are defined in this file.

.PHONY: grpc
grpc:
	protoc --proto_path=. \
	--go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    pkg/server/server.proto \
	pkg/testing/remote/loader.proto \
    pkg/runner/monitor/monitor.proto

.PHONY: grpc.gw
grpc.gw:
	protoc -I . --grpc-gateway_out . \
    --grpc-gateway_opt logtostderr=true \
    --grpc-gateway_opt paths=source_relative \
    --grpc-gateway_opt generate_unbound_methods=true \
	--openapiv2_out . \
	--openapiv2_opt logtostderr=true \
	--openapiv2_opt generate_unbound_methods=true \
    pkg/server/server.proto

.PHONY: grpc.java
grpc.java:
	protoc --plugin=protoc-gen-grpc-java \
    --grpc-java_out=bin --proto_path=. \
	pkg/server/server.proto \
	pkg/testing/remote/loader.proto

.PHONY: grpc.js
grpc.js:
	protoc -I=pkg/server server.proto \
    --js_out=import_style=commonjs:bin \
    --grpc-web_out=import_style=commonjs,mode=grpcwebtext:bin

.PHONY: grpc.ts
# https://github.com/grpc/grpc-web
grpc.ts:
	protoc -I=pkg/server server.proto \
    --js_out=import_style=commonjs,binary:console/atest-ui/src \
    --grpc-web_out=import_style=typescript,mode=grpcwebtext:console/atest-ui/src

# grpc-java:
# 	protoc --plugin=protoc-gen-grpc-java=/usr/local/bin/protoc-gen-grpc-java \
#     --grpc-java_out=bin --proto_path=pkg/server server.proto

.PHONY: grpc.decs
grpc.decs:
	protoc --proto_path=. \
	--descriptor_set_out=.github/testing/server.pb \
    pkg/server/server.proto

.PHONY: grpc.testproto
grpc.testproto:
	protoc -I . \
	--descriptor_set_out=pkg/runner/grpc_test/test.pb \
	pkg/runner/grpc_test/test.proto

	protoc -I . \
	--go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	pkg/runner/grpc_test/test.proto

##@ Grpc Protobufs

.PHONY: grpc
grpc-go: ## Generate Go gRPC code
grpc-go: grpc

.PHONY: grpc-gw
grpc-gw: ## Generate Go gRPC Gateway code
grpc-gw: grpc.gw

.PHONY: grpc-java
grpc-java: ## Generate Java gRPC code
grpc-java: grpc.java

.PHONY: grpc-js
grpc-js: ## Generate JavaScript gRPC code
grpc-js: grpc.js

.PHONY: grpc-ts
grpc-ts: ## Generate TypeScript gRPC code
grpc-ts: grpc.ts

.PHONY: grpc-decs
grpc-decs: ## Generate DescriptorSet
grpc-decs: grpc.decs

.PHONY: grpc-all
grpc-all: ## Generate all gRPC code
grpc-all: grpc grpc-gw grpc-java grpc-js grpc-ts grpc-decs ## Generate all gRPC code

.PHONY: proto-test
proto-test: ## Test the protobuf files
proto-test: grpc.testproto
