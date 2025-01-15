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
package testing

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/linuxsuren/api-testing/pkg/util"
	"gopkg.in/yaml.v3"
)

// TestSuite represents a set of test cases
type TestSuite struct {
	Name  string            `yaml:"name,omitempty" json:"name,omitempty"`
	API   string            `yaml:"api,omitempty" json:"api,omitempty"`
	Spec  APISpec           `yaml:"spec,omitempty" json:"spec,omitempty"`
	Param map[string]string `yaml:"param,omitempty" json:"param,omitempty"`
	Items []TestCase        `yaml:"items,omitempty" json:"items,omitempty"`
	Proxy *Proxy            `yaml:"proxy,omitempty" json:"proxy,omitempty"`
}

type APISpec struct {
	Kind   string   `yaml:"kind,omitempty" json:"kind,omitempty"`
	URL    string   `yaml:"url,omitempty" json:"url,omitempty"`
	RPC    *RPCDesc `yaml:"rpc,omitempty" json:"rpc,omitempty"`
	Secure *Secure  `yaml:"secure,omitempty" json:"secure,omitempty"`
	Metric *Metric  `yaml:"metric,omitempty" json:"metric,omitempty"`
}

type HistoryTestSuite struct {
	HistorySuiteName string            `yaml:"name,omitempty" json:"name,omitempty"`
	Items            []HistoryTestCase `yaml:"items,omitempty" json:"items,omitempty"`
}

type HistoryTestCase struct {
	ID               string            `yaml:"id,omitempty" json:"id,omitempty"`
	CaseName         string            `yaml:"caseName,omitempty" json:"name,omitempty"`
	SuiteName        string            `yaml:"suiteName,omitempty" json:"suiteName,omitempty"`
	HistorySuiteName string            `yaml:"historySuiteName,omitempty" json:"historySuiteName,omitempty"`
	CreateTime       time.Time         `yaml:"createTime,omitempty" json:"createTime,omitempty"`
	SuiteAPI         string            `yaml:"api,omitempty" json:"api,omitempty"`
	SuiteSpec        APISpec           `yaml:"spec,omitempty" json:"spec,omitempty"`
	SuiteParam       map[string]string `yaml:"param,omitempty" json:"param,omitempty"`
	Data             TestCase          `yaml:"data,omitempty" json:"data,omitempty"`
	HistoryHeader    map[string]string `yaml:"historyHeader,omitempty" json:"historyHeader,omitempty"`
}

type HistoryTestResult struct {
	Message        string           `yaml:"message,omitempty" json:"message,omitempty"`
	Error          string           `yaml:"error,omitempty" json:"error,omitempty"`
	TestCaseResult []TestCaseResult `yaml:"testCaseResult,omitempty" json:"testCaseResult,omitempty"`
	Data           HistoryTestCase  `yaml:"data,omitempty" json:"data,omitempty"`
	CreateTime     time.Time        `yaml:"createTime,omitempty" json:"createTime,omitempty"`
}

type RPCDesc struct {
	ImportPath       []string `yaml:"import,omitempty" json:"import,omitempty"`
	ServerReflection bool     `yaml:"serverReflection,omitempty" json:"serverReflection,omitempty"`
	ProtoFile        string   `yaml:"protofile,omitempty" json:"protofile,omitempty"`
	ProtoSet         string   `yaml:"protoset,omitempty" json:"protoset,omitempty"`
	Raw              string   `yaml:"raw,omitempty" json:"raw,omitempty"`
}

type Secure struct {
	Insecure   bool   `yaml:"insecure,omitempty" json:"insecure,omitempty"`
	CertFile   string `yaml:"cert,omitempty" json:"cert,omitempty"`
	CAFile     string `yaml:"ca,omitempty" json:"ca,omitempty"`
	KeyFile    string `yaml:"key,omitempty" json:"key,omitempty"`
	ServerName string `yaml:"serverName,omitempty" json:"serverName,omitempty"`
}

type Metric struct {
	Type string `yaml:"type,omitempty" json:"type,omitempty"`
	URL  string `yaml:"url,omitempty" json:"url,omitempty"`
}

