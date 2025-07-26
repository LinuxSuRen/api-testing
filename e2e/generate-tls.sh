#!/bin/bash
set -e

# Generate private key
openssl genrsa -out server.key 2048
# Generate self-signed certificate
openssl req -new -x509 -key server.key -out server.crt -days 36500 \
    -subj "/C=US/ST=Denial/L=Springfield/O=Dis/CN=atest" \
# Generate Certificate Signing Request (CSR)
openssl req -new -key server.key -out server.csr \
    -subj "/C=US/ST=Denial/L=Springfield/O=Dis/CN=atest" \
# Generate a new private key
openssl genpkey -algorithm RSA -out test.key
# Generate a new CSR
openssl req -new -nodes -key test.key -out test.csr -days 3650 \
    -subj "/C=US/ST=Denial/L=Springfield/O=Dis/CN=atest" \
    -config "openssl.cnf" -extensions v3_req
# Sign the new CSR with the self-signed certificate
openssl x509 -req -days 365 -in test.csr \
    -out test.pem -CA server.crt -CAkey server.key \
    -CAcreateserial -extfile "openssl.cnf" -extensions v3_req
