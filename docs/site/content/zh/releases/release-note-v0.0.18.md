+++
title = "v0.0.18"
+++

`atest` 发布 `v0.0.18`

`atest` 是致力于帮助开发者持续保持高质量 API 的开源接口工具。

你可以在命令行终端或者容器中启动：

```shell
docker run -p 8080:8080 ghcr.io/linuxsuren/api-testing:v0.0.18
```

## 亮点

* 在开源之夏 2024 中 `atest` 增加了基于 MySQL 的测试用例历史的支持
* HTTP API Mock 功能的支持

在系统和平台的开发过程中，我们通常会采用前后端分离的开发模式。在后端API尚未开发完成、稳定化，并且未部署到公共集成测试环境之前，前端开发者往往需要通过硬编码数据来推进页面开发。待后端开发完成后，会进入所谓的“联调”阶段，这时可能会遇到以下问题：

* 前端可能需要调整数据结构、页面布局和逻辑，并重新进行测试
* 在实际查看页面后，可能会发现后端的数据结构和API的请求与响应需要调整

在最坏的情况下，前后端的联调可能会耗费远超预期的时间。为了更有效地解决这一问题，`atest` 提供了HTTP API Mock功能。

在设计评审阶段，我们可以根据API设计提供相应的Mock服务配置，从而快速模拟后端API的响应数据。例如：

```yaml
objects:
  - name: users
    sample: |
      {
        "name": "LinuxSuRen",
        "age": 18
        "gender": "male"
      }
proxies:
  - path: /api/v1/projects/{projectID}
    target: http://localhost:8080
```

把上面的内容放到 `mock.yaml` 文件中，然后使用 `atest mock --prefix /api/v1 --port 6060 mock.yaml` 命令即可启动一个 HTTP Mock 服务。

此时，Mock 服务就会把**代理**模块指定的 API 转发到已有服务的的接口上，并同时提供了 `users` 对象的增删查改（CRUD）的标准 API。你可以用 `atest` 或者 `curl` 命令来调用这些 API。

```shell
curl -X POST -d '{"name": "Rick"}' http://localhost:6060/api/v1/users
curl -X GET http://localhost:6060/api/v1/users
curl -X PUT -d '{"name": "Rick", "age": 20}' http://localhost:6060/api/v1/users/Rick
curl -X GET http://localhost:6060/api/v1/users/Rick
curl -X DELETE http://localhost:6060/api/v1/users/Rick
```

非常期待 `atest` 可以帮助更多的项目持续提升、保持 API 稳定性。

## 🚀 主要的新功能

* Mock 功能的增强，包含对象、原始、代理三种模式 (#552) @LinuxSuRen
* 支持重命名测试用例、测试集 (#550) @LinuxSuRen
* 支持给定频率下重复执行测试用例 (#548) @LinuxSuRen
* 下载插件文件时显示进度信息 (#544) @LinuxSuRen
* 支持生成随机图片并上传 (#541) @LinuxSuRen
* 支持上传嵌入式文件（基于 base64 编码） (#538) @LinuxSuRen
* 支持导入其他 atest 实例的用例数据 (#539) @LinuxSuRen
* UI 上显示响应体的大小 (#536) @LinuxSuRen
* 增加基于 MySQL 位存储的测试用例执行历史记录  (#524) @SamYSF
* 支持设置插件下载的“前缀”信息 (#532) @SamYSF
* 优化存储插件管理界面 (#518) @LinuxSuRen
* 在 UI 上增加快捷键支持 (#510) @LinuxSuRen
* 重构 API 风格为 restFul (#497) @LinuxSuRen
* 增加 Mock 配置的 JSON schema (#499) @LinuxSuRen
* 增加了对 JSON 兼容性的响应格式的支持 (#496) @LinuxSuRen

## 🐛 缺陷修复

* 修复测试用例重复时被覆盖的问题 (#531) @LinuxSuRen

## 致谢

本次版本发布，包含了以下 3 位 contributor 的努力：

* [@LinuxSuRen](https://github.com/LinuxSuRen)
* [@SamYSF](https://github.com/SamYSF)
* [@yuluo-yx](https://github.com/yuluo-yx)

## 相关数据

下面是 `atest` 截止到 `v0.0.18` 的部分数据：

* watch 9
* fork 50
* star 249 (+40)
* contributor 25 (+1)
* 二进制文件下载量 6.3k (+3.2k)
* 部分镜像 6.4k (+0.9k)
* 单元测试覆盖率 76% (+2%)

想了解完整信息的话，请访问 https://github.com/LinuxSuRen/api-testing/releases/tag/v0.0.18
