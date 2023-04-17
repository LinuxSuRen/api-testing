package runner

import (
	"sort"
	"time"
)

type memoryTestReporter struct {
	records []*ReportRecord
}

// NewMemoryTestReporter creates a memory based test reporter
func NewMemoryTestReporter() TestReporter {
	return &memoryTestReporter{
		records: []*ReportRecord{},
	}
}

// ReportResultWithTotal holds the total duration base on ReportResult
type ReportResultWithTotal struct {
	ReportResult
	Total time.Duration
	First time.Time
	Last  time.Time
}

// PutRecord puts the record to memory
func (r *memoryTestReporter) PutRecord(record *ReportRecord) {
	r.records = append(r.records, record)
}

// GetAllRecords returns all the records
func (r *memoryTestReporter) GetAllRecords() []*ReportRecord {
	return r.records
}

func getMaxAndMin(max, min, duration time.Duration) (time.Duration, time.Duration) {
	if max < duration {
		max = duration
	}

	if min > duration {
		min = duration
	}
	return max, min
}

// ExportAllReportResults exports all the report results
func (r *memoryTestReporter) ExportAllReportResults() (result ReportResultSlice, err error) {
	resultWithTotal := map[string]*ReportResultWithTotal{}
	for _, record := range r.records {
		api := record.Method + " " + record.API
		duration := record.Duration()

		if item, ok := resultWithTotal[api]; ok {
			item.Max, item.Min = getMaxAndMin(item.Max, item.Min, duration)
			item.Error += record.ErrorCount()
			item.Total += duration
			item.Count += 1

			if record.EndTime.After(item.Last) {
				item.Last = record.EndTime
			}
			if record.BeginTime.Before(item.First) {
				item.First = record.BeginTime
			}
		} else {
			resultWithTotal[api] = &ReportResultWithTotal{
				ReportResult: ReportResult{
					API:   api,
					Count: 1,
					Max:   duration,
					Min:   duration,
					Error: record.ErrorCount(),
				},
				First: record.BeginTime,
				Last:  record.EndTime,
				Total: duration,
			}
		}
	}

	for _, r := range resultWithTotal {
		r.Average = r.Total / time.Duration(r.Count)
		if duration := int(r.Last.Sub(r.First).Seconds()); duration > 0 {
			r.QPS = r.Count / duration
		}
		result = append(result, r.ReportResult)
	}

	sort.Sort(result)
	return
}
