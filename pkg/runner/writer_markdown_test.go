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
	writer.WithAPIConverage(nil)

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
	assert.Equal(t, `There are 2 test cases:

| API | Average | Max | Min | Count | Error |
|---|---|---|---|---|---|
| api | 3ns | 4ns | 2ns | 3 | 0 |
| api | 3ns | 4ns | 2ns | 3 | 0 |`, buf.String())
}
