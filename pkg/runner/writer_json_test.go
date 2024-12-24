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

package runner_test

import (
	"bytes"
	_ "embed"
	"testing"

	"github.com/linuxsuren/api-testing/pkg/runner"
	"github.com/stretchr/testify/assert"
)

func TestJSONResultWriter(t *testing.T) {
	buf := new(bytes.Buffer)
	writer := runner.NewJSONResultWriter(buf)
	writer.WithAPICoverage(nil)

	err := writer.Output([]runner.ReportResult{{
		Name:    "foo",
		API:     "api",
		Average: 3,
		Max:     4,
		Min:     2,
		Count:   3,
		Error:   0,
	}, {
		Name:    "bar",
		API:     "api",
		Average: 3,
		Max:     4,
		Min:     2,
		Count:   3,
		Error:   0,
	}})
	assert.Nil(t, err)
	assert.JSONEq(t, jsonResult, buf.String())
	assert.NotNil(t, writer.WithResourceUsage(nil))
}

//go:embed testdata/json-result.json
var jsonResult string
