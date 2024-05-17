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

package apispec_test

import (
	"github.com/go-openapi/spec"
	"net/http"
	"testing"

	_ "embed"

	"github.com/h2non/gock"
	"github.com/linuxsuren/api-testing/pkg/apispec"
	"github.com/stretchr/testify/assert"
)

func TestParseURLToSwagger(t *testing.T) {
	tests := []struct {
		name       string
		swaggerURL string
		verify     func(t *testing.T, swagger *spec.Swagger, err error)
	}{{
		name:       "normal",
		swaggerURL: urlFoo,
		verify: func(t *testing.T, swagger *spec.Swagger, err error) {
			assert.NoError(t, err)
			assert.Equal(t, "2.0", swagger.Swagger)
			assert.Equal(t, spec.InfoProps{
				Description: "sample",
				Title:       "sample",
				Version:     "1.0.0",
			}, swagger.Info.InfoProps)
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gock.New(tt.swaggerURL).Get("/").Reply(200).BodyString(testdataSwaggerJSON)
			defer gock.Off()

			s, err := apispec.ParseURLToSwagger(tt.swaggerURL)
			tt.verify(t, s, err)
		})
	}
}

func TestHaveAPI(t *testing.T) {
	tests := []struct {
		name         string
		swaggerURL   string
		path, method string
		expectExist  bool
	}{{
		name:        "normal, exist",
		swaggerURL:  urlFoo,
		path:        "/api/v1/users",
		method:      http.MethodGet,
		expectExist: true,
	}, {
		name:        "create user, exist",
		swaggerURL:  urlFoo,
		path:        "/api/v1/users",
		method:      http.MethodPost,
		expectExist: true,
	}, {
		name:        "get a user, exist",
		swaggerURL:  urlFoo,
		path:        "/api/v1/users/linuxsuren",
		method:      http.MethodGet,
		expectExist: true,
	}, {
		name:        "normal, not exist",
		swaggerURL:  urlFoo,
		path:        "/api/v1/users",
		method:      http.MethodDelete,
		expectExist: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gock.New(tt.swaggerURL).Get("/").Reply(200).BodyString(testdataSwaggerJSON)
			defer gock.Off()

			swagger, err := apispec.ParseURLToSwagger(tt.swaggerURL)
			swaggerAPI := apispec.NewSwaggerAPI(swagger)
			assert.NoError(t, err)
			exist := swaggerAPI.HaveAPI(tt.path, tt.method)
			assert.Equal(t, tt.expectExist, exist)
		})
	}
}

func TestAPICount(t *testing.T) {
	tests := []struct {
		name        string
		swaggerURL  string
		expectCount int
	}{{
		name:        "normal",
		swaggerURL:  urlFoo,
		expectCount: 5,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gock.New(tt.swaggerURL).Get("/").Reply(200).BodyString(testdataSwaggerJSON)
			defer gock.Off()

			swagger, err := apispec.ParseURLToSwagger(tt.swaggerURL)
			swaggerAPI := apispec.NewSwaggerAPI(swagger)
			assert.NoError(t, err)
			count := swaggerAPI.APICount()
			assert.Equal(t, tt.expectCount, count)
		})
	}
}

//go:embed testdata/swagger.json
var testdataSwaggerJSON string

const urlFoo = "http://foo"
