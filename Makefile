IMG_TOOL?=podman

build:
	mkdir -p bin
	rm -rf bin/atest
	go build -o bin/atest main.go
build-embed-ui:
	cd console/atest-ui && npm i && npm run build-only
	cp console/atest-ui/dist/index.html cmd/data/index.html
	cp console/atest-ui/dist/assets/*.js cmd/data/index.js
	cp console/atest-ui/dist/assets/*.css cmd/data/index.css
	go build -ldflags "-w -s -X github.com/linuxsuren/api-testing/pkg/version.version=$(git rev-parse --short HEAD)" -o bin/atest main.go
	echo -n '' > cmd/data/index.html
	echo -n '' > cmd/data/index.js
	echo -n '' > cmd/data/index.css
goreleaser:
	goreleaser build --rm-dist --snapshot
build-image:
	${IMG_TOOL} build -t ghcr.io/linuxsuren/api-testing:master . --build-arg GOPROXY=https://goproxy.cn,direct
run-image:
	docker run -p 7070:7070 -p 8080:8080 ghcr.io/linuxsuren/api-testing:master
run-server:
	go run . server --local-storage 'sample/*.yaml' --console-path console/atest-ui/dist
run-console:
	cd console/atest-ui && npm run dev
copy:
	sudo cp bin/atest /usr/local/bin/
copy-restart: build-embed-ui
	atest service stop
	make copy
	atest service restart

test:
	go test ./... -cover -v -coverprofile=coverage.out
	go tool cover -func=coverage.out
test-ui:
	cd console/atest-ui && npm run test:unit
test-e2e:
	cd console/atest-ui && npm run test:e2e
test-collector:
	go test github.com/linuxsuren/api-testing/extensions/collector/./... -cover -v -coverprofile=collector-coverage.out
	go tool cover -func=collector-coverage.out
test-store-orm:
	go test github.com/linuxsuren/api-testing/extensions/store-orm/./... -cover -v -coverprofile=store-orm-coverage.out
	go tool cover -func=store-orm-coverage.out
test-all-backend: test test-collector test-store-orm
test-all: test-all-backend test-ui

install-precheck:
	cp .github/pre-commit .git/hooks/pre-commit

grpc:
	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    pkg/server/server.proto

	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    pkg/testing/remote/loader.proto
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
	hd i protoc-gen-grpc-gateway
