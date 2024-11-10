FROM docker.io/library/node:20-alpine3.17 AS ui

WORKDIR /workspace
COPY console/atest-ui .
RUN npm install --ignore-scripts --registry=https://registry.npmmirror.com
RUN npm run build-only

FROM docker.io/golang:1.22.4 AS builder

ARG VERSION
ARG GOPROXY
WORKDIR /workspace
RUN mkdir -p console/atest-ui

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
COPY console/atest-ui/ui.go console/atest-ui/ui.go
COPY README.md README.md
COPY LICENSE LICENSE

COPY --from=ui /workspace/dist/index.html cmd/data/index.html
COPY --from=ui /workspace/dist/assets/*.js cmd/data/index.js
COPY --from=ui /workspace/dist/assets/*.css cmd/data/index.css

# RUN go mod download
RUN CGO_ENABLED=0 go build -v -a -ldflags "-w -s -X github.com/linuxsuren/api-testing/pkg/version.version=${VERSION}\
    -X github.com/linuxsuren/api-testing/pkg/version.date=$(date +%Y-%m-%d)" -o atest .

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
COPY --from=builder /workspace/LICENSE /LICENSE
COPY --from=builder /workspace/README.md /README.md

RUN apt update -y && \
    # required for atest-store-git
    apt install -y --no-install-recommends ssh-client ca-certificates curl && \
    rm -rf /var/lib/apt/lists/*

EXPOSE 8080
CMD ["atest", "server", "--local-storage=/var/data/api-testing/*.yaml"]
