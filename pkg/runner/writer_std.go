/*
Copyright 2023-2024 API Testing Authors.

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
	"fmt"
	"io"

	"github.com/linuxsuren/api-testing/pkg/apispec"
)

type stdResultWriter struct {
	writer       io.Writer
	apiConverage apispec.APICoverage
}

// NewResultWriter creates a result writer with the specific io.Writer
func NewResultWriter(writer io.Writer) ReportResultWriter {
	return &stdResultWriter{writer: writer}
}

// NewDiscardResultWriter creates a report result writer which discard everything
func NewDiscardResultWriter() ReportResultWriter {
	return &stdResultWriter{writer: io.Discard}
}

// Output writer the report to target writer
func (w *stdResultWriter) Output(results []ReportResult) error {
	var errResults []ReportResult
	_, _ = fmt.Fprintf(w.writer, "Name Average Max Min QPS Count Error\n")
	for _, r := range results {
		_, _ = fmt.Fprintf(w.writer, "%s %v %v %v %d %d %d\n", r.Name, r.Average, r.Max,
			r.Min, r.QPS, r.Count, r.Error)
		if r.Error > 0 && r.LastErrorMessage != "" {
			errResults = append(errResults, r)
		}
	}

	for _, r := range errResults {
		_, _ = fmt.Fprintf(w.writer, "%s error: %s\n", r.API, r.LastErrorMessage)
	}

	_, _ = fmt.Fprintf(w.writer, "Test case count: %d\n", len(results))
	apiConveragePrint(results, w.apiConverage, w.writer)
	return nil
}

// WithAPIConverage sets the api coverage
func (w *stdResultWriter) WithAPICoverage(apiConverage apispec.APICoverage) ReportResultWriter {
	w.apiConverage = apiConverage
	return w
}

func (w *stdResultWriter) WithResourceUsage([]ResourceUsage) ReportResultWriter {
	return w
}

func (w *stdResultWriter) GetWriter() io.Writer {
	return w.writer
}

func apiConveragePrint(result []ReportResult, apiConverage apispec.APICoverage, w io.Writer) {
	covered, total := apiConverageCount(result, apiConverage)
	if total > 0 {
		fmt.Fprintf(w, "\nAPI Coverage: %d/%d\n", covered, total)
	}
}

func apiConverageCount(result []ReportResult, apiConverage apispec.APICoverage) (covered, total int) {
	if apiConverage == nil {
		return
	}

	for _, item := range result {
		if apiConverage.HaveAPI(item.API, "GET") {
			covered++
		}
	}
	total = apiConverage.APICount()
	if covered > total {
		covered = total
	}
	return covered, total
}
