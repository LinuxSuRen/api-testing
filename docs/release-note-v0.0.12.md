`atest` 版本发布 v0.0.12

`atest` 是一款用 Golang 编写的、基于 YAML 格式的开源接口测试工具，可以方便地在本地、服务端、持续集成等场景中使用。
我们希望提供一个简单、强大、高质量的测试工具，方便测试、研发人员快速、低成本地借助接口测试为产品研发质量保驾护航。

通过以下命令启动 HTTP 代理服务器后，给您的浏览器配置该代理，打开业务系统就会自动录制：

```shell
docker run -p 1234:8080 -v /var/tmp:/var/tmp \
  ghcr.io/linuxsuren/api-testing atest-collector \
  --filter-path /api \
  -o /var/tmp/sample.yaml
# --filter-path /api 会过滤所有以 /api 为前缀的 HTTP 请求
# 关闭服务后，您可以在 /var/tmp/sample 这个目录中找到生成的测试用文件
```

## 更新重点

* 支持通过基于 HTTP 代理服务生成测试用例
* 支持根据 Swagger 数据生成接口测试覆盖率
* 增加 HTML、Markdown 等格式的测试报告
* 代码重构，包括：包结构、原文件名整理，逻辑抽象为接口以及不同实现
* 支持打印所有支持的模板函数
* 优化 Kubernetes 的部署清单文件
* 修复已知缺陷

本次版本发布，包含了以下三位 contributor 的努力：

* [@LinuxSuRen](https://github.com/LinuxSuRen)
* [@wongearl](https://github.com/wongearl)
* [@yJunS](https://github.com/yJunS)

## 相关数据

下面是 `atest` 截止到 v0.0.12 的部分数据：

* watch 3
* fork 9
* star 33
* contributor 4
* 二进制文件下载量 561
* 代码行数 7.6k
* 单元测试覆盖率 94%

想了解完整信息的话，请访问 https://github.com/LinuxSuRen/api-testing/releases/tag/v0.0.12
