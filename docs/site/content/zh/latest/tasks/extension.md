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
