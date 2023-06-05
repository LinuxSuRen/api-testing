package pkg

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	atestpkg "github.com/linuxsuren/api-testing/pkg/testing"
	"gopkg.in/yaml.v2"
)

// SampleExporter is a sample exporter
type SampleExporter struct {
	TestSuite atestpkg.TestSuite
}

// NewSampleExporter creates a new exporter
func NewSampleExporter() *SampleExporter {
	return &SampleExporter{
		TestSuite: atestpkg.TestSuite{
			Name: "sample",
		},
	}
}

// Add adds a request to the exporter
func (e *SampleExporter) Add(r *http.Request) {

	fmt.Println("receive", r.URL.Path)
	req := atestpkg.Request{
		API:    r.URL.String(),
		Method: r.Method,
		Header: map[string]string{},
	}

	if body := r.Body; body != nil {
		if data, err := io.ReadAll(body); err == nil {
			req.Body = string(data)
		}
	}

	testCase := atestpkg.TestCase{
		Request: req,
		Expect: atestpkg.Response{
			StatusCode: http.StatusOK,
		},
	}

	specs := strings.Split(r.URL.Path, "/")
	if len(specs) > 0 {
		testCase.Name = specs[len(specs)-1]
	}

	if val := r.Header.Get("Content-Type"); val != "" {
		req.Header["Content-Type"] = val
	}
	if val := r.Header.Get("Authorization"); val != "" {
		req.Header["Authorization"] = val
	}

	e.TestSuite.Items = append(e.TestSuite.Items, testCase)
}

var prefix = `#!api-testing
# yaml-language-server: $schema=https://gitee.com/linuxsuren/api-testing/raw/master/sample/api-testing-schema.json
`

// Export exports the test suite
func (e *SampleExporter) Export() (string, error) {
	marker := map[string]int{}

	for i, item := range e.TestSuite.Items {
		if _, ok := marker[item.Name]; ok {
			marker[item.Name]++
			e.TestSuite.Items[i].Name = fmt.Sprintf("%s-%d", item.Name, marker[item.Name])
		} else {
			marker[item.Name] = 0
		}
	}

	data, err := yaml.Marshal(e.TestSuite)
	return prefix + string(data), err
}
