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
