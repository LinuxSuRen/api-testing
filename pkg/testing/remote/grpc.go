package remote

import (
	context "context"
	"fmt"

	"github.com/linuxsuren/api-testing/pkg/testing"
	"google.golang.org/grpc"
)

type gRPCLoader struct {
	address string
	client  LoaderClient
}

func NewGRPCLoader(address string) (writer testing.Writer, err error) {
	var conn *grpc.ClientConn
	if conn, err = grpc.Dial(address, grpc.WithInsecure()); err == nil {
		writer = &gRPCLoader{
			address: address,
			client:  NewLoaderClient(conn),
		}
	}
	return
}

func (g *gRPCLoader) HasMore() bool {
	// nothing to do
	return false
}

func (g *gRPCLoader) Load() ([]byte, error) {
	// nothing to do
	return nil, nil
}

func (g *gRPCLoader) Put(path string) error {
	// nothing to do
	return nil
}

func (g *gRPCLoader) GetContext() string {
	// nothing to do
	return ""
}

func (g *gRPCLoader) GetCount() int {
	// nothing to do
	return 0
}

func (g *gRPCLoader) Reset() {
	// nothing to do
}

func convertToGRPCTestCase(testcase testing.TestCase) (result *TestCase) {
	result = &TestCase{
		Name: testcase.Name,
		Request: &Request{
			Api:    testcase.Request.API,
			Method: testcase.Request.Method,
			Body:   testcase.Request.Body,
			Header: mapToPair(testcase.Request.Header),
			Query:  mapToPair(testcase.Request.Query),
			Form:   mapToPair(testcase.Request.Form),
		},
		Response: &Response{
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

func convertToNormalTestCase(testcase *TestCase) (result testing.TestCase) {
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

func mapToPair(data map[string]string) (pairs []*Pair) {
	pairs = make([]*Pair, 0)
	for k, v := range data {
		pairs = append(pairs, &Pair{
			Key:   k,
			Value: v,
		})
	}
	return
}

func mapInterToPair(data map[string]interface{}) (pairs []*Pair) {
	pairs = make([]*Pair, 0)
	for k, v := range data {
		pairs = append(pairs, &Pair{
			Key:   k,
			Value: fmt.Sprintf("%v", v),
		})
	}
	return
}

func pairToMap(pairs []*Pair) (data map[string]string) {
	data = make(map[string]string)
	for _, pair := range pairs {
		data[pair.Key] = pair.Value
	}
	return
}

func pairToInterMap(pairs []*Pair) (data map[string]interface{}) {
	data = make(map[string]interface{})
	for _, pair := range pairs {
		data[pair.Key] = pair.Value
	}
	return
}

func (g *gRPCLoader) ListTestCase(suite string) (testcases []testing.TestCase, err error) {
	var testCases *TestCases
	testCases, err = g.client.ListTestCases(context.Background(), &TestSuite{
		Name: suite,
	})

	if err == nil && testCases.Data != nil {
		for _, item := range testCases.Data {
			if item.Name == "" {
				continue
			}
			testcases = append(testcases, convertToNormalTestCase(item))
		}
	}
	return
}
func (g *gRPCLoader) GetTestCase(suite, name string) (testcase testing.TestCase, err error) {
	var result *TestCase
	result, err = g.client.GetTestCase(context.Background(), &TestCase{
		Name:      name,
		SuiteName: suite,
	})
	if err == nil && result != nil {
		testcase = convertToNormalTestCase(result)
	}
	return
}

func (g *gRPCLoader) CreateTestCase(suite string, testcase testing.TestCase) (err error) {
	payload := convertToGRPCTestCase(testcase)
	payload.SuiteName = suite
	_, err = g.client.CreateTestCase(context.Background(), payload)
	return
}

func (g *gRPCLoader) UpdateTestCase(suite string, testcase testing.TestCase) (err error) {
	payload := convertToGRPCTestCase(testcase)
	payload.SuiteName = suite
	_, err = g.client.UpdateTestCase(context.Background(), payload)
	return
}

func (g *gRPCLoader) DeleteTestCase(suite, testcase string) (err error) {
	_, err = g.client.DeleteTestCase(context.Background(), &TestCase{
		Name:      testcase,
		SuiteName: suite,
	})
	return
}

func (g *gRPCLoader) ListTestSuite() (suites []testing.TestSuite, err error) {
	var items *TestSuites
	items, err = g.client.ListTestSuite(context.Background(), &Empty{})
	if err == nil && items != nil {
		for _, item := range items.Data {
			suites = append(suites, testing.TestSuite{
				Name: item.Name,
				API:  item.Api,
			})
		}
	}
	return
}

func (g *gRPCLoader) GetTestSuite(name string, full bool) (suite testing.TestSuite, err error) {
	var result *TestSuite
	if result, err = g.client.GetTestSuite(context.Background(),
		&TestSuite{Name: name, Full: full}); err == nil {
		suite = testing.TestSuite{
			Name: result.Name,
			API:  result.Api,
		}

		if result.Items != nil {
			for i := range result.Items {
				suite.Items = append(suite.Items, convertToNormalTestCase(result.Items[i]))
			}
		}
	}
	return
}

func (g *gRPCLoader) CreateSuite(name, api string) (err error) {
	_, err = g.client.CreateTestSuite(context.Background(), &TestSuite{
		Name: name,
		Api:  api,
	})
	return
}

func (g *gRPCLoader) GetSuite(name string) (reply *testing.TestSuite, _ string, err error) {
	var suite *TestSuite
	if suite, err = g.client.GetTestSuite(context.Background(),
		&TestSuite{Name: name}); err != nil {
		return
	}

	reply = &testing.TestSuite{
		Name: suite.Name,
		API:  suite.Api,
	}
	return
}

func (g *gRPCLoader) UpdateSuite(name, api string) (err error) {
	_, err = g.client.UpdateTestSuite(context.Background(), &TestSuite{
		Name: name,
		Api:  api,
	})
	return
}

func (g *gRPCLoader) DeleteSuite(name string) (err error) {
	_, err = g.client.DeleteTestSuite(context.Background(), &TestSuite{
		Name: name,
	})
	return
}
