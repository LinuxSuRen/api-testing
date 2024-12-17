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
package testing_test

import (
    "github.com/linuxsuren/api-testing/pkg/util"
    "testing"

    atesting "github.com/linuxsuren/api-testing/pkg/testing"
    "github.com/stretchr/testify/assert"

    "gopkg.in/yaml.v3"
)

func TestInScope(t *testing.T) {
    testCase := &atesting.TestCase{Name: "foo"}
    assert.True(t, testCase.InScope(nil))
    assert.True(t, testCase.InScope([]string{"foo"}))
    assert.False(t, testCase.InScope([]string{"bar"}))
}

func TestRequestBody(t *testing.T) {
    req := &atesting.Request{}
    graphqlBody := `api: /api
body:
    query: query
    operationName: ""
    variables:
        name: rick
`

    err := yaml.Unmarshal([]byte(graphqlBody), req)
    assert.Nil(t, err)
    assert.Equal(t, `{"query":"query","operationName":"","variables":{"name":"rick"}}`, req.Body.String())

    var data []byte
    data, err = yaml.Marshal(req)
    assert.Nil(t, err)
    assert.Equal(t, graphqlBody, string(data))

    err = yaml.Unmarshal([]byte(`body: plain`), req)
    assert.Nil(t, err)
    assert.Equal(t, "plain", req.Body.String())
}

func TestResponse(t *testing.T) {
    resp := &atesting.Response{
        Body: "body",
        BodyFieldsExpect: map[string]interface{}{
            "name": "rick",
        },
    }
    assert.Equal(t, "body", resp.GetBody())
    assert.Equal(t, map[string]interface{}{"name": "rick"}, resp.GetBodyFieldsExpect())
}

func TestSortedKeysStringMap(t *testing.T) {
    obj := atesting.SortedKeysStringMap{
        "c": "d",
        "f": map[string]interface{}{
            "value": "f",
        },
        "e": &atesting.Verifier{
            Value: "e",
        },
        "a": "b",
    }
    assert.Equal(t, []string{"a", "c", "e", "f"}, obj.Keys())
    assert.Equal(t, "b", obj.GetValue("a"))
    assert.Nil(t, obj.GetVerifier("b"))
    assert.Equal(t, "e", obj.GetValue("e"))
    assert.Equal(t, "f", obj.GetValue("f"))
    assert.Equal(t, "f", obj.GetVerifier("f").Value)
    assert.Empty(t, obj.GetValue("not-found"))
}

func TestBodyBytes(t *testing.T) {
    const defaultPlainText = "hello"
    const defaultBase64Text = "aGVsbG8="

    tt := []struct {
        name    string
        rawBody string
        expect  []byte
    }{{
        name:    "image base64",
        rawBody: util.ImageBase64Prefix + defaultBase64Text,
        expect:  []byte(defaultPlainText),
    }, {
        name:    "pdf",
        rawBody: util.PDFBase64Prefix + defaultBase64Text,
        expect:  []byte(defaultPlainText),
    }, {
        name:    "zip",
        rawBody: util.ZIPBase64Prefix + defaultBase64Text,
        expect:  []byte(defaultPlainText),
    }, {
        name:    "binary",
        rawBody: util.BinaryBase64Prefix + defaultBase64Text,
        expect:  []byte(defaultPlainText),
    }, {
        name:    "raw",
        rawBody: defaultPlainText,
        expect:  []byte(defaultPlainText),
    }}
    for _, tc := range tt {
        t.Run(tc.name, func(t *testing.T) {
            body := atesting.RequestBody{
                Value: tc.rawBody,
            }
            data := body.Bytes()
            assert.Equal(t, tc.expect, data)
            assert.False(t, body.IsEmpty())
        })
    }
}
