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
