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

## 数组值检查

### 检查数组中是否有元素的字段包含特定值

示例数据：

```json
{
  "data": [{
    "key": "Content-Type"
  }]
}
```

校验配置：

```yaml
- name: popularHeaders
  request:
    api: /popularHeaders
  expect:
    verify:
      - any(data.data, {.key == "Content-Type"})
```

### 检查数组中是否有元素的字段只包含特定值

校验配置：

```yaml
- name: popularHeaders
  request:
    api: /popularHeaders
  expect:
    verify:
      - all(data.data, {.key == "Content-Type" or .key == "Target"})
```

[更多用法](https://expr-lang.org/docs/language-definition#any).

## 字符串判断

```yaml
- name: metrics
  request:
    api: |
      {{.param.server}}/metrics
  expect:
    verify:
      - indexOf(data, "atest_execution_count") != -1
```

[更多用法](https://expr-lang.org/docs/language-definition#indexOf).
