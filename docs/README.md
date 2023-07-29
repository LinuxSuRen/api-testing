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
* Podman, and Docker

Please see the following example usage:

```shell
atest service start -m podman --version master
```

## Storage
There are multiple storage backends supported. See the status from the list:

| Name | Status |
|---|---|
| Local Storage | Ready |
| ORM DataBase | Developing |
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

## Extensions
Developers could have a storage extension. Implement a gRPC server according to [loader.proto](../pkg/testing/remote/loader.proto) is required.

## Official Images
You could find the official images from both [Docker Hub](https://hub.docker.com/r/linuxsuren/api-testing) and [GitHub Images](https://github.com/users/LinuxSuRen/packages/container/package/api-testing). See the image path:

* `ghcr.io/linuxsuren/api-testing:latest`
* `linuxsuren/api-testing:latest`
* `docker.m.daocloud.io/linuxsuren/api-testing` (mirror)

The tag `latest` represents the latest release version. The tag `master` represents the image of the latest master branch. We highly recommend you using a fixed version instead of those in a production environment.

## Release Notes
* [v0.0.12](release-note-v0.0.12.md)

## Articles
* [Introduction](introduce-zh.md)
* [GLCC 2023 announccement](glcc-2023-announce.md)
