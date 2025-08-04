+++
title = "Mock 服务"
weight = 104
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

### 面向对象

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
      path: /api/v1/repos/{repo}/prs
    response:
      header:
        server: mock
        Content-Type: application/json
      body: |
        {
          "count": 1,
          "repo": "{{.Param.repo}}",
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
curl http://localhost:6060/mock/api/v1/repos/atest/prs -v
```

另外，为了满足复杂的场景，还可以对 Response Body 做特定的解码，目前支持：`base64`、`url`、`raw`：

> encoder 为 `raw` 时，表示不进行处理

```yaml
#!api-testing-mock
# yaml-language-server: $schema=https://linuxsuren.github.io/api-testing/api-testing-mock-schema.json
items:
  - name: base64
    request:
      path: /api/v1/base64
    response:
      body: aGVsbG8=
      encoder: base64
```

上面 Body 的内容是经过 `base64` 编码的，这可以用于不希望直接明文显示，或者是图片的场景：

```shell
curl http://localhost:6060/mock/api/v1/base64
```

如果你的 Body 内容可以通过另外一个 HTTP 请求（GET）获得，那么你可以这么写：

```yaml
#!api-testing-mock
# yaml-language-server: $schema=https://linuxsuren.github.io/api-testing/api-testing-mock-schema.json
items:
  - name: baidu
    request:
      path: /api/v1/baidu
    response:
      body: https://baidu.com
      encoder: url
```

如果你的响应内容比较大，或者保存在一个本地文件中，那么你可以这么写：

```yaml
#!api-testing-mock
# yaml-language-server: $schema=https://linuxsuren.github.io/api-testing/api-testing-mock-schema.json
items:
  - name: baidu
    request:
      path: /api/v1/baidu
    response:
      bodyFromFile: /tmp/baidu.html
```

在实际情况中，往往是向已有系统或平台添加新的 API，此时要 Mock 所有已经存在的 API 就既没必要也需要很多工作量。因此，我们提供了一种简单的方式，即可以增加**代理**的方式把已有的 API 请求转发到实际的地址，只对新增的 API 进行 Mock 处理。如下所示：

```yaml
#!api-testing-mock
# yaml-language-server: $schema=https://linuxsuren.github.io/api-testing/api-testing-mock-schema.json
proxies:
  - path: /api/v1/{part}
    target: http://atest.localhost:8080
```

当我们发起如下的请求时，实际请求的地址为 `http://atest.localhost:8080/api/v1/projects`

```shell
curl http://localhost:6060/mock/api/v1/projects
```

如何希望把所有的请求都转发到某个地址，则可以使用通配符的方式：

```yaml
proxies:
  - path: /{path:.*}
    target: http://192.168.123.58:9200
```

## TCP 协议代理

```yaml
proxies:
  - protocol: tcp
    port: 3306
    path: /
    target: 192.168.123.58:33060
```

## 代理多个服务

```shell
atest mock-compose bin/compose.yaml
```

执行上面的命令，会启动多个 Mock 代理服务，分别以不同的端口代理了 Elasticsearch 和 Eureka 服务：

```yaml
proxies:
  - prefix: /
    port: 9200
    path: /{path:.*}
    target: http://192.168.123.58:9200
  - prefix: /
    port: 17001
    path: /{path:.*}
    target: http://192.168.123.58:17001
  - protocol: tcp
    port: 33060
    path: /
    target: 192.168.123.58:33060
```

当前代理支持 HTTP 和 TCP 协议，上面的例子中代理了 MySQL 的 `33060` 端口。

## Webhook

有些场景下，需要定时向服务器发送请求，这时可以使用 Webhook。当前支持的协议包括：

* HTTP
* Syslog

> 更多 URL 中通配符的用法，请参考 https://github.com/gorilla/mux
