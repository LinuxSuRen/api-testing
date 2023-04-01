package runner

type discardTestReporter struct {
}

// NewDiscardTestReporter creates a test reporter which discard everything
func NewDiscardTestReporter() TestReporter {
	return &discardTestReporter{}
}

func (r *discardTestReporter) PutRecord(*ReportRecord) {}
func (r *discardTestReporter) GetAllRecords() []*ReportRecord {
	return nil
}
func (r *discardTestReporter) ExportAllReportResults() (ReportResultSlice, error) {
	return nil, nil
}
