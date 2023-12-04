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

// Package runner provides the common expr style functions
package runner

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/ast"
	"github.com/expr-lang/expr/builtin"
	"github.com/expr-lang/expr/vm"
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
			log.Println("checking the response")
			exprText := params[2].(string)

			// check the response
			var data []byte
			if data, err = io.ReadAll(resp.Body); err == nil {
				unstruct := make(map[string]interface{})

				if err = json.Unmarshal(data, &unstruct); err != nil {
					log.Printf("failed to unmarshal the response data: %v\n", err)
					return
				}

				unstruct["data"] = unstruct
				var program *vm.Program
				if program, err = expr.Compile(exprText, expr.Env(unstruct)); err != nil {
					log.Printf("failed to compile: %s, %v\n", exprText, err)
					return
				}

				var result interface{}
				if result, err = expr.Run(program, unstruct); err != nil {
					log.Printf("failed to Run: %s, %v\n", exprText, err)
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

		log.Println("waiting for", api)
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
