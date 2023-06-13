package runner_test

import (
	"bytes"
	"testing"

	_ "embed"

	"github.com/linuxsuren/api-testing/pkg/runner"
	"github.com/stretchr/testify/assert"
)

func TestHTMLResultWriter(t *testing.T) {
	tests := []struct {
		name    string
		buf     *bytes.Buffer
		results []runner.ReportResult
		expect  string
	}{{
		name: "simple",
		buf:  new(bytes.Buffer),
		results: []runner.ReportResult{{
			API:     "/foo",
			Max:     3,
			Min:     3,
			Average: 3,
			Error:   0,
			Count:   1,
		}},
		expect: htmlReportExpect,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := runner.NewHTMLResultWriter(tt.buf)
			w.WithAPIConverage(nil)
			err := w.Output(tt.results)
			assert.NoError(t, err)
			assert.Equal(t, tt.expect, tt.buf.String())
		})
	}
}

//go:embed testdata/report.html
var htmlReportExpect string
