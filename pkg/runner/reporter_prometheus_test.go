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

	"github.com/h2non/gock"
	"github.com/linuxsuren/api-testing/pkg/runner"
	"github.com/stretchr/testify/assert"
)

func TestPutRecordToPrometheus(t *testing.T) {
	reporter := runner.NewPrometheusWriter(urlFoo, true)

	assert.Nil(t, reporter.GetAllRecords())
	_, err := reporter.ExportAllReportResults()
	assert.NoError(t, err)

	defer gock.Clean()
	gock.New(urlFoo).Put("/metrics/job/api-testing").Reply(http.StatusBadRequest)

	now := time.Now()
	reporter.PutRecord(&runner.ReportRecord{
		Group:     "foo",
		Name:      "bar",
		BeginTime: now,
		EndTime:   now.Add(time.Second * 3),
	})
	reporter.PutRecord(&runner.ReportRecord{
		Group:     "foo",
		Name:      "bar",
		BeginTime: now,
		EndTime:   now.Add(time.Second * 3),
		Error:     errors.New("fake"),
	})
}
