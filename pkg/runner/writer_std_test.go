/*
MIT License

Copyright (c) 2023 API Testing Authors.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
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
		apiConverage apispec.APIConverage
		results      []runner.ReportResult
		expect       string
	}{{
		name:    "result is nil",
		buf:     new(bytes.Buffer),
		results: nil,
		expect: `API Average Max Min QPS Count Error
`,
	}, {
		name: "have one item",
		buf:  new(bytes.Buffer),
		apiConverage: apispec.NewFakeAPISpec([][]string{{
			"/api", "GET",
		}}),
		results: []runner.ReportResult{{
			API:     "/api",
			Average: 1,
			Max:     1,
			Min:     1,
			QPS:     10,
			Count:   1,
			Error:   0,
		}},
		expect: `API Average Max Min QPS Count Error
/api 1ns 1ns 1ns 10 1 0

API Coverage: 1/1
`,
	}, {
		name: "have errors",
		buf:  new(bytes.Buffer),
		results: []runner.ReportResult{{
			API:              "api",
			Average:          1,
			Max:              1,
			Min:              1,
			QPS:              10,
			Count:            1,
			Error:            1,
			LastErrorMessage: "error",
		}},
		expect: `API Average Max Min QPS Count Error
api 1ns 1ns 1ns 10 1 1
api error: error
`,
	}, {
		name: "have no errors but with message",
		buf:  new(bytes.Buffer),
		results: []runner.ReportResult{{
			API:              "api",
			Average:          1,
			Max:              1,
			Min:              1,
			QPS:              10,
			Count:            1,
			Error:            0,
			LastErrorMessage: "message",
		}},
		expect: `API Average Max Min QPS Count Error
api 1ns 1ns 1ns 10 1 0
`,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := runner.NewResultWriter(tt.buf)
			if !assert.NotNil(t, writer) {
				return
			}

			writer.WithAPIConverage(tt.apiConverage)
			err := writer.Output(tt.results)
			assert.Nil(t, err)
			assert.Equal(t, tt.expect, tt.buf.String())

			assert.NotNil(t, writer.WithResourceUsage(nil))
		})
	}

	discardResultWriter := runner.NewDiscardResultWriter()
	assert.NotNil(t, discardResultWriter)
}
