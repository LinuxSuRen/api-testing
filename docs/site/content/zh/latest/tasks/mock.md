+++
title = "Mock 服务"
+++

Mock 服务在前后端并行开发、系统对接、设备对接场景下能起到非常好的作用，可以极大地降低团队之间、系统之间的耦合度。

用户可以通过命令行终端（CLI）、Web UI 的方式来使用 Mock 服务。

## 命令行

```shell
atest mock --prefix / --port 9090 mock.yaml
```

## Web

在 UI 上可以实现和命令行相同的功能，并可以通过页面编辑的方式修改、加载 Mock 服务配置。

## 语法

从整体上来看，我们的写法和 HTTP 的名称基本保持一致，用户无需再了解额外的名词。此外，提供两种描述 Mock 服务的方式：

* 针对某个数据对象的 CRUD
* 任意 HTTP 服务

下面是一个具体的例子：

```yaml
#!api-testing-mock
# yaml-language-server: $schema=https://linuxsuren.github.io/api-testing/api-testing-mock-schema.json
objects:
  - name: repo
    fields:
      - name: name
        kind: string
      - name: url
        kind: string
  - name: projects
    initCount: 3
    sample: |
      {
        "name": "api-testing",
        "color": "{{ randEnum "blue" "read" "pink" }}"
      }
items:
  - name: base64
    request:
      path: /v1/base64
    response:
      body: aGVsbG8=
      encoder: base64
  - name: prList
    request:
      path: /v1/repos/{repo}/prs
      header:
        name: rick
    response:
      header:
        server: mock
      body: |
        {
          "count": 1,
          "items": [{
            "title": "fix: there is a bug on page {{ randEnum "one" }}",
            "number": 123,
            "message": "{{.Response.Header.server}}",
            "author": "someone",
            "status": "success"
          }]
        }
```
