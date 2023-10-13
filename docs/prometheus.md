## Pushing the test results into Prometheus

You can use the following command to do it:

```shell
atest run --report prometheus --report-file http://localhost:9091 \
    -p sample/testsuite-gitee.yaml --duration 30m --qps 1
```

It will push the test results data into Prometheus [PushGateway](https://github.com/prometheus/pushgateway).
Then Prometheus could get the metrics from it.

Skip the following instructions if you are familiar with Prometheus:
```shell
docker run \
    -p 9090:9090 \
    -v /etc/timezone:/etc/timezone:ro \
    -v /etc/localtime:/etc/localtime:ro \
    -v /root/prometheus.yml:/etc/prometheus/prometheus.yml \
    prom/prometheus

docker run -p 9091:9091 \
    -v /etc/timezone:/etc/timezone:ro \
    -v /etc/localtime:/etc/localtime:ro \
    prom/pushgateway

docker run -p 3000:3000 docker.io/grafana/grafana
```
