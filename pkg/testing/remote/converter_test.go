/*
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
