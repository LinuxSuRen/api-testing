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
	"bytes"
	"testing"

	"github.com/linuxsuren/api-testing/pkg/apispec"
	"github.com/linuxsuren/api-testing/pkg/runner"
	"github.com/stretchr/testify/assert"
)

func TestMarkdownWriter(t *testing.T) {
	sample := runner.ReportResult{
		API:     "api",
		Average: 3,
		Max:     4,
		Min:     2,
		Count:   3,
		Error:   0,
	}
	errSample := runner.ReportResult{
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
		writer.WithAPIConverage(nil)
		err := writer.Output(createSlice(sample, 2))
		assert.Nil(t, err)
		assert.Equal(t, `There are 2 test cases:

| API | Average | Max | Min | Count | Error |
|---|---|---|---|---|---|
| api | 3ns | 4ns | 2ns | 3 | 0 |
| api | 3ns | 4ns | 2ns | 3 | 0 |`, buf.String())
	})

	t.Run("long", func(t *testing.T) {
		buf := new(bytes.Buffer)
		writer := runner.NewMarkdownResultWriter(buf)
		writer.WithAPIConverage(nil)
		err := writer.Output(createSlice(sample, 8))
		assert.Nil(t, err)
		assert.Equal(t, `There are 8 test cases:

<details>
  <summary><b>See all test records</b></summary>

| API | Average | Max | Min | Count | Error |
|---|---|---|---|---|---|
| api | 3ns | 4ns | 2ns | 3 | 0 |
| api | 3ns | 4ns | 2ns | 3 | 0 |
| api | 3ns | 4ns | 2ns | 3 | 0 |
| api | 3ns | 4ns | 2ns | 3 | 0 |
| api | 3ns | 4ns | 2ns | 3 | 0 |
| api | 3ns | 4ns | 2ns | 3 | 0 |
| api | 3ns | 4ns | 2ns | 3 | 0 |
| api | 3ns | 4ns | 2ns | 3 | 0 |
</details>`, buf.String())
	})

	t.Run("long, there are error cases", func(t *testing.T) {
		buf := new(bytes.Buffer)
		writer := runner.NewMarkdownResultWriter(buf)
		writer.WithAPIConverage(nil)
		err := writer.Output(append(createSlice(sample, 8), errSample))
		assert.Nil(t, err)
		assert.Equal(t, `There are 9 test cases:

| API | Average | Max | Min | Count | Error |
|---|---|---|---|---|---|
| api | 3ns | 4ns | 2ns | 3 | 1 |

<details>
  <summary><b>See all test records</b></summary>

| API | Average | Max | Min | Count | Error |
|---|---|---|---|---|---|
| api | 3ns | 4ns | 2ns | 3 | 0 |
| api | 3ns | 4ns | 2ns | 3 | 0 |
| api | 3ns | 4ns | 2ns | 3 | 0 |
| api | 3ns | 4ns | 2ns | 3 | 0 |
| api | 3ns | 4ns | 2ns | 3 | 0 |
| api | 3ns | 4ns | 2ns | 3 | 0 |
| api | 3ns | 4ns | 2ns | 3 | 0 |
| api | 3ns | 4ns | 2ns | 3 | 0 |
| api | 3ns | 4ns | 2ns | 3 | 1 |
</details>`, buf.String())
	})

	t.Run("with resource usage", func(t *testing.T) {
		buf := new(bytes.Buffer)
		writer := runner.NewMarkdownResultWriter(buf)
		writer.WithAPIConverage(nil)
		writer.WithResourceUsage([]runner.ResourceUsage{{
			CPU:    1,
			Memory: 1,
		}})
		err := writer.Output(createSlice(sample, 2))
		assert.Nil(t, err)
		assert.Equal(t, `There are 2 test cases:

| API | Average | Max | Min | Count | Error |
|---|---|---|---|---|---|
| api | 3ns | 4ns | 2ns | 3 | 0 |
| api | 3ns | 4ns | 2ns | 3 | 0 |

Resource usage:
* CPU: 1
* Memory: 1`, buf.String())
	})

	t.Run("have error message", func(t *testing.T) {
		buf := new(bytes.Buffer)
		writer := runner.NewMarkdownResultWriter(buf)
		writer.WithAPIConverage(nil)
		result := sample
		result.LastErrorMessage = "error happend"
		err := writer.Output(createSlice(result, 2))
		assert.Nil(t, err)
		assert.Equal(t, `There are 2 test cases:

| API | Average | Max | Min | Count | Error |
|---|---|---|---|---|---|
| api | 3ns | 4ns | 2ns | 3 | 0 |
| api | 3ns | 4ns | 2ns | 3 | 0 |

<details>
  <summary><b>See the error message</b></summary>
* error happend
* error happend
</details>`, buf.String())
	})

	t.Run("with api converage", func(t *testing.T) {
		buf := new(bytes.Buffer)
		writer := runner.NewMarkdownResultWriter(buf)
		writer.WithAPIConverage(apispec.NewFakeAPISpec([][]string{{
			"api", "GET",
		}}))
		err := writer.Output(createSlice(sample, 2))
		assert.Nil(t, err)
		assert.Equal(t, `There are 2 test cases:

| API | Average | Max | Min | Count | Error |
|---|---|---|---|---|---|
| api | 3ns | 4ns | 2ns | 3 | 0 |
| api | 3ns | 4ns | 2ns | 3 | 0 |

API Coverage: 1/1`, buf.String())
	})
}

func createSlice(sample runner.ReportResult, count int) (result []runner.ReportResult) {
	for i := 0; i < count; i++ {
		result = append(result, sample)
	}
	return
}
