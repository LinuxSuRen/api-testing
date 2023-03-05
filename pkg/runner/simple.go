package runner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strings"

	"github.com/andreyvit/diff"
	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	"github.com/linuxsuren/api-testing/pkg/exec"
	"github.com/linuxsuren/api-testing/pkg/testing"
	unstructured "github.com/linuxsuren/unstructured/pkg"
)

// RunTestCase runs the test case
func RunTestCase(testcase *testing.TestCase, ctx interface{}) (output interface{}, err error) {
	fmt.Printf("start to run: '%s'\n", testcase.Name)
	if err = doPrepare(testcase); err != nil {
		err = fmt.Errorf("failed to prepare, error: %v", err)
		return
	}

	defer func() {
		if testcase.Clean.CleanPrepare {
			if err = doCleanPrepare(testcase); err != nil {
				return
			}
		}
	}()

	client := http.Client{}
	var requestBody io.Reader
	if testcase.Request.Body != "" {
		requestBody = bytes.NewBufferString(testcase.Request.Body)
	} else if testcase.Request.BodyFromFile != "" {
		var data []byte
		if data, err = os.ReadFile(testcase.Request.BodyFromFile); err != nil {
			return
		}
		requestBody = bytes.NewBufferString(string(data))
	}

	if err = testcase.Request.Render(ctx); err != nil {
		return
	}

	var request *http.Request
	if request, err = http.NewRequest(testcase.Request.Method, testcase.Request.API, requestBody); err != nil {
		return
	}

	// set headers
	for key, val := range testcase.Request.Header {
		request.Header.Add(key, val)
	}

	fmt.Println("start to send request to", testcase.Request.API)

	// send the HTTP request
	var resp *http.Response
	if resp, err = client.Do(request); err != nil {
		return
	}

	if testcase.Expect.StatusCode != 0 {
		if err = expectInt(testcase.Name, testcase.Expect.StatusCode, resp.StatusCode); err != nil {
			return
		}
	}

	for key, val := range testcase.Expect.Header {
		actualVal := resp.Header.Get(key)
		if err = expectString(testcase.Name, val, actualVal); err != nil {
			return
		}
	}

	var responseBodyData []byte
	if responseBodyData, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}
	if testcase.Expect.Body != "" {
		if string(responseBodyData) != strings.TrimSpace(testcase.Expect.Body) {
			err = fmt.Errorf("case: %s, got different response body, diff: \n%s", testcase.Name,
				diff.LineDiff(testcase.Expect.Body, string(responseBodyData)))
			return
		}
	}

	mapOutput := map[string]interface{}{}
	if err = json.Unmarshal(responseBodyData, &mapOutput); err != nil {
		switch b := err.(type) {
		case *json.UnmarshalTypeError:
			if b.Value != "array" {
				return
			} else {
				arrayOutput := []interface{}{}
				if err = json.Unmarshal(responseBodyData, &arrayOutput); err != nil {
					return
				}
				output = arrayOutput
			}
		default:
			return
		}
	} else {
		output = mapOutput
	}

	for key, expectVal := range testcase.Expect.BodyFieldsExpect {
		var val interface{}
		var ok bool
		if val, ok, err = unstructured.NestedField(mapOutput, strings.Split(key, "/")...); err != nil {
			err = fmt.Errorf("failed to get field: %s, %v", key, err)
			return
		} else if !ok {
			err = fmt.Errorf("not found field: %s", key)
			return
		} else if !reflect.DeepEqual(expectVal, val) {
			err = fmt.Errorf("field[%s] expect value: %v, actual: %v", key, expectVal, val)
			return
		}
	}

	for _, verify := range testcase.Expect.Verify {
		var program *vm.Program
		if program, err = expr.Compile(verify, expr.Env(output), expr.AsBool()); err != nil {
			return
		}

		var result interface{}
		if result, err = expr.Run(program, output); err != nil {
			return
		}

		if !result.(bool) {
			err = fmt.Errorf("faild to verify: %s", verify)
			break
		}
	}
	return
}

func doPrepare(testcase *testing.TestCase) (err error) {
	for i := range testcase.Prepare.Kubernetes {
		item := testcase.Prepare.Kubernetes[i]

		if err = exec.RunCommand("kubectl", "apply", "-f", item); err != nil {
			return
		}
	}
	return
}

func doCleanPrepare(testcase *testing.TestCase) (err error) {
	count := len(testcase.Prepare.Kubernetes)
	for i := count - 1; i >= 0; i-- {
		item := testcase.Prepare.Kubernetes[i]

		if err = exec.RunCommand("kubectl", "delete", "-f", item); err != nil {
			return
		}
	}
	return
}

func expectInt(name string, expect, actual int) (err error) {
	if expect != actual {
		err = fmt.Errorf("case: %s, expect %d, actual %d", name, expect, actual)
	}
	return
}

func expectString(name, expect, actual string) (err error) {
	if expect != actual {
		err = fmt.Errorf("case: %s, expect %s, actual %s", name, expect, actual)
	}
	return
}
