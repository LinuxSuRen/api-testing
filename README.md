[![Codacy Badge](https://app.codacy.com/project/badge/Grade/3f16717cd6f841118006f12c346e9341)](https://www.codacy.com/gh/LinuxSuRen/api-testing/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=LinuxSuRen/api-testing&amp;utm_campaign=Badge_Grade)
[![Codacy Badge](https://app.codacy.com/project/badge/Coverage/3f16717cd6f841118006f12c346e9341)](https://www.codacy.com/gh/LinuxSuRen/api-testing/dashboard?utm_source=github.com&utm_medium=referral&utm_content=LinuxSuRen/api-testing&utm_campaign=Badge_Coverage)
![GitHub All Releases](https://img.shields.io/github/downloads/linuxsuren/api-testing/total)

This is a API testing tool.

## Feature
* Response Body fields equation check  
* Response Body [eval](https://expr.medv.io/)  
* Output reference between TestCase  
* Run in server mode, and provide the gRPC endpoint
* [VS Code extension](https://github.com/LinuxSuRen/vscode-api-testing) support

## Template
The following fields are templated with [sprig](http://masterminds.github.io/sprig/):

* API  
* Request Body  

## TODO
* Reduce the size of context  
* Support customized context  

## Limit
* Only support to parse the response body when it's a map or array  
