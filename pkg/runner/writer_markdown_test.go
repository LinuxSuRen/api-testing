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
	"bytes"
	"strings"
	"testing"

	"github.com/linuxsuren/api-testing/pkg/apispec"
	"github.com/linuxsuren/api-testing/pkg/runner"
	"github.com/stretchr/testify/assert"
)

const (
	winLineEnd = "\r\n"
	unixLineEnd = "\n"
)

func TestMarkdownWriter(t *testing.T) {
	sample := runner.ReportResult{
		Name:    "api",
		API:     "api",
		Average: 3,
		Max:     4,
		Min:     2,
		Count:   3,
		Error:   0,
	}
	errSample := runner.ReportResult{
		Name:    "foo",
		API:     "api",
		Average: 3,
		Max:     4,
		Min:     2,
		Count:   3,
		Error:   1,
	}

	t.Run("short", func(t *testing.T) {
		buf := new(bytes.Buffer)
		writer := runner.NewMarkdownResultWriter(buf)
		writer.WithAPICoverage(nil)
		err := writer.Output(createSlice(sample, 2))
		assert.Nil(t, err)
		// Normalize line endings for Windows compatibility
		actual := buf.String()
		actual = normalizeLineEndings(actual)
		assert.Equal(t, `There are 2 test cases, failed count 0:

| Name | Average | Max | Min | Count | Error |
|---|---|---|---|---|---|
| api | 3ns | 4ns | 2ns | 3 | 0 |
| api | 3ns | 4ns | 2ns | 3 | 0 |`, actual)
	})

	t.Run("long", func(t *testing.T) {
		buf := new(bytes.Buffer)
		writer := runner.NewMarkdownResultWriter(buf)
		writer.WithAPICoverage(nil)
		err := writer.Output(createSlice(sample, 8))
		assert.Nil(t, err)
		actual := normalizeLineEndings(buf.String())
		assert.Equal(t, `There are 8 test cases, failed count 0:

<details>
  <summary><b>See all test records</b></summary>

| Name | Average | Max | Min | Count | Error |
|---|---|---|---|---|---|
| api | 3ns | 4ns | 2ns | 3 | 0 |
| api | 3ns | 4ns | 2ns | 3 | 0 |
| api | 3ns | 4ns | 2ns | 3 | 0 |
| api | 3ns | 4ns | 2ns | 3 | 0 |
| api | 3ns | 4ns | 2ns | 3 | 0 |
| api | 3ns | 4ns | 2ns | 3 | 0 |
| api | 3ns | 4ns | 2ns | 3 | 0 |
| api | 3ns | 4ns | 2ns | 3 | 0 |
</details>`, actual)
	})

	t.Run("long, there are error cases", func(t *testing.T) {
		buf := new(bytes.Buffer)
		writer := runner.NewMarkdownResultWriter(buf)
		writer.WithAPICoverage(nil)
		err := writer.Output(append(createSlice(sample, 8), errSample))
		assert.Nil(t, err)
		actual := normalizeLineEndings(buf.String())
		assert.Equal(t, `There are 9 test cases, failed count 1:

| Name | Average | Max | Min | Count | Error |
|---|---|---|---|---|---|
| foo | 3ns | 4ns | 2ns | 3 | 1 |

<details>
  <summary><b>See all test records</b></summary>

| Name | Average | Max | Min | Count | Error |
|---|---|---|---|---|---|
| api | 3ns | 4ns | 2ns | 3 | 0 |
| api | 3ns | 4ns | 2ns | 3 | 0 |
| api | 3ns | 4ns | 2ns | 3 | 0 |
| api | 3ns | 4ns | 2ns | 3 | 0 |
| api | 3ns | 4ns | 2ns | 3 | 0 |
| api | 3ns | 4ns | 2ns | 3 | 0 |
| api | 3ns | 4ns | 2ns | 3 | 0 |
| api | 3ns | 4ns | 2ns | 3 | 0 |
| foo | 3ns | 4ns | 2ns | 3 | 1 |
</details>`, actual)
	})

	t.Run("with resource usage", func(t *testing.T) {
		buf := new(bytes.Buffer)
		writer := runner.NewMarkdownResultWriter(buf)
		writer.WithAPICoverage(nil)
		writer.WithResourceUsage([]runner.ResourceUsage{{
			CPU:    1,
			Memory: 1,
		}})
		err := writer.Output(createSlice(sample, 2))
		assert.Nil(t, err)
		actual := normalizeLineEndings(buf.String())
		assert.Equal(t, `There are 2 test cases, failed count 0:

| Name | Average | Max | Min | Count | Error |
|---|---|---|---|---|---|
| api | 3ns | 4ns | 2ns | 3 | 0 |
| api | 3ns | 4ns | 2ns | 3 | 0 |

Resource usage:
* CPU: 1
* Memory: 1`, actual)
	})

	t.Run("have error message", func(t *testing.T) {
		buf := new(bytes.Buffer)
		writer := runner.NewMarkdownResultWriter(buf)
		writer.WithAPICoverage(nil)
		result := sample
		result.LastErrorMessage = "error happend"
		err := writer.Output(createSlice(result, 2))
		assert.Nil(t, err)
		actual := normalizeLineEndings(buf.String())
		assert.Equal(t, `There are 2 test cases, failed count 0:

| Name | Average | Max | Min | Count | Error |
|---|---|---|---|---|---|
| api | 3ns | 4ns | 2ns | 3 | 0 |
| api | 3ns | 4ns | 2ns | 3 | 0 |

<details>
  <summary><b>See the error message</b></summary>
* error happend
* error happend
</details>`, actual)
	})

	t.Run("with api converage", func(t *testing.T) {
		buf := new(bytes.Buffer)
		writer := runner.NewMarkdownResultWriter(buf)
		writer.WithAPICoverage(apispec.NewFakeAPISpec([][]string{{
			"api", "GET",
		}}))
		err := writer.Output(createSlice(sample, 2))
		assert.Nil(t, err)
		actual := normalizeLineEndings(buf.String())
		assert.Equal(t, `There are 2 test cases, failed count 0:

| Name | Average | Max | Min | Count | Error |
|---|---|---|---|---|---|
| api | 3ns | 4ns | 2ns | 3 | 0 |
| api | 3ns | 4ns | 2ns | 3 | 0 |

API Coverage: 1/1`, actual)
	})
}

func normalizeLineEndings(s string) string {
	return strings.ReplaceAll(s, winLineEnd, unixLineEnd)
}

func createSlice(sample runner.ReportResult, count int) (result []runner.ReportResult) {
	for i := 0; i < count; i++ {
		result = append(result, sample)
	}
	return
}
