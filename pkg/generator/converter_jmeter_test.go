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
	"testing"

	_ "embed"

	atest "github.com/linuxsuren/api-testing/pkg/testing"
	"github.com/stretchr/testify/assert"
)

func TestJmeterConvert(t *testing.T) {
	t.Run("common", func(t *testing.T) {
		jmeterConvert := GetTestSuiteConverter("jmeter")
		assert.NotNil(t, jmeterConvert)

		converters := GetTestSuiteConverters()
		assert.Equal(t, 2, len(converters))
	})

	converter := &jmeterConverter{}
	output, err := converter.Convert(createTestSuiteForTest())
	assert.NoError(t, err)

	assert.Equal(t, expectedJmeter, output, output)
}

func createTestSuiteForTest() *atest.TestSuite {
	return &atest.TestSuite{
		Name: "API Testing",
		API:  `{{default "http://localhost:8080/server.Runner" (env "SERVER")}}`,
		Items: []atest.TestCase{{
			Name: "hello-jmeter",
			Request: atest.Request{
				Method: "POST",
				API:    "/GetSuites",
				Body:   atest.NewRequestBody(`sample`),
			},
		}},
	}
}

//go:embed testdata/expected_jmeter.jmx
var expectedJmeter string
