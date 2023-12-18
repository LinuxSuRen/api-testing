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

package compare

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

func TestElement(t *testing.T) {
	exp := `{
		"data": [
		  {
			"key": "hell",
			"value": "func() strin"
		  }
		]
	  }
	`
	act := `
	  {
		"data": [
		  {
			"key": "hello",
			"value": "func() string"
		  }
		]
	  }`
	expect := gjson.Parse(exp)
	actul := gjson.Parse(act)

	err := Element("TestElement", expect, actul)

	expmsg1 := "compare: field TestElement.data.0.value: expect func() strin but got func() string"
	expmsg2 := "compare: field TestElement.data.0.key: expect hell but got hello"
	assert.Contains(t, err.Error(), expmsg1)
	assert.Contains(t, err.Error(), expmsg2)
}
