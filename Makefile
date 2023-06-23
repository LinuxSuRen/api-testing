build:
	mkdir -p bin
	rm -rf bin/atest
	go build -o bin/atest main.go
goreleaser:
	goreleaser build --rm-dist --snapshot
build-image:
	docker build -t ghcr.io/linuxsuren/api-testing:dev .
run-image:
	docker run ghcr.io/linuxsuren/api-testing:dev
copy: build
	sudo cp bin/atest /usr/local/bin/
copy-restart: build
	atest service stop
	make copy
	atest service restart
test:
	go test ./... -cover -v -coverprofile=coverage.out
	go tool cover -func=coverage.out
test-collector:
	go test github.com/linuxsuren/api-testing/extensions/collector/./... -cover -v -coverprofile=collector-coverage.out
	go tool cover -func=collector-coverage.out
grpc:
	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    pkg/server/server.proto
grpc-gw:
	protoc -I . --grpc-gateway_out . \
    --grpc-gateway_opt logtostderr=true \
    --grpc-gateway_opt paths=source_relative \
    --grpc-gateway_opt generate_unbound_methods=true \
    pkg/server/server.proto
grpc-js:
	protoc -I=pkg/server server.proto \
    --js_out=import_style=commonjs:bin \
    --grpc-web_out=import_style=commonjs,mode=grpcwebtext:bin
# https://github.com/grpc/grpc-web
grpc-ts:
	protoc -I=pkg/server server.proto \
    --js_out=import_style=commonjs,binary:console/atest-ui/src \
    --grpc-web_out=import_style=typescript,mode=grpcwebtext:console/atest-ui/src
grpc-java:
	protoc --plugin=protoc-gen-grpc-java=/usr/local/bin/protoc-gen-grpc-java \
    --grpc-java_out=bin --proto_path=pkg/server server.proto
install-tool:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	hd i protoc-gen-grpc-web
