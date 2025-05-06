+++
title = "关于安全"
weight = 8
+++

通常在不使用 TLS 证书认证时，gRPC 客户端与服务端之间进行的是明文通信，信息易被第三方监听或篡改。所以大多数情况下均推荐使用 SSL/TLS 保护 gRPC 服务。目前`atest`已实现服务端 TLS，双向 TLS(mTLS) 需等待后续实现。

默认情况下`atest`不使用任何安全策略，等价于`spec.secure.insecure = true`。启用 TLS 仅需 yaml 中添加以下内容：

```yaml
spec:
  secure:
    cert: server.pem
    serverName: atest
```

## 字段说明

`secure`共有以下五个字段：

| 字段名     | 类型   | 是否可选 |
| ---------- | ------ | -------- |
| cert       | string | x        |
| ca         | string | √        |
| key        | string | √        |
| serverName | string | x        |
| insecure   | bool   | √        |

`cert`为客户端需要配置的证书的文件路径，格式为`PEM`。

`serverName`为 TLS 所需的服务名，通常为签发证书时使用的 x509 SAN。

`ca`为 CA 证书的路径，`key`为与`cert`对应的私钥，这两项填写后代表启用 mTLS。(mTLS 尚未实现)

当`insecure`为`false`时，`cert`和`serverName`为必填项。
