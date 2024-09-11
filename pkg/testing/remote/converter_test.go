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
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"

	server "github.com/linuxsuren/api-testing/pkg/server"
	atest "github.com/linuxsuren/api-testing/pkg/testing"
	"github.com/stretchr/testify/assert"
)

func TestConvert(t *testing.T) {
	now := time.Now().UTC()
	t.Run("convertToNormalTestSuite, empty object", func(t *testing.T) {
		assert.Equal(t, &atest.TestSuite{
			Param: map[string]string{},
		}, ConvertToNormalTestSuite(&TestSuite{}))
	})

	t.Run("convertToNormalTestSuite, normal object", func(t *testing.T) {
		assert.Equal(t, &atest.TestSuite{
			Param: defaultMap,
			Spec: atest.APISpec{
				Kind: "http",
				URL:  "/v1",
				RPC: &atest.RPCDesc{
					Raw: "fake",
				},
				Secure: &atest.Secure{
					KeyFile: "fake",
				},
			},
			Items: []atest.TestCase{{
				Name: "fake",
			}},
		}, ConvertToNormalTestSuite(&TestSuite{
			Param: defaultPairs,
			Spec: &server.APISpec{
				Url:  "/v1",
				Kind: "http",
				Rpc: &server.RPC{
					Raw: "fake",
				},
				Secure: &server.Secure{
					Key: "fake",
				},
			},
			Items: []*server.TestCase{{
				Name: "fake",
			}},
		}))
	})

	t.Run("convertToGRPCTestSuite, normal object", func(t *testing.T) {
		result := ConvertToGRPCTestSuite(&atest.TestSuite{
			API:   "v1",
			Param: defaultMap,
			Spec: atest.APISpec{
				RPC: &atest.RPCDesc{
					Raw: "fake",
				},
				Secure: &atest.Secure{
					KeyFile: "fake",
				},
			},
			Items: []atest.TestCase{{
				Name: "fake",
			}},
		})
		assert.Equal(t, "v1", result.Api)
		assert.Equal(t, defaultPairs, result.Param)
		assert.Equal(t, "fake", result.Spec.Secure.Key)
	})

	t.Run("convertToNormalTestCase", func(t *testing.T) {
		assert.Equal(t, atest.TestCase{
			Request: atest.Request{
				API:    "/v1",
				Header: defaultMap,
				Query:  map[string]interface{}{},
				Form:   map[string]string{},
			},
			Expect: atest.Response{
				BodyFieldsExpect: defaultInterMap,
				Header:           map[string]string{},
			},
		}, ConvertToNormalTestCase(&server.TestCase{
			Request: &server.Request{
				Api:    "/v1",
				Query:  nil,
				Header: defaultPairs,
			},
			Response: &server.Response{
				BodyFieldsExpect: defaultPairs,
			},
		}))
	})

	t.Run("convertToGRPCTestCase", func(t *testing.T) {
		result := ConvertToGRPCTestCase(atest.TestCase{
			Expect: atest.Response{
				BodyFieldsExpect: defaultInterMap,
				Header:           defaultMap,
			},
		})
		if !assert.NotNil(t, result) {
			return
		}
		assert.Equal(t, defaultPairs, result.Response.BodyFieldsExpect)
		assert.Equal(t, defaultPairs, result.Response.Header)
	})

	t.Run("convertHistoryToGRPCTestCase", func(t *testing.T) {
		result := ConvertHistoryToGRPCTestCase(&server.HistoryTestCase{
			CaseName: "fake",
			Request: &server.Request{
				Header: defaultPairs,
			},
			Response: &server.Response{
				BodyFieldsExpect: defaultPairs,
			},
		})
		if !assert.NotNil(t, result) {
			return
		}
		assert.Equal(t, defaultMap, result.Request.Header)
		assert.Equal(t, defaultInterMap, result.Expect.BodyFieldsExpect)
		assert.Equal(t, "fake", result.Name)
	})

	t.Run("convertToNormalHistoryTestCase", func(t *testing.T) {
		assert.Equal(t, atest.HistoryTestCase{
			CreateTime: now,
			SuiteParam: defaultMap,
			SuiteSpec: atest.APISpec{
				Kind: "http",
				URL:  "/v1",
				RPC: &atest.RPCDesc{
					Raw: "fake",
				},
				Secure: &atest.Secure{
					KeyFile: "fake",
				},
			},
			Data: atest.TestCase{
				Request: atest.Request{
					API:    "/v1",
					Header: defaultMap,
					Query:  map[string]interface{}{},
					Form:   map[string]string{},
				},
				Expect: atest.Response{
					BodyFieldsExpect: defaultInterMap,
					Header:           map[string]string{},
				},
			},
		}, ConvertToNormalHistoryTestCase(&server.HistoryTestCase{
			CreateTime: timestamppb.New(now),
			SuiteParam: defaultPairs,
			SuiteSpec: &server.APISpec{
				Url:  "/v1",
				Kind: "http",
				Rpc: &server.RPC{
					Raw: "fake",
				},
				Secure: &server.Secure{
					Key: "fake",
				},
			},
			Request: &server.Request{
				Header: defaultPairs,
				Query:  nil,
				Api:    "/v1",
			},
			Response: &server.Response{
				BodyFieldsExpect: defaultPairs,
			},
		}))
	})

	t.Run("convertToNormalHistoryTestSuite, empty object", func(t *testing.T) {
		assert.Equal(t, &atest.HistoryTestSuite{}, ConvertToNormalHistoryTestSuite(&HistoryTestSuite{}))
	})

	t.Run("convertToNormalHistoryTestSuite, normal object", func(t *testing.T) {
		assert.Equal(t, &atest.HistoryTestSuite{
			HistorySuiteName: "fake",
			Items: []atest.HistoryTestCase{
				{
					CreateTime: now,
					SuiteParam: defaultMap,
					SuiteSpec: atest.APISpec{
						Kind: "http",
						URL:  "/v1",
						RPC: &atest.RPCDesc{
							Raw: "fake",
						},
						Secure: &atest.Secure{
							KeyFile: "fake",
						},
					},
				},
			},
		}, ConvertToNormalHistoryTestSuite(&HistoryTestSuite{
			HistorySuiteName: "fake",
			Items: []*server.HistoryTestCase{
				{
					CreateTime: timestamppb.New(now),
					SuiteParam: defaultPairs,
					SuiteSpec: &server.APISpec{
						Url:  "/v1",
						Kind: "http",
						Rpc: &server.RPC{
							Raw: "fake",
						},
						Secure: &server.Secure{
							Key: "fake",
						},
					},
				},
			},
		}))
	})

	t.Run("convertToGRPCHistoryTestCase", func(t *testing.T) {
		result := ConvertToGRPCHistoryTestCase(atest.HistoryTestCase{
			SuiteParam: defaultMap,
			SuiteSpec: atest.APISpec{
				Secure: &atest.Secure{
					KeyFile: "fake",
				},
			},
			Data: atest.TestCase{
				Request: atest.Request{
					Header: defaultMap,
				},
				Expect: atest.Response{
					BodyFieldsExpect: defaultInterMap,
				},
			},
		})
		assert.Equal(t, defaultPairs, result.SuiteParam)
		assert.Equal(t, defaultPairs, result.Request.Header)
		assert.Equal(t, defaultPairs, result.Response.BodyFieldsExpect)
		assert.Equal(t, "fake", result.SuiteSpec.Secure.Key)
	})

	t.Run("convertToGRPCHistoryTestCaseResult", func(t *testing.T) {
		result := ConvertToGRPCHistoryTestCaseResult(atest.TestCaseResult{
			Body:   "fake body",
			Output: "fake output",
		}, &atest.TestSuite{
			Param: defaultMap,
			Spec: atest.APISpec{
				Secure: &atest.Secure{
					KeyFile: "fake",
				},
			},
			Items: []atest.TestCase{
				{
					Request: atest.Request{
						Header: defaultMap,
					},
					Expect: atest.Response{
						BodyFieldsExpect: defaultInterMap,
					},
				},
			},
		})
		assert.Equal(t, defaultPairs, result.Data.SuiteParam)
		assert.Equal(t, defaultPairs, result.Data.Request.Header)
		assert.Equal(t, defaultPairs, result.Data.Response.BodyFieldsExpect)
		assert.Equal(t, "fake", result.Data.SuiteSpec.Secure.Key)
		assert.Equal(t, "fake output", result.TestCaseResult[0].Output)
		assert.Equal(t, "fake body", result.TestCaseResult[0].Body)
	})

	t.Run("convertToNormalTestCaseResult", func(t *testing.T) {
		assert.Equal(t, atest.HistoryTestResult{
			CreateTime: now,
			Data: atest.HistoryTestCase{
				SuiteParam: defaultMap,
				CreateTime: now,
			},
			TestCaseResult: []atest.TestCaseResult{
				{
					Body:   "fake body",
					Output: "fake output",
					Header: defaultMap,
				},
				{
					Body:   "fake body 2",
					Output: "fake output 2",
					Header: defaultMap,
				},
			},
		}, ConvertToNormalTestCaseResult(&server.HistoryTestResult{
			CreateTime: timestamppb.New(now),
			Data: &server.HistoryTestCase{
				SuiteParam: defaultPairs,
				CreateTime: timestamppb.New(now),
			},
			TestCaseResult: []*server.TestCaseResult{
				{
					Body:   "fake body",
					Output: "fake output",
					Header: defaultPairs,
				},
				{
					Body:   "fake body 2",
					Output: "fake output 2",
					Header: defaultPairs,
				},
			},
		}))
	})
}

var defaultInterMap = map[string]interface{}{"foo": "bar"}
var defaultMap map[string]string = map[string]string{"foo": "bar"}
var defaultPairs []*server.Pair = []*server.Pair{{Key: "foo", Value: "bar"}}
