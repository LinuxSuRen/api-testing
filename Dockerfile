FROM docker.io/library/node:20-alpine3.17 AS ui

WORKDIR /workspace
COPY console/atest-ui .
RUN npm install --ignore-scripts --registry=https://registry.npmmirror.com
RUN npm run build-only

# FROM docker.io/apache/skywalking-go:0.4.0-go1.19 AS sk
# use above tag once https://github.com/apache/skywalking-go/pull/134 got released
FROM ghcr.io/apache/skywalking-go/skywalking-go:74b68861aed04b4d78fcc5b4bcd925113f7de81d-go1.19 AS sk

FROM docker.io/golang:1.19 AS builder

ARG VERSION
ARG GOPROXY
WORKDIR /workspace
COPY cmd/ cmd/
COPY pkg/ pkg/
COPY extensions/ extensions/
COPY operator/ operator/
COPY .github/testing/*.yaml sample/
COPY .github/ .github/
COPY sample/ sample/
COPY docs/ docs/
COPY go.mod go.mod
COPY go.sum go.sum
COPY go.work go.work
COPY go.work.sum go.work.sum
COPY main.go main.go
COPY README.md README.md
COPY LICENSE LICENSE

COPY --from=ui /workspace/dist/index.html cmd/data/index.html
COPY --from=ui /workspace/dist/assets/*.js cmd/data/index.js
COPY --from=ui /workspace/dist/assets/*.css cmd/data/index.css

COPY --from=sk /usr/local/bin/skywalking-go-agent /usr/local/bin/skywalking-go-agent

RUN GOPROXY=${GOPROXY} go mod download
RUN GOPROXY=${GOPROXY} CGO_ENABLED=0 go build -toolexec="skywalking-go-agent" -a -ldflags "-w -s -X github.com/linuxsuren/api-testing/pkg/version.version=${VERSION}" -o atest .
RUN GOPROXY=${GOPROXY} CGO_ENABLED=0 go build -ldflags "-w -s" -o atest-collector extensions/collector/main.go
RUN GOPROXY=${GOPROXY} CGO_ENABLED=0 go build -ldflags "-w -s" -o atest-store-orm extensions/store-orm/main.go
RUN GOPROXY=${GOPROXY} CGO_ENABLED=0 go build -ldflags "-w -s" -o atest-store-s3 extensions/store-s3/main.go
RUN GOPROXY=${GOPROXY} CGO_ENABLED=0 go build -ldflags "-w -s" -o atest-store-etcd extensions/store-etcd/main.go
RUN GOPROXY=${GOPROXY} CGO_ENABLED=0 go build -toolexec="skywalking-go-agent" -a -ldflags "-w -s" -o atest-store-git extensions/store-git/main.go

FROM docker.io/library/ubuntu:23.04

LABEL "com.github.actions.name"="API testing"
LABEL "com.github.actions.description"="API testing"
LABEL "com.github.actions.icon"="home"
LABEL "com.github.actions.color"="red"

LABEL "repository"="https://github.com/linuxsuren/api-testing"
LABEL "homepage"="https://github.com/linuxsuren/api-testing"
LABEL "maintainer"="Rick <linuxsuren@gmail.com>"

LABEL "Name"="API testing"

COPY --from=builder /workspace/atest /usr/local/bin/atest
COPY --from=builder /workspace/atest-collector /usr/local/bin/atest-collector
COPY --from=builder /workspace/atest-store-orm /usr/local/bin/atest-store-orm
COPY --from=builder /workspace/atest-store-s3 /usr/local/bin/atest-store-s3
COPY --from=builder /workspace/atest-store-etcd /usr/local/bin/atest-store-etcd
COPY --from=builder /workspace/atest-store-git /usr/local/bin/atest-store-git
COPY --from=builder /workspace/LICENSE /LICENSE
COPY --from=builder /workspace/README.md /README.md

COPY --from=builder /workspace/sample /var/www/sample

RUN apt update -y && \
    # required for atest-store-git
    apt install -y --no-install-recommends ssh-client ca-certificates && \
    apt install -y curl

CMD ["atest", "server", "--local-storage=/var/www/sample/*.yaml"]
