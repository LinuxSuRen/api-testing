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

package runner

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"

	"github.com/linuxsuren/api-testing/pkg/apispec"
)

type jsonResultWriter struct {
	writer io.Writer
}

// NewJSONResultWriter creates a new jsonResultWriter
func NewJSONResultWriter(writer io.Writer) ReportResultWriter {
	return &jsonResultWriter{writer: writer}
}

// Output writes the JSON base report to target writer
func (w *jsonResultWriter) Output(result []ReportResult) (err error) {
	jsonData, err := json.Marshal(result)
	if err == nil {
		_, err = fmt.Fprint(w.writer, string(jsonData))
	}
	return
}

// WithAPIConverage sets the api coverage
func (w *jsonResultWriter) WithAPICoverage(apiConverage apispec.APICoverage) ReportResultWriter {
	return w
}

func (w *jsonResultWriter) WithResourceUsage([]ResourceUsage) ReportResultWriter {
	return w
}
