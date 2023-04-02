package testing

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/linuxsuren/api-testing/pkg/render"
	"gopkg.in/yaml.v2"
)

// Parse parses a file and returns the test suite
func Parse(configFile string) (testSuite *TestSuite, err error) {
	var data []byte
	if data, err = os.ReadFile(configFile); err == nil {
		testSuite, err = ParseFromData(data)
	}
	return
}

// ParseFromData parses data and returns the test suite
func ParseFromData(data []byte) (testSuite *TestSuite, err error) {
	testSuite = &TestSuite{}
	err = yaml.Unmarshal(data, testSuite)
	return
}

// ParseTestCaseFromData parses the data to a test case
func ParseTestCaseFromData(data []byte) (testCase *TestCase, err error) {
	testCase = &TestCase{}
	err = yaml.Unmarshal(data, testCase)
	return
}

// Render injects the template based context
func (r *Request) Render(ctx interface{}) (err error) {
	// template the API
	var result string
	if result, err = render.Render("api", r.API, ctx); err == nil {
		r.API = result
	} else {
		err = fmt.Errorf("failed render '%s', %v", r.API, err)
		return
	}

	// read body from file
	if r.BodyFromFile != "" {
		var data []byte
		if data, err = os.ReadFile(r.BodyFromFile); err != nil {
			return
		}
		r.Body = strings.TrimSpace(string(data))
	}

	// template the header
	for key, val := range r.Header {
		if result, err = render.Render("header", val, ctx); err == nil {
			r.Header[key] = result
		} else {
			return
		}
	}

	// template the body
	if result, err = render.Render("body", r.Body, ctx); err == nil {
		r.Body = result
	} else {
		return
	}

	// template the form
	for key, val := range r.Form {
		if result, err = render.Render("form", val, ctx); err == nil {
			r.Form[key] = result
		} else {
			return
		}
	}

	// setting default values
	r.Method = emptyThenDefault(r.Method, http.MethodGet)
	return
}

// Render renders the response
func (r *Response) Render(ctx interface{}) (err error) {
	r.StatusCode = zeroThenDefault(r.StatusCode, http.StatusOK)
	return
}

func zeroThenDefault(val, defVal int) int {
	if val == 0 {
		val = defVal
	}
	return val
}

func emptyThenDefault(val, defVal string) string {
	if strings.TrimSpace(val) == "" {
		val = defVal
	}
	return val
}
