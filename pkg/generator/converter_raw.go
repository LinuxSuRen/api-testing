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
package generator

import (
	"github.com/linuxsuren/api-testing/pkg/testing"
	"gopkg.in/yaml.v3"
)

type rawConverter struct {
}

func init() {
	RegisterTestSuiteConverter("raw", &rawConverter{})
}

func (c *rawConverter) Convert(testSuite *testing.TestSuite) (result string, err error) {
	if err = testSuite.Render(make(map[string]interface{})); err == nil {
		for _, item := range testSuite.Items {
			item.Request.RenderAPI(testSuite.API)
		}

		var data []byte
		if data, err = yaml.Marshal(testSuite); err == nil {
			result = string(data)
		}
	}
	return
}
