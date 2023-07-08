package pkg

import (
	"encoding/json"

	"github.com/linuxsuren/api-testing/pkg/testing/remote"
)

func ConverToDBTestCase(testcase *remote.TestCase) (result *TestCase) {
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
		result.ExpectVerify = sliceToJSON(resp.Verify)
	}
	return
}

func ConvertToRemoteTestCase(testcase *TestCase) (result *remote.TestCase) {
	result = &remote.TestCase{
		Name: testcase.Name,

		Request: &remote.Request{
			Api:    testcase.API,
			Method: testcase.Method,
			Body:   testcase.Body,
			Header: jsonToPair(testcase.Header),
			Query:  jsonToPair(testcase.Query),
			Form:   jsonToPair(testcase.Form),
		},

		Response: &remote.Response{
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
	return
}

func ConvertToGRPCTestSuite(suite *TestSuite) (result *remote.TestSuite) {
	result = &remote.TestSuite{
		Name: suite.Name,
		Api:  suite.API,
	}
	return
}

func sliceToJSON(slice []string) (result string) {
	var data []byte
	var err error
	if data, err = json.Marshal(slice); err == nil {
		result = string(data)
	}
	return
}

func pairToJSON(pair []*remote.Pair) (result string) {
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

func jsonToPair(jsonStr string) (pairs []*remote.Pair) {
	pairMap := make(map[string]string, 0)
	err := json.Unmarshal([]byte(jsonStr), &pairMap)
	if err == nil {
		for k, v := range pairMap {
			pairs = append(pairs, &remote.Pair{
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
