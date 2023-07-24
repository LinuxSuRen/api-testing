package remote

import (
	context "context"

	"github.com/linuxsuren/api-testing/pkg/testing"
	"google.golang.org/grpc"
)

type gRPCLoader struct {
	address string
	client  LoaderClient
}

// NewGRPCLoader creates a new gRPC loader
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
			suites = append(suites, *convertToNormalTestSuite(item))
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

	reply = convertToNormalTestSuite(suite)
	return
}

func (g *gRPCLoader) UpdateSuite(suite testing.TestSuite) (err error) {
	_, err = g.client.UpdateTestSuite(context.Background(), convertToGRPCTestSuite(&suite))
	return
}

func (g *gRPCLoader) DeleteSuite(name string) (err error) {
	_, err = g.client.DeleteTestSuite(context.Background(), &TestSuite{
		Name: name,
	})
	return
}
