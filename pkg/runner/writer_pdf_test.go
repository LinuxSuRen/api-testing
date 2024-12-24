/*
Copyright 2024 API Testing Authors.

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

package runner

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestPDF(t *testing.T) {
	tests := []struct {
		name    string
		buf     io.Writer
		results []ReportResult
		verify  func(t *testing.T)
		hasErr  bool
	}{{
		name: "normal",
		buf:  new(bytes.Buffer),
		results: []ReportResult{{
			Name:    "/api",
			API:     "/api",
			Average: 1,
			Max:     1,
			Min:     1,
			QPS:     10,
			Count:   1,
			Error:   0,
		}},
		verify: func(t *testing.T) {

		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := NewPDFResultWriter(tt.buf)
			if !assert.NotNil(t, writer) {
				return
			}

			err := writer.Output(tt.results)
			if tt.hasErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			tt.verify(t)
			writer.WithResourceUsage(nil)
			writer.WithAPICoverage(nil)
		})
	}
}
