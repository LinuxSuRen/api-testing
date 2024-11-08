---
title: "快速入门"
weight: 1
description: 只需几个简单的步骤即可开始使用 API Testing。
---

本指南将帮助您通过几个简单的步骤开始使用 API Testing。

## 执行部分测试用例

下面的命令会执行名称中包含 `sbom` 的所有测试用例：

```shell
atest run -p test-suite.yaml --case-filter sbom
```
