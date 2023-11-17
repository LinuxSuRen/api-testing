Ports in extensions:

| Name | Port |
|------|------|
| orm  | 4071 |
| s3   | 4072 |
| etcd | 4073 |
| git  | 4074 |
| mongodb | 4075 |

## Contribute a new extension

* First, create a directory in current directory. And please keep the same naming convertion.
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
