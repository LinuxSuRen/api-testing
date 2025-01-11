+++
title = "用例模板"
+++

`atest` 采用 [sprig](https://masterminds.github.io/sprig/) 作为测试用例的模板引擎。通过模板函数可以生成很多随机数据：

## 手机号

下面的代码可以生成 `182` 开头的手机号：

```
182{{shuffle "09876543"}}
```