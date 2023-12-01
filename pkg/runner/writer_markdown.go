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

type markdownResultWriter struct {
	writer        io.Writer
	apiConverage  apispec.APIConverage
	resourceUsage []ResourceUsage
}

// NewMarkdownResultWriter creates the Markdown writer
func NewMarkdownResultWriter(writer io.Writer) ReportResultWriter {
	return &markdownResultWriter{writer: writer}
}

// Output writes the Markdown based report to target writer
func (w *markdownResultWriter) Output(result []ReportResult) (err error) {
	report := &markdownReport{
		Total: len(result),
		Items: result,
	}
	if len(w.resourceUsage) > 0 {
		report.LastResourceUsage = w.resourceUsage[len(w.resourceUsage)-1]
	}

	for _, item := range result {
		if item.Error > 0 {
			report.Error++
		}
		if item.LastErrorMessage != "" {
			report.Errors = append(report.Errors, item.LastErrorMessage)
		}
	}

	return render.RenderThenPrint("md-report", markdownReportTpl, report, w.writer)
}

// WithAPIConverage sets the api coverage
func (w *markdownResultWriter) WithAPIConverage(apiConverage apispec.APIConverage) ReportResultWriter {
	w.apiConverage = apiConverage
	return w
}

func (w *markdownResultWriter) WithResourceUsage(resurceUage []ResourceUsage) ReportResultWriter {
	w.resourceUsage = resurceUage
	return w
}

type markdownReport struct {
	Total             int
	Error             int
	Items             []ReportResult
	LastResourceUsage ResourceUsage
	Errors            []string
}

//go:embed data/report.md
var markdownReportTpl string
