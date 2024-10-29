+++
title = "gRPC测试用例编写指南"
+++

本文档将介绍如何编写`api-testing`的 gRPC API 的测试用例。

阅读本文档之前，您需要先安装并配置好`api-testing`，具体操作可以参考[安装](../install/_index.md)章节。如果您已经完成了这些步骤，可以继续阅读本文档的后续部分。

## 创建测试项目

创建一个基于服务反射的 gRPC 测试用例仅需在 yaml 文件的`spec`路径下加入以下内容：

```yaml
spec:
  rpc:
    serverReflection: true
```

`rpc`字段一共有五个子字段

| 字段名           | 类型     | 是否可选 |
| ---------------- | -------- | -------- |
| import           | []string | √        |
| protofile        | string   | √        |
| protoset         | string   | √        |
| serverReflection | bool     | √        |

### 字段`import`和`protofile`

`protofile`是一个文件路径，指向`api-testing`查找描述符的`.proto`文件的位置。

`import`字段与`protc`编译器中的`--import_path`参数类似，用于确定查找`proto`文件的位置与解析依赖的目录。与`protoc`一样，您不需要在此处指定某些`proto`文件的位置（以`google.protobuf`开头的`Protocol Buffers Well-Known Types`），它们已经内置在了`api-testing`中。

### 字段`protoset`

`protoset`字段既可以是一个文件路径，也可以是`http(s)://`开头的网络地址。

当您的`proto`数量繁多或引用关系复杂，可以尝试使用`protoc --descriptor_set_out=set.pb`生成`proto`描述符集合。本质上它是一个使用了`wire`编码的二进制文件，其中包括了所有需要的描述符。

### 字段`serverReflection`

若目标服务器支持服务反射，将此项设为`true`则不再需要提供上述三个字段。

---
注：`api-testing`对三种描述符来源的优先级顺序为

`serverReflection` > `protoset` > `protofile`

## 编写gRPC API测试

与编写`HTTP`测试用例类型，您需要在根节点的`api`字段定义服务器的地址。

```yaml
api: 127.0.0.1:7070
```

默认情况下`api-testing`使用不安全的方式连接到目标服务器。若您想配置TLS证书，请参考文档[关于安全](./secure-zh.md)

---

编写`gRPC API`测试的方式与编写`HTTP API`测试的方式基本相同。

```yaml
- name: FunctionsQuery
  request:
    api: /server.Runner/FunctionsQuery
    body: |
      {
        "name": "hello"
      }
  expect:
    body: |
      {
        "data": [
          {
            "key": "hello",
            "value": "func() string"
          }
        ]
      }
```

`api`字段的格式为`/package.service/method`，支持 gRPC 一元调用、客户端流、服务端流和双向流调用。

与`api`字段同级的`body`字段是以`JSON`格式表示的`Protocol Buffers`消息，代表将要调用的`api`的入参。特别的，当您需要调用客户端流或双向流 API 时，请使用`JSON Array`格式编写字段`body`，如：

```yaml
body: |
  [
    {
      "name": "hello"
    },
    {
      "name": "title"
    }
  ]
```

## 编写返回内容验证

编写`gRPC API`返回内容验证的方式与`HTTP API`基本相同。对与`gRPC API`来说，一切返回值都被视为`map`类型，被放入`api testing`特定的返回结构中：

```yaml
expect:
  body: |
    {
      "data": [
        {
          "key": "hello",
          "value": "func() string"
        }
      ]
    }
```

`api-testing`为`JSON`比对编写了一套对比库，请参考此处[compare](https://pkg.go.dev/github.com/linuxsuren/api-testing/pkg/compare)

请注意，对于服务端流和双向流模式，服务器发送多条消息的情况下，此处的`data`字段内的填写的目标数组，需同时满足与待验证数组长度相，两个数组同一下标的内容完全相同。

`gRPC API`的`verify`功能与`HTTP API`保持一致，此处不再赘述。