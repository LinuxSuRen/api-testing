/**
MIT License

Copyright (c) 2023 API Testing Authors.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package server

import (
	"strings"

	"github.com/linuxsuren/api-testing/pkg/testing"
	"github.com/linuxsuren/api-testing/pkg/util"
)

// ToGRPCStore convert the normal store to GRPC store
func ToGRPCStore(store testing.Store) (result *Store) {
	result = &Store{
		Name: store.Name,
		Kind: &StoreKind{
			Name:    store.Kind.Name,
			Url:     store.Kind.URL,
			Enabled: store.Kind.Enabled,
		},
		Description: store.Description,
		Url:         store.URL,
		Username:    store.Username,
		Password:    store.Password,
		Properties:  mapToPair(store.Properties),
	}
	return
}

// ToNormalStore convert the GRPC store to normal store
func ToNormalStore(store *Store) (result testing.Store) {
	result = testing.Store{
		Name:        store.Name,
		Description: store.Description,
		URL:         store.Url,
		Username:    store.Username,
		Password:    store.Password,
		Properties:  pairToMap(store.Properties),
	}
	if store.Kind != nil {
		result.Kind = testing.StoreKind{
			Name: store.Kind.Name,
			URL:  store.Kind.Url,
		}
	}
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
				Raw:        suite.Spec.Rpc.Raw,
				ProtoFile:  suite.Spec.Rpc.Protofile,
				ImportPath: suite.Spec.Rpc.Import,
			}
		}
	}
	return
}

func ToGRPCTestCase(testCase testing.TestCase) (result *TestCase) {
	req := &Request{
		Api:    testCase.Request.API,
		Method: testCase.Request.Method,
		Query:  mapToPair(testCase.Request.Query),
		Header: mapToPair(testCase.Request.Header),
		Form:   mapToPair(testCase.Request.Form),
		Body:   testCase.Request.Body,
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
		result.Request.Body = req.Body
		result.Request.Header = pairToMap(req.Header)
		result.Request.Form = pairToMap(req.Form)
		result.Request.Query = pairToMap(req.Query)
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
