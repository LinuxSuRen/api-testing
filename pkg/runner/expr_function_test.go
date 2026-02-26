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
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
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
		defer gock.Off()
		gock.New(urlFoo).Reply(http.StatusOK)

		_, err := runner.ExprFuncHTTPReady(urlFoo, 1)
		assert.NoError(t, err)
	})

	t.Run("failed", func(t *testing.T) {
		defer gock.Off()
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

	t.Run("check the response", func(t *testing.T) {
		defer gock.Off()
		gock.New(urlFoo).Reply(http.StatusOK).BodyString(`{"name": "test"}`)
		_, err := runner.ExprFuncHTTPReady(urlFoo, 1, `data.name == "test"`)
		assert.NoError(t, err)
	})

	t.Run("response is not JSON", func(t *testing.T) {
		defer gock.Off()
		gock.New(urlFoo).Reply(http.StatusOK).BodyString(`name: test`)
		_, err := runner.ExprFuncHTTPReady(urlFoo, 1, `data.name == "test"`)
		assert.Error(t, err)
	})

	t.Run("response checking result is failed", func(t *testing.T) {
		defer gock.Off()
		gock.New(urlFoo).Reply(http.StatusOK).BodyString(`{"name": "test"}`)
		_, err := runner.ExprFuncHTTPReady(urlFoo, 1, `data.name == "test"`)
		assert.NoError(t, err)
	})

	t.Run("not a bool expr", func(t *testing.T) {
		defer gock.Off()
		gock.New(urlFoo).Reply(http.StatusOK).BodyString(`{"name": "test"}`)
		_, err := runner.ExprFuncHTTPReady(urlFoo, 1, `name + "test"`)
		assert.Error(t, err)
	})

	t.Run("failed to compile", func(t *testing.T) {
		defer gock.Off()
		gock.New(urlFoo).Reply(http.StatusOK).BodyString(`{"name": "test"}`)
		_, err := runner.ExprFuncHTTPReady(urlFoo, 1, `1~!@`)
		assert.Error(t, err)
	})
}

func TestFunctions(t *testing.T) {
	tmpFile, err := os.CreateTemp(os.TempDir(), "test")
	if err != nil {
		t.Fatal("failed to create temp file")
	}
	defer os.Remove(tmpFile.Name())

	tests := []struct {
		name      string
		expr      string
		syntaxErr bool
		verify    func(t *testing.T, result any, resultErr error)
	}{{
		name:      "invalid syntax",
		expr:      "sleep 1",
		syntaxErr: true,
	}, {
		name: "command",
		expr: `command("echo 1")`,
		verify: func(t *testing.T, result any, resultErr error) {
			assert.NoError(t, resultErr)
			assert.Equal(t, "1", strings.TrimSpace(result.(string)))
		},
	}, {
		name: "writeFile",
		expr: fmt.Sprintf(`writeFile("%s", "hello")`, filepath.ToSlash(tmpFile.Name())),
		verify: func(t *testing.T, result any, resultErr error) {
			assert.NoError(t, resultErr)

			data, err := io.ReadAll(tmpFile)
			assert.NoError(t, err, "failed to read file: %v", err)
			assert.Equal(t, "hello", string(data))
		},
	}}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var program *vm.Program
			program, err = expr.Compile(tt.expr, expr.Env(tt))
			if tt.syntaxErr {
				assert.Error(t, err, "%q %d", tt.name, i)
				return
			}
			if !assert.NotNil(t, program, "%q, index: %d, expr: %s: error: %v", tt.name, i, tt.expr, err) {
				return
			}

			var result any
			result, err = expr.Run(program, tt)
			if tt.verify != nil {
				tt.verify(t, result, err)
			}
		})
	}
}
