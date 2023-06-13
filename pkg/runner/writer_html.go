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

//go:embed data/html.html
var htmlReport string
