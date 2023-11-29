#!/bin/bash
set -e

mkdir -p /root/.config/atest

nohup atest server&

echo "start to run testing"
kind=orm target=mysql:3306 driver=mysql atest run -p test-suite-common.yaml
kind=orm target=mariadb:3306 driver=mysql atest run -p test-suite-common.yaml
kind=etcd target=etcd:2379 atest run -p test-suite-common.yaml
kind=mongodb target=mongo:27017 atest run -p test-suite-common.yaml

# TODO online git repository is unstable, need to fix
# if [ -z "$GITEE_TOKEN" ]
# then
#     atest run -p git.yaml
# else
#     echo "found gitee token"
#     kind=git target=https://gitee.com/linuxsuren/test username=linuxsuren password=$GITEE_TOKEN atest run -p test-suite-common.yaml
# fi

# TODO need to fix below cases
kind=orm target=postgres:5432 driver=postgres atest run -p test-suite-common.yaml
# kind=s3 target=minio:9000 atest run -p test-suite-common.yaml

cat /root/.config/atest/stores.yaml
exit 0
