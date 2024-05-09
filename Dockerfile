FROM docker.io/library/node:20-alpine3.17 AS ui

WORKDIR /workspace
COPY console/atest-ui .
RUN npm install --ignore-scripts --registry=https://registry.npmmirror.com
RUN npm run build-only

FROM docker.io/golang:1.22.2 AS builder

ARG VERSION
ARG GOPROXY
WORKDIR /workspace
COPY cmd/ cmd/
COPY pkg/ pkg/
COPY operator/ operator/
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

# RUN go mod download
RUN CGO_ENABLED=0 go build -v -a -ldflags "-w -s -X github.com/linuxsuren/api-testing/pkg/version.version=${VERSION}\
    -X github.com/linuxsuren/api-testing/pkg/version.date=$(date +%Y-%m-%d)" -o atest .

FROM ghcr.io/linuxsuren/atest-ext-store-mongodb:master as mango
FROM ghcr.io/linuxsuren/atest-ext-store-git:master as git
FROM ghcr.io/linuxsuren/atest-ext-store-s3:master as s3
FROM ghcr.io/linuxsuren/atest-ext-store-etcd:master as etcd
FROM ghcr.io/linuxsuren/atest-ext-store-orm:master as orm
FROM ghcr.io/linuxsuren/atest-ext-monitor-docker:master as docker
FROM ghcr.io/linuxsuren/atest-ext-collector:master as collector
FROM ghcr.io/linuxsuren/api-testing-vault-extension:v0.0.1 as vault

FROM docker.io/library/ubuntu:23.10

LABEL "com.github.actions.name"="API testing"
LABEL "com.github.actions.description"="API testing"
LABEL "com.github.actions.icon"="home"
LABEL "com.github.actions.color"="red"
LABEL org.opencontainers.image.description "This is an API testing tool that supports HTTP, gRPC, and GraphQL." 

LABEL "repository"="https://github.com/linuxsuren/api-testing"
LABEL "homepage"="https://github.com/linuxsuren/api-testing"
LABEL "maintainer"="Rick <linuxsuren@gmail.com>"

LABEL "Name"="API testing"

COPY --from=builder /workspace/atest /usr/local/bin/atest
COPY --from=collector /usr/local/bin/atest-collector /usr/local/bin/atest-collector
COPY --from=orm /usr/local/bin/atest-store-orm /usr/local/bin/atest-store-orm
COPY --from=s3 /usr/local/bin/atest-store-s3 /usr/local/bin/atest-store-s3
COPY --from=etcd /usr/local/bin/atest-store-etcd /usr/local/bin/atest-store-etcd
COPY --from=git /usr/local/bin/atest-store-git /usr/local/bin/atest-store-git
COPY --from=mango /usr/local/bin/atest-store-mongodb /usr/local/bin/atest-store-mongodb
COPY --from=docker /usr/local/bin/atest-monitor-docker /usr/local/bin/atest-monitor-docker
COPY --from=vault /usr/local/bin/atest-vault-ext /usr/local/bin
COPY --from=builder /workspace/LICENSE /LICENSE
COPY --from=builder /workspace/README.md /README.md

RUN apt update -y && \
    # required for atest-store-git
    apt install -y --no-install-recommends ssh-client ca-certificates && \
    apt install -y curl

EXPOSE 8080
CMD ["atest", "server", "--local-storage=/var/data/api-testing/*.yaml"]
