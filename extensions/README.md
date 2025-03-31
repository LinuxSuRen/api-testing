Ports in extensions:

| Type | Name                                                                     | Port |
|------|--------------------------------------------------------------------------|------|
| Store | [orm](https://github.com/LinuxSuRen/atest-ext-store-orm)                 | 4071 |
| Store | [s3](https://github.com/LinuxSuRen/atest-ext-store-s3)                   | 4072 |
| Store | [etcd](https://github.com/LinuxSuRen/atest-ext-store-etcd)               | 4073 |
| Store | [git](https://github.com/LinuxSuRen/atest-ext-store-git)                 | 4074 |
| Store | [mongodb](https://github.com/LinuxSuRen/atest-ext-store-mongodb)         | 4075 |
| Store | [redis](https://github.com/LinuxSuRen/atest-ext-store-redis)             |  |
| Store | [iotdb](https://github.com/LinuxSuRen/atest-ext-store-iotdb) | |
| Store | [Cassandra](https://github.com/LinuxSuRen/atest-ext-store-cassandra) | |
| Monitor | [docker-monitor](https://github.com/LinuxSuRen/atest-ext-monitor-docker) |  |
| Agent | [collector](https://github.com/LinuxSuRen/atest-ext-collector)           |  |
| Secret | [Vault](https://github.com/LinuxSuRen/api-testing-vault-extension)       | |
| Data | [Swagger](https://github.com/LinuxSuRen/atest-ext-data-swagger) | |

## Contribute a new extension

* First, create a repository. And please keep the same naming convertion.
* Second, implement the `Loader` gRPC service which defined by [this proto](../pkg/testing/remote/loader.proto).
* Finally, add the extension's name into function [SupportedExtensions](../console/atest-ui/src/views/store.ts).

## Naming conventions

Please follow the following conventions if you want to add a new store extension:

`store-xxx`

`xxx` should be a type of a backend storage.

## Test

First, build and copy the binary file into the system path. You can run the following
command in the root directory of this project:

```shell
make build-ext-etcd copy-ext
```
