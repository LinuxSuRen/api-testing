[![CLA assistant](https://cla-assistant.io/readme/badge/LinuxSuRen/api-testing)](https://cla-assistant.io/LinuxSuRen/api-testing)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/3f16717cd6f841118006f12c346e9341)](https://app.codacy.com/gh/LinuxSuRen/api-testing/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade)
[![Codacy Badge](https://app.codacy.com/project/badge/Coverage/3f16717cd6f841118006f12c346e9341)](https://app.codacy.com/gh/LinuxSuRen/api-testing/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade)
[![GitHub All Releases](https://img.shields.io/github/downloads/linuxsuren/api-testing/total)](https://tooomm.github.io/github-release-stats/?username=linuxsuren&repository=api-testing)
[![Docker Pulls](https://img.shields.io/docker/pulls/linuxsuren/api-testing)](https://hub.docker.com/r/linuxsuren/api-testing)
[![LinuxSuRen/open-source-best-practice](https://img.shields.io/static/v1?label=OSBP&message=%E5%BC%80%E6%BA%90%E6%9C%80%E4%BD%B3%E5%AE%9E%E8%B7%B5&color=blue)](https://github.com/LinuxSuRen/open-source-best-practice)
![GitHub Created At](https://img.shields.io/github/created-at/linuxsuren/api-testing)

> 中文 | [English](README.md)

一个开源的 API 测试工具。🚀

## 功能特性

* 支持的协议: HTTP, gRPC, tRPC
* 支持多种格式的测试结果导出: Markdown, HTML, PDF, Stdout
* 简单易用的 Mock 服务，支持 OpenAPI
* 支持转换为 [JMeter](https://jmeter.apache.org/) 文件格式
* 支持响应体字段检查或 [eval](https://expr.medv.io/)
* 使用 [JSON schema] 校验响应参数(https://json-schema.org/)
* 支持预处理和后处理 API 请求
* 支持以服务器模式运行并支持 [gRPC](pkg/server/server.proto) 和 HTTP endpoint
* [VS Code 扩展支持](https://github.com/LinuxSuRen/vscode-api-testing)
* [Github 扩展支持](https://github.com/marketplace/actions/api-testing-with-kubernetes)
* 支持多种存储方式 (Local, ORM Database, S3, Git, Etcd, etc.)
* [HTTP API record](https://github.com/LinuxSuRen/atest-ext-collector)
* 支持多种安装方式(CLI, Container, Native-Service, Operator, Helm, etc.)
* 整合 Prometheus, SkyWalking 监控

## 快速开始

[![Try in PWD](https://github.com/play-with-docker/stacks/raw/cff22438cb4195ace27f9b15784bbb497047afa7/assets/images/button.png)](http://play-with-docker.com?stack=https://raw.githubusercontent.com/LinuxSuRen/api-testing/master/docs/manifests/docker-compose.yml)

通过 [hd](https://github.com/LinuxSuRen/http-downloader) 安装，或从 [releases](https://github.com/LinuxSuRen/api-testing/releases) 下载安装:

```shell
hd install atest
```

您也可以通过 kubernetes 安装，更多细节请参考： [manifests](docs/manifests/kubernetes/default/manifest.yaml).

用法如下：

```shell
API testing tool

Usage:
  atest [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  func        Print all the supported functions
  help        Help about any command
  json        Print the JSON schema of the test suites struct
  run         Run the test suite
  sample      Generate a sample test case YAML file
  server      Run as a server mode
  service     Install atest as a Linux service

Flags:
  -h, --help      help for atest
  -v, --version   version for atest

Use "atest [command] --help" for more information about a command.
```

API Testing 使用示例，在此示例中，您将通过 md 格式阅览生成的接口测试报告：

`atest run -p sample/testsuite-gitlab.yaml --duration 1m --thread 3  --report md`

| API | Average | Max | Min | Count | Error |
|---|---|---|---|---|---|
| GET https://gitlab.com/api/v4/projects | 1.152777167s | 2.108680194s | 814.928496ms | 99 | 0 |
| GET https://gitlab.com/api/v4/projects/45088772 | 840.761064ms | 1.487285371s | 492.583066ms | 10 | 0 |
consume: 1m2.153686448s

## 在 Docker 中使用

在 Docker 中以服务器模式运行 `atest`，您可以通过 `8080` 访问 `atest` 的 UI 控制台：

```bash
docker run --pull always -p 8080:8080 ghcr.io/linuxsuren/api-testing:master
```

在 Docker 中使用 `atest-collector`:

```shell
docker run -p 1234:8080 -v /var/tmp:/var/tmp \
  ghcr.io/linuxsuren/api-testing atest-collector \
  --filter-path /api \
  -o /var/tmp/sample.yaml
# you could find the test cases file from /var/tmp/sample
# cat /var/tmp/sample
```

## 模板

以下字段的模板配置参考：[sprig](http://masterminds.github.io/sprig/):

* API
* Request Body
* Request Header

### Functions

您可以使用 [sprig](http://masterminds.github.io/sprig/) 中的所有常用函数。此外，还有一些特殊函数可以在 `atest` 中使用：

| Name | Usage |
|---|---|
| `randomKubernetesName` | `{{randomKubernetesName}}` to generate Kubernetes resource name randomly, the name will have 8  chars |
| `sleep` | `{{sleep(1)}}` in the pre and post request handle |

## 验证 Kuberntes 资源

`atest` 可以验证任何类型的 Kubernetes 资源。使用前请先设置 Kubernetes 相关的环境变量：

* `KUBERNETES_SERVER`
* `KUBERNETES_TOKEN`

另请参考 [example](sample/kubernetes.yaml)。

## 待办事项

* 减少上下文的大小
* 支持自定义上下文

## 功能限制

* 仅支持解析 map 或 array 类型的响应体。

## 社区交流

欢迎使用以下联系方式，探讨有关 API Testing 的任何问题！

### 邮件列表

`api-testing-tech@googlegroups.com`, 欢迎通过此邮件列表讨论与 API Testing 相关的任何问题。

### GitHub Discussion

[GitHub Discussion](https://github.com/LinuxSuRen/api-testing/discussions/new/choose)
