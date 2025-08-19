/*
Copyright 2023-2025 API Testing Authors.

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
package server

import (
	"strings"

	"github.com/linuxsuren/api-testing/pkg/testing"
	"github.com/linuxsuren/api-testing/pkg/util"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ToGRPCStore convert the normal store to GRPC store
func ToGRPCStore(store testing.Store) (result *Store) {
	result = &Store{
		Name:  store.Name,
		Owner: store.Owner,
		Kind: &StoreKind{
			Name:       store.Kind.Name,
			Url:        store.Kind.URL,
			Enabled:    store.Kind.Enabled,
			Categories: store.Kind.Categories,
			Link:       store.Kind.Link,
		},
		Description: store.Description,
		Url:         store.URL,
		Username:    store.Username,
		Password:    store.Password,
		Disabled:    store.Disabled,
		Properties:  mapToPair(store.Properties),
	}
	return
}

// ToNormalStore convert the GRPC store to normal store
func ToNormalStore(store *Store) (result testing.Store) {
	result = testing.Store{
		Name:        store.Name,
		Owner:       store.Owner,
		Description: store.Description,
		URL:         store.Url,
		Username:    store.Username,
		Password:    store.Password,
		Disabled:    store.Disabled,
		Properties:  pairToMap(store.Properties),
	}
	if store.Kind != nil {
		result.Kind = testing.StoreKind{
			Name:       store.Kind.Name,
			URL:        store.Kind.Url,
			Link:       store.Kind.Link,
			Categories: store.Kind.Categories,
		}
	}
	return
}

func ConvertToGRPCHistoryTestCase(historyTestcase testing.HistoryTestCase) (result *HistoryTestCase) {
	req := historyTestcase.Data.Request
	res := historyTestcase.Data.Expect
	result = &HistoryTestCase{
		ID:               historyTestcase.ID,
		CreateTime:       timestamppb.New(historyTestcase.CreateTime),
		CaseName:         historyTestcase.CaseName,
		SuiteName:        historyTestcase.SuiteName,
		HistorySuiteName: historyTestcase.HistorySuiteName,
		SuiteSpec:        ToGRPCTestSuiteSpec(historyTestcase.SuiteSpec),
		SuiteApi:         historyTestcase.SuiteAPI,
		SuiteParam:       mapToPair(historyTestcase.SuiteParam),
		HistoryHeader:    mapToPair(historyTestcase.HistoryHeader),

		Request: &Request{
			Api:    req.API,
			Method: req.Method,
			Body:   req.Body.String(),
			Header: mapToPair(req.Header),
			Query:  mapInterToPair(req.Query),
			Form:   mapToPair(req.Form),
		},

		Response: &Response{
			Body:             res.Body,
			StatusCode:       int32(res.StatusCode),
			Schema:           res.Schema,
			Verify:           res.Verify,
			Header:           mapToPair(res.Header),
			BodyFieldsExpect: mapInterToPair(res.BodyFieldsExpect),
		},
	}
	result.SuiteSpec = ToGRPCTestSuiteSpec(historyTestcase.SuiteSpec)
	return
}

func ToGRPCSuite(suite *testing.TestSuite) (result *TestSuite) {
	result = &TestSuite{
		Name:  suite.Name,
		Api:   suite.API,
		Param: mapToPair(suite.Param),
		Spec: &APISpec{
			Kind: suite.Spec.Kind,
			Url:  suite.Spec.URL,
		},
	}
	if suite.Proxy != nil {
		result.Proxy = &ProxyConfig{
			Http:  suite.Proxy.HTTP,
			Https: suite.Proxy.HTTPS,
			No:    suite.Proxy.No,
		}
	}
	if suite.Spec.Secure != nil {
		result.Spec.Secure = &Secure{
			Insecure:   suite.Spec.Secure.Insecure,
			Cert:       suite.Spec.Secure.CertFile,
			Ca:         suite.Spec.Secure.CAFile,
			ServerName: suite.Spec.Secure.ServerName,
			Key:        suite.Spec.Secure.KeyFile,
		}
	}
	if suite.Spec.RPC != nil {
		result.Spec.Rpc = &RPC{
			Import:           suite.Spec.RPC.ImportPath,
			ServerReflection: suite.Spec.RPC.ServerReflection,
			Protofile:        suite.Spec.RPC.ProtoFile,
			Protoset:         suite.Spec.RPC.ProtoSet,
			Raw:              suite.Spec.RPC.Raw,
		}
	}
	return
}

func ToNormalSuite(suite *TestSuite) (result *testing.TestSuite) {
	result = &testing.TestSuite{
		Name:  suite.Name,
		API:   suite.Api,
		Param: pairToMap(suite.Param),
	}
	if suite.Proxy != nil {
		result.Proxy = &testing.Proxy{
			HTTP:  suite.Proxy.Http,
			HTTPS: suite.Proxy.Https,
			No:    suite.Proxy.No,
		}
	}
	if suite.Spec != nil {
		result.Spec = testing.APISpec{
			Kind: suite.Spec.Kind,
			URL:  suite.Spec.Url,
		}
		if suite.Spec.Secure != nil {
			result.Spec.Secure = &testing.Secure{
				Insecure:   suite.Spec.Secure.Insecure,
				CertFile:   suite.Spec.Secure.Cert,
				CAFile:     suite.Spec.Secure.Ca,
				ServerName: suite.Spec.Secure.ServerName,
				KeyFile:    suite.Spec.Secure.Key,
			}
		}
		if suite.Spec.Rpc != nil {
			result.Spec.RPC = &testing.RPCDesc{
				Raw:              suite.Spec.Rpc.Raw,
				ProtoFile:        suite.Spec.Rpc.Protofile,
				ProtoSet:         suite.Spec.Rpc.Protoset,
				ImportPath:       suite.Spec.Rpc.Import,
				ServerReflection: suite.Spec.Rpc.ServerReflection,
			}
		}
	}
	return
}

func ToNormalSuiteYAML(suite *TestSuite) ([]byte, error) {
	result := ToNormalSuite(suite)
	return testing.ToYAML(result)
}

func ToGRPCTestCase(testCase testing.TestCase) (result *TestCase) {
	req := &Request{
		Api:    testCase.Request.API,
		Method: testCase.Request.Method,
		Query:  mapInterToPair(testCase.Request.Query),
		Header: mapToPair(testCase.Request.Header),
		Cookie: mapToPair(testCase.Request.Cookie),
		Form:   mapToPair(testCase.Request.Form),
		Body:   testCase.Request.Body.String(),
	}

	resp := &Response{
		StatusCode:       int32(testCase.Expect.StatusCode),
		Body:             testCase.Expect.Body,
		Header:           mapToPair(testCase.Expect.Header),
		BodyFieldsExpect: mapInterToPair(testCase.Expect.BodyFieldsExpect),
		Verify:           testCase.Expect.Verify,
		Schema:           testCase.Expect.Schema,
	}

	result = &TestCase{
		Name:     testCase.Name,
		Request:  req,
		Response: resp,
	}
	return
}

func ToNormalTestCase(in *TestCase) (result testing.TestCase) {
	result = testing.TestCase{
		Name: in.Name,
	}
	req := in.Request
	resp := in.Response

	if req != nil {
		result.Request.API = req.Api
		result.Request.Method = req.Method
		result.Request.Body = testing.NewRequestBody(req.Body)
		result.Request.Header = pairToMap(req.Header)
		result.Request.Cookie = pairToMap(req.Cookie)
		result.Request.Form = pairToMap(req.Form)
		result.Request.Query = pairToInterMap(req.Query)
	}

	if resp != nil {
		result.Expect.Body = strings.TrimSpace(resp.Body)
		result.Expect.Schema = strings.TrimSpace(resp.Schema)
		result.Expect.StatusCode = int(resp.StatusCode)
		result.Expect.Verify = util.RemoeEmptyFromSlice(resp.Verify)
		result.Expect.ConditionalVerify = convertConditionalVerify(resp.ConditionalVerify)
		result.Expect.BodyFieldsExpect = pairToInterMap(resp.BodyFieldsExpect)
		result.Expect.Header = pairToMap(resp.Header)
	}
	return
}

func ToNormalTestCaseResult(testCaseResult *TestCaseResult) (result testing.TestCaseResult) {
	result = testing.TestCaseResult{
		StatusCode: int(testCaseResult.StatusCode),
		Error:      testCaseResult.Error,
		Body:       testCaseResult.Body,
		Header:     pairToMap(testCaseResult.Header),
		Id:         testCaseResult.Id,
		Output:     testCaseResult.Output,
	}
	return result
}

func ToGRPCHistoryTestCaseResult(historyTestResult testing.HistoryTestResult) (result *HistoryTestResult) {
	convertedHistoryTestCase := ConvertToGRPCHistoryTestCase(historyTestResult.Data)

	result = &HistoryTestResult{
		Message:    historyTestResult.Message,
		Error:      historyTestResult.Error,
		CreateTime: timestamppb.New(historyTestResult.CreateTime),
		Data:       convertedHistoryTestCase,
	}

	for _, testCaseResult := range historyTestResult.TestCaseResult {
		result.TestCaseResult = append(result.TestCaseResult, &TestCaseResult{
			StatusCode: int32(testCaseResult.StatusCode),
			Error:      testCaseResult.Error,
			Body:       testCaseResult.Body,
			Header:     mapToPair(testCaseResult.Header),
			Output:     testCaseResult.Output,
			Id:         testCaseResult.Id,
		})
	}

	return result
}

func ToGRPCTestSuiteSpec(spec testing.APISpec) (result *APISpec) {
	result = &APISpec{
		Kind: spec.Kind,
		Url:  spec.URL,
	}
	if spec.RPC != nil {
		result.Rpc = &RPC{
			Raw:              spec.RPC.Raw,
			Protofile:        spec.RPC.ProtoFile,
			Import:           spec.RPC.ImportPath,
			ServerReflection: spec.RPC.ServerReflection,
		}
	}
	if spec.Secure != nil {
		result.Secure = &Secure{
			Insecure:   spec.Secure.Insecure,
			Cert:       spec.Secure.CertFile,
			Ca:         spec.Secure.CAFile,
			ServerName: spec.Secure.ServerName,
			Key:        spec.Secure.KeyFile,
		}
	}
	return
}
