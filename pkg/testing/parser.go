package testing

import (
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
		testSuite = &TestSuite{}
		err = yaml.Unmarshal(data, testSuite)
	}
	return
}

// Render injects the template based context
func (r *Request) Render(ctx interface{}) (err error) {
	// template the API
	var result string
	if result, err = render.Render("api", r.API, ctx); err == nil {
		r.API = result
	} else {
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
