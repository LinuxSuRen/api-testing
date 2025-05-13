ARG LAN_ENV=docker.io/library/golang:1.21

FROM ghcr.io/linuxsuren/api-testing:master AS atest
FROM docker.io/demisto/jq:1.0.0.117497 AS jq
FROM $LAN_ENV

WORKDIR /workspace
COPY . .
COPY --from=jq /usr/local/bin/jq /usr/local/bin/jq
COPY --from=atest /usr/local/bin/atest /usr/local/bin/atest

CMD [ "/workspace/entrypoint.sh" ]
