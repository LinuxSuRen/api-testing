package runner

import "time"

// TestReporter is the interface of the report
type TestReporter interface {
	PutRecord(*ReportRecord)
	GetAllRecords() []*ReportRecord
	ExportAllReportResults() (ReportResultSlice, error)
}

// ReportRecord represents the raw data of a request
type ReportRecord struct {
	Group     string
	Name      string
	Method    string
	API       string
	Body      string
	BeginTime time.Time
	EndTime   time.Time
	Error     error
}

// Duration returns the duration between begin and end time
func (r *ReportRecord) Duration() time.Duration {
	return r.EndTime.Sub(r.BeginTime)
}

// ErrorCount returns the count number of errors
func (r *ReportRecord) ErrorCount() int {
	if r.Error == nil {
		return 0
	}
	return 1
}

// GetErrorMessage returns the error message
func (r *ReportRecord) GetErrorMessage() string {
	if r.ErrorCount() > 0 {
		return r.Body
	} else {
		return ""
	}
}

// NewReportRecord creates a record, and set the begin time to be now
func NewReportRecord() *ReportRecord {
	return &ReportRecord{
		BeginTime: time.Now(),
	}
}
