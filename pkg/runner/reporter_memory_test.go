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

package runner_test

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/linuxsuren/api-testing/pkg/runner"
	"github.com/stretchr/testify/assert"
)

const urlFoo = "http://foo"
const urlBar = "http://bar"
const urlFake = "http://fake"

func TestExportAllReportResults(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name    string
		records []*runner.ReportRecord
		expect  runner.ReportResultSlice
	}{{
		name:    "no records",
		records: []*runner.ReportRecord{},
		expect:  nil,
	}, {
		name: "normal",
		records: []*runner.ReportRecord{{
			API:       urlFoo,
			Method:    http.MethodGet,
			BeginTime: now,
			EndTime:   now.Add(time.Second * 3),
		}, {
			API:       urlFoo,
			Method:    http.MethodGet,
			BeginTime: now,
			EndTime:   now.Add(time.Second * 4),
			Error:     errors.New("fake"),
			Body:      "fake",
		}, {
			API:       urlFoo,
			Method:    http.MethodGet,
			BeginTime: now,
			EndTime:   now.Add(time.Second * 2),
		}, {
			API:       urlBar,
			Method:    http.MethodGet,
			BeginTime: now,
			EndTime:   now.Add(time.Second),
		}, {
			API:       urlFake,
			Method:    http.MethodGet,
			BeginTime: now,
			EndTime:   now.Add(time.Second * 5),
		}},
		expect: runner.ReportResultSlice{{
			API:     "GET " + urlFake,
			Average: time.Second * 5,
			Max:     time.Second * 5,
			Min:     time.Second * 5,
			Count:   1,
			Error:   0,
		}, {
			API:              "GET http://foo",
			Average:          time.Second * 3,
			Max:              time.Second * 4,
			Min:              time.Second * 2,
			Count:            3,
			Error:            1,
			LastErrorMessage: "Case: . error: fake. body: fake",
		}, {
			API:     "GET http://bar",
			Average: time.Second,
			Max:     time.Second,
			Min:     time.Second,
			QPS:     1,
			Count:   1,
			Error:   0,
		}},
	}, {
		name: "first record has error",
		records: []*runner.ReportRecord{{
			Name:      "fake",
			API:       urlFoo,
			Method:    http.MethodGet,
			BeginTime: now,
			EndTime:   now.Add(time.Second * 4),
			Error:     errors.New("fake"),
			Body:      "fake",
		}},
		expect: runner.ReportResultSlice{{
			API:              "GET http://foo",
			Average:          time.Second * 4,
			Max:              time.Second * 4,
			Min:              time.Second * 4,
			Count:            1,
			Error:            1,
			LastErrorMessage: "Case: fake. error: fake. body: fake",
		}},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reporter := runner.NewMemoryTestReporter(nil, "")
			assert.NotNil(t, reporter)

			for i := range tt.records {
				reporter.PutRecord(tt.records[i])
			}
			assert.Equal(t, tt.records, reporter.GetAllRecords())

			result, err := reporter.ExportAllReportResults()
			assert.Nil(t, err)
			assert.Equal(t, tt.expect, result)

			assert.Equal(t, len(tt.records), len(reporter.GetResourceUsage()))
		})
	}
}
