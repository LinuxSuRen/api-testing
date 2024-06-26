#!/bin/bash
set -e

sleep 6
echo "Running k8s.sh"

ls -hal
cd api-testing
echo "build helm dependency"
helm dependency build

echo "install helm chart"
helm install --kube-apiserver https://server:6443 --kube-token abcd --kube-insecure-skip-tls-verify \
    api-testing . \
    --set service.type=NodePort \
    --set service.nodePort=30000 \
    --set persistence.enabled=false \
    --set image.registry=ghcr.io \
    --set image.repository=linuxsuren/api-testing \
    --set image.tag=master \
    --set extension.registry=ghcr.io

SERVER=http://server:30000 atest run -p git.yaml
