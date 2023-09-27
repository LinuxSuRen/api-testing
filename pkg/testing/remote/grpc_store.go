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
	context "context"
	"errors"

	server "github.com/linuxsuren/api-testing/pkg/server"
	"github.com/linuxsuren/api-testing/pkg/testing"
	"google.golang.org/grpc"
)

type gRPCLoader struct {
	store  *testing.Store
	client LoaderClient
	ctx    context.Context
}

func NewGRPCloaderFromStore() testing.StoreWriterFactory {
	return &gRPCLoader{}
}

func (g *gRPCLoader) NewInstance(store testing.Store) (writer testing.Writer, err error) {
	address := store.Kind.URL

	var conn *grpc.ClientConn
	if conn, err = grpc.Dial(address, grpc.WithInsecure()); err == nil {
		writer = &gRPCLoader{
			store:  &store,
			ctx:    WithStoreContext(context.Background(), &store),
			client: NewLoaderClient(conn),
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
	var testCases *server.TestCases
	testCases, err = g.client.ListTestCases(g.ctx, &TestSuite{
		Name: suite,
	})

	if err == nil && testCases.Data != nil {
		for _, item := range testCases.Data {
			if item.Name == "" {
				continue
			}
			testcases = append(testcases, ConvertToNormalTestCase(item))
		}
	}
	return
}
func (g *gRPCLoader) GetTestCase(suite, name string) (testcase testing.TestCase, err error) {
	var result *server.TestCase
	result, err = g.client.GetTestCase(g.ctx, &server.TestCase{
		Name:      name,
		SuiteName: suite,
	})
	if err == nil && result != nil {
		testcase = ConvertToNormalTestCase(result)
	}
	return
}

func (g *gRPCLoader) CreateTestCase(suite string, testcase testing.TestCase) (err error) {
	payload := ConvertToGRPCTestCase(testcase)
	payload.SuiteName = suite
	_, err = g.client.CreateTestCase(g.ctx, payload)
	return
}

func (g *gRPCLoader) UpdateTestCase(suite string, testcase testing.TestCase) (err error) {
	payload := ConvertToGRPCTestCase(testcase)
	payload.SuiteName = suite
	_, err = g.client.UpdateTestCase(g.ctx, payload)
	return
}

func (g *gRPCLoader) DeleteTestCase(suite, testcase string) (err error) {
	_, err = g.client.DeleteTestCase(g.ctx, &server.TestCase{
		Name:      testcase,
		SuiteName: suite,
	})
	return
}

func (g *gRPCLoader) ListTestSuite() (suites []testing.TestSuite, err error) {
	var items *TestSuites
	items, err = g.client.ListTestSuite(g.ctx, &server.Empty{})
	if err == nil && items != nil {
		for _, item := range items.Data {
			suites = append(suites, *ConvertToNormalTestSuite(item))
		}
	}
	return
}

func (g *gRPCLoader) GetTestSuite(name string, full bool) (suite testing.TestSuite, err error) {
	var result *TestSuite
	if result, err = g.client.GetTestSuite(g.ctx,
		&TestSuite{Name: name, Full: full}); err == nil {
		suite = *ConvertToNormalTestSuite(result)
	}
	return
}

func (g *gRPCLoader) CreateSuite(name, api string) (err error) {
	_, err = g.client.CreateTestSuite(g.ctx, &TestSuite{
		Name: name,
		Api:  api,
	})
	return
}

func (g *gRPCLoader) GetSuite(name string) (reply *testing.TestSuite, _ string, err error) {
	var suite *TestSuite
	if suite, err = g.client.GetTestSuite(g.ctx,
		&TestSuite{Name: name}); err != nil {
		return
	}

	reply = ConvertToNormalTestSuite(suite)
	return
}

func (g *gRPCLoader) UpdateSuite(suite testing.TestSuite) (err error) {
	_, err = g.client.UpdateTestSuite(g.ctx, ConvertToGRPCTestSuite(&suite))
	return
}

func (g *gRPCLoader) DeleteSuite(name string) (err error) {
	_, err = g.client.DeleteTestSuite(g.ctx, &TestSuite{
		Name: name,
	})
	return
}

func (g *gRPCLoader) Verify() (err error) {
	var result *server.CommonResult
	if result, err = g.client.Verify(g.ctx, &server.Empty{}); err == nil {
		if !result.Success {
			err = errors.New(result.Message)
		}
	}
	return
}
