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

package runner

import (
	"fmt"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	"github.com/linuxsuren/api-testing/pkg/runner/kubernetes"
	"github.com/linuxsuren/api-testing/pkg/testing"
)

// Verify if the data satisfies the expression.
func Verify(expect testing.Response, data map[string]any) (err error) {
	for _, verifyExpr := range expect.Verify {
		var ok bool
		if ok, err = verify(verifyExpr, data); !ok {
			err = fmt.Errorf("failed to verify: %q, %v", verifyExpr, err)
			return
		}
	}

	for _, verifyCon := range expect.ConditionalVerify {
		pass := true
		for _, con := range verifyCon.Condition {
			if ok, _ := verify(con, data); !ok {
				pass = false
				break
			}
		}

		if pass {
			for _, verifyExpr := range verifyCon.Verify {
				var ok bool
				if ok, err = verify(verifyExpr, data); !ok {
					err = fmt.Errorf("failed to verify: %q, %v", verifyExpr, err)
					return
				}
			}
		}
	}
	return
}

func verify(verify string, data map[string]any) (ok bool, err error) {
	var program *vm.Program
	if program, err = expr.Compile(verify, expr.Env(data),
		expr.AsBool(), kubernetes.PodValidatorFunc(),
		kubernetes.KubernetesValidatorFunc()); err != nil {
		return
	}

	var result interface{}
	if result, err = expr.Run(program, data); err != nil {
		return
	}

	ok = result.(bool)
	return
}
