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
