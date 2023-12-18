/*
Copyright 2023 API Testing Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package testing_test

import (
	"testing"

	atesting "github.com/linuxsuren/api-testing/pkg/testing"
	"github.com/stretchr/testify/assert"

	"gopkg.in/yaml.v3"
)

func TestInScope(t *testing.T) {
	testCase := &atesting.TestCase{Name: "foo"}
	assert.True(t, testCase.InScope(nil))
	assert.True(t, testCase.InScope([]string{"foo"}))
	assert.False(t, testCase.InScope([]string{"bar"}))
}

func TestRequestBody(t *testing.T) {
	req := &atesting.Request{}
	graphqlBody := `api: /api
body:
    query: query
    operationName: ""
    variables:
        name: rick
`

	err := yaml.Unmarshal([]byte(graphqlBody), req)
	assert.Nil(t, err)
	assert.Equal(t, `{"query":"query","operationName":"","variables":{"name":"rick"}}`, req.Body.String())

	var data []byte
	data, err = yaml.Marshal(req)
	assert.Nil(t, err)
	assert.Equal(t, graphqlBody, string(data))

	err = yaml.Unmarshal([]byte(`body: plain`), req)
	assert.Nil(t, err)
	assert.Equal(t, "plain", req.Body.String())
}
