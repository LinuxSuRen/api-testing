/**
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

package testing

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/linuxsuren/api-testing/docs"
	"github.com/linuxsuren/api-testing/pkg/render"
	"github.com/linuxsuren/api-testing/pkg/util"
	"github.com/xeipuuv/gojsonschema"
)

const (
	ContextKeyGlobalParam = "param"
)

// Parse parses a file and returns the test suite
func Parse(data []byte) (testSuite *TestSuite, err error) {
	testSuite, err = ParseFromData(data)

	// schema validation
	if err == nil {
		// convert YAML to JSON
		var jsonData []byte
		if jsonData, err = yaml.YAMLToJSON(data); err == nil {
			schemaLoader := gojsonschema.NewStringLoader(docs.Schema)
			documentLoader := gojsonschema.NewBytesLoader(jsonData)

			var result *gojsonschema.Result
			if result, err = gojsonschema.Validate(schemaLoader, documentLoader); err == nil {
				if !result.Valid() {
					err = fmt.Errorf("%v", result.Errors())
				}
			}
		}
	}
	return
}

// ParseFromStream parses the stream and returns the test suite
func ParseFromStream(stream io.Reader) (testSuite *TestSuite, err error) {
	var data []byte
	if data, err = io.ReadAll(stream); err == nil {
		testSuite, err = ParseFromData(data)
	}
	return
}

// ParseFromData parses data and returns the test suite
func ParseFromData(data []byte) (testSuite *TestSuite, err error) {
	testSuite = &TestSuite{}
	if err = yaml.Unmarshal(data, testSuite); err != nil {
		return
	}

	names := map[string]struct{}{}
	for _, item := range testSuite.Items {
		if _, ok := names[item.Name]; !ok {
			names[item.Name] = struct{}{}
		} else {
			err = fmt.Errorf("having duplicated name '%s'", item.Name)
			break
		}
	}
	return
}

// ParseTestCaseFromData parses the data to a test case
func ParseTestCaseFromData(data []byte) (testCase *TestCase, err error) {
	testCase = &TestCase{}
	err = yaml.Unmarshal(data, testCase)
	return
}

// ParseTestSuiteFromFile parses from suite path
func ParseTestSuiteFromFile(suitePath string) (testSuite *TestSuite, err error) {
	var data []byte
	if data, err = os.ReadFile(suitePath); err == nil {
		testSuite = &TestSuite{}
		yaml.Unmarshal(data, testSuite)
	}
	return
}

// GetHeader returns the header of the YAML config file
func GetHeader() string {
	return `#!api-testing
# yaml-language-server: $schema=https://linuxsuren.github.io/api-testing/api-testing-schema.json
`
}

func ToYAML(suite *TestSuite) ([]byte, error) {
	data, err := yaml.Marshal(suite)
	return data, err
}

// SaveTestSuiteToFile saves the test suite to file
func SaveTestSuiteToFile(suite *TestSuite, suitePath string) (err error) {
	var data []byte
	if data, err = ToYAML(suite); err == nil {
		// add header
		data = append([]byte(GetHeader()), data...)
		err = os.WriteFile(suitePath, data, 0644)
	}
	return
}

// Render injects the template based context
func (s *TestSuite) Render(dataContext map[string]interface{}) (err error) {
	// render the API
	var result string
	if result, err = render.Render("base api", s.API, dataContext); err == nil {
		s.API = result
		s.API = strings.TrimSuffix(s.API, "/")
		// render the parameters
		s.Param, err = renderMap(dataContext, s.Param, "parameter")
		dataContext[ContextKeyGlobalParam] = s.Param
	}
	return
}

// Render injects the template based context
func (r *Request) Render(ctx interface{}, dataDir string) (err error) {
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
		if data, err = os.ReadFile(path.Join(dataDir, r.BodyFromFile)); err != nil {
			return
		}
		r.Body = strings.TrimSpace(string(data))
	}

	// template the header
	if r.Header, err = renderMap(ctx, r.Header, "header"); err != nil {
		return
	}

	// template the body
	if result, err = render.Render("body", r.Body, ctx); err == nil {
		r.Body = result
	} else {
		return
	}

	// template the form
	if r.Form, err = renderMap(ctx, r.Form, "form"); err != nil {
		return
	}

	// setting default values
	r.Method = util.EmptyThenDefault(r.Method, http.MethodGet)
	return
}

// RenderAPI will combine with the base API
func (r *Request) RenderAPI(base string) {
	// reuse the API prefix
	if strings.HasPrefix(r.API, "/") {
		r.API = fmt.Sprintf("%s%s", base, r.API)
	}
}

// GetBody returns the request body
func (r *Request) GetBody() (reader io.Reader, err error) {
	if len(r.Form) > 0 {
		if r.Header[util.ContentType] == util.MultiPartFormData {
			multiBody := &bytes.Buffer{}
			writer := multipart.NewWriter(multiBody)
			for key, val := range r.Form {
				writer.WriteField(key, val)
			}

			_ = writer.Close()
			reader = multiBody
			r.Header[util.ContentType] = writer.FormDataContentType()
		} else if r.Header[util.ContentType] == util.Form {
			data := url.Values{}
			for key, val := range r.Form {
				data.Set(key, val)
			}
			reader = strings.NewReader(data.Encode())
		}
	} else if r.Body != "" {
		reader = bytes.NewBufferString(r.Body)
	} else if r.BodyFromFile != "" {
		var data []byte
		if data, err = os.ReadFile(r.BodyFromFile); err == nil {
			reader = bytes.NewBufferString(string(data))
		}
	}
	return
}

// Render renders the response
func (r *Response) Render(ctx interface{}) (err error) {
	r.StatusCode = util.ZeroThenDefault(r.StatusCode, http.StatusOK)

	toDel := []string{}
	for k, v := range r.BodyFieldsExpect {
		var keyStr string
		if keyStr, err = render.Render("bodyFieldsExpect key", k, ctx); err == nil {
			if k != keyStr {
				// means the key is a template string
				toDel = append(toDel, k)
				k = keyStr
			}
		} else {
			return
		}

		valStr, ok := v.(string)
		if !ok {
			continue
		}

		if valStr, err = render.Render("bodyFieldsExpect value", valStr, ctx); err == nil {
			r.BodyFieldsExpect[k] = valStr
		} else {
			return
		}
	}

	for _, k := range toDel {
		delete(r.BodyFieldsExpect, k)
	}
	return
}

func renderMap(ctx interface{}, data map[string]string, title string) (result map[string]string, err error) {
	var tmpVal string
	for key, val := range data {
		if tmpVal, err = render.Render(title, val, ctx); err == nil {
			data[key] = tmpVal
		} else {
			break
		}
	}
	result = data
	return
}
