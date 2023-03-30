package runner_test

import (
	"bytes"
	"testing"

	"github.com/linuxsuren/api-testing/pkg/runner"
	"github.com/stretchr/testify/assert"
)

func TestMarkdownWriter(t *testing.T) {
	buf := new(bytes.Buffer)
	writer := runner.NewMarkdownResultWriter(buf)

	err := writer.Output([]runner.ReportResult{{
		API:     "api",
		Average: 3,
		Max:     4,
		Min:     2,
		Count:   3,
		Error:   0,
	}, {
		API:     "api",
		Average: 3,
		Max:     4,
		Min:     2,
		Count:   3,
		Error:   0,
	}})
	assert.Nil(t, err)
	assert.Equal(t, `| API | Average | Max | Min | Count | Error |
|---|---|---|---|---|---|
| api | 3ns | 4ns | 2ns | 3 | 0 |
| api | 3ns | 4ns | 2ns | 3 | 0 |
`, buf.String())
}

func TestNewStdResultWriter(t *testing.T) {
	tests := []struct {
		name    string
		buf     *bytes.Buffer
		results []runner.ReportResult
		expect  string
	}{{
		name:    "result is nil",
		buf:     new(bytes.Buffer),
		results: nil,
		expect: `API Average Max Min Count Error
`,
	}, {
		name: "have one item",
		buf:  new(bytes.Buffer),
		results: []runner.ReportResult{{
			API:     "api",
			Average: 1,
			Max:     1,
			Min:     1,
			Count:   1,
			Error:   0,
		}},
		expect: `API Average Max Min Count Error
api 1ns 1ns 1ns 1 0
`,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := runner.NewResultWriter(tt.buf)
			if !assert.NotNil(t, writer) {
				return
			}

			err := writer.Output(tt.results)
			assert.Nil(t, err)
			assert.Equal(t, tt.expect, tt.buf.String())
		})
	}

	discardResultWriter := runner.NewDiscardResultWriter()
	assert.NotNil(t, discardResultWriter)
}
