You could install api-testing via Helm chart:

```shell
helm install atest oci://registry-1.docker.io/surenpi/api-testing \
    --version v0.0.2 \
    --set service.type=NodePort
```

or upgrade it:

```shell
helm upgrade atest oci://registry-1.docker.io/surenpi/api-testing --version v0.0.3
```
