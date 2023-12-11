`atest` 版本发布 `v0.0.14`

`atest` 是一款用 Golang 编写的、开源的接口测试工具。

你可以在容器中启动：

```shell
docker run --network host \
  linuxsuren/api-testing:v0.0.14
```

或者，直接[下载二进制文件](https://github.com/LinuxSuRen/api-testing/releases/tag/v0.0.14)后启动：

```shell
atest server --local-storage /var/www/sample
```

## 主要的新功能

* 增加了对 `tRPC` 和 `gRPC` 协议的（命令行与 Web 界面）支持
* 新增了 Helm Chart 的安装方式
* 支持通过按钮切换暗模式
* 支持启动启动插件
* 支持在 Web 界面中参数化执行
* 支持生成 `curl` 与 `Golang` 代码
* 支持从 Postman 中导入测试用例
* 可观测方便，增加了对 Apache SkyWalking 和 Prometheus 的支持
* 一些 Web 界面操作的优化（例如：多语言、测试结果缓存、自动保存）

本次版本发布，包含了以下 5 位 contributor 的努力：

* [@Ink-33](https://github.com/Ink-33)
* [@LinuxSuRen](https://github.com/LinuxSuRen)
* [@hellojukay](https://github.com/hellojukay)
* [@kuv2707](https://github.com/kuv2707)
* [@yuluo-yx](https://github.com/yuluo-yx)

## 相关数据

下面是 `atest` 截止到 `v0.0.14` 的部分数据：

* watch 7
* fork 23
* star 104
* contributor 12
* 二进制文件下载量 1.1k
* 代码行数 45k
* 单元测试覆盖率 88%

想了解完整信息的话，请访问 https://github.com/LinuxSuRen/api-testing/releases/tag/v0.0.14
