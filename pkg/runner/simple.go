package runner

import (
	"bytes"
	"fmt"
	"github.com/andreyvit/diff"
	"github.com/linuxsuren/api-testing/pkg/exec"
	"github.com/linuxsuren/api-testing/pkg/testing"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

func RunTestCase(testcase *testing.TestCase) (err error) {
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

	if testcase.Expect.Body != "" {
		var data []byte
		if data, err = ioutil.ReadAll(resp.Body); err != nil {
			return
		}

		if string(data) != strings.TrimSpace(testcase.Expect.Body) {
			err = fmt.Errorf("case: %s, got different response body, diff: \n%s", testcase.Name,
				diff.LineDiff(testcase.Expect.Body, string(data)))
			return
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
