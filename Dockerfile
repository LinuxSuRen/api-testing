FROM golang:1.18 AS builder

WORKDIR /workspace
COPY cmd/ cmd/
COPY pkg/ pkg/
COPY extensions/ extensions/
COPY console/atest-ui atest-ui/
COPY sample/ sample/
COPY go.mod go.mod
COPY go.sum go.sum
COPY go.work go.work
COPY go.work.sum go.work.sum
COPY main.go main.go
COPY README.md README.md
COPY LICENSE LICENSE

RUN go mod download
RUN CGO_ENABLE=0 go build -ldflags "-w -s" -o atest .
RUN CGO_ENABLE=0 go build -ldflags "-w -s" -o atest-collector extensions/collector/main.go

FROM node:20-alpine3.17 AS ui

WORKDIR /workspace
COPY --from=builder /workspace/atest-ui .
RUN npm install --ignore-scripts
RUN npm run build-only

FROM ubuntu:23.04

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
COPY --from=builder /workspace/LICENSE /LICENSE
COPY --from=builder /workspace/README.md /README.md

RUN mkdir -p /var/www
COPY --from=builder /workspace/sample /var/www/sample
COPY --from=ui /workspace/dist /var/www/html

CMD ["atest", "server", "--console-path=/var/www/html", "--local-storage=/var/www/sample/*.yaml"]
