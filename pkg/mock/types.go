/*
Copyright 2024 API Testing Authors.

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
package mock

type Object struct {
	Name      string  `yaml:"name"`
	InitCount *int    `yaml:"initCount"`
	Sample    string  `yaml:"sample"`
	Fields    []Field `yaml:"fields"`
}

type Field struct {
	Name string `yaml:"name"`
	Kind string `yaml:"kind"`
}

type Item struct {
	Name     string   `yaml:"name"`
	Request  Request  `yaml:"request"`
	Response Response `yaml:"response"`
	Param    map[string]string
}

type Request struct {
	Path   string            `yaml:"path"`
	Method string            `yaml:"method"`
	Header map[string]string `yaml:"header"`
	Body   string            `yaml:"body"`
}

type Response struct {
	Encoder    string            `yaml:"encoder"`
	Body       string            `yaml:"body"`
	Header     map[string]string `yaml:"header"`
	StatusCode int               `yaml:"statusCode"`
	BodyData   []byte
}

type Webhook struct {
	Name    string  `yaml:"name"`
	Timer   string  `yaml:"timer"`
	Request Request `yaml:"request"`
}

type Server struct {
	Objects  []Object  `yaml:"objects"`
	Items    []Item    `yaml:"items"`
	Webhooks []Webhook `yaml:"webhooks"`
}
