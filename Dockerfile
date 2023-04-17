FROM golang:1.18 AS builder

WORKDIR /workspace
COPY cmd/ cmd/
COPY pkg/ pkg/
COPY sample/ sample/
COPY go.mod go.mod
COPY go.sum go.sum
COPY main.go main.go
COPY README.md README.md
COPY LICENSE LICENSE

RUN go mod download
RUN CGO_ENABLE=0 go build -ldflags "-w -s" -o atest .

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
COPY --from=builder /workspace/LICENSE /LICENSE
COPY --from=builder /workspace/README.md /README.md

CMD ["atest", "server"]