// Proxy configuration for the test suite
type Proxy struct {
	HTTP  string `yaml:"http,omitempty" json:"http,omitempty"`   // HTTP proxy URL
	HTTPS string `yaml:"https,omitempty" json:"https,omitempty"` // HTTPS proxy URL
	No    string `yaml:"no,omitempty" json:"no,omitempty"`       // Comma-separated list of hosts to exclude from proxying
}

// TestCase represents a test case
type TestCase struct {
	ID      string   `yaml:"id,omitempty" json:"id,omitempty"`
	Name    string   `yaml:"name,omitempty" json:"name,omitempty"`
	Group   string   `yaml:"group,omitempty" json:"group,omitempty"`
	Before  *Job     `yaml:"before,omitempty" json:"before,omitempty"`
	After   *Job     `yaml:"after,omitempty" json:"after,omitempty"`
	Request Request  `yaml:"request" json:"request"`
	Expect  Response `yaml:"expect,omitempty" json:"expect,omitempty"`
}

// InScope returns true if the test case is in scope with the given items.
// Returns true if the items is empty.
func (c *TestCase) InScope(items []string) bool {
	if len(items) == 0 {
		return true
	}
	for _, item := range items {
		if item == c.Name {
			return true
		}
	}
	return false
}

// Job contains a list of jobs
type Job struct {
	Items []string `yaml:"items,omitempty" json:"items,omitempty"`
}

// Request represents a HTTP request
type Request struct {
	API          string              `yaml:"api" json:"api"`
	Method       string              `yaml:"method,omitempty" json:"method,omitempty" jsonschema:"enum=GET,enum=POST,enum=PUT,enum=DELETE"`
	Query        SortedKeysStringMap `yaml:"query,omitempty" json:"query,omitempty"`
	Header       map[string]string   `yaml:"header,omitempty" json:"header,omitempty"`
	Cookie       map[string]string   `yaml:"cookie,omitempty" json:"cookie,omitempty"`
	Form         map[string]string   `yaml:"form,omitempty" json:"form,omitempty"`
	Body         RequestBody         `yaml:"body,omitempty" json:"body,omitempty"`
	BodyFromFile string              `yaml:"bodyFromFile,omitempty" json:"bodyFromFile,omitempty"`
}

type RequestBody struct {
	Value  string `json:"value" yaml:"value"`
	isJson bool
}

func NewRequestBody(val string) RequestBody {
	return RequestBody{Value: val}
}

func (e *RequestBody) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	gql := &GraphQLRequestBody{}
	err = unmarshal(gql)
	if err != nil {
		val := ""
		if err = unmarshal(&val); err == nil {
			e.Value = val
		}
	} else {
		var data []byte
		if data, err = json.Marshal(gql); err == nil {
			e.Value = string(data)
			e.isJson = true
		}
	}
	return
}

func (e RequestBody) MarshalYAML() (val interface{}, err error) {
	val = e.Value
	if e.isJson {
		gql := &GraphQLRequestBody{}
		if err = json.Unmarshal([]byte(e.Value), gql); err == nil {
			val = gql
		}
	}
	return
}

var _ yaml.Marshaler = &RequestBody{}

// var _ yaml.Unmarshaler = &RequestBody{}

func (e RequestBody) String() string {
	return e.Value
}

func (e RequestBody) IsEmpty() bool {
	return e.Value == ""
}

func (e RequestBody) Bytes() (data []byte) {
	var err error
	if strings.HasPrefix(e.Value, util.ImageBase64Prefix) {
		data, err = decodeBase64Body(e.Value, util.ImageBase64Prefix)
	} else if strings.HasPrefix(e.Value, util.PDFBase64Prefix) {
		data, err = decodeBase64Body(e.Value, util.PDFBase64Prefix)
	} else if strings.HasPrefix(e.Value, util.ZIPBase64Prefix) {
		data, err = decodeBase64Body(e.Value, util.ZIPBase64Prefix)
	} else if strings.HasPrefix(e.Value, util.BinaryBase64Prefix) {
		data, err = decodeBase64Body(e.Value, util.BinaryBase64Prefix)
	} else {
		data = []byte(e.Value)
	}

	if err != nil {
		log.Printf("Error decoding: %v", err)
	}
	return
}

