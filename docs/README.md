Welcome to use `atest` to improve your code quality.

## Get started
TODO

## Storage
There are multiple storage backend supported: See the status from the list:

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
# create a network
# start the server with gRPC storage
podman run -p 8080:8080 ghcr.io/linuxsuren/api-testing:master \
    atest server --storage grpc --grpc-storage 192.168.1.98:7071 --console-path=/var/www/html
# start the gRPC storage which connect to an ORM database
podman run -p 7071:7071 ghcr.io/linuxsuren/api-testing:master \
    atest-store-orm --address 192.168.1.98:4000 --user root --database test
```

> Please don't forget to replace `192.168.1.98` to your own IP address.

## Extensions
Developers could have a storage extension. Implement a gRPC server according to [loader.proto](../pkg/testing/remote/loader.proto) is required.

## Official Images
You could find the official images from both [Docker Hub](https://hub.docker.com/r/linuxsuren/api-testing) and [GitHub Images](https://github.com/users/LinuxSuRen/packages/container/package/api-testing). See the image path:

* `ghcr.io/linuxsuren/api-testing:latest`
* `linuxsuren/api-testing:latest`

The tag `latest` represents the latest release version. The tag `master` represents the image of the latest master branch. We highly recommend you using a fixed version instead of those in a production environment.

## Release Notes
* [v0.0.12](release-note-v0.0.12.md)

## Articles
* [Introduction](introduce-zh.md)
* [GLCC 2023 announccement](glcc-2023-announce.md)
