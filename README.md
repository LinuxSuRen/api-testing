[![Codacy Badge](https://app.codacy.com/project/badge/Grade/3f16717cd6f841118006f12c346e9341)](https://www.codacy.com/gh/LinuxSuRen/api-testing/dashboard?utm_source=github.com\&utm_medium=referral\&utm_content=LinuxSuRen/api-testing\&utm_campaign=Badge_Grade)
[![Codacy Badge](https://app.codacy.com/project/badge/Coverage/3f16717cd6f841118006f12c346e9341)](https://www.codacy.com/gh/LinuxSuRen/api-testing/dashboard?utm_source=github.com\&utm_medium=referral\&utm_content=LinuxSuRen/api-testing\&utm_campaign=Badge_Coverage)
![GitHub All Releases](https://img.shields.io/github/downloads/linuxsuren/api-testing/total)
![](https://sloc.xyz/github/linuxsuren/api-testing)
[![LinuxSuRen/open-source-best-practice](https://img.shields.io/static/v1?label=OSBP&message=%E5%BC%80%E6%BA%90%E6%9C%80%E4%BD%B3%E5%AE%9E%E8%B7%B5&color=blue)](https://github.com/LinuxSuRen/open-source-best-practice)

This is a API testing tool.

## Features

*   Multiple test report formats: Markdown, HTML, Stdout
*   Response Body fields equation check
*   Response Body [eval](https://expr.medv.io/)
*   Verify the Kubernetes resources
*   Validate the response body with [JSON schema](https://json-schema.org/)
*   Pre and post handle with the API request
*   Output reference between TestCase
*   Run in server mode, and provide the [gRPC](pkg/server/server.proto) and HTTP endpoint
*   [VS Code extension](https://github.com/LinuxSuRen/vscode-api-testing) support
*   Multiple storage backends supported(Local, ORM Database, S3, etc)
*   [HTTP API record](extensions/collector)
*   Install in mutiple use cases(CLI, Container, Native-Service, Operator, etc)

## Get started

Install it via [hd](https://github.com/LinuxSuRen/http-downloader) or download from [releases](https://github.com/LinuxSuRen/api-testing/releases):

```shell
hd install atest
```

or, you can install it in Kubernetes. See also the [manifests](sample/kubernetes/manifest.yaml).

see the following usage:

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

below is an example of the usage, and you could see the report as well:

`atest run -p sample/testsuite-gitlab.yaml --duration 1m --thread 3  --report md`

| API | Average | Max | Min | Count | Error |
|---|---|---|---|---|---|
| GET https://gitlab.com/api/v4/projects | 1.152777167s | 2.108680194s | 814.928496ms | 99 | 0 |
| GET https://gitlab.com/api/v4/projects/45088772 | 840.761064ms | 1.487285371s | 492.583066ms | 10 | 0 |
consume: 1m2.153686448s

## Use in Docker

Use `atest` as server mode in Docker, then you could visit the UI from `8080`:
```
docker run --pull always -p 8080:8080 ghcr.io/linuxsuren/api-testing:master
```

Use `atest-collector` in Docker:
```shell
docker run -p 1234:8080 -v /var/tmp:/var/tmp \
  ghcr.io/linuxsuren/api-testing atest-collector \
  --filter-path /api \
  -o /var/tmp/sample.yaml
# you could find the test cases file from /var/tmp/sample
# cat /var/tmp/sample
```

## Template

The following fields are templated with [sprig](http://masterminds.github.io/sprig/):

*   API
*   Request Body
*   Request Header

### Functions

You could use all the common functions which comes from [sprig](http://masterminds.github.io/sprig/). Besides some specific functions are available:

| Name | Usage |
|---|---|
| `randomKubernetesName` | `{{randomKubernetesName}}` to generate Kubernetes resource name randomly, the name will have 8  chars |
| `sleep` | `{{sleep(1)}}` in the pre and post request handle |

## Verify against Kubernetes

It could verify any kinds of Kubernetes resources. Please set the environment variables before using it:

*   `KUBERNETES_SERVER`
*   `KUBERNETES_TOKEN`

See also the [example](sample/kubernetes.yaml).

## TODO

*   Reduce the size of context
*   Support customized context

## Limit

*   Only support to parse the response body when it's a map or array
