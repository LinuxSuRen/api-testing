+++
title = "插件"
+++

`atest` 会把非核心、可扩展的功能以插件（extension）的形式实现。下面介绍有哪些插件，以及如何使用：

> 在不同的系统中，插件有着不同的表述，例如：extension、plugin 等。

| 类型 | 名称 | 描述 |
|------|------|------|
| 存储 | [orm](https://github.com/LinuxSuRen/atest-ext-store-orm)  | 保存数据到关系型数据库中，例如：MySQL |
| 存储 | [s3](https://github.com/LinuxSuRen/atest-ext-store-s3)   | 保存数据到对象存储中 |
| 存储 | [etcd](https://github.com/LinuxSuRen/atest-ext-store-etcd) | 保存数据到 Etcd 数据库中 |
| 存储 | [git](https://github.com/LinuxSuRen/atest-ext-store-git)  | 保存数据到 Git 仓库中 |
| 存储 | [mongodb](https://github.com/LinuxSuRen/atest-ext-store-mongodb) | 保存数据到 MongDB 中 |

> `atest` 也是唯一支持如此丰富的存储的接口开发、测试的开源工具。

## 下载插件

我们建议通过如下的命令来下载插件：

```shell
atest extension orm
```

上面的命令，会识别当前的操作系统，自动下载最新版本的插件。当然，用户可以通过自行编译、手动下载的方式获取插件二进制文件。

`atest` 可以从任意支持 OCI 的镜像仓库中（命令参数说明中给出了支持的镜像服务地址）下载插件，也可以指定下载超时时间：

```shell
atest extension orm --registry ghcr.io --timeout 2ms
```

想要下载其他类型的插件的话，可以使用下面的命令：

```shell
atest extension --kind data swagger
```
