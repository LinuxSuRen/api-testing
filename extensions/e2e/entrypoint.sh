#!/bin/bash
set -e

mkdir -p /root/.config/atest

nohup etcd&
nohup atest server&

atest run -p etcd.yaml
atest run -p git.yaml

cat /root/.config/atest/stores.yaml
