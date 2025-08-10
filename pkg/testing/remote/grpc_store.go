/*
Copyright 2023-2025 API Testing Authors.

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
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/linuxsuren/api-testing/pkg/logging"
	"github.com/linuxsuren/api-testing/pkg/mock"
	server "github.com/linuxsuren/api-testing/pkg/server"
	"github.com/linuxsuren/api-testing/pkg/testing"

	"google.golang.org/grpc"
)

var (
	grpcLogger = logging.DefaultLogger(logging.LogLevelInfo).WithName("grpc")
)

type gRPCLoader struct {
	store    *testing.Store
	client   LoaderClient
	ctx      context.Context
	conn     *grpc.ClientConn
	mockData []byte
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
	var testSuite testing.TestSuite
	if testSuite, err = g.GetTestSuite(suite, true); err == nil {
		testSuiteYaml, err = testing.ToYAML(&testSuite)
	}
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

func (g *gRPCLoader) GetHistoryTestCaseWithResult(id string) (result testing.HistoryTestResult, err error) {
	var historyTestResult *server.HistoryTestResult
	historyTestResult, err = g.client.GetHistoryTestCaseWithResult(g.ctx, &server.HistoryTestCase{
		ID: id,
	})
	if err == nil && historyTestResult != nil {
		result = ConvertToNormalTestCaseResult(historyTestResult)
	}
	return
}

func (g *gRPCLoader) GetHistoryTestCase(id string) (result testing.HistoryTestCase, err error) {
	var historyTestCase *server.HistoryTestCase
	historyTestCase, err = g.client.GetHistoryTestCase(g.ctx, &server.HistoryTestCase{
		ID: id,
	})
	if err == nil && historyTestCase != nil {
		result = ConvertToNormalHistoryTestCase(historyTestCase)
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

func (g *gRPCLoader) RenameTestCase(suite, oldName, newName string) (err error) {
	_, err = g.client.RenameTestCase(g.ctx, &server.TestCaseDuplicate{
		SourceCaseName:  oldName,
		SourceSuiteName: suite,
		TargetCaseName:  newName,
	})
	return
}

func (g *gRPCLoader) CreateHistoryTestCase(testcaseResult testing.TestCaseResult, testSuite *testing.TestSuite, historyHeader map[string]string) (err error) {
	payload := ConvertToGRPCHistoryTestCaseResult(testcaseResult, testSuite, historyHeader)
	_, err = g.client.CreateTestCaseHistory(g.ctx, payload)
	return
}

func (g *gRPCLoader) ListHistoryTestSuite() (suites []testing.HistoryTestSuite, err error) {
	var items *HistoryTestSuites
	items, err = g.client.ListHistoryTestSuite(g.ctx, &server.Empty{})
	if err == nil && items != nil {
		for _, item := range items.Data {
			suites = append(suites, *ConvertToNormalHistoryTestSuite(item))
		}
	}
	return
}

func (g *gRPCLoader) DeleteHistoryTestCase(id string) (err error) {
	_, err = g.client.DeleteHistoryTestCase(g.ctx, &server.HistoryTestCase{
		ID: id,
	})
	return
}

func (g *gRPCLoader) DeleteAllHistoryTestCase(suite, name string) (err error) {
	_, err = g.client.DeleteAllHistoryTestCase(g.ctx, &server.HistoryTestCase{
		SuiteName: suite,
		CaseName:  name,
	})
	return
}

func (g *gRPCLoader) GetTestCaseAllHistory(suite, name string) (historyTestcases []testing.HistoryTestCase, err error) {
	var historyTestCases *server.HistoryTestCases
	historyTestCases, err = g.client.GetTestCaseAllHistory(g.ctx, &server.TestCase{
		Name:      name,
		SuiteName: suite,
	})
	if err == nil && historyTestCases.Data != nil {
		for _, item := range historyTestCases.Data {
			historyTestcases = append(historyTestcases, ConvertToNormalHistoryTestCase(item))
		}
	}
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

func (g *gRPCLoader) RenameTestSuite(oldName, newName string) (err error) {
	_, err = g.client.RenameTestSuite(g.ctx, &server.TestSuiteDuplicate{
		SourceSuiteName: oldName,
		TargetSuiteName: newName,
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

func (g *gRPCLoader) Query(query map[string]string) (result testing.DataResult, err error) {
	var dataResult *server.DataQueryResult
	offset, _ := strconv.ParseInt(query["offset"], 10, 64)
	limit, _ := strconv.ParseInt(query["limit"], 10, 64)
	if dataResult, err = g.client.Query(g.ctx, &server.DataQuery{
		Sql:    query["sql"],
		Key:    query["key"],
		Offset: offset,
		Limit:  limit,
	}); err == nil {
		result.Pairs = pairToMap(dataResult.Data)
		for _, item := range dataResult.Items {
			result.Rows = append(result.Rows, pairToMap(item.Data))
		}

		if dataResult.Meta != nil {
			result.Databases = dataResult.Meta.Databases
			result.Tables = dataResult.Meta.Tables
			result.CurrentDatabase = dataResult.Meta.CurrentDatabase
			result.Duration = dataResult.Meta.Duration
			result.Labels = pairToMap(dataResult.Meta.Labels)
		}
	}
	return
}

func (g *gRPCLoader) GetThemes() (result []string, err error) {
	var simpleList *server.SimpleList
	if simpleList, err = g.client.GetThemes(g.ctx, &server.Empty{}); err == nil && simpleList.Data != nil {
		for _, item := range simpleList.Data {
			result = append(result, item.Key)
		}
	}
	return
}

func (g *gRPCLoader) GetTheme(name string) (result string, err error) {
	var themeData *server.CommonResult
	if themeData, err = g.client.GetTheme(g.ctx, &server.SimpleName{
		Name: name,
	}); err == nil && themeData != nil {
		result = themeData.Message
	}
	return
}

func (g *gRPCLoader) GetBindings() (result []string, err error) {
	var simpleList *server.SimpleList
	if simpleList, err = g.client.GetBindings(g.ctx, &server.Empty{}); err == nil && simpleList.Data != nil {
		for _, item := range simpleList.Data {
			result = append(result, item.Key)
		}
	}
	return
}

func (g *gRPCLoader) GetBinding(name string) (result string, err error) {
	var themeData *server.CommonResult
	if themeData, err = g.client.GetBinding(g.ctx, &server.SimpleName{
		Name: name,
	}); err == nil && themeData != nil {
		result = themeData.Message
	}
	return
}

func (g *gRPCLoader) GetMenus() (result []*testing.Menu, err error) {
	fmt.Println("getting menus from grpc server", g.store.Kind)

	var menuList *server.MenuList
	if menuList, err = g.client.GetMenus(g.ctx, &server.Empty{}); err == nil {
		for _, item := range menuList.Data {
			result = append(result, &testing.Menu{
				Name:  item.Name,
				Icon:  item.Icon,
				Index: item.Index,
			})
		}
	}
	return
}

func (g *gRPCLoader) GetPageOfJS(name string) (result string, err error) {
	var themeData *server.CommonResult
	if themeData, err = g.client.GetPageOfJS(g.ctx, &server.SimpleName{
		Name: name,
	}); err == nil && themeData != nil {
		result = themeData.Message
	}
	return
}

func (g *gRPCLoader) GetPageOfCSS(name string) (result string, err error) {
	var themeData *server.CommonResult
	if themeData, err = g.client.GetPageOfCSS(g.ctx, &server.SimpleName{
		Name: name,
	}); err == nil && themeData != nil {
		result = themeData.Message
	}
	return
}

func (g *gRPCLoader) Parse() (server *mock.Server, err error) {
	return
}

func (g *gRPCLoader) GetData() []byte {
	return g.mockData
}

func (g *gRPCLoader) Write(data []byte) {
	g.mockData = data
}

func (g *gRPCLoader) Close() {
	if g.conn != nil {
		g.conn.Close()
	}
}
