HTTP API record tool.

## Usage

```shell
atest-collector --filter-path /answer/api/v1
```

It will start a HTTP proxy server, and set the server address to your browser proxy (such as: [SwitchyOmega](https://github.com/FelisCatus/SwitchyOmega)).

`atest-collector` will record all HTTP requests which has prefix `/answer/api/v1`, and 
save it to file `sample.yaml` once you close the server.

## Features

* Basic authorization
* Upstream proxy
* URL path filter
* Support save response body or not
