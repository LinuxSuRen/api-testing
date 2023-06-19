package runner_test

import (
	"bytes"
	"testing"

	"github.com/linuxsuren/api-testing/pkg/runner"
	"github.com/stretchr/testify/assert"
)

func TestJSONResultWriter(t *testing.T) {
	buf := new(bytes.Buffer)
	writer := runner.NewJSONResultWriter(buf)
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
	assert.Equal(t,
		"[{\"API\":\"api\",\"Count\":3,\"Average\":3,\"Max\":4,\"Min\":2,\"QPS\":0,\"Error\":0,\"LastErrorMessage\":\"\"},{\"API\":\"api\",\"Count\":3,\"Average\":3,\"Max\":4,\"Min\":2,\"QPS\":0,\"Error\":0,\"LastErrorMessage\":\"\"}]",
		buf.String())
}
