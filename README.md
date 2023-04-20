[![Codacy Badge](https://app.codacy.com/project/badge/Grade/3f16717cd6f841118006f12c346e9341)](https://www.codacy.com/gh/LinuxSuRen/api-testing/dashboard?utm_source=github.com\&utm_medium=referral\&utm_content=LinuxSuRen/api-testing\&utm_campaign=Badge_Grade)
[![Codacy Badge](https://app.codacy.com/project/badge/Coverage/3f16717cd6f841118006f12c346e9341)](https://www.codacy.com/gh/LinuxSuRen/api-testing/dashboard?utm_source=github.com\&utm_medium=referral\&utm_content=LinuxSuRen/api-testing\&utm_campaign=Badge_Coverage)
![GitHub All Releases](https://img.shields.io/github/downloads/linuxsuren/api-testing/total)

This is a API testing tool.

## Feature

*   Response Body fields equation check
*   Response Body [eval](https://expr.medv.io/)
*   Verify the Kubernetes resources
*   Validate the response body with [JSON schema](https://json-schema.org/)
*   Output reference between TestCase
*   Run in server mode, and provide the gRPC endpoint
*   [VS Code extension](https://github.com/LinuxSuRen/vscode-api-testing) support

## Get started

Install it via [hd](https://github.com/LinuxSuRen/http-downloader) or download from [releases](https://github.com/LinuxSuRen/api-testing/releases):

```shell
hd install atest
```

see the following usage:

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

## Template

The following fields are templated with [sprig](http://masterminds.github.io/sprig/):

*   API
*   Request Body
*   Request Header

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
