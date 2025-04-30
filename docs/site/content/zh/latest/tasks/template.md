+++
title = "用例模板"
weight = 100
+++

`atest` 采用 [sprig](https://masterminds.github.io/sprig/) 作为测试用例的模板引擎。通过模板函数可以生成很多随机数据：

## 手机号

下面的代码可以生成 `182` 开头的手机号：

```
182{{shuffle "09876543"}}
```

## 带权重的随机枚举

下面的代码以 80% 的概率返回 `open`，以 20% 的概率返回 `closed`：

```
{{randWeightEnum (weightObject 4 "open") (weightObject 1 "closed")}}
```

## 时间

下面的代码可以生成当前时间，并制定时间格式：

```
{{ now.Format "2006-01-02T15:04:05Z07:00" }}
```

如果想要获取其他时间，则可以参考如下写法：

```gotemplate
{{(now | date_modify "2h") | date "2006-01-02T15:04:05"}}
```

## 环境变量

下面的代码可以获取环境变量 `SHELL` 的值，在需要使用一个全局变量的时候，可以使用这个模板函数：

```
{{ env "SHELL" }}
```
