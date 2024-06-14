+++
title = "gRPC testsuite writing manual"
+++

This document will introduce how to write testsuite for the gRPC API of `api-testing`.

Before reading this document, you need to install and configure `api-testing`. For specific operations, please refer to [Install Document](../install/_index.md). If you have completed these steps, you can continue reading the rest of this document.

## Create testsuite

To create a gRPC testsuite based on service reflection, just add the following content to the `spec` path of the yaml file:

```yaml
spec:
  rpc:
    serverReflection: true
```

Field `rpc` has five subfields in total:

| Name             | Type     | Optional |
| ---------------- | -------- | -------- |
| import           | []string | √        |
| protofile        | string   | √        |
| protoset         | string   | √        |
| serverReflection | bool     | √        |


### Field `import` and `protofile`

`protofile` is a file path pointing to the location of the `.proto` file where `api-testing` looks for descriptors.

The `import` field is similar to the `--import_path` parameter in the `protc` compiler, and is used to determine the location of the `proto` file and the directory for parsing dependencies. Like `protoc`, you don't need to specify the location of certain `proto` files here (such as `Protocol Buffers Well-Known Types` starting with `google.protobuf`), they are already built into `api-testing` binary file.

### Field `protoset`

Field `protoset` can be either a file path or a network address starting with `http(s)://`.

When you have a large number of `proto` or complex dependencies, you can try to use `protoc --descriptor_set_out=set.pb` to generate a `proto descriptor set`. Essentially it is a wire-encoded binary file that includes all the required descriptors.


### Field `serverReflection`

If the target server supports service reflection, setting this to `true` will no longer need to provide the above three fields.

---
Note: The priority order of `api-testing` for the three descriptor sources is

`serverReflection` > `protoset` > `protofile`

## Write gRPC API testsuite

Like writing the `HTTP` testsuite, you need to define the address of the server in the `api` field of the root node.

```yaml
api: 127.0.0.1:7070
```

By default, `api-testing` uses an insecure way to connect to the target server. If you want to configure a TLS certificate, please refer to the document [About Security](./secure.md)

---

Writing testsuite for `gRPC API` is basically the same as writing testsuite for `HTTP API`.

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

The format of the `api` field is `/package.service/method`, which supports gRPC unary calls, client streams, server streams and bidirectional stream calls.

The `body` field at the same level as the `api` field is a `Protocol Buffers` message expressed in `JSON` format, representing the input parameters of the `api` to be called. Especially, when you need to call the client stream or bidirectional stream API, please use the `JSON Array` format to write the field `body`, such as:

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

## Write return content verification

Writing `gRPC API` to return content validation is basically the same as `HTTP API`. For the `gRPC API`, all return values are treated as `map` types and put into the `api testing` specific return structure: 

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

`api-testing` has written a comparison lib for `JSON` comparison, please refer to here [compare](https://pkg.go.dev/github.com/linuxsuren/api-testing/pkg/compare)

Please note that for server-side streaming and bi-directional streaming modes where the server sends multiple messages, the target array in the `data` field must be the same length as the array to be validated, and both arrays must have the same contents under the same index.

The `verify` functionality of the `gRPC API` is consistent with the `HTTP API` and will not be repeated here.