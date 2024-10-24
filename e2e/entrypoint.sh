#!/bin/bash
set -e

SCRIPT_DIR=$(dirname "$(readlink -f "$0")")
mkdir -p /root/.config/atest
mkdir -p /var/data

# Generate private key
openssl genrsa -out server.key 2048
# Generate self-signed certificate
openssl req -new -x509 -key server.key -out server.crt -days 36500 \
    -subj "/C=US/ST=Denial/L=Springfield/O=Dis/CN=www.example.com"
# Generate Certificate Signing Request (CSR)
openssl req -new -key server.key -out server.csr \
    -subj "/C=US/ST=Denial/L=Springfield/O=Dis/CN=www.example.com"
# Generate a new private key
openssl genpkey -algorithm RSA -out test.key
# Generate a new CSR
openssl req -new -nodes -key test.key -out test.csr -days 3650 \
    -subj "/C=US/ST=Denial/L=Springfield/O=Dis/CN=www.example.com" \
    -config "openssl.cnf" -extensions v3_req
# Sign the new CSR with the self-signed certificate
openssl x509 -req -days 365 -in test.csr \
    -out test.pem -CA server.crt -CAkey server.key \
    -CAcreateserial -extfile "openssl.cnf" -extensions v3_req

echo "start to download extenions"
atest extension --output /usr/local/bin --registry ghcr.io git
atest extension --output /usr/local/bin --registry ghcr.io orm
atest extension --output /usr/local/bin --registry ghcr.io etcd
atest extension --output /usr/local/bin --registry ghcr.io mongodb

echo "start to run server"
nohup atest server --tls-grpc --cert-file test.pem --key-file test.key&
cmd="atest run -p test-suite-common.yaml --request-ignore-error"

echo "start to run testing: $cmd"
kind=orm target=mysql:3306 driver=mysql $cmd
kind=orm target=mariadb:3306 driver=mysql $cmd
kind=etcd target=etcd:2379 $cmd
kind=mongodb target=mongo:27017 $cmd
kind=orm target=postgres:5432 driver=postgres $cmd

# TODO online git repository is unstable, need to fix
# if [ -z "$GITEE_TOKEN" ]
# then
#     atest run -p git.yaml
# else
#     echo "found gitee token"
#     kind=git target=https://gitee.com/linuxsuren/test username=linuxsuren password=$GITEE_TOKEN atest run -p test-suite-common.yaml
# fi

# TODO need to fix below cases
# kind=s3 target=minio:9000 atest run -p test-suite-common.yaml

cat /root/.config/atest/stores.yaml
