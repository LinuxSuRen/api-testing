package pkg

import (
	"encoding/json"

	"github.com/linuxsuren/api-testing/pkg/server"
	"github.com/linuxsuren/api-testing/pkg/testing/remote"
)

func ConverToDBTestCase(testcase *server.TestCase) (result *TestCase) {
	result = &TestCase{
		Name:      testcase.Name,
		SuiteName: testcase.SuiteName,
	}
	request := testcase.Request
	if request != nil {
		result.API = request.Api
		result.Method = request.Method
		result.Body = request.Body
		result.Header = pairToJSON(request.Header)
		result.Form = pairToJSON(request.Form)
		result.Query = pairToJSON(request.Query)
	}

	resp := testcase.Response
	if resp != nil {
		result.ExpectBody = resp.Body
		result.ExpectSchema = resp.Schema
		result.ExpectStatusCode = int(resp.StatusCode)
		result.ExpectHeader = pairToJSON(resp.Header)
		result.ExpectBodyFields = pairToJSON(resp.BodyFieldsExpect)
		result.ExpectVerify = SliceToJSON(resp.Verify)
	}
	return
}

func ConvertToRemoteTestCase(testcase *TestCase) (result *server.TestCase) {
	result = &server.TestCase{
		Name: testcase.Name,

		Request: &server.Request{
			Api:    testcase.API,
			Method: testcase.Method,
			Body:   testcase.Body,
			Header: jsonToPair(testcase.Header),
			Query:  jsonToPair(testcase.Query),
			Form:   jsonToPair(testcase.Form),
		},

		Response: &server.Response{
			StatusCode:       int32(testcase.ExpectStatusCode),
			Body:             testcase.ExpectBody,
			Schema:           testcase.ExpectSchema,
			Verify:           jsonToSlice(testcase.ExpectVerify),
			BodyFieldsExpect: jsonToPair(testcase.ExpectBodyFields),
			Header:           jsonToPair(testcase.ExpectHeader),
		},
	}
	return
}

func ConvertToDBTestSuite(suite *remote.TestSuite) (result *TestSuite) {
	result = &TestSuite{
		Name: suite.Name,
		API:  suite.Api,
	}
	if suite.Spec != nil {
		result.SpecKind = suite.Spec.Kind
		result.SpecURL = suite.Spec.Url
	}
	if suite.Param != nil {
		result.Param = pairToJSON(suite.Param)
	}
	return
}

func ConvertToGRPCTestSuite(suite *TestSuite) (result *remote.TestSuite) {
	result = &remote.TestSuite{
		Name: suite.Name,
		Api:  suite.API,
		Spec: &server.APISpec{
			Kind: suite.SpecKind,
			Url:  suite.SpecURL,
		},
		Param: jsonToPair(suite.Param),
	}
	return
}

func SliceToJSON(slice []string) (result string) {
	var data []byte
	var err error
	if slice != nil {
		if data, err = json.Marshal(slice); err == nil {
			result = string(data)
		}
	}
	if result == "" {
		result = "[]"
	}
	return
}

func pairToJSON(pair []*server.Pair) (result string) {
	var obj = make(map[string]string)
	for i := range pair {
		k := pair[i].Key
		v := pair[i].Value
		obj[k] = v
	}

	var data []byte
	var err error
	if data, err = json.Marshal(obj); err == nil {
		result = string(data)
	}
	return
}

func jsonToPair(jsonStr string) (pairs []*server.Pair) {
	pairMap := make(map[string]string, 0)
	err := json.Unmarshal([]byte(jsonStr), &pairMap)
	if err == nil {
		for k, v := range pairMap {
			pairs = append(pairs, &server.Pair{
				Key: k, Value: v,
			})
		}
	}
	return
}

func jsonToSlice(jsonStr string) (result []string) {
	_ = json.Unmarshal([]byte(jsonStr), &result)
	return
}
