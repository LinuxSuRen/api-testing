Welcome to use `atest` to improve your code quality.

## Get started
You can use `atest` as a CLI or other ways:

* Web UI
* [VS Code Extension](https://marketplace.visualstudio.com/items?itemName=linuxsuren.api-testing)

See also the screenshots below:

![image](https://github.com/LinuxSuRen/api-testing/assets/1450685/e3404c53-34bc-4bf0-8f6a-c1873e2263a2)

![image](https://github.com/LinuxSuRen/api-testing/assets/1450685/e959f560-1fb5-4592-9f45-ec883c385785)

## Installation
You can install in various methods:

* CLI via `hd i atest`
* Web server
* [Kubernetes](https://github.com/LinuxSuRen/api-testing/tree/master/sample/kubernetes)
* [Argo CD](https://github.com/LinuxSuRen/api-testing/blob/master/sample/argocd/simple.yaml)

If you're developing APIs locally, the best way is installing it as a container service.
Then you can access it via your browser.

Currently, it supports the following kinds of services:

* Operate System services
  * Linux, and Darwin
* [Podman](https://github.com/containers/podman), and Docker

Please see the following example usage:

```shell
sudo atest service install -m podman --version master
```

or run in Docker:
```shell
docker run -v /var/www/sample:/var/www/sample \
  --network host \
  linuxsuren/api-testing:master
```

the default web server port is `8080`. So you can visit it via: http://localhost:8080

## Run in k3s

```shell
sudo k3s server --write-kubeconfig-mode 666

k3s kubectl apply -k sample/kubernetes/default

kustomize build sample/kubernetes/docker.io/ | k3s kubectl apply -f -
```

## Run your test cases
The test suite file could be in local, or in the HTTP server. See the following different ways:

* `atest run -p your-local-file.yaml`
* `atest run -p https://gitee.com/linuxsuren/api-testing/raw/master/sample/testsuite-gitee.yaml`
* `atest run -p http://localhost:8080/server.Runner/ConvertTestSuite?suite=sample`

For the last one, it represents the API Testing server.

## Convert to JMeter
[JMeter](https://jmeter.apache.org/) is a load test tool. You can run the following commands from the root directory of this repository:

```shell
atest convert --converter jmeter -p sample/testsuite-gitee.yaml --target bin/gitee.jmx

jmeter -n -t bin/gitee.jmx
```

Please feel free to bring more test tool converters.

## Run in Jenkins
You can run the API testings in Jenkins, see also the following example:

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

## Storage
There are multiple storage backends supported. See the status from the list:

| Name | Status |
|---|---|
| Local Storage | Ready |
| S3 | Ready |
| ORM DataBase | Developing |
| Git Repository | Developing |
| Etcd DataBase | Developing |

### Local Storage
Local storage is the built-in solution. You can run it with the following command:

```shell
podman run --pull always -p 8080:8080 ghcr.io/linuxsuren/api-testing:master

# The default local storage directory is: /var/www/sample
# You can find the test case YAML files in it.
# Visit it from http://localhost:8080 once it's ready.
```

or, you can run the CLI in terminal like this:

```shell
atest server --local-storage 'sample/*.yaml' --console-path console/atest-ui/dist
```

using the host network mode if you want to connect to your local environment:
```shell
podman run --pull always --network host ghcr.io/linuxsuren/api-testing:master
```

### ORM DataBase Storage
Start a database with the following command if you don't have a database already. You can install [tiup](https://tiup.io/) via `hd i tiup`.

```shell
tiup playground --db.host 0.0.0.0
```

```shell
# create a config file
mkdir bin
echo "- name: db
  kind:
    name: database
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

See also the expected configuration below:

```yaml
- name: s3
  url: http://172.11.0.13:30999   # address of the s3 server
  kind:
    name: s3
    url: localhost:7072           # address of the s3 storage extension
  properties:
    accessKeyID: 6e03rIMChrsZ6YZl
    secretAccessKey: F0xH6o2qRYTyAUyRuXO81B4gj7zUrSaj
    disableSSL:  true
    forcepathstyle: true
    bucket: vm1
    region: cn
```

### Git Storage
You can use a git repository as the storage backend.

```shell
# The default port is 7074
podman run --network host \
    ghcr.io/linuxsuren/api-testing:master atest-store-git
```

See also the expected configuration below:

```yaml
- name: git
  url: http://172.11.0.13:30999   # address of the git repository
  username: linuxsuren
  password: linuxsuren
  kind:
    name: git
    url: localhost:7074           # address of the git storage extension
  properties:
    targetPath: .
```

## Secret Server
You can put the sensitive information into a secret server. For example, [Vault](https://www.github.com/hashicorp/vault).

Connect to [a vault extension](https://github.com/LinuxSuRen/api-testing-secret-extension) via flag: `--secret-server`. Such as:

```shell
atest server --secret-server localhost:7073
```

## Extensions
Developers could have storage, secret extensions. Implement a gRPC server according to [loader.proto](../pkg/testing/remote/loader.proto) is required.

## Official Images
You could find the official images from both [Docker Hub](https://hub.docker.com/r/linuxsuren/api-testing) and [GitHub Images](https://github.com/users/LinuxSuRen/packages/container/package/api-testing). See the image path:

* `ghcr.io/linuxsuren/api-testing:latest`
* `linuxsuren/api-testing:latest`
* `docker.m.daocloud.io/linuxsuren/api-testing` (mirror)

The tag `latest` represents the latest release version. The tag `master` represents the image of the latest master branch. We highly recommend you using a fixed version instead of those in a production environment.

## Release Notes
* [v0.0.13](release-note-v0.0.13.md)
* [v0.0.12](release-note-v0.0.12.md)

## Articles
* [Introduction](introduce-zh.md)
* [GLCC 2023 announccement](glcc-2023-announce.md)
