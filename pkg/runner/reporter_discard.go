package runner

type discardTestReporter struct {
}

// NewDiscardTestReporter creates a test reporter which discard everything
func NewDiscardTestReporter() TestReporter {
	return &discardTestReporter{}
}

// PutRecord does nothing
func (r *discardTestReporter) PutRecord(*ReportRecord) {
	// Do nothing which is the design purpose
}

// GetAllRecords does nothing
func (r *discardTestReporter) GetAllRecords() []*ReportRecord {
	return nil
}

// ExportAllReportResults does nothing
func (r *discardTestReporter) ExportAllReportResults() (ReportResultSlice, error) {
	return nil, nil
}
