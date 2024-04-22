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

// Package runner provides the common expr style functions
package runner

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/ast"
	"github.com/expr-lang/expr/builtin"
	"github.com/expr-lang/expr/vm"
	"github.com/linuxsuren/api-testing/pkg/logging"
)

var (
	runnerLogger = logging.DefaultLogger(logging.LogLevelInfo).WithName("runner")
)

// ExprFuncSleep is an expr function for sleeping
func ExprFuncSleep(params ...interface{}) (res interface{}, err error) {
	if len(params) < 1 {
		err = fmt.Errorf("the duration param is required")
		return
	}

	switch duration := params[0].(type) {
	case int:
		time.Sleep(time.Duration(duration) * time.Second)
	case string:
		var dur time.Duration
		if dur, err = time.ParseDuration(duration); err == nil {
			time.Sleep(dur)
		}
	}
	return
}

// ExprFuncHTTPReady is an expr function for reading status from a HTTP server
func ExprFuncHTTPReady(params ...interface{}) (res interface{}, err error) {
	if len(params) < 2 {
		err = fmt.Errorf("usage: api retry")
		return
	}

	api, ok := params[0].(string)
	if !ok {
		err = fmt.Errorf("the API param should be a string")
		return
	}

	retry, ok := params[1].(int)
	if !ok {
		err = fmt.Errorf("the retry param should be a integer")
		return
	}

	var resp *http.Response
	for i := 0; i < retry; i++ {
		resp, err = http.Get(api)
		alive := err == nil && resp != nil && resp.StatusCode == http.StatusOK

		if alive && len(params) >= 3 {
			runnerLogger.Info("checking the response")
			exprText := params[2].(string)

			// check the response
			var data []byte
			if data, err = io.ReadAll(resp.Body); err == nil {
				unstruct := make(map[string]interface{})

				if err = json.Unmarshal(data, &unstruct); err != nil {
					runnerLogger.Info("failed to unmarshal the response data ", err)
					return
				}

				unstruct["data"] = unstruct
				var program *vm.Program
				if program, err = expr.Compile(exprText, expr.Env(unstruct)); err != nil {
					runnerLogger.Info("failed to compile ", exprText, "error: ", err)
					return
				}

				var result interface{}
				if result, err = expr.Run(program, unstruct); err != nil {
					runnerLogger.Info("failed to Run ", exprText, "error: ", err)
					return
				}

				if val, ok := result.(bool); ok {
					if val {
						return
					}
				} else {
					err = fmt.Errorf("the result of %s should be a bool", exprText)
					return
				}
			}
		} else if alive {
			return
		}

		runnerLogger.Info("waiting for", api)
		time.Sleep(1 * time.Second)
	}
	err = fmt.Errorf("failed to wait for the API ready in %d times", retry)
	return
}

func init() {
	builtin.Builtins = append(builtin.Builtins, []*ast.Function{
		{
			Name: "sleep",
			Func: ExprFuncSleep,
		},
		{
			Name: "httpReady",
			Func: ExprFuncHTTPReady,
		},
		{
			Name: "command",
			Func: func(params ...interface{}) (res any, err error) {
				var output []byte
				output, err = exec.Command("sh", "-c", params[0].(string)).CombinedOutput()
				if output != nil {
					res = string(output)
				}
				return
			},
		},
		{
			Name: "writeFile",
			Func: func(params ...interface{}) (res any, err error) {
				filename := params[0]
				content := params[1]

				err = os.WriteFile(filename.(string), []byte(content.(string)), 0644)
				return
			},
		},
	}...)
}
