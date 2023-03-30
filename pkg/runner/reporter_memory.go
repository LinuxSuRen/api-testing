package runner

import (
	"sort"
	"time"
)

type memoryTestReporter struct {
	records []*ReportRecord
}

// NewmemoryTestReporter creates a memory based test reporter
func NewmemoryTestReporter() TestReporter {
	return &memoryTestReporter{
		records: []*ReportRecord{},
	}
}

type ReportResultWithTotal struct {
	ReportResult
	Total time.Duration
}

func (r *memoryTestReporter) PutRecord(record *ReportRecord) {
	r.records = append(r.records, record)
}
func (r *memoryTestReporter) GetAllRecords() []*ReportRecord {
	return r.records
}
func (r *memoryTestReporter) ExportAllReportResults() (result ReportResultSlice, err error) {
	resultWithTotal := map[string]*ReportResultWithTotal{}
	for _, record := range r.records {
		api := record.Method + " " + record.API
		duration := record.Duration()

		if item, ok := resultWithTotal[api]; ok {
			if item.Max < duration {
				item.Max = duration
			}

			if item.Min > duration {
				item.Min = duration
			}
			item.Error += record.ErrorCount()
			item.Total += duration
			item.Count += 1
		} else {
			resultWithTotal[api] = &ReportResultWithTotal{
				ReportResult: ReportResult{
					API:   api,
					Count: 1,
					Max:   duration,
					Min:   duration,
					Error: record.ErrorCount(),
				},
				Total: duration,
			}
		}
	}

	for _, r := range resultWithTotal {
		r.Average = r.Total / time.Duration(r.Count)
		result = append(result, r.ReportResult)
	}

	sort.Sort(result)
	return
}
