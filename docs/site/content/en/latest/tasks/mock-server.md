+++
title = "Mock server get started"
weight = -99
+++

## Get started

You can start a mock server of [container registry](https://distribution.github.io/distribution/) with below command:

```shell
atest mock --prefix / mock/image-registry.yaml
```

then, you can pull images from it:

```shell
docker pull localhost:6060/repo/name:tag
```
