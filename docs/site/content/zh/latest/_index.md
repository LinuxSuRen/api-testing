+++
title = "欢迎访问 API Testing"
linktitle = "文档"
description = "API Testing 文档"

[[cascade]]
type = "docs"
+++

{{% alert title="记录" color="primary" %}}

该项目正在**积极**开发中，很多功能尚待补充，我们希望您[参与其中](contributions/)！

{{% /alert %}}

API Testing 一个基于 YAML 文件的开源接口测试工具，同时支持运行在本地、服务端。

在选择工具时，可以从很多方面进行考量、对比，以下几点是该工具的特色或者优点：

* 开源与否，atest 采用 MIT 开源协议，是最流行的宽松开源协议之一。有些工具也许有非常丰富的功能、漂亮的界面，但相比于开源项目，免费的工具不定什么时候就有可能变为收费的；而且，你的使用感受几乎很难直接反馈到产品中，只能被动接受。
* 质量、可靠性，作为一款用于测试场景的工具，atest 本身的单元测试覆盖率达 89%，单测代码与业务逻辑代码量平分秋色；另外，每次代码改动都需要通过代码扫描、单元测试等流水线。
* 身材小巧，整个工具大小为 18M，支持 Windows、Linux、macOS 平台。
* 只有简单的可执行二进制文件，不像部分工具会给你的操作系统安装莫名其妙的系统启动项目、系统服务等。
* 基于 YAML 文件，提交到 Git 仓库后，天生支持团队协作，无需注册额外账号。
* 同时提供简单、高级两种模式的返回值断言，还包括 JSON Schema 以及针对 Kubernetes 资源的校验判断。
* 支持性能测试。
* 直接在 VS Code 中直接触发执行单个或整个测试文件。

## 如何使用？

那么，这个工具长什么样子呢，下面是命令行 `atest` 的参数说明：

```shell
API testing tool

Usage:
  atest [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  json        Print the JSON schema of the test suites struct
  run         Run the test suite
  sample      Generate a sample test case YAML file
  server      Run as a server mode
  service     Install atest as a Linux service

Flags:
  -h, --help   help for atest

Use "atest [command] --help" for more information about a command.
```

### 本地模式

执行一个测试用例集文件：`atest run -p sample/testsuite-gitlab.yaml`，其中的参数 `-p` 支持模糊匹配多个文件。

如果希望对测试用例集执行性能测试的话，可以增加响应的参数：

`atest run -p sample/testsuite-gitlab.yaml --duration 1m --thread 3  --report md`

其中的参数 `--report` 可以指定性能测试输出报告，目前支持 Markdown 以及控制台输出。效果如下所示：

| API | Average | Max | Min | Count | Error |
|---|---|---|---|---|---|
| GET https://gitlab.com/api/v4/projects | 1.152777167s | 2.108680194s | 814.928496ms | 99 | 0 |
| GET https://gitlab.com/api/v4/projects/45088772 | 840.761064ms | 1.487285371s | 492.583066ms | 10 | 0 |
consume: 1m2.153686448s

### 服务端模式

除了本地执行外，`atest` 还提供了基于 `gRPC` 协议服务端，通过下面的命令即可启动：

```shell
atest server
```

对于 Linux 操作系统，用户还可以通过下面的命令安装后台服务：

```shell
atest service (install | start | stop | restart)
```

当然，如果你对容器、Kubernetes 比较熟悉的话，本项目也提供了对应的支持。

这种模式，对于想要集成的用户而言，可以通过调用 `gRPC` 来执行测试。也可以安装 [VS Code](https://marketplace.visualstudio.com/items?itemName=linuxsuren.api-testing) 插件，在编码与接口测试之间无缝切换，您可以搜索 `api-testing` 找到该插件。

插件会识别所有第一行是 `#!api-testing` 的 YAML 文件，并提供快速的执行操作，请参考如下截图：

![](atest-vscode.png)

如图所示，会有四个快捷执行操作：

* `run suite` 会执行整个文件
* `run suite with env` 会加载 `env.yaml` 文件并执行整个文件
* `run` 执行单个测试用例（包括所依赖的用例）
* `debug` 执行单个测试用例，并输出接口返回值

当你安装了 VS Code 插件后，会自动下载并安装 `atest` 及其服务。当然，你也可以配置不同的远端服务地址。

### 文件格式

`atest` 定义的 YAML 格式，基本遵循 HTTP 的语义，熟悉 HTTP 协议的同学即可快速上手。下面是一个范例，更多例子[请参考这里](https://github.com/LinuxSuRen/api-testing/blob/master/sample/)：


```yaml
#!api-testing
name: Kubernetes
api: https://192.168.123.121:6443
items:
- name: pods
  request:
    api: /api/v1/namespaces/kube-system/pods
    header:
      Authorization: Bearer token
  expect:
    verify:
    - pod("kube-system", "kube-ovn-cni-55bz9").Exist()
    - k8s("deployments", "kube-system", "coredns").Exist()
    - k8s("deployments", "kube-system", "coredns").ExpectField(2, "spec", "replicas")
    - k8s({"kind":"virtualmachines","group":"kubevirt.io"}, "vm-test", "vm-win10-dkkhl").Exist()
```

用户可以自定义请求的 Header、Payload 等，可以对响应体做全面的断言判断。

## 后续计划

如果您已经耐心阅读到这里的话，可以再顺便了解下这个项目后续的一些想法。

通过更多的实际场景来打磨、优化 `atest` 对接口测试的便利性、可扩展性，以不丢失易用性为前提增强功能。例如：

* 优化错误提示、反馈
* 提供与 CICD 集成的最佳实践
* 增加 gRPC 等协议的支持
* 增加测试记录信息的持久化
* VS Code 插件支持测试用例编写的提示、格式校验
* 提供插件机制，增加对数据库等数据源的格式校验

最后期待您的反馈 https://github.com/LinuxSuRen/api-testing/issues

## 准备好开始了吗？ {#ready-to-get-started}
