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

type markdownResultWriter struct {
	writer        io.Writer
	apiConverage  apispec.APICoverage
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
	report.Converage.Covered, report.Converage.Total = apiConverageCount(result, w.apiConverage)

	return render.RenderThenPrint("md-report", markdownReportTpl, report, w.writer)
}

// WithAPIConverage sets the api coverage
func (w *markdownResultWriter) WithAPICoverage(apiConverage apispec.APICoverage) ReportResultWriter {
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
	Converage         converage
}

type converage struct {
	Covered int
	Total   int
}

//go:embed data/report.md
var markdownReportTpl string
