+++
title = "End-to-End"
weight = 103
+++

`atest` 非常适合针对（HTTP）接口做 E2E（端到端）测试，E2E 测试可以确保后端接口在完整的环境中持续地保持正确运行。下面采用 Docker compose 给出一个使用事例：

```yaml
version: '3.1'
services:
  testing:
    image: ghcr.io/linuxsuren/api-testing:latest
    environment:
      SERVER: http://server:8080
    volumes:
      - ./testsuite.yaml:/work/testsuite.yaml
    command: atest run -p /work/testsuite.yaml
    depends_on:
      server:
        condition: service_healthy
    links:
      - server
  server:
    image: ghcr.io/devops-ws/learn-springboot:master
    healthcheck:
      test: ["CMD", "bash", "-c", "cat < /dev/null > /dev/tcp/127.0.0.1/8080"]
      interval: 3s
      timeout: 60s
      retries: 10
      start_period: 3s
```

从 Docker compose `v2.36.0` 开始，可以采用如下简化的写法：

```yaml
services:
  testing:
    scale: 0
    provider:
      type: atest
      options:
        pattern: testsuite.yaml
    environment:
      SERVER: http://server:8080
    depends_on:
      server:
        condition: service_healthy
    links:
      - server
  server:
    image: ghcr.io/devops-ws/learn-springboot:master
    healthcheck:
      test: ["CMD", "bash", "-c", "cat < /dev/null > /dev/tcp/127.0.0.1/8080"]
      interval: 3s
      timeout: 60s
      retries: 10
      start_period: 3s
```
