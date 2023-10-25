You could install `api-testing` via Helm chart:

```shell
helm install atest oci://docker.io/linuxsuren/api-testing \
    --version v0.0.1-helm \
    --set service.type=NodePort
```

or upgrade it:

```shell
helm upgrade atest oci://docker.io/surenpi/api-testing \
    --version v0.0.1-helm \
    --set image.tag=master \
    --set replicaCount=3
```

## SkyWalking

```shell
helm install atest oci://docker.io/linuxsuren/api-testing \
    --version v0.0.1-helm \
    --set skywalking.endpoint.http=http://skywalking-skywalking-helm-oap.skywalking.svc:12800
    --set skywalking.endpoint.grpc=skywalking-skywalking-helm-oap.skywalking.svc:11800
    --set service.type=NodePort
```
