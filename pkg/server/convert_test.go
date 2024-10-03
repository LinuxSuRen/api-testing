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
package server

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"

	atest "github.com/linuxsuren/api-testing/pkg/testing"
	"github.com/stretchr/testify/assert"
)

func TestToGRPCStore(t *testing.T) {
	assert.Equal(t, &Store{
		Name:        "test",
		Owner:       "rick",
		Description: "desc",
		Kind: &StoreKind{
			Name: "test",
			Url:  urlFoo,
		},
		Url:      urlFoo,
		Username: "user",
		Password: "pass",
		Disabled: true,
		Properties: []*Pair{{
			Key: "foo", Value: "bar",
		}},
	}, ToGRPCStore(atest.Store{
		Name:        "test",
		Owner:       "rick",
		Description: "desc",
		Kind: atest.StoreKind{
			Name: "test",
			URL:  urlFoo,
		},
		URL:      urlFoo,
		Username: "user",
		Password: "pass",
		Disabled: true,
		Properties: map[string]string{
			"foo": "bar",
		},
	}))
}

func TestToNormalStore(t *testing.T) {
	assert.Equal(t, atest.Store{
		Name:        "test",
		Description: "desc",
		Kind: atest.StoreKind{
			Name: "test",
			URL:  urlFoo,
		},
		URL:      urlFoo,
		Username: "user",
		Password: "pass",
		Properties: map[string]string{
			"foo": "bar",
		},
	}, ToNormalStore(&Store{
		Name:        "test",
		Description: "desc",
		Kind: &StoreKind{
			Name: "test",
			Url:  urlFoo,
		},
		Url:      urlFoo,
		Username: "user",
		Password: "pass",
		Properties: []*Pair{{
			Key: "foo", Value: "bar",
		}},
	}))
}

func TestToGRPCSuite(t *testing.T) {
	assert.Equal(t, &TestSuite{
		Name: "test",
		Api:  "api",
		Param: []*Pair{{
			Key: "foo", Value: "bar",
		}},
		Spec: &APISpec{
			Secure: &Secure{
				Insecure: true,
			},
			Rpc: &RPC{
				Raw: "raw",
			},
		},
	}, ToGRPCSuite(&atest.TestSuite{
		Name: "test",
		API:  "api",
		Param: map[string]string{
			"foo": "bar",
		},
		Spec: atest.APISpec{
			Secure: &atest.Secure{
				Insecure: true,
			},
			RPC: &atest.RPCDesc{
				Raw: "raw",
			},
		},
	}))
}

func TestToNormalSuite(t *testing.T) {
	assert.Equal(t, &atest.TestSuite{
		Name: "test",
		API:  "api",
		Param: map[string]string{
			"foo": "bar",
		},
		Spec: atest.APISpec{
			Secure: &atest.Secure{
				Insecure: true,
			},
			RPC: &atest.RPCDesc{
				Raw: "raw",
			},
		},
	}, ToNormalSuite(&TestSuite{
		Name: "test",
		Api:  "api",
		Param: []*Pair{{
			Key: "foo", Value: "bar",
		}},
		Spec: &APISpec{
			Secure: &Secure{
				Insecure: true,
			},
			Rpc: &RPC{
				Raw: "raw",
			},
		},
	}))
}

func TestConvertToGRPCHistoryTestCase(t *testing.T) {
	now := time.Now().UTC()
	result := ConvertToGRPCHistoryTestCase(atest.HistoryTestCase{
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
				Header: defaultMap,
			},
			Expect: atest.Response{
				BodyFieldsExpect: defaultInterMap,
			},
		},
	})
	assert.Equal(t, result.Request.Header, defaultPairs)
	assert.Equal(t, result.SuiteParam, defaultPairs)
	assert.Equal(t, result.Response.BodyFieldsExpect, defaultPairs)
	assert.Equal(t, "fake", result.SuiteSpec.Secure.Key)
	assert.Equal(t, timestamppb.New(now), result.CreateTime)
}

func TestToNormalTestCaseResult(t *testing.T) {
	assert.Equal(t, atest.TestCaseResult{
		Body:       "body",
		Error:      "error",
		Header:     defaultMap,
		Id:         "id",
		Output:     "output",
		StatusCode: 200,
	}, ToNormalTestCaseResult(&TestCaseResult{
		Body:       "body",
		Error:      "error",
		Header:     defaultPairs,
		Id:         "id",
		Output:     "output",
		StatusCode: 200,
	}))
}

func TestToGRPCHistoryTestCaseResult(t *testing.T) {
	t.Run("TestCaseResult is empty", func(t *testing.T) {
		historyTestResult := atest.HistoryTestResult{
			Message:    "test message",
			Error:      "test error",
			CreateTime: time.Now(),
			Data: atest.HistoryTestCase{
				ID: "test-id",
			},
			TestCaseResult: nil,
		}

		result := ToGRPCHistoryTestCaseResult(historyTestResult)

		assert.Equal(t, 0, len(result.TestCaseResult))
		assert.Equal(t, historyTestResult.Message, result.Message)
		assert.Equal(t, historyTestResult.Error, result.Error)
	})

	t.Run("TestCaseResult is not empty", func(t *testing.T) {
		now := time.Now().UTC()

		result := ToGRPCHistoryTestCaseResult(atest.HistoryTestResult{
			Message:    "fake message",
			CreateTime: now,
			Data: atest.HistoryTestCase{
				ID: "fake id",
			},
			TestCaseResult: []atest.TestCaseResult{
				{
					StatusCode: 200,
					Output:     "fake output",
				},
				{
					Output: "fake output 2",
				},
			},
		})

		assert.Equal(t, 2, len(result.TestCaseResult))
		assert.Equal(t, "fake message", result.Message)
		assert.Equal(t, now, result.CreateTime.AsTime())
		assert.Equal(t, "fake output", result.TestCaseResult[0].Output)
		assert.Equal(t, "fake output 2", result.TestCaseResult[1].Output)
	})
}

func TestToGRPCTestSuiteSpec(t *testing.T) {

	t.Run("empty", func(t *testing.T) {
		assert.Equal(t, &APISpec{}, ToGRPCTestSuiteSpec(atest.APISpec{}))
	})

	t.Run("fields", func(t *testing.T) {
		assert.Equal(t, &APISpec{
			Url:  "/v1",
			Kind: "http",
			Rpc: &RPC{
				Raw: "fake",
			},
			Secure: &Secure{
				Key: "fake",
			},
		}, ToGRPCTestSuiteSpec(atest.APISpec{
			Kind: "http",
			URL:  "/v1",
			RPC: &atest.RPCDesc{
				Raw: "fake",
			},
			Secure: &atest.Secure{
				KeyFile: "fake",
			},
		}))
	})
}

var defaultInterMap = map[string]interface{}{"foo": "bar"}
var defaultMap map[string]string = map[string]string{"foo": "bar"}
var defaultPairs []*Pair = []*Pair{{Key: "foo", Value: "bar"}}
