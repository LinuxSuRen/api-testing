FROM bitnami/etcd:3.5.10 AS etcd

FROM ghcr.io/linuxsuren/api-testing:master

WORKDIR /workspace
COPY --from=etcd /opt/bitnami/etcd/bin/etcd /usr/local/bin/etcd
COPY . .
RUN chmod +x entrypoint.sh

CMD [ "/workspace/entrypoint.sh" ]
