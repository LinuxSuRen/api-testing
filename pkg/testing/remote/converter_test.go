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
	"testing"

	server "github.com/linuxsuren/api-testing/pkg/server"
	atest "github.com/linuxsuren/api-testing/pkg/testing"
	"github.com/stretchr/testify/assert"
)

func TestConvert(t *testing.T) {
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
}

var defaultInterMap = map[string]interface{}{"foo": "bar"}
var defaultMap map[string]string = map[string]string{"foo": "bar"}
var defaultPairs []*server.Pair = []*server.Pair{{Key: "foo", Value: "bar"}}
