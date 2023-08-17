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
