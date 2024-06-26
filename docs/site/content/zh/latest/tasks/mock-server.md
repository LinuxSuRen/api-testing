+++
title = "Mock Server 功能使用"
+++

## Get started

您可以通过执行下面的命令 mock 一个容器仓库服务[container registry](https://distribution.github.io/distribution/):

```shell
atest mock --prefix / mock/image-registry.yaml
```

之后，您可以通过使用如下的命令使用 mock 功能。

```shell
docker pull localhost:6060/repo/name:tag
```
