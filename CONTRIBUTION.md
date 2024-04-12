Please join us to improve this project.

The backend is written by [Golang](https://go.dev/), and the front-end is written by [Vue](https://vuejs.org/).

## For beginner
You might need to know the following tech before get started.

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

## Setup development environment
It's highly recommended you to configure the git pre-commit hook. It will force to run unit tests before commit.
Run the following command:

```shell
make install-precheck
```

## Print the code of lines:

```shell
git ls-files | xargs cloc
```

## pprof

```
go tool pprof -http=:9999 http://localhost:8080/debug/pprof/heap
```

Other usage of this:
* `/debug/pprof/heap?gc=1`
* `/debug/pprof/heap?seconds=10`
* `/debug/pprof/goroutine/?debug=0`

## Skywalking

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

Run SkyWalking with BanYanDB
```shell
docker run -p 17912:17912 -p 17913:17913 apache/skywalking-banyandb:latest  standalone

docker run -p 12800:12800 -p 9412:9412 \
    -e SW_STORAGE=banyandb \
    -e SW_STORAGE_BANYANDB_HOST=192.168.1.98 \
    docker.io/apache/skywalking-oap-server
```

## FAQ

* Got sum missing match error of go.
  * Run command: `go clean -modcache && go mod tidy`
