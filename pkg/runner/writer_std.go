package runner

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"text/template"
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
func (w *stdResultWriter) Output(result []ReportResult) error {
	fmt.Fprintf(w.writer, "API Average Max Min Count Error\n")
	for _, r := range result {
		fmt.Fprintf(w.writer, "%s %v %v %v %d %d\n", r.API, r.Average, r.Max,
			r.Min, r.Count, r.Error)
	}
	return nil
}

type markdownResultWriter struct {
	writer io.Writer
}

// NewMarkdownResultWriter creates the Markdown writer
func NewMarkdownResultWriter(writer io.Writer) ReportResultWriter {
	return &markdownResultWriter{writer: writer}
}

// Output writer the Markdown based report to target writer
func (w *markdownResultWriter) Output(result []ReportResult) (err error) {
	var tpl *template.Template
	if tpl, err = template.New("report").Parse(markDownReport); err == nil {
		buf := new(bytes.Buffer)

		if err = tpl.Execute(buf, result); err == nil {
			fmt.Fprint(w.writer, buf.String())
		}
	}
	return
}

//go:embed data/report.md
var markDownReport string
