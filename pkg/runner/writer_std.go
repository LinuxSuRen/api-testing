package runner

import (
	_ "embed"
	"fmt"
	"io"
)

type stdResultWriter struct {
	writer io.Writer
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
	return nil
}
