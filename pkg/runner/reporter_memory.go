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
	"context"
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
		runnerLogger.Info("failed to get resource usage", "error", err)
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
		id := record.Name
		api := record.Method + " " + record.API
		duration := record.Duration()

		if item, ok := resultWithTotal[id]; ok {
			item.Max, item.Min = getMaxAndMin(item.Max, item.Min, duration)
			item.Error += record.ErrorCount()
			item.Total += duration
			item.Count += 1

			item.Last = getLaterTime(record.EndTime, item.Last)
			item.LastErrorMessage = getOriginalStringWhenEmpty(item.LastErrorMessage, record.GetErrorMessage())
		} else {
			resultWithTotal[id] = &ReportResultWithTotal{
				ReportResult: ReportResult{
					Name:  record.Name,
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
			resultWithTotal[id].LastErrorMessage = record.GetErrorMessage()
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
