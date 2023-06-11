package runner

import (
	_ "embed"
	"io"

	"github.com/linuxsuren/api-testing/pkg/render"
)

type markdownResultWriter struct {
	writer io.Writer
}

// NewMarkdownResultWriter creates the Markdown writer
func NewMarkdownResultWriter(writer io.Writer) ReportResultWriter {
	return &markdownResultWriter{writer: writer}
}

// Output writes the Markdown based report to target writer
func (w *markdownResultWriter) Output(result []ReportResult) (err error) {
	return render.RenderThenPrint("md-report", markdownReport, result, w.writer)
}

//go:embed data/report.md
var markdownReport string
