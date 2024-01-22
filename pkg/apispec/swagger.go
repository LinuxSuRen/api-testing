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

package apispec

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"strings"
)

type Swagger struct {
	Swagger string `json:"swagger"`
	// Paths includes all the API requests.
	// The keys is the HTTP request method which as lower-case, for example: get, post.
	Paths map[string]map[string]SwaggerAPI `json:"paths"`
	Info  SwaggerInfo                      `json:"info"`
}

type SwaggerAPI struct {
	OperationId string      `json:"operationId"`
	Parameters  []Parameter `json:"parameters"`
	Summary     string      `json:"summary"`
}

type Parameter struct {
	Name string `json:"name"`
	// In represents the parameter type, supported values: query, path
	In       string `json:"in"`
	Required bool   `json:"required"`
	Schema   Schema `json:"schema"`
}

type Schema struct {
	Type   string `json:"type"`
	Format string `json:"format"`
}

type SwaggerInfo struct {
	Description string `json:"description"`
	Title       string `json:"title"`
	Version     string `json:"version"`
}

type APIConverage interface {
	HaveAPI(path, method string) (exist bool)
	APICount() (count int)
}

// HaveAPI check if the swagger has the API.
// If the path is /api/v1/names/linuxsuren, then will match /api/v1/names/{name}
func (s *Swagger) HaveAPI(path, method string) (exist bool) {
	method = strings.ToLower(method)
	for item := range s.Paths {
		if matchAPI(path, item) {
			for m := range s.Paths[item] {
				if strings.ToLower(m) == method {
					exist = true
					return
				}
			}
		}
	}
	return
}

func matchAPI(particularAPI, swaggerAPI string) (matched bool) {
	result := swaggerAPIConvert(swaggerAPI)
	reg, err := regexp.Compile(result)
	if err == nil {
		matched = reg.MatchString(particularAPI)
	}
	return
}

func swaggerAPIConvert(text string) (result string) {
	result = text
	reg, err := regexp.Compile("{.*}")
	if err == nil {
		result = reg.ReplaceAllString(text, ".*")
	}
	return
}

// APICount return the count of APIs
func (s *Swagger) APICount() (count int) {
	for path := range s.Paths {
		for range s.Paths[path] {
			count++
		}
	}
	return
}

func ParseToSwagger(data []byte) (swagger *Swagger, err error) {
	swagger = &Swagger{}
	err = json.Unmarshal(data, swagger)
	return
}

func ParseURLToSwagger(swaggerURL string) (swagger *Swagger, err error) {
	var resp *http.Response
	if resp, err = http.Get(swaggerURL); err == nil && resp != nil && resp.StatusCode == http.StatusOK {
		swagger, err = ParseStreamToSwagger(resp.Body)
	}
	return
}

func ParseStreamToSwagger(stream io.Reader) (swagger *Swagger, err error) {
	var data []byte
	if data, err = io.ReadAll(stream); err == nil {
		swagger, err = ParseToSwagger(data)
	}
	return
}
