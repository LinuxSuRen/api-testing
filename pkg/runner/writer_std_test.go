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
		})
	}

	discardResultWriter := runner.NewDiscardResultWriter()
	assert.NotNil(t, discardResultWriter)
}
