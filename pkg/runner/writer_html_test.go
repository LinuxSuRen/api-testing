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
			assert.NotNil(t, w.WithResourceUsage(nil))
		})
	}
}

//go:embed testdata/report.html
var htmlReportExpect string
