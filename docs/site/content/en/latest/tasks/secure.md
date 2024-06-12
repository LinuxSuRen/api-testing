+++
title = "Secure"
+++

Usually, when TLS certificate authentication is not used, the gRPC client and server communicate in plain text, and the information is easily eavesdropped or tampered by a third party. Therefore, it is recommended to use SSL/TLS to protect gRPC services in most cases. Currently, `atest` has implemented server-side TLS, and mutual TLS (mTLS) needs to wait for implementation.

By default `atest` does not use any security policy, which is equivalent to `spec.secure.insecure = true`. Enabling TLS only requires adding the following content to your yaml:

```yaml
spec:
  secure:
    cert: server.pem
    serverName: atest
```

## Field description

`secure` has the following five fields:

| Name       | Type   | Optional |
| ---------- | ------ | -------- |
| cert       | string | x        |
| ca         | string | √        |
| key        | string | √        |
| serverName | string | x        |
| insecure   | bool   | √        |

`cert` is the path to the certificate that the client needs to configure, in the format of `PEM`.

`serverName` is the service name required by TLS, usually the x509 SAN used when issuing certificates.

`ca` is the path to the CA certificate, and `key` is the private key corresponding to `cert`. After filling in these two items, mTLS is enabled. (mTLS is not implemented yet)

When `insecure` is `false`, `cert` and `serverName` are required.