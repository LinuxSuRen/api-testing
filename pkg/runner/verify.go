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
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
	"github.com/linuxsuren/api-testing/pkg/runner/kubernetes"
	"github.com/linuxsuren/api-testing/pkg/testing"
	"github.com/linuxsuren/api-testing/pkg/util"
	"github.com/tidwall/gjson"
	"gopkg.in/yaml.v3"
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
	if result, err = expr.Run(program, data); err == nil {
		ok = result.(bool)
	}
	return
}

type BodyVerifier interface {
	Parse(data []byte) (interface{}, error)
	Verify(data []byte) error
}

type BodyGetter interface {
	GetBody() string
	GetBodyFieldsExpect() map[string]interface{}
}

func NewBodyVerify(contentType string, body BodyGetter) BodyVerifier {
	switch contentType {
	case util.JSON:
		return &jsonBodyVerifier{body: body}
	case util.YAML:
		return &yamlBodyVerifier{body: body}
	default:
		return nil
	}
}

type jsonBodyVerifier struct {
	body BodyGetter
}

func (v *jsonBodyVerifier) Parse(data []byte) (obj interface{}, err error) {
	mapOutput := map[string]interface{}{}
	if err = json.Unmarshal(data, &mapOutput); err != nil {
		switch b := err.(type) {
		case *json.UnmarshalTypeError:
			if b.Value != "array" {
				return
			}

			var arrayOutput []interface{}
			if err = json.Unmarshal(data, &arrayOutput); err == nil {
				obj = arrayOutput
			}
		}
	} else {
		obj = mapOutput
	}
	return
}

func (v *jsonBodyVerifier) Verify(data []byte) (err error) {
	for key, expectVal := range v.body.GetBodyFieldsExpect() {
		result := gjson.Get(string(data), key)
		if result.Exists() {
			err = valueCompare(expectVal, result, key)
		} else {
			err = fmt.Errorf("not found field: %s", key)
		}

		if err != nil {
			break
		}
	}
	return
}

type yamlBodyVerifier struct {
	body BodyGetter
}

func (v *yamlBodyVerifier) Parse(data []byte) (obj interface{}, err error) {
	obj = map[string]interface{}{}
	err = yaml.Unmarshal(data, &obj)
	return
}

func (v *yamlBodyVerifier) Verify(data []byte) (err error) {
	// TODO need to implement
	return
}

func valueCompare(expect interface{}, acutalResult gjson.Result, key string) (err error) {
	var actual interface{}
	actual = acutalResult.Value()

	if !reflect.DeepEqual(expect, actual) {
		switch acutalResult.Type {
		case gjson.Number:
			expect = fmt.Sprintf("%v", expect)
			actual = fmt.Sprintf("%v", actual)

			if strings.Compare(expect.(string), actual.(string)) == 0 {
				return
			}
		}
		err = fmt.Errorf("field[%s] expect value: '%v', actual: '%v'", key, expect, actual)
	}
	return
}
