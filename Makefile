IMG_TOOL?=docker
BINARY?=atest
TOOLEXEC?= #-toolexec="skywalking-go-agent"
BUILD_FLAG?=-ldflags "-w -s -X github.com/linuxsuren/api-testing/pkg/version.version=$(shell git describe --tags) \
	-X github.com/linuxsuren/api-testing/pkg/version.date=$(shell date +%Y-%m-%d)"
GOPROXY?=direct
HELM_VERSION?=v0.0.2
APP_VERSION?=v0.0.13
HELM_REPO?=docker.io/linuxsuren

fmt:
	go mod tidy
	go fmt ./...
	cd extensions/store-etcd && go mod tidy && go fmt ./...
	cd extensions/store-git && go mod tidy && go fmt ./...
	cd extensions/store-orm && go mod tidy && go fmt ./...
	cd extensions/store-s3 && go mod tidy && go fmt ./...
	cd extensions/store-mongodb && go mod tidy && go fmt ./...
build:
	mkdir -p bin
	rm -rf bin/atest
	CGO_ENABLED=0 go build ${TOOLEXEC} -a ${BUILD_FLAG} -o bin/${BINARY} main.go
build-ext: build-ext-git build-ext-orm build-ext-s3 build-ext-etcd build-ext-mongodb build-ext-monitor-docker
build-ext-git:
	CGO_ENABLED=0 go build -ldflags "-w -s" -o bin/atest-store-git extensions/store-git/main.go
build-ext-orm:
	CGO_ENABLED=0 go build -ldflags "-w -s" -o bin/atest-store-orm extensions/store-orm/main.go
build-ext-etcd:
	CGO_ENABLED=0 go build -ldflags "-w -s" -o bin/atest-store-etcd extensions/store-etcd/main.go
build-ext-s3:
	CGO_ENABLED=0 go build -ldflags "-w -s" -o bin/atest-store-s3 extensions/store-s3/main.go
build-ext-mongodb:
	CGO_ENABLED=0 go build -ldflags "-w -s" -o bin/atest-store-mongodb extensions/store-mongodb/main.go
build-ext-monitor-docker:
	CGO_ENABLED=0 go build -ldflags "-w -s" -o bin/atest-monitor-docker extensions/monitor-docker/main.go
build-ui:
	cd console/atest-ui && npm i && npm run build-only
