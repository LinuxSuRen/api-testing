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
	Item    []PostmanItem  `json:"item"`
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
		if len(item.Item) == 0 {
			suite.Items[i] = testing.TestCase{
				Name: item.Name,
				Request: testing.Request{
					Method: item.Request.Method,
					API:    item.Request.URL.Raw,
					Body:   testing.NewRequestBody(item.Request.Body.Raw),
					Header: item.Request.Header.ToMap(),
				},
			}
		} else {
			for _, sub := range item.Item {
				suite.Items[i] = testing.TestCase{
					Name: item.Name + " " + sub.Name,
					Request: testing.Request{
						Method: sub.Request.Method,
						API:    sub.Request.URL.Raw,
						Body:   testing.NewRequestBody(item.Request.Body.Raw),
						Header: sub.Request.Header.ToMap(),
					},
				}
			}
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
