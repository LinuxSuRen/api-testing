#!/bin/bash
set -e

sleep 6
echo "Running k8s.sh"
helm install --kube-apiserver https://server:6443 --kube-token abcd --kube-insecure-skip-tls-verify \
    api-testing ./api-testing \
    --set service.type=NodePort \
    --set service.nodePort=30000 \
    --set persistence.enabled=false \
    --set image.registry=linuxsuren.docker.scarf.sh \
	--set image.repository=linuxsuren/api-testing \
    --set image.tag=master

SERVER=http://server:30000 atest run -p git.yaml
