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
