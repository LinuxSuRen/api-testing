package apispec

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"strings"
)

type Swagger struct {
	Swagger string                           `json:"swagger"`
	Paths   map[string]map[string]SwaggerAPI `json:"paths"`
	Info    SwaggerInfo                      `json:"info"`
}

type SwaggerAPI struct {
	OperationId string `json:"operationId"`
	Summary     string `json:"summary"`
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
