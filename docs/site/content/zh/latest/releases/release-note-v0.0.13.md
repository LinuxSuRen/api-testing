+++
title = "v0.0.13"
+++

`atest` 版本发布 `v0.0.13`

`atest` 是一款用 Golang 编写的、开源的接口测试工具。

你可以在容器中启动：

```shell
docker run -v /var/www/sample:/var/www/sample \
  --network host \
  linuxsuren/api-testing:master
```

或者，直接[下载二进制文件](https://github.com/LinuxSuRen/api-testing/releases/tag/v0.0.13)后启动：

```shell
atest server --local-storage /var/www/sample
```

对于持续集成（CI）场景，可以通过在流水线中执行命令的方式：

```shell
# 执行本地文件
atest run -p your-test-suite.yaml
# 执行远程文件
atest run -p https://gitee.com/linuxsuren/api-testing/raw/master/sample/testsuite-gitee.yaml
# 容器中执行
docker run linuxsuren/api-testing:master atest run -p https://gitee.com/linuxsuren/api-testing/raw/master/sample/testsuite-gitee.yaml
```

你也可以把测试用例转为 JMeter 文件并执行：

```shell
# 格式转换
atest convert --converter jmeter -p https://gitee.com/linuxsuren/api-testing/raw/master/sample/testsuite-gitee.yaml --target gitee.jmx
# 执行
jmeter -n -t gitee.jmx
```

## 主要的新功能

* 增加了插件扩展机制，支持以 Git、S3、关系型数据为后端存储，支持从 [Vault](https://github.com/hashicorp/vault) 获取密码等敏感信息
* 新增对 gRPC 接口的用例支持 @Ink-33
* 支持导出 [JMeter](https://github.com/apache/jmeter) 文件
* 支持通过 [Operator](https://operatorhub.io/operator/api-testing-operator) 的方式安装，并上架 OperatorHub.io
* 提供了基本的 Web UI
* 支持导出 PDF 格式的测试报告 @wjsvec

本次版本发布，包含了以下 5 位 contributor 的努力：

* [@Ink-33](https://github.com/Ink-33)
* [@LinuxSuRen](https://github.com/LinuxSuRen)
* [@chan158](https://github.com/chan158)
* [@setcy](https://github.com/setcy)
* [@wjsvec](https://github.com/wjsvec)

## 相关数据

下面是 `atest` 截止到 `v0.0.13` 的部分数据：

* watch 7
* fork 18
* star 69
* contributor 8
* 二进制文件下载量 872
* 代码行数 45k
* 单元测试覆盖率 84%

想了解完整信息的话，请访问 https://github.com/LinuxSuRen/api-testing/releases/tag/v0.0.13
