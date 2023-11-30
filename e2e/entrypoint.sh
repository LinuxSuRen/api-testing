#!/bin/bash
set -e

mkdir -p /root/.config/atest
mkdir -p /var/data

nohup atest server&
cmd="atest run -p test-suite-common.yaml --report github --report-github-identity e2e-testing --report-file /var/data/report.json --report-github-repo linuxsuren/api-testing --report-github-pr ${PULL_REQUEST:-0}"

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
exit 0
