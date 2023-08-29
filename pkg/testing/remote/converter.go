package remote

import (
	"fmt"

	server "github.com/linuxsuren/api-testing/pkg/server"
	"github.com/linuxsuren/api-testing/pkg/testing"
)

func ConvertToNormalTestSuite(suite *TestSuite) (result *testing.TestSuite) {
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
	}

	for _, testcase := range suite.Items {
		result.Items = append(result.Items, ConvertToNormalTestCase(testcase))
	}
	return
}

func ConvertToGRPCTestSuite(suite *testing.TestSuite) (result *TestSuite) {
	result = &TestSuite{
		Name:  suite.Name,
		Api:   suite.API,
		Param: mapToPair(suite.Param),
		Spec: &server.APISpec{
			Kind: suite.Spec.Kind,
			Url:  suite.Spec.URL,
		},
	}

	for _, testcase := range suite.Items {
		result.Items = append(result.Items, ConvertToGRPCTestCase(testcase))
	}
	return
}

func ConvertToNormalTestCase(testcase *server.TestCase) (result testing.TestCase) {
	result = testing.TestCase{
		Name: testcase.Name,
	}
	if testcase.Request != nil {
		result.Request = testing.Request{
			API:    testcase.Request.Api,
			Method: testcase.Request.Method,
			Body:   testcase.Request.Body,
			Header: pairToMap(testcase.Request.Header),
			Query:  pairToMap(testcase.Request.Query),
			Form:   pairToMap(testcase.Request.Form),
		}
	}
	if testcase.Response != nil {
		result.Expect = testing.Response{
			Body:             testcase.Response.Body,
			StatusCode:       int(testcase.Response.StatusCode),
			Schema:           testcase.Response.Schema,
			Verify:           testcase.Response.Verify,
			Header:           pairToMap(testcase.Response.Header),
			BodyFieldsExpect: pairToInterMap(testcase.Response.BodyFieldsExpect),
		}
	}
	return
}

func ConvertToGRPCTestCase(testcase testing.TestCase) (result *server.TestCase) {
	result = &server.TestCase{
		Name: testcase.Name,
		Request: &server.Request{
			Api:    testcase.Request.API,
			Method: testcase.Request.Method,
			Body:   testcase.Request.Body,
			Header: mapToPair(testcase.Request.Header),
			Query:  mapToPair(testcase.Request.Query),
			Form:   mapToPair(testcase.Request.Form),
		},
		Response: &server.Response{
			Body:             testcase.Expect.Body,
			StatusCode:       int32(testcase.Expect.StatusCode),
			Schema:           testcase.Expect.Schema,
			Verify:           testcase.Expect.Verify,
			Header:           mapToPair(testcase.Expect.Header),
			BodyFieldsExpect: mapInterToPair(testcase.Expect.BodyFieldsExpect),
		},
	}
	return
}

func mapToPair(data map[string]string) (pairs []*server.Pair) {
	pairs = make([]*server.Pair, 0)
	for k, v := range data {
		pairs = append(pairs, &server.Pair{
			Key:   k,
			Value: v,
		})
	}
	return
}

func mapInterToPair(data map[string]interface{}) (pairs []*server.Pair) {
	pairs = make([]*server.Pair, 0)
	for k, v := range data {
		pairs = append(pairs, &server.Pair{
			Key:   k,
			Value: fmt.Sprintf("%v", v),
		})
	}
	return
}

func pairToMap(pairs []*server.Pair) (data map[string]string) {
	data = make(map[string]string)
	for _, pair := range pairs {
		data[pair.Key] = pair.Value
	}
	return
}

func pairToInterMap(pairs []*server.Pair) (data map[string]interface{}) {
	data = make(map[string]interface{})
	for _, pair := range pairs {
		data[pair.Key] = pair.Value
	}
	return
}
