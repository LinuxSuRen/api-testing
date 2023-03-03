package testing

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"os"

	"github.com/Masterminds/sprig/v3"
	"gopkg.in/yaml.v2"
)

func Parse(configFile string) (testSuite *TestSuite, err error) {
	var data []byte
	if data, err = ioutil.ReadFile(configFile); err != nil {
		return
	}

	testSuite = &TestSuite{}
	if err = yaml.Unmarshal(data, testSuite); err != nil {
		return
	}

	// for i, testCase := range testSuite.Items {
	// 	// template the API
	// 	var tpl *template.Template
	// 	if tpl, err = template.New("base").Funcs(sprig.FuncMap()).Parse(testCase.Request.API); err != nil {
	// 		return
	// 	}
	// 	buf := new(bytes.Buffer)
	// 	if err = tpl.Execute(buf, os.Environ()); err != nil {
	// 		return
	// 	} else {
	// 		testCase.Request.API = buf.String()
	// 	}

	// 	// read body from file
	// 	if testCase.Request.BodyFromFile != "" {
	// 		if data, err = os.ReadFile(testCase.Request.BodyFromFile); err != nil {
	// 			return
	// 		}
	// 		testCase.Request.Body = string(data)
	// 	}

	// 	// template the body
	// 	if tpl, err = template.New("base").Funcs(sprig.FuncMap()).Parse(testCase.Request.Body); err != nil {
	// 		return
	// 	}
	// 	buf = new(bytes.Buffer)
	// 	if err = tpl.Execute(buf, os.Environ()); err != nil {
	// 		return
	// 	} else {
	// 		testCase.Request.Body = buf.String()
	// 	}

	// 	testSuite.Items[i] = testCase
	// }
	return
}

func (r *Request) Render(ctx interface{}) (err error) {
	// template the API
	var tpl *template.Template
	if tpl, err = template.New("base").Funcs(sprig.FuncMap()).Parse(r.API); err != nil {
		return
	}
	buf := new(bytes.Buffer)
	if err = tpl.Execute(buf, ctx); err != nil {
		return
	} else {
		r.API = buf.String()
	}

	// read body from file
	if r.BodyFromFile != "" {
		var data []byte
		if data, err = os.ReadFile(r.BodyFromFile); err != nil {
			return
		}
		r.Body = string(data)
	}

	// template the body
	if tpl, err = template.New("base").Funcs(sprig.FuncMap()).Parse(r.Body); err != nil {
		return
	}
	buf = new(bytes.Buffer)
	if err = tpl.Execute(buf, ctx); err != nil {
		return
	} else {
		r.Body = buf.String()
	}
	return
}
