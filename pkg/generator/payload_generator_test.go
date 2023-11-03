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

package generator

import (
	"encoding/json"
	"testing"

	_ "embed"

	atest "github.com/linuxsuren/api-testing/pkg/testing"
	"github.com/linuxsuren/api-testing/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestGenerateJSON(t *testing.T) {
	t.Run("TestAdvancedType", func(t *testing.T) {
		result, err := generateGRPCPayloadAsJSON(&atest.RPCDesc{
			ProtoFile: "testdata/test.proto",
		}, "grpctest.Main.TestAdvancedType")
		if err != nil {
			t.Fatal(err)
		}

		obj := map[string]interface{}{}
		err = json.Unmarshal([]byte(result), &obj)
		if assert.NoError(t, err) {
			keys := util.Keys(obj)
			assert.ElementsMatch(t, keys,
				[]string{"Int32Array", "Int64Array", "Uint32Array", "Uint64Array",
					"Float32Array", "Float64Array", "StringArray", "BoolArray", "HelloReplyMap", "Protocol"})
		}
	})

	t.Run("TestBasicType", func(t *testing.T) {
		result, err := generateGRPCPayloadAsJSON(&atest.RPCDesc{
			ProtoFile: "testdata/test.proto",
		}, "grpctest.Main.TestBasicType")
		if err != nil {
			t.Fatal(err)
		}

		obj := map[string]interface{}{}
		err = json.Unmarshal([]byte(result), &obj)
		if assert.NoError(t, err) {
			keys := util.Keys(obj)
			assert.ElementsMatch(t, keys,
				[]string{"Int32", "Int64", "Uint32", "Uint64",
					"Float32", "Float64", "String", "Bool"})
		}
	})

	t.Run("ClientStream", func(t *testing.T) {
		result, err := generateGRPCPayloadAsJSON(&atest.RPCDesc{
			ProtoFile: "testdata/test.proto",
			Raw:       protoContent,
		}, "grpctest.Main.ClientStream")
		if err != nil {
			t.Fatal(err)
		}

		obj := map[string]interface{}{}
		err = json.Unmarshal([]byte(result), &obj)
		if assert.NoError(t, err) {
			keys := util.Keys(obj)
			assert.ElementsMatch(t, keys,
				[]string{"MsgID", "ExpectLen"})
		}
	})

	t.Run("BidStream, no proto file give, only file content", func(t *testing.T) {
		result, err := generateGRPCPayloadAsJSON(&atest.RPCDesc{
			Raw: protoContent,
		}, "grpctest.Main.BidStream")
		if err != nil {
			t.Fatal(err)
		}

		obj := map[string]interface{}{}
		err = json.Unmarshal([]byte(result), &obj)
		if assert.NoError(t, err) {
			keys := util.Keys(obj)
			assert.ElementsMatch(t, keys,
				[]string{"MsgID", "ExpectLen"})
		}
	})

	t.Run("call the generate method", func(t *testing.T) {
		generator := NewGrpcPayloadGenerator()
		result, err := generator.Generate(&atest.TestSuite{
			API: "localhost:7070",
			Spec: atest.APISpec{
				RPC: &atest.RPCDesc{
					ProtoFile: "testdata/test.proto",
				},
			},
		}, &atest.TestCase{
			Request: atest.Request{
				API: "/grpctest.Main/ServerStream",
			},
		})
		if err != nil {
			t.Fatal(err)
		}

		obj := map[string]interface{}{}
		err = json.Unmarshal([]byte(result), &obj)
		if assert.NoError(t, err) {
			keys := util.Keys(obj)
			assert.ElementsMatch(t, keys,
				[]string{"data"})
		}
	})
}

//go:embed testdata/test.proto
var protoContent string
