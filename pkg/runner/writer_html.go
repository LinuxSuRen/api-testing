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
	"io"

	"github.com/linuxsuren/api-testing/pkg/apispec"
	"github.com/linuxsuren/api-testing/pkg/render"
)

type htmlResultWriter struct {
	writer       io.Writer
	apiConverage apispec.APIConverage
}

// NewHTMLResultWriter creates a new htmlResultWriter
func NewHTMLResultWriter(writer io.Writer) ReportResultWriter {
	return &htmlResultWriter{writer: writer}
}

// Output writes the HTML base report to target writer
func (w *htmlResultWriter) Output(result []ReportResult) (err error) {
	return render.RenderThenPrint("html-report", htmlReport, result, w.writer)
}

// WithAPIConverage sets the api coverage
func (w *htmlResultWriter) WithAPIConverage(apiConverage apispec.APIConverage) ReportResultWriter {
	w.apiConverage = apiConverage
	return w
}

func (w *htmlResultWriter) WithResourceUsage([]ResourceUsage) ReportResultWriter {
	return w
}

//go:embed data/html.html
var htmlReport string
