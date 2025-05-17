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
	"testing"

	"github.com/linuxsuren/api-testing/pkg/apispec"
	"github.com/linuxsuren/api-testing/pkg/runner"
	"github.com/stretchr/testify/assert"
)

func TestNewStdResultWriter(t *testing.T) {
	tests := []struct {
		name         string
		buf          *bytes.Buffer
		apiConverage apispec.APICoverage
		results      []runner.ReportResult
		expect       string
	}{{
		name:    "result is nil",
		buf:     new(bytes.Buffer),
		results: nil,
		expect: `Name Average Max Min QPS Count Error
Test case count: 0
`,
	}, {
		name: "have one item",
		buf:  new(bytes.Buffer),
		apiConverage: apispec.NewFakeAPISpec([][]string{{
			"/api", "GET",
		}}),
		results: []runner.ReportResult{{
			Name:    "/api",
			API:     "/api",
			Average: 1,
			Max:     1,
			Min:     1,
			QPS:     10,
			Count:   1,
			Error:   0,
		}},
		expect: `Name Average Max Min QPS Count Error
/api 1ns 1ns 1ns 10 1 0
Test case count: 1

API Coverage: 1/1
`,
	}, {
		name: "have errors",
		buf:  new(bytes.Buffer),
		results: []runner.ReportResult{{
			Name:             "api",
			API:              "api",
			Average:          1,
			Max:              1,
			Min:              1,
			QPS:              10,
			Count:            1,
			Error:            1,
			LastErrorMessage: "error",
		}},
		expect: `Name Average Max Min QPS Count Error
api 1ns 1ns 1ns 10 1 1
api error: error
Test case count: 1
`,
	}, {
		name: "have no errors but with message",
		buf:  new(bytes.Buffer),
		results: []runner.ReportResult{{
			Name:             "api",
			API:              "api",
			Average:          1,
			Max:              1,
			Min:              1,
			QPS:              10,
			Count:            1,
			Error:            0,
			LastErrorMessage: "message",
		}},
		expect: `Name Average Max Min QPS Count Error
api 1ns 1ns 1ns 10 1 0
Test case count: 1
`,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := runner.NewResultWriter(tt.buf)
			if !assert.NotNil(t, writer) {
				return
			}

			writer.WithAPICoverage(tt.apiConverage)
			err := writer.Output(tt.results)
			assert.Nil(t, err)
			assert.Equal(t, tt.expect, tt.buf.String())
			assert.NotNil(t, writer.GetWriter())
			assert.NotNil(t, writer.WithResourceUsage(nil))
		})
	}

	discardResultWriter := runner.NewDiscardResultWriter()
	assert.NotNil(t, discardResultWriter)
}
