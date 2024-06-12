+++
title = "v0.0.17"
+++

`atest` 发布 `v0.0.17`

`atest` 是致力于帮助开发者持续保持高质量 API 的开源接口工具。

你可以在命令行终端或者容器中启动：

```shell
docker run -p 8080:8080 ghcr.io/linuxsuren/api-testing:v0.0.17
```

## 亮点

* 我们提供了基于 Electron 的桌面应用，会极大地方便开发者在桌面环境中测试 API。
* 为缩减镜像的体积（40M），我们把插件全部以 OCI 的格式单独存储，并在启用时自动下载。
* 诞生了第二位项目 Committer [@yuluo-yx](https://github.com/LinuxSuRen/api-testing/discussions/479)

非常期待 `atest` 可以帮助更多的项目持续提升、保持 API 稳定性。

## 🚀 主要的新功能

* 支持通过 HTTP 请求执行测试套件 (#478) @LinuxSuRen
* 增加 gRPC 接口对 TLS 的支持 (#477) @DWJ-Squirtle
* 支持自动下载插件 (#471) @LinuxSuRen
* 补充代码生成器的 e2e 测试 (#458) @LinuxSuRen
* 支持复制测试用例和测试套件 (#455) @LinuxSuRen
* Web 界面上添加切换语言的按钮 (#447) @SamYSF
* 支持通过 Web 界面查看 YAML 格式的测试套件 (#438) @SamYSF
* 支持发送测试报告到 gRPC 服务 (#431) @lizzy-0323
* 支持发送测试报告到 HTTP 服务 (#367) @hahahashen
* 增加基于 Electron 的桌面应用 (#428) @LinuxSuRen
* 实现了镜像 Registry 的 Mock 服务 (#425) @LinuxSuRen
* 支持在 Web 界面启动、刷新 Mock 服务 (#410) @LinuxSuRen
* 支持根据测试用例生成 JavaScript 代码 (#400) @YukiCoco
* 支持根据测试用例生成 Python 代码 (#398) @zhouzhou1017
* 支持根据测试用例生成 Java 代码 (#369) @Agility6
* 增加日志框架的支持 (#389) @yuluo-yx
* 生成 Golang 代码时支持 Cookie 的设置 (#363) @SLOWDOWNO
* 测试用例支持 Cookie 设置 (#355) @LinuxSuRen

## 🐛 缺陷修复

* 解决测试用例页面徽章显示的问题 (#462) @SamYSF
* 解决无法导入 Postman 子集的问题 (#426) @SamYSF
* 优化 gRPC 消息超过默认值的处理 (#399) @acceleratorssr
* 解决 golang.org/x/net 的安全漏洞 CVE-2023-45288 (#401) @yuluo-yx
* 修复生成 Golang 代码时对 HTTP 请求体的设置 (#383) @Agility6

## 📝 文档

* 增加行为准则说明 (#379) @yuluo-yx
* 增加安全漏洞相关的说明 (#391) @yuluo-yx
* 更新贡献文档说明 (#380) @yuluo-yx

## 👻 维护

* 用 openapi 官方的依赖库替换当前实现 (#439) @dshyjtdes8888
* 增加 issue comment github actions (#382) @yuluo-yx

## 致谢

本次版本发布，包含了以下 13 位 contributor 的努力：

* [@Agility6](https://github.com/Agility6)
* [@DWJ-Squirtle](https://github.com/DWJ-Squirtle)
* [@LinuxSuRen](https://github.com/LinuxSuRen)
* [@SLOWDOWNO](https://github.com/SLOWDOWNO)
* [@SamYSF](https://github.com/SamYSF)
* [@YukiCoco](https://github.com/YukiCoco)
* [@acceleratorssr](https://github.com/acceleratorssr)
* [@dshyjtdes8888](https://github.com/dshyjtdes8888)
* [@hahahashen](https://github.com/hahahashen)
* [@lizzy-0323](https://github.com/lizzy-0323)
* [@wt-goodluck](https://github.com/wt-goodluck)
* [@yuluo-yx](https://github.com/yuluo-yx)
* [@zhouzhou1017](https://github.com/zhouzhou1017)

## 相关数据

下面是 `atest` 截止到 `v0.0.17` 的部分数据：

* watch 8
* fork 47
* star 209 (+86)
* contributor 24 (+11)
* 二进制文件下载量 3.1k (+1.8k)
* 部分镜像 5.5k (+3.3k)
* 单元测试覆盖率 74% (-8%)

想了解完整信息的话，请访问 https://github.com/LinuxSuRen/api-testing/releases/tag/v0.0.17
