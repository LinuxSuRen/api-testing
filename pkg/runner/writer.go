package runner

import "github.com/linuxsuren/api-testing/pkg/apispec"

// ReportResultWriter is the interface of the report writer
type ReportResultWriter interface {
	Output([]ReportResult) error
	WithAPIConverage(apiConverage apispec.APIConverage) ReportResultWriter
}
