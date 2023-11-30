/*
MIT License

Copyright (c) 2023 API Testing Authors.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package runner

import (
	"context"
	"log"
	"sort"
	"time"

	"github.com/linuxsuren/api-testing/pkg/runner/monitor"
)

type memoryTestReporter struct {
	records        []*ReportRecord
	resourceUsages []ResourceUsage
	resMonitor     monitor.MonitorClient
	monitorTarget  string
}

// NewMemoryTestReporter creates a memory based test reporter
func NewMemoryTestReporter(resMonitor monitor.MonitorClient, monitorTarget string) TestReporter {
	if resMonitor == nil {
		resMonitor = monitor.NewDumyMonitor()
	}
	return &memoryTestReporter{
		records:       []*ReportRecord{},
		resMonitor:    resMonitor,
		monitorTarget: monitorTarget,
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
	usage, err := r.resMonitor.GetResourceUsage(context.TODO(), &monitor.Target{
		Name: r.monitorTarget,
	})
	if err != nil {
		log.Println("failed to get resource usage:", err)
	} else {
		r.resourceUsages = append(r.resourceUsages, ResourceUsage{
			Memory: usage.Memory,
			CPU:    usage.Cpu,
			Time:   time.Now(),
		})
	}
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

			item.Last = getLaterTime(record.EndTime, item.Last)
			item.LastErrorMessage = getOriginalStringWhenEmpty(item.LastErrorMessage, record.GetErrorMessage())
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
			resultWithTotal[api].LastErrorMessage = record.GetErrorMessage()
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

func (r *memoryTestReporter) GetResourceUsage() []ResourceUsage {
	return r.resourceUsages
}

func getLaterTime(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}

func getOriginalStringWhenEmpty(a, b string) string {
	if b == "" {
		return a
	}
	return b
}
