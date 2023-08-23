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

package generator

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/linuxsuren/api-testing/pkg/testing"
)

type PostmanCollection struct {
	Collection Postman `json:"collection"`
}

type Postman struct {
	Info PostmanInfo   `json:"info"`
	Item []PostmanItem `json:"item"`
}

type PostmanInfo struct {
	Name string
}

type PostmanItem struct {
	Name    string         `json:"name"`
	Request PostmanRequest `json:"request"`
}

type PostmanRequest struct {
	Method string      `json:"method"`
	URL    PostmanURL  `json:"url"`
	Header Paris       `json:"header"`
	Body   PostmanBody `json:"body"`
}

type PostmanBody struct {
	Mode string `json:"mode"`
	Raw  string `json:"raw"`
}

type PostmanURL struct {
	Raw   string   `json:"raw"`
	Path  []string `json:"path"`
	Query Paris    `json:"query"`
}

type Paris []Pair
type Pair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (p Paris) ToMap() (result map[string]string) {
	count := len(p)
	if count == 0 {
		return
	}
	result = make(map[string]string, count)
	for _, item := range p {
		result[item.Key] = item.Value
	}
	return
}

type Importer interface {
	Convert(data []byte) (*testing.TestSuite, error)
	ConvertFromFile(dataFile string) (*testing.TestSuite, error)
	ConvertFromURL(dataURL string) (*testing.TestSuite, error)
}

type postmanImporter struct {
}

// NewPostmanImporter returns a new postman importer
func NewPostmanImporter() Importer {
	return &postmanImporter{}
}

// Convert converts the postman data to test suite
func (p *postmanImporter) Convert(data []byte) (suite *testing.TestSuite, err error) {
	postman := &Postman{}
	if err = json.Unmarshal(data, postman); err != nil {
		return
	}
	if postman.Info.Name == "" {
		postmanCollection := &PostmanCollection{}
		if err = json.Unmarshal(data, postmanCollection); err != nil {
			return
		}
		postman = &postmanCollection.Collection
	}

	suite = &testing.TestSuite{}
	suite.Name = postman.Info.Name
	suite.Items = make([]testing.TestCase, len(postman.Item))

	for i, item := range postman.Item {
		suite.Items[i] = testing.TestCase{
			Name: item.Name,
			Request: testing.Request{
				Method: item.Request.Method,
				API:    item.Request.URL.Raw,
				Body:   item.Request.Body.Raw,
				Header: item.Request.Header.ToMap(),
			},
		}
	}
	return
}

func (p *postmanImporter) ConvertFromFile(dataFile string) (suite *testing.TestSuite, err error) {
	var data []byte
	if data, err = os.ReadFile(dataFile); err == nil {
		suite, err = p.Convert(data)
	}
	return
}

func (p *postmanImporter) ConvertFromURL(dataURL string) (suite *testing.TestSuite, err error) {
	var resp *http.Response
	if resp, err = http.Get(dataURL); err == nil {
		var data []byte
		if data, err = io.ReadAll(resp.Body); err == nil {
			suite, err = p.Convert(data)
		}
	}
	return
}
