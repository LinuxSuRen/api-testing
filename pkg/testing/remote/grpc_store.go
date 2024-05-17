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

package remote

import (
	context "context"
	"errors"
	"time"

	"github.com/linuxsuren/api-testing/pkg/logging"
	server "github.com/linuxsuren/api-testing/pkg/server"
	"github.com/linuxsuren/api-testing/pkg/testing"

	"google.golang.org/grpc"
)

var (
	grpcLogger = logging.DefaultLogger(logging.LogLevelInfo).WithName("grpc")
)

type gRPCLoader struct {
	store  *testing.Store
	client LoaderClient
	ctx    context.Context
	conn   *grpc.ClientConn
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
			conn:   conn,
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

func (g *gRPCLoader) GetTestSuiteYaml(suite string) (testSuiteYaml []byte, err error) {
	return
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

func (g *gRPCLoader) Verify() (readOnly bool, err error) {
	// avoid to long to wait the response
	ctx, cancel := context.WithTimeout(g.ctx, time.Second*5)
	defer cancel()

	var result *server.ExtensionStatus
	if result, err = g.client.Verify(ctx, &server.Empty{}); err == nil {
		readOnly = result.ReadOnly
		if !result.Ready {
			err = errors.New(result.Message)
		}
	}
	return
}

func (g *gRPCLoader) PProf(name string) []byte {
	data, err := g.client.PProf(context.Background(), &server.PProfRequest{
		Name: name,
	})
	if err != nil {
		grpcLogger.Info("failed to get pprof:", "error", err)
	}
	return data.Data
}

func (g *gRPCLoader) Close() {
	if g.conn != nil {
		g.conn.Close()
	}
}
