# API Testing

Welcome to use `atest` to improve your code quality.

## Get started

You can use `atest` as a CLI or as:

* Web UI
* [VS Code Extension](https://marketplace.visualstudio.com/items?itemName=linuxsuren.api-testing)

See also the screenshots below:

![image](https://github.com/LinuxSuRen/api-testing/assets/1450685/e3404c53-34bc-4bf0-8f6a-c1873e2263a2)

![image](https://github.com/LinuxSuRen/api-testing/assets/1450685/e959f560-1fb5-4592-9f45-ec883c385785)

## Installation

There are various ways of installing `atest`:

* CLI via `hd i atest`
* Web server
* [Kubernetes](https://github.com/LinuxSuRen/api-testing/tree/master/docs/manifests/kubernetes)
* [Argo CD](https://github.com/LinuxSuRen/api-testing/blob/master/docs/manifests/argocd/simple.yaml)
* [Helm](helm.md)

If you're developing APIs locally, the best way is to install it as a container service.
Then you can access it via your browser.

Currently, it supports the following kinds of services:

* Operate System services
  * Linux and Darwin
* [Podman](https://github.com/containers/podman) and Docker

[![Deployed on Zeabur](https://zeabur.com/deployed-on-zeabur-dark.svg)](https://zeabur.com?referralCode=LinuxSuRen&utm_source=LinuxSuRen&utm_campaign=oss)

### Have a look at the following example usage

#### Podman

```shell
sudo atest service install -m podman --version master
```

#### Docker

```shell
docker run -v /var/www/sample:/var/www/sample \
  --network host \
  linuxsuren/api-testing:master
```

The default web server port is `8080`. So you can visit it via: <http://localhost:8080>

## Run in k3s

```shell
sudo k3s server --write-kubeconfig-mode 666

k3s kubectl apply -k sample/kubernetes/default

kustomize build sample/kubernetes/docker.io/ | k3s kubectl apply -f -
```

## Run your test cases

The test suite file could be local, or in the HTTP server. Have a look at some examples:

* `atest run -p your-local-file.yaml`
* `atest run -p https://gitee.com/linuxsuren/api-testing/raw/master/sample/testsuite-gitee.yaml`
* `atest run -p http://localhost:8080/server.Runner/ConvertTestSuite?suite=sample`

The last example pertains to the API Testing server.

## Functions

There are two kinds of functions for two situations: template rendering and test results verification.

* Template rendering functions are based on [the Go template](https://pkg.go.dev/text/template)
* The verification functions are based on [expr library](https://expr.medv.io/)

You can query the supported functions on the UI page or by using the command:

```shell
atest func
```

## Hooks

In some cases you may want to run the test cases after an HTTP server is ready. Then you can use the hooks feature as follows:

```yaml
name: Gitlab
api: https://gitlab.com/api/v4
param:
  user: linuxsuren
items:
- name: projects
  request:
    api: /projects
  before:
    items:
      - "sleep(1)"
  after:
    items:
      - "sleep(1)"
```

You can use all the functions that are available in the expr library.

## Convert to JMeter

[JMeter](https://jmeter.apache.org/) is a load test tool. You can run the following commands from the root directory of this repository:

```shell
atest convert --converter jmeter -p sample/testsuite-gitee.yaml --target bin/gitee.jmx

jmeter -n -t bin/gitee.jmx
```

Please feel free to bring more test tool converters.

## Run in Jenkins

You can run the API testings in Jenkins, as demonstrated in the example below:

```Jenkinsfile
pipeline {
    agent any
    
    stages() {
        stage('test') {
            steps {
                sh '''
                curl http://localhost:9090/get -o atest
                chmod u+x atest
                
                ./atest run -p http://localhost:9090/server.Runner/ConvertTestSuite?suite=api-testing
                '''
            }
        }
    }
}
```

## Report

You can see the test results in [Grafana](prometheus.md).

## Monitoring

It can monitor the server and browser via the [Apache SkyWalking](https://skywalking.apache.org/).
Please add the following flag if you want to get the browser tracing data:

```shell
# 12800 is the HTTP port of Apache SkyWalking
# 11800 is the gRPC port of Apache SkyWalking
export SW_AGENT_REPORTER_GRPC_BACKEND_SERVICE=localhost:11800
export SW_AGENT_REPORTER_DISCARD=false

atest server --skywalking http://localhost:12800
```

## Storage

There are multiple storage backends supported. See the status from the list:

| Name | Status |
|---|---|
| Local Storage | Ready |
| S3 | Ready |
| ORM DataBase | Ready |
| Git Repository | Ready |
| Etcd | Ready |
| MongoDB | Devloping |

### Local Storage

Local storage is the built-in solution. You can run it with the following command:

```shell
podman run --pull always -p 8080:8080 ghcr.io/linuxsuren/api-testing:master

# The default local storage directory is: /var/www/sample
# You can find the test case YAML files in it.
# Visit it from http://localhost:8080 once it's ready.
```

Or, you can run the CLI in the terminal like this:

```shell
atest server --local-storage 'sample/*.yaml' --console-path console/atest-ui/dist
```

Use the host network mode if you want to connect to your local environment:

```shell
podman run --pull always --network host ghcr.io/linuxsuren/api-testing:master
```

### ORM Database Storage

Start a database with the following command if you don't have a database already. You can install [tiup](https://tiup.io/) via `hd i tiup`.

```shell
tiup playground --db.host 0.0.0.0
```

```shell
# create a config file
mkdir bin
echo "- name: db
  kind:
    name: atest-store-orm
    url: localhost:7071
  url: localhost:4000
  username: root
  properties:
    database: test" > bin/stores.yaml

# start the server with gRPC storage
podman run -p 8080:8080 -v bin:var/data/atest \
    --network host \
    ghcr.io/linuxsuren/api-testing:master \
    atest server --console-path=/var/www/html \
    --config-dir=/var/data/atest

# start the gRPC storage which ready to connect to an ORM database
podman run -p 7071:7071 \
    --network host \
    ghcr.io/linuxsuren/api-testing:master atest-store-orm
```

### S3 Storage

You can use a S3 compatible storage as the storage backend.

```shell
# The default port is 7072
podman run --network host \
    ghcr.io/linuxsuren/api-testing:master atest-store-s3
```

Have a look at the expected configuration below:

```yaml
- name: s3
  url: http://172.11.0.13:30999   # address of the s3 server
  kind:
    name: atest-store-s3
    url: localhost:7072           # address of the s3 storage extension
  properties:
    accessKeyID: 6e03rIMChrsZ6YZl
    secretAccessKey: F0xH6o2qRYTyAUyRuXO81B4gj7zUrSaj
    sessiontoken: ""
    region: cn
    disableSSL:  true
    forcepathstyle: true
    bucket: vm1
```

### Git Storage

You can use a git repository as the storage backend.

```shell
# The default port is 7074
podman run --network host \
    ghcr.io/linuxsuren/api-testing:master atest-store-git
```

Have a look at the expected configuration below:

```yaml
- name: git
  url: http://172.11.0.13:30999   # address of the git repository
  username: linuxsuren
  password: linuxsuren
  kind:
    name: atest-store-git         # the extension binary file name
    url: localhost:7074           # address of the git storage extension
  properties:                     # optional properties for specific features
    targetPath: .                 # target path to find YAML files
    name: linuxsuren              # the name for git commit
    email: linuxsuren@github.com  # the email address for git commit
    insecure: false               # whether to use insecure
```

### MongoDB Storage

You can use a MongoDB as the storage backend.

Have a look at the expected configuration below:

```yaml
- name: mongodb
  url: 172.11.0.13:27017   # address of the mongodb
  username: linuxsuren
  password: linuxsuren
  kind:
    name: atest-store-mongodb     # the extension binary file name
  properties:                     # optional properties for specific features
    database: testing             # the database name
    collection: atest             # the collection name
```

## Secret Server

You can put sensitive information into a secret server. For example, [Vault](https://www.github.com/hashicorp/vault).

Connect to [a vault extension](https://github.com/LinuxSuRen/api-testing-secret-extension) via flag: `--secret-server`. Such as:

```shell
atest server --secret-server localhost:7073
```

## Application monitor

You can get the resource usage in the report through Docker:

```shell
atest run -p sample/testsuite-gitlab.yaml --monitor-docker test --report md
```

## Verify

| Item | Description |
|---|---|
| `expect.bodyFieldsExpect` | See also the syntax from [gjson](https://github.com/tidwall/gjson) |

## OAuth

It support GitHub, [Dex](https://github.com/dexidp/dex) as OAuth provider. See also the following usage:

```shell
atest server --auth oauth --client-id your-id --client-secret your-secret
```

## Extensions

Developers can have storage, secret extensions. Implementing a gRPC server according to [loader.proto](../pkg/testing/remote/loader.proto) is required.

## Official Images

You can find the official images from both [Docker Hub](https://hub.docker.com/r/linuxsuren/api-testing) and others. See the image path:

* `ghcr.io/linuxsuren/api-testing:master`
* `docker.io/linuxsuren/api-testing:master`
* `registry.aliyuncs.com/linuxsuren/api-testing:master`
* `ccr.ccs.tencentyun.com/linuxsuren/api-testing:master`
* `docker.m.daocloud.io/linuxsuren/api-testing:master` (mirror)

The tag `latest` represents the latest release version. The tag `master` represents the image of the latest master branch. We highly recommend you to use a fixed version instead of those in a production environment.

## Release Notes

* [v0.0.15](release-note-v0.0.15.md)
* [v0.0.14](release-note-v0.0.14.md)
* [v0.0.13](release-note-v0.0.13.md)
* [v0.0.12](release-note-v0.0.12.md)

## Articles

* [Introduction](introduce-zh.md)
* [GLCC 2023 announccement](glcc-2023-announce.md)
