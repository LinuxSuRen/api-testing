Please join us to improve this project.

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
docker run -p 8080:8080 -e SW_OAP_ADDRESS=http://localhost:12800 -e SW_ZIPKIN_ADDRESS=http://localhost:9412 docker.io/apache/skywalking-ui:9.0.0

make build

export SW_AGENT_NAME=atest
export SW_AGENT_REPORTER_GRPC_BACKEND_SERVICE=10.121.218.184:31065
export SW_AGENT_PLUGIN_CONFIG_HTTP_SERVER_COLLECT_PARAMETERS=true
export SW_AGENT_METER_COLLECT_INTERVAL=3
export SW_AGENT_LOG_TYPE=std
./bin/atest server --local-storage 'bin/*.yaml' --http-port 8082 --port 7072 --console-path console/atest-ui/dist/
```