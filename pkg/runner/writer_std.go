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
	"fmt"
	"io"

	"github.com/linuxsuren/api-testing/pkg/apispec"
)

type stdResultWriter struct {
	writer       io.Writer
	apiConverage apispec.APIConverage
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
	fmt.Fprintf(w.writer, "API Average Max Min QPS Count Error\n")
	for _, r := range results {
		fmt.Fprintf(w.writer, "%s %v %v %v %d %d %d\n", r.API, r.Average, r.Max,
			r.Min, r.QPS, r.Count, r.Error)
		if r.Error > 0 && r.LastErrorMessage != "" {
			errResults = append(errResults, r)
		}
	}

	for _, r := range errResults {
		fmt.Fprintf(w.writer, "%s error: %s\n", r.API, r.LastErrorMessage)
	}

	apiConveragePrint(results, w.apiConverage, w.writer)
	return nil
}

// WithAPIConverage sets the api coverage
func (w *stdResultWriter) WithAPIConverage(apiConverage apispec.APIConverage) ReportResultWriter {
	w.apiConverage = apiConverage
	return w
}

func apiConveragePrint(result []ReportResult, apiConverage apispec.APIConverage, w io.Writer) {
	if apiConverage == nil {
		return
	}

	var covered int
	for _, item := range result {
		if apiConverage.HaveAPI(item.API, "GET") {
			covered++
		}
	}
	fmt.Fprintf(w, "\nAPI Coverage: %d/%d\n", covered, apiConverage.APICount())
}
