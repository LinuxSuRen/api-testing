ARG LAN_ENV=docker.io/library/golang:1.21

FROM ghcr.io/linuxsuren/api-testing:master AS atest
FROM ghcr.io/linuxsuren/hd:v0.0.42 as downloader
RUN hd install jq
FROM $LAN_ENV

WORKDIR /workspace
COPY . .
COPY --from=downloader /usr/local/bin/jq /usr/local/bin/jq
COPY --from=atest /usr/local/bin/atest /usr/local/bin/atest

CMD [ "/workspace/entrypoint.sh" ]
