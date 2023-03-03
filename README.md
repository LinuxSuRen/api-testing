This is a API testing tool.

## Feature
* Response Body fields equation check
* Response Body [eval](https://expr.medv.io/)
* Output reference between TestCase

## Template
The following fields are templated with [sprig](http://masterminds.github.io/sprig/):

* API
* Request Body

## Limit
* Only support to parse the response body when it's a map
