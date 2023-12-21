/*
Copyright 2023 API Testing Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package runner

import (
	"fmt"
	"time"
)

// TestReporter is the interface of the report
type TestReporter interface {
	PutRecord(*ReportRecord)
	GetAllRecords() []*ReportRecord
	ExportAllReportResults() (ReportResultSlice, error)
	GetResourceUsage() []ResourceUsage
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
		return fmt.Sprintf("Case: %s. error: %v. body: %s", r.Name, r.Error, r.Body)
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

type ResourceUsage struct {
	Memory uint64
	CPU    uint64
	Time   time.Time
}
