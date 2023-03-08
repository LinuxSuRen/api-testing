package testing

import (
	"bytes"
	"html/template"
	"os"
	"strings"

	"github.com/Masterminds/sprig/v3"
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
	var tpl *template.Template
	if tpl, err = template.New("api").Funcs(sprig.FuncMap()).Parse(r.API); err != nil {
		return
	}
	buf := new(bytes.Buffer)
	if err = tpl.Execute(buf, ctx); err != nil {
		return
	}
	r.API = buf.String()

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
		if tpl, err = template.New("header").Funcs(sprig.FuncMap()).Parse(val); err == nil {
			buf = new(bytes.Buffer)
			if err = tpl.Execute(buf, ctx); err == nil {
				r.Header[key] = buf.String()
			}
		}
	}

	// template the body
	if tpl, err = template.New("body").Funcs(sprig.FuncMap()).Parse(r.Body); err == nil {
		buf = new(bytes.Buffer)
		if err = tpl.Execute(buf, ctx); err == nil {
			r.Body = buf.String()
		}
	}

	// template the form
	for key, val := range r.Form {
		if tpl, err = template.New("form").Funcs(sprig.FuncMap()).Parse(val); err == nil {
			buf = new(bytes.Buffer)
			if err = tpl.Execute(buf, ctx); err == nil {
				r.Form[key] = buf.String()
			}
		}
	}
	return
}
