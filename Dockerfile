FROM golang:1.17 as builder

WORKDIR /workspace
COPY . .
RUN go mod download
RUN CGO_ENABLE=0 go build -ldflags "-w -s" -o atest cmd/*.go

FROM ghcr.io/linuxsuren/hd:v0.0.67 as hd

FROM alpine:3.10

LABEL "com.github.actions.name"="API testing"
LABEL "com.github.actions.description"="API testing"
LABEL "com.github.actions.icon"="home"
LABEL "com.github.actions.color"="red"

LABEL "repository"="https://github.com/linuxsuren/api-testing"
LABEL "homepage"="https://github.com/linuxsuren/api-testing"
LABEL "maintainer"="Rick <linuxsuren@gmail.com>"

LABEL "Name"="API testing"

ENV LC_ALL C.UTF-8
ENV LANG en_US.UTF-8
ENV LANGUAGE en_US.UTF-8

RUN apk add --no-cache \
        git \
        openssh-client \
        libc6-compat \
        libstdc++

COPY entrypoint.sh /entrypoint.sh
COPY --from=builder /workspace/atest /usr/bin/atest
COPY --from=hd /usr/local/bin/hd /usr/local/bin/hd

ENTRYPOINT ["/entrypoint.sh"]
