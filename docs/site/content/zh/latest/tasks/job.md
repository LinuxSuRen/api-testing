+++
title = "任务"
weight = 101
+++

在特定情况下，执行接口测试用例前需要执行相应的任务，例如：数据初始化、等待服务就绪等等。 `atest` 的任务功能，就是为了满足这类场景而设计的。

> 任务的执行引擎为 [expr](https://expr.medv.io)，如果当前页面给出的示例无法满足你的需求，可以查阅相关的官方文档。

## 等待接口响应码为 200

以下的例子会每隔一秒请求一次指定的接口，并检查响应码（Status Code）是否为 200，如果不是的话，则最多会重试 3 次：

```yaml
name: demo
api: http://localhost
items:
  - name: login
    before:
      items:
        - httpReady("http://localhost/health", 3)
    request:
      api: /demo
```

如果希望检查响应体的话，则可以用下面的表达式：

```
httpReady("http://localhost:17001/actuator/health", 3, 'components.discoveryComposite.components.eureka.details.applications["AUTH"] == 1')
```
