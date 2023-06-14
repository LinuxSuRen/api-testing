package runner

import (
	_ "embed"
	"io"

	"github.com/linuxsuren/api-testing/pkg/apispec"
	"github.com/linuxsuren/api-testing/pkg/render"
)

type markdownResultWriter struct {
	writer       io.Writer
	apiConverage apispec.APIConverage
}

// NewMarkdownResultWriter creates the Markdown writer
func NewMarkdownResultWriter(writer io.Writer) ReportResultWriter {
	return &markdownResultWriter{writer: writer}
}

// Output writes the Markdown based report to target writer
func (w *markdownResultWriter) Output(result []ReportResult) (err error) {
	return render.RenderThenPrint("md-report", markdownReport, result, w.writer)
}

// WithAPIConverage sets the api coverage
func (w *markdownResultWriter) WithAPIConverage(apiConverage apispec.APIConverage) ReportResultWriter {
	w.apiConverage = apiConverage
	return w
}

//go:embed data/report.md
var markdownReport string
