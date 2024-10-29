+++
title = "测试用例验证"
+++

`atest` 采用 https://expr.medv.io 对 HTTP 请求响应的验证，比如：返回的数据列表长度验证、具体值的验证等等。下面给出一些例子：

> 需要注意的是，`data` 指的是 HTTP Response Body（响应体）的 JSON 对象。

## 数组长度判断

```yaml
  - name: projectKinds
    request:
      api: /api/resources/projectKinds
    expect:
      verify:
        - len(data.data)  == 6
```
