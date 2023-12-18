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
			Name:      "foo",
			API:       urlFoo,
			Method:    http.MethodGet,
			BeginTime: now,
			EndTime:   now.Add(time.Second * 3),
		}, {
			Name:      "foo",
			API:       urlFoo,
			Method:    http.MethodGet,
			BeginTime: now,
			EndTime:   now.Add(time.Second * 4),
			Error:     errors.New("fake"),
			Body:      "fake",
		}, {
			Name:      "foo",
			API:       urlFoo,
			Method:    http.MethodGet,
			BeginTime: now,
			EndTime:   now.Add(time.Second * 2),
		}, {
			Name:      "bar",
			API:       urlBar,
			Method:    http.MethodGet,
			BeginTime: now,
			EndTime:   now.Add(time.Second),
		}, {
			Name:      "fake",
			API:       urlFake,
			Method:    http.MethodGet,
			BeginTime: now,
			EndTime:   now.Add(time.Second * 5),
		}},
		expect: runner.ReportResultSlice{{
			Name:    "fake",
			API:     "GET " + urlFake,
			Average: time.Second * 5,
			Max:     time.Second * 5,
			Min:     time.Second * 5,
			Count:   1,
			Error:   0,
		}, {
			Name:             "foo",
			API:              "GET http://foo",
			Average:          time.Second * 3,
			Max:              time.Second * 4,
			Min:              time.Second * 2,
			Count:            3,
			Error:            1,
			LastErrorMessage: "Case: foo. error: fake. body: fake",
		}, {
			Name:    "bar",
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
			Name:             "fake",
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
