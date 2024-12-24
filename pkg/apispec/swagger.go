/*
Copyright 2023-2024 API Testing Authors.

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
	"github.com/go-openapi/spec"
	"io"
	"net/http"
	"regexp"
	"strings"
)

type APICoverage interface {
	HaveAPI(path, method string) (exist bool)
	APICount() (count int)
}

type SwaggerAPI struct {
	Swagger *spec.Swagger
	ApiMap  map[string][]string
}

func NewSwaggerAPI(swagger *spec.Swagger) *SwaggerAPI {
	return &SwaggerAPI{
		Swagger: swagger,
		ApiMap:  buildAPIMap(swagger),
	}
}

func buildAPIMap(swagger *spec.Swagger) map[string][]string {
	apiMap := make(map[string][]string)
	for path, pathItem := range swagger.Paths.Paths {
		var methods []string
		if pathItem.Get != nil {
			methods = append(methods, "get")
		}
		if pathItem.Put != nil {
			methods = append(methods, "put")
		}
		if pathItem.Post != nil {
			methods = append(methods, "post")
		}
		if pathItem.Delete != nil {
			methods = append(methods, "delete")
		}
		if pathItem.Options != nil {
			methods = append(methods, "options")
		}
		if pathItem.Head != nil {
			methods = append(methods, "head")
		}
		if pathItem.Patch != nil {
			methods = append(methods, "patch")
		}
		apiMap[path] = methods
	}
	return apiMap
}

// HaveAPI check if the swagger has the API.
// If the path is /api/v1/names/linuxsuren, then will match /api/v1/names/{name}
func (s *SwaggerAPI) HaveAPI(path, method string) (exist bool) {
	method = strings.ToLower(method)
	for p := range s.ApiMap {
		if matchAPI(path, p) {
			if methods, ok := s.ApiMap[p]; ok {
				for _, m := range methods {
					if m == method {
						exist = true
						return
					}
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
func (s *SwaggerAPI) APICount() (count int) {
	for _, methods := range s.ApiMap {
		count += len(methods)
	}
	return
}

func ParseToSwagger(data []byte) (swagger *spec.Swagger, err error) {
	swagger = &spec.Swagger{}
	err = swagger.UnmarshalJSON(data)
	return
}

func ParseURLToSwagger(swaggerURL string) (swagger *spec.Swagger, err error) {
	var resp *http.Response
	if resp, err = http.Get(swaggerURL); err == nil && resp != nil && resp.StatusCode == http.StatusOK {
		swagger, err = ParseStreamToSwagger(resp.Body)
	}
	return
}

func ParseStreamToSwagger(stream io.Reader) (swagger *spec.Swagger, err error) {
	var data []byte
	if data, err = io.ReadAll(stream); err == nil {
		swagger, err = ParseToSwagger(data)
	}
	return
}