func decodeBase64Body(raw, prefix string) ([]byte, error) {
	rawStr := strings.TrimPrefix(raw, prefix)
	return base64.StdEncoding.DecodeString(rawStr)
}

type GraphQLRequestBody struct {
	Query         string            `yaml:"query" json:"query"`
	OperationName string            `yaml:"operationName" json:"operationName"`
	Variables     map[string]string `yaml:"variables" json:"variables"`
}

// Response is the expected response
type Response struct {
	StatusCode        int                    `yaml:"statusCode,omitempty" json:"statusCode,omitempty"`
	Body              string                 `yaml:"body,omitempty" json:"body,omitempty"`
	Header            map[string]string      `yaml:"header,omitempty" json:"header,omitempty"`
	BodyFieldsExpect  map[string]interface{} `yaml:"bodyFieldsExpect,omitempty" json:"bodyFieldsExpect,omitempty"`
	Verify            []string               `yaml:"verify,omitempty" json:"verify,omitempty"`
	ConditionalVerify []ConditionalVerify    `yaml:"conditionalVerify,omitempty" json:"conditionalVerify,omitempty"`
	Schema            string                 `yaml:"schema,omitempty" json:"schema,omitempty"`
}

func (r Response) GetBody() string {
	return r.Body
}

func (r Response) GetBodyFieldsExpect() map[string]interface{} {
	return r.BodyFieldsExpect
}

type ConditionalVerify struct {
	Condition []string `yaml:"condition,omitempty" json:"condition,omitempty"`
	Verify    []string `yaml:"verify,omitempty" json:"verify,omitempty"`
}

type SortedKeysStringMap map[string]interface{}

func (m SortedKeysStringMap) Keys() (keys []string) {
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return
}

func (m SortedKeysStringMap) GetValue(key string) string {
	val := m[key]

	switch o := any(val).(type) {
	case string:
		return val.(string)
	case map[string]interface{}:
		verifier := convertToVerifier(o)
		return verifier.Value
	case *Verifier:
		return o.Value
	}

	return ""
}

func (m SortedKeysStringMap) GetVerifier(key string) (verifier *Verifier) {
	val := m[key]

	switch o := any(val).(type) {
	case map[string]interface{}:
		verifier = convertToVerifier(o)
	}

	return
}

func convertToVerifier(data map[string]interface{}) (verifier *Verifier) {
	verifier = &Verifier{}

	if data, err := yaml.Marshal(data); err == nil {
		if err = yaml.Unmarshal(data, verifier); err != nil {
			verifier = nil
		}
	}
	return
}

type Verifier struct {
	Value     string `yaml:"value,omitempty" json:"value,omitempty"`
	Required  bool   `yaml:"required,omitempty" json:"required,omitempty"`
	Max       int    `yaml:"max"`
	Min       int    `yaml:"min"`
	MaxLength int    `yaml:"maxLength"`
	MinLength int    `yaml:"minLength"`
}

type TestResult struct {
	Message        string            `yaml:"message,omitempty" json:"message,omitempty"`
	Error          string            `yaml:"error,omitempty" json:"error,omitempty"`
	TestCaseResult []*TestCaseResult `yaml:"testCaseResult,omitempty" json:"testCaseResult,omitempty"`
}

type TestCaseResult struct {
	StatusCode int               `yaml:"statusCode,omitempty" json:"statusCode,omitempty"`
	Body       string            `yaml:"body,omitempty" json:"body,omitempty"`
	Header     map[string]string `yaml:"header,omitempty" json:"header,omitempty"`
	Error      string            `yaml:"error,omitempty" json:"error,omitempty"`
	Id         string            `yaml:"id,omitempty" json:"id,omitempty"`
	Output     string            `yaml:"output,omitempty" json:"output,omitempty"`
}
