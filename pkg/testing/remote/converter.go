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

package remote

import (
	"fmt"

	server "github.com/linuxsuren/api-testing/pkg/server"
	"github.com/linuxsuren/api-testing/pkg/testing"
)

func ConvertToNormalTestSuite(suite *TestSuite) (result *testing.TestSuite) {
	result = &testing.TestSuite{
		Name:  suite.Name,
		API:   suite.Api,
		Param: pairToMap(suite.Param),
		Spec:  ConvertToNormalTestSuiteSpec(suite.Spec),
	}

	for _, testcase := range suite.Items {
		result.Items = append(result.Items, ConvertToNormalTestCase(testcase))
	}
	return
}

func ConvertToNormalTestSuiteSpec(spec *server.APISpec) (result testing.APISpec) {
	if spec != nil {
		result = testing.APISpec{
			Kind: spec.Kind,
			URL:  spec.Url,
		}
		if spec.Rpc != nil {
			result.RPC = &testing.RPCDesc{
				Raw:              spec.Rpc.Raw,
				ProtoFile:        spec.Rpc.Protofile,
				ImportPath:       spec.Rpc.Import,
				ServerReflection: spec.Rpc.ServerReflection,
			}
		}
		if spec.Secure != nil {
			result.Secure = &testing.Secure{
				Insecure:   spec.Secure.Insecure,
				CertFile:   spec.Secure.Cert,
				CAFile:     spec.Secure.Ca,
				ServerName: spec.Secure.ServerName,
				KeyFile:    spec.Secure.Key,
			}
		}
	}
	return
}

func ConvertToGRPCTestSuite(suite *testing.TestSuite) (result *TestSuite) {
	result = &TestSuite{
		Name:  suite.Name,
		Api:   suite.API,
		Param: mapToPair(suite.Param),
		Spec: &server.APISpec{
			Kind: suite.Spec.Kind,
			Url:  suite.Spec.URL,
		},
	}
	if suite.Spec.Secure != nil {
		result.Spec.Secure = &server.Secure{
			Insecure:   suite.Spec.Secure.Insecure,
			Cert:       suite.Spec.Secure.CertFile,
			Ca:         suite.Spec.Secure.CAFile,
			ServerName: suite.Spec.Secure.ServerName,
			Key:        suite.Spec.Secure.KeyFile,
		}
	}
	if suite.Spec.RPC != nil {
		result.Spec.Rpc = &server.RPC{
			Import:           suite.Spec.RPC.ImportPath,
			ServerReflection: suite.Spec.RPC.ServerReflection,
			Protofile:        suite.Spec.RPC.ProtoFile,
			Protoset:         suite.Spec.RPC.ProtoSet,
			Raw:              suite.Spec.RPC.Raw,
		}
	}

	for _, testcase := range suite.Items {
		result.Items = append(result.Items, ConvertToGRPCTestCase(testcase))
	}
	return
}

func ConvertToNormalTestCase(testcase *server.TestCase) (result testing.TestCase) {
	result = testing.TestCase{
		Name: testcase.Name,
	}
	if testcase.Request != nil {
		result.Request = testing.Request{
			API:    testcase.Request.Api,
			Method: testcase.Request.Method,
			Body:   testing.NewRequestBody(testcase.Request.Body),
			Header: pairToMap(testcase.Request.Header),
			Query:  pairToMapInter(testcase.Request.Query),
			Form:   pairToMap(testcase.Request.Form),
		}
	}
	if testcase.Response != nil {
		result.Expect = testing.Response{
			Body:             testcase.Response.Body,
			StatusCode:       int(testcase.Response.StatusCode),
			Schema:           testcase.Response.Schema,
			Verify:           testcase.Response.Verify,
			Header:           pairToMap(testcase.Response.Header),
			BodyFieldsExpect: pairToInterMap(testcase.Response.BodyFieldsExpect),
		}
	}
	return
}

func ConvertToGRPCTestCase(testcase testing.TestCase) (result *server.TestCase) {
	result = &server.TestCase{
		Name: testcase.Name,
		Request: &server.Request{
			Api:    testcase.Request.API,
			Method: testcase.Request.Method,
			Body:   testcase.Request.Body.String(),
			Header: mapToPair(testcase.Request.Header),
			Query:  mapInterToPair(testcase.Request.Query),
			Form:   mapToPair(testcase.Request.Form),
		},
		Response: &server.Response{
			Body:             testcase.Expect.Body,
			StatusCode:       int32(testcase.Expect.StatusCode),
			Schema:           testcase.Expect.Schema,
			Verify:           testcase.Expect.Verify,
			Header:           mapToPair(testcase.Expect.Header),
			BodyFieldsExpect: mapInterToPair(testcase.Expect.BodyFieldsExpect),
		},
	}
	return
}

func mapToPair(data map[string]string) (pairs []*server.Pair) {
	pairs = make([]*server.Pair, 0)
	for k, v := range data {
		pairs = append(pairs, &server.Pair{
			Key:   k,
			Value: v,
		})
	}
	return
}

func mapInterToPair(data map[string]interface{}) (pairs []*server.Pair) {
	pairs = make([]*server.Pair, 0)
	for k, v := range data {
		pairs = append(pairs, &server.Pair{
			Key:   k,
			Value: fmt.Sprintf("%v", v),
		})
	}
	return
}

func pairToMap(pairs []*server.Pair) (data map[string]string) {
	data = make(map[string]string)
	for _, pair := range pairs {
		data[pair.Key] = pair.Value
	}
	return
}

func pairToMapInter(pairs []*server.Pair) (data map[string]interface{}) {
	data = make(map[string]interface{})
	for _, pair := range pairs {
		data[pair.Key] = pair.Value
	}
	return
}

func pairToInterMap(pairs []*server.Pair) (data map[string]interface{}) {
	data = make(map[string]interface{})
	for _, pair := range pairs {
		data[pair.Key] = pair.Value
	}
	return
}
