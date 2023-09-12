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
	"net/http"
	"testing"

	"github.com/h2non/gock"
	"github.com/linuxsuren/api-testing/pkg/runner"
	"github.com/stretchr/testify/assert"
)

func TestExprFuncSleep(t *testing.T) {
	tests := []struct {
		name   string
		params []interface{}
		hasErr bool
	}{{
		name:   "string format duration",
		params: []interface{}{"0.01s"},
		hasErr: false,
	}, {
		name:   "without params",
		hasErr: true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := runner.ExprFuncSleep(tt.params...)
			if tt.hasErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestExprFuncHTTPReady(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		defer gock.Clean()
		gock.New(urlFoo).Reply(http.StatusOK)

		_, err := runner.ExprFuncHTTPReady(urlFoo, 1)
		assert.NoError(t, err)
	})

	t.Run("failed", func(t *testing.T) {
		defer gock.Clean()
		gock.New(urlFoo).Reply(http.StatusNotFound)

		_, err := runner.ExprFuncHTTPReady(urlFoo, 1)
		assert.Error(t, err)
	})

	t.Run("params less than 2", func(t *testing.T) {
		_, err := runner.ExprFuncHTTPReady()
		assert.Error(t, err)
	})

	t.Run("API param is not a string", func(t *testing.T) {
		_, err := runner.ExprFuncHTTPReady(1, 2)
		assert.Error(t, err)
	})

	t.Run("retry param is not an integer", func(t *testing.T) {
		_, err := runner.ExprFuncHTTPReady(urlFoo, "two")
		assert.Error(t, err)
	})
}
