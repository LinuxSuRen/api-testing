+++
title = "通过 Helm 安装的方式使用 API Testing"
weight = -98
+++

You could install `api-testing` via Helm chart:

```shell
helm install atest oci://docker.io/linuxsuren/api-testing \
    --version v0.0.2-helm \
    --set service.type=NodePort
```

or upgrade it:

```shell
helm upgrade atest oci://docker.io/surenpi/api-testing \
    --version v0.0.2-helm \
    --set image.tag=master \
    --set replicaCount=3
```

## SkyWalking

```shell
helm install atest oci://docker.io/linuxsuren/api-testing \
    --version v0.0.2-helm \
    --set image.tag=master \
    --set service.type=NodePort \
    --set service.nodePort=30154 \
    --set skywalking.endpoint.http=http://skywalking-skywalking-helm-oap.skywalking.svc:12800 \
    --set skywalking.endpoint.grpc=skywalking-skywalking-helm-oap.skywalking.svc:11800
```
