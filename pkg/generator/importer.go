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
	Mode       string      `json:"mode"`
	Raw        string      `json:"raw"`
	FormData   []TypeField `json:"formdata"`
	URLEncoded []TypeField `json:"urlencoded"`
	Disabled   bool        `json:"disabled"`
}

type TypeField struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	Disabled    bool   `json:"disabled"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Src         string `json:"src"`
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
	nativeImporter
}

// NewPostmanImporter returns a new postman importer
func NewPostmanImporter() Importer {
	return &postmanImporter{
		nativeImporter: nativeImporter{},
	}
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
	if err = p.convertItems(postman.Item, "", suite); err != nil {
		return
	}

	return
}

func (p *postmanImporter) convertItems(items []PostmanItem, prefix string, suite *testing.TestSuite) (err error) {
	for _, item := range items {
		itemName := prefix + item.Name
		if len(item.Item) == 0 {
			suite.Items = append(suite.Items, testing.TestCase{
				Name: itemName,
				Request: testing.Request{
					Method: item.Request.Method,
					API:    item.Request.URL.Raw,
					Body:   testing.NewRequestBody(item.Request.Body.Raw),
					Header: item.Request.Header.ToMap(),
				},
			})
		} else {
			if err = p.convertItems(item.Item, itemName+" ", suite); err != nil {
				return
			}
		}
	}
	return
}
