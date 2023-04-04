FROM golang:1.17 as builder

WORKDIR /workspace
COPY . .
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

CMD ["atest", "server"]