embed-ui:
	cd console/atest-ui && npm i && npm run build-only
	cp console/atest-ui/dist/index.html cmd/data/index.html
	cp console/atest-ui/dist/assets/*.js cmd/data/index.js
	cp console/atest-ui/dist/assets/*.css cmd/data/index.css
clean-embed-ui:
	git checkout cmd/data/index.html
	git checkout cmd/data/index.js
	git checkout cmd/data/index.css
build-embed-ui: embed-ui
	GOOS=${OS} go build ${TOOLEXEC} -a -ldflags "-w -s -X github.com/linuxsuren/api-testing/pkg/version.version=$(shell git rev-parse --short HEAD)" -o bin/${BINARY} main.go
	make clean-embed-ui
build-darwin:
	BINARY=atest_darwin GOOS=darwin make build
build-win:
	BINARY=atest.exe GOOS=windows make build
build-win-embed-ui:
	BINARY=atest.exe GOOS=windows make build-embed-ui
goreleaser:
	goreleaser build --rm-dist --snapshot
	make clean-embed-ui
build-image:
	${IMG_TOOL} build -t ghcr.io/linuxsuren/api-testing:master . \
		--build-arg GOPROXY=${GOPROXY} \
		--build-arg VERSION=$(shell git describe --abbrev=0 --tags)-$(shell git rev-parse --short HEAD)
run-image:
	docker run -p 7070:7070 -p 8080:8080 ghcr.io/linuxsuren/api-testing:master
run-server: build-ui
	go run . server --local-storage 'bin/*.yaml' --console-path console/atest-ui/dist
run-console:
	cd console/atest-ui && npm run dev
copy:
	sudo cp bin/atest /usr/local/bin/
copy-ext:
	sudo cp bin/atest-* /usr/local/bin/
copy-restart: build-embed-ui
	atest service stop
	make copy
	atest service restart

# helm
helm-package:
	helm package helm/api-testing --version ${HELM_VERSION}-helm --app-version ${APP_VERSION} -d bin
helm-push:
	helm push bin/api-testing-${HELM_VERSION}-helm.tgz oci://${HELM_REPO}
helm-lint:
	helm lint helm/api-testing

# plugins
plugin-git:
	GOOS=${OS} go build -ldflags "-w -s" -o bin/atest-store-git extensions/store-git/main.go

test:
	go test ./... -cover -v -coverprofile=coverage.out
	go tool cover -func=coverage.out
testlong:
	go test pkg/limit/limiter_long_test.go -v
test-ui:
	cd console/atest-ui && npm run test:unit
test-ui-e2e:
	cd console/atest-ui && npm run test:e2e
test-collector:
	go test github.com/linuxsuren/api-testing/extensions/collector/./... -cover -v -coverprofile=collector-coverage.out
	go tool cover -func=collector-coverage.out
test-store-orm:
	go test github.com/linuxsuren/api-testing/extensions/store-orm/./... -cover -v -coverprofile=store-orm-coverage.out
	go tool cover -func=store-orm-coverage.out
test-store-s3:
	go test github.com/linuxsuren/api-testing/extensions/store-s3/./... -cover -v -coverprofile=store-s3-coverage.out
	go tool cover -func=store-s3-coverage.out
test-store-git:
	go test github.com/linuxsuren/api-testing/extensions/store-git/./... -cover -v -coverprofile=store-git-coverage.out
	go tool cover -func=store-git-coverage.out
test-store-etcd:
	go test github.com/linuxsuren/api-testing/extensions/store-etcd/./... -cover -v -coverprofile=store-etcd-coverage.out
	go tool cover -func=store-etcd-coverage.out
test-operator:
	cd operator && make test # converage file path: operator/cover.out
test-all-backend: test test-collector test-store-orm test-store-s3 test-store-git test-store-etcd
test-all: test-all-backend test-ui
test-e2e:
	cd e2e && ./start.sh && ./start.sh compose-k8s.yaml
install-precheck:
	cp .github/pre-commit .git/hooks/pre-commit

grpc:
	protoc --proto_path=. \
	--go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    pkg/server/server.proto \
	pkg/testing/remote/loader.proto \
    pkg/runner/monitor/monitor.proto
grpc-gw:
	protoc -I . --grpc-gateway_out . \
    --grpc-gateway_opt logtostderr=true \
    --grpc-gateway_opt paths=source_relative \
    --grpc-gateway_opt generate_unbound_methods=true \
	--openapiv2_out . \
	--openapiv2_opt logtostderr=true \
	--openapiv2_opt generate_unbound_methods=true \
    pkg/server/server.proto
grpc-java:
	protoc --plugin=protoc-gen-grpc-java \
    --grpc-java_out=bin --proto_path=. \
	pkg/server/server.proto \
	pkg/testing/remote/loader.proto
grpc-js:
	protoc -I=pkg/server server.proto \
    --js_out=import_style=commonjs:bin \
    --grpc-web_out=import_style=commonjs,mode=grpcwebtext:bin
# https://github.com/grpc/grpc-web
grpc-ts:
	protoc -I=pkg/server server.proto \
    --js_out=import_style=commonjs,binary:console/atest-ui/src \
    --grpc-web_out=import_style=typescript,mode=grpcwebtext:console/atest-ui/src
# grpc-java:
# 	protoc --plugin=protoc-gen-grpc-java=/usr/local/bin/protoc-gen-grpc-java \
#     --grpc-java_out=bin --proto_path=pkg/server server.proto
grpc-decs:
	protoc --proto_path=. \
	--descriptor_set_out=.github/testing/server.pb \
    pkg/server/server.proto 

grpc-testproto:
	protoc -I . \
	--descriptor_set_out=pkg/runner/grpc_test/test.pb \
	pkg/runner/grpc_test/test.proto

	protoc -I . \
	--go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	pkg/runner/grpc_test/test.proto

hd:
	curl https://linuxsuren.github.io/tools/install.sh|bash
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
init-env: hd
	hd i cli/cli
	gh extension install linuxsuren/gh-dev
