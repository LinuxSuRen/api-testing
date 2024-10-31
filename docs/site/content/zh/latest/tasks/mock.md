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

## Mock Docker Registry

您可以通过执行下面的命令 mock 一个容器仓库服务[container registry](https://distribution.github.io/distribution/):

```shell
atest mock --prefix / mock/image-registry.yaml
```

之后，您可以通过使用如下的命令使用 mock 功能。

```shell
docker pull localhost:6060/repo/name:tag
```

## 语法

从整体上来看，我们的写法和 HTTP 的名称基本保持一致，用户无需再了解额外的名词。此外，提供两种描述 Mock 服务的方式：

* 面向对象的 CRUD
* 自定义 HTTP 服务

### 面对对象

```yaml
#!api-testing-mock
# yaml-language-server: $schema=https://linuxsuren.github.io/api-testing/api-testing-mock-schema.json
objects:
  - name: projects
    initCount: 3
    sample: |
      {
        "name": "atest",
        "color": "{{ randEnum "blue" "read" "pink" }}"
      }
```

上面 `projects` 的配置，会自动提供该对象的 CRUD（创建、查找、更新、删除）的 API，你可以通过 `atest` 或类似工具发出 HTTP 请求。例如：

```shell
curl http://localhost:6060/mock/projects

curl http://localhost:6060/mock/projects/atest

curl http://localhost:6060/mock/projects -X POST -d '{"name": "new"}'

curl http://localhost:6060/mock/projects -X PUT -d '{"name": "new", "remark": "this is a project"}'

curl http://localhost:6060/mock/projects/atest -X DELETE
```

> `initCount` 是指按照 `sample` 给定的数据初始化多少个对象；如果没有指定的话，则默认值为 1.

### 自定义

```yaml
#!api-testing-mock
# yaml-language-server: $schema=https://linuxsuren.github.io/api-testing/api-testing-mock-schema.json
items:
  - name: prList
    request:
      path: /v1/repos/{repo}/prs
    response:
      header:
        server: mock
        Content-Type: application/json
      body: |
        {
          "count": 1,
          "items": [{
            "title": "fix: there is a bug on page {{ randEnum "one", "two" }}",
            "number": 123,
            "message": "{{.Response.Header.server}}",
            "author": "someone",
            "status": "success"
          }]
        }
```

启动 Mock 服务后，我们就可以发起如下的请求：

```shell
curl http://localhost:6060/mock/v1/repos/atest/prs -v
```

另外，为了满足复杂的场景，还可以对 Response Body 做特定的解码，目前支持：`base64`、`url`：

```yaml
#!api-testing-mock
# yaml-language-server: $schema=https://linuxsuren.github.io/api-testing/api-testing-mock-schema.json
items:
  - name: base64
    request:
      path: /v1/base64
    response:
      body: aGVsbG8=
      encoder: base64
```

上面 Body 的内容是经过 `base64` 编码的，这可以用于不希望直接明文显示，或者是图片的场景：

```shell
curl http://localhost:6060/mock/v1/base64
```

如果你的 Body 内容可以通过另外一个 HTTP 请求（GET）获得，那么你可以这么写：

```
#!api-testing-mock
# yaml-language-server: $schema=https://linuxsuren.github.io/api-testing/api-testing-mock-schema.json
items:
  - name: baidu
    request:
      path: /v1/baidu
    response:
      body: https://baidu.com
      encoder: url
```
