---
title: "贡献指南"
weight: -1
description: "API Testing 贡献指南."
---

请加入我们，共同完善这个项目。

后端由 [Golang](https://go.dev/) 编写，前端由 [Vue](https://vuejs.org/) 编写。

### 对于初学者

在开始之前，您可能需要了解以下技术:

| Name                                                                        | Domain                                                                 |
|-----------------------------------------------------------------------------|------------------------------------------------------------------------|
| [HTTP](https://developer.mozilla.org/en-US/docs/Web/HTTP/Overview) Protocol | Core                                                                   |
| [RESTful](https://en.wikipedia.org/wiki/REST)                               | Core                                                                   |
| [gRPC](https://grpc.io/)                                                    | `gRPC` runner extension                                                |
| [Prometheus](https://prometheus.io/)                                        | Application monitor                                                    |
| [Cobra](https://github.com/spf13/cobra)                                     | The Go CLI framework                                                   |
| [Element Plus](https://element-plus.org/)                                   | The front-end framework                                                |
| [Docker](https://www.docker.com/get-started/)                               | The container image build                                              |
| [Helm chart](https://helm.sh/)                                              | The [Kubernetes](https://kubernetes.io/docs/home/) application package |
| [GitHub Actions](https://docs.github.com/en/actions)                        | The continuous integration                                             |
| [make](https://www.gnu.org/software/make/)                                  | The automated Build Tools                                              |
| [Docs Guide](https://github.com/LinuxSuRen/api-testing.git) | 文档编写指南 |

## 设置开发环境

> 本项目使用 `make` 作为构建工具，并设计了非常强大的 make 指令系统。您可以通过运行 `make help` 查看所有可用的命令。

强烈建议您配置 `git pre-commit` 钩子。它会强制在提交前运行单元测试。
运行以下命令：

```shell
make install-precheck
```

## 打印各行代码：

```shell
git ls-files | xargs cloc
```

## pprof

```shell
go tool pprof -http=:9999 http://localhost:8080/debug/pprof/heap
```

其他用法：

* `/debug/pprof/heap?gc=1`
* `/debug/pprof/heap?seconds=10` 
* `/debug/pprof/goroutine/?debug=0` 

## SkyWalking

```shell
docker run -p 12800:12800 -p 9412:9412 docker.io/apache/skywalking-oap-server:9.0.0
docker run -p 8080:8080 -e SW_OAP_ADDRESS=http://172.11.0.6:12800 -e SW_ZIPKIN_ADDRESS=http://172.11.0.6:9412 docker.io/apache/skywalking-ui:9.0.0

make build

export SW_AGENT_NAME=atest
export SW_AGENT_REPORTER_GRPC_BACKEND_SERVICE=172.11.0.6:30689
export SW_AGENT_PLUGIN_CONFIG_HTTP_SERVER_COLLECT_PARAMETERS=true
export SW_AGENT_METER_COLLECT_INTERVAL=3
export SW_AGENT_LOG_TYPE=std
export SW_AGENT_REPORTER_DISCARD=true
./bin/atest server --local-storage 'bin/*.yaml' --http-port 8082 --port 7072 --console-path console/atest-ui/dist/
```

通过 BanYanDB 运行 SkyWalking：

```shell
docker run -p 17912:17912 -p 17913:17913 apache/skywalking-banyandb:latest  standalone

docker run -p 12800:12800 -p 9412:9412 \
    -e SW_STORAGE=banyandb \
    -e SW_STORAGE_BANYANDB_HOST=192.168.1.98 \
    docker.io/apache/skywalking-oap-server
```

## 第一次贡献

对于第一次对此项目贡献代码的开发者，您应该在本地开发环境运行如下命令：

```shell
make test
```

以确保通过项目测试，这会有助于您检查并解决在提交时遇到的错误，同时减少 review 的复杂度。

## FAQ

* Got sum missing match error of go.
  * 运行命令： `go clean -modcache && go mod tidy`.
