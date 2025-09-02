/*
Copyright 2024-2025 API Testing Authors.

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
	Name      string `yaml:"name" json:"name"`
	InitCount *int   `yaml:"initCount" json:"initCount"`
	Sample    string `yaml:"sample" json:"sample"`
}

type Item struct {
	Name     string   `yaml:"name" json:"name"`
	Request  Request  `yaml:"request" json:"request"`
	Response Response `yaml:"response" json:"response"`
	Param    map[string]interface{}
}

type Request struct {
	Protocol string            `yaml:"protocol" json:"protocol"`
	Path     string            `yaml:"path" json:"path"`
	Method   string            `yaml:"method" json:"method"`
	Header   map[string]string `yaml:"header" json:"header"`
	Body     string            `yaml:"body" json:"body"`
}

type RequestWithAuth struct {
	Request   `yaml:",inline"`
	BearerAPI string `yaml:"bearerAPI" json:"bearerAPI"`
	Username  string `yaml:"username" json:"username"`
	Password  string `yaml:"password" json:"password"`
}

type Response struct {
	Encoder      string            `yaml:"encoder" json:"encoder"`
	Body         string            `yaml:"body" json:"body"`
	BodyFromFile string            `yaml:"bodyFromFile" json:"bodyFromFile"`
	Header       map[string]string `yaml:"header" json:"header"`
	StatusCode   int               `yaml:"statusCode" json:"statusCode"`
	BodyData     []byte
}

type Webhook struct {
	Name    string            `yaml:"name" json:"name"`
	Timer   string            `yaml:"timer" json:"timer"`
	Param   map[string]string `yaml:"param" json:"param"`
	Request RequestWithAuth   `yaml:"request" json:"request"`
}

type Proxy struct {
	Prefix       string       `yaml:"prefix" json:"prefix"`
	Port         int          `yaml:"port" json:"port"`
	Path         string       `yaml:"path" json:"path"`
	Target       string       `yaml:"target" json:"target"`
	RequestAmend RequestAmend `yaml:"requestAmend" json:"requestAmend"`
	Protocol     string       `yaml:"protocol" json:"protocol"`
	Echo         bool         `yaml:"echo" json:"echo"`
}

type RequestAmend struct {
	BodyPatch string `yaml:"bodyPatch" json:"bodyPatch"`
}

type Server struct {
	Objects  []Object  `yaml:"objects" json:"objects"`
	Items    []Item    `yaml:"items" json:"items"`
	Proxies  []Proxy   `yaml:"proxies" json:"proxies"`
	Webhooks []Webhook `yaml:"webhooks" json:"webhooks"`
}
