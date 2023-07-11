Print the code of lines:

```shell
git ls-files | xargs cloc
```

## Setup development environment
It's highly recommended you to configure the git pre-commit hook. It will force to run unit tests before commit.
Run the following command:

```shell
make install-precheck
```
