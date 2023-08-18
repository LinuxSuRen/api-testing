// Package server provides a GRPC based server
package server

import (
	"bytes"
	context "context"
	"errors"
	"fmt"
	"io"
	"os"
	reflect "reflect"
	"regexp"
	"strings"

	_ "embed"

	"github.com/linuxsuren/api-testing/pkg/apispec"
	"github.com/linuxsuren/api-testing/pkg/generator"
	"github.com/linuxsuren/api-testing/pkg/render"
	"github.com/linuxsuren/api-testing/pkg/runner"
	"github.com/linuxsuren/api-testing/pkg/testing"
	"github.com/linuxsuren/api-testing/pkg/util"
	"github.com/linuxsuren/api-testing/pkg/version"
	"github.com/linuxsuren/api-testing/sample"
	"google.golang.org/grpc/metadata"
	"gopkg.in/yaml.v3"
)

type server struct {
	UnimplementedRunnerServer
	loader             testing.Writer
	storeWriterFactory testing.StoreWriterFactory
	configDir          string

	secretServer SecretServiceServer
}

type SecretServiceServer interface {
	GetSecrets(context.Context, *Empty) (*Secrets, error)
	CreateSecret(context.Context, *Secret) (*CommonResult, error)
	DeleteSecret(context.Context, *Secret) (*CommonResult, error)
	UpdateSecret(context.Context, *Secret) (*CommonResult, error)
}

type SecertServiceGetable interface {
	GetSecret(context.Context, *Secret) (*Secret, error)
}

type fakeSecretServer struct{}

var errNoSecretService = errors.New("no secret service found")

func (f *fakeSecretServer) GetSecrets(ctx context.Context, in *Empty) (reply *Secrets, err error) {
	err = errNoSecretService
	return
}

func (f *fakeSecretServer) CreateSecret(ctx context.Context, in *Secret) (reply *CommonResult, err error) {
	err = errNoSecretService
	return
}

func (f *fakeSecretServer) DeleteSecret(ctx context.Context, in *Secret) (reply *CommonResult, err error) {
	err = errNoSecretService
	return
}

func (f *fakeSecretServer) UpdateSecret(ctx context.Context, in *Secret) (reply *CommonResult, err error) {
	err = errNoSecretService
	return
}

// NewRemoteServer creates a remote server instance
func NewRemoteServer(loader testing.Writer, storeWriterFactory testing.StoreWriterFactory, secretServer SecretServiceServer, configDir string) RunnerServer {
	if secretServer == nil {
		secretServer = &fakeSecretServer{}
	}

	return &server{
		loader:             loader,
		storeWriterFactory: storeWriterFactory,
		configDir:          configDir,
		secretServer:       secretServer,
	}
}

func withDefaultValue(old, defVal any) any {
	if old == "" || old == nil {
		old = defVal
	}
	return old
}

func parseSuiteWithItems(data []byte) (suite *testing.TestSuite, err error) {
	suite, err = testing.ParseFromData(data)
	if err == nil && (suite == nil || suite.Items == nil) {
		err = errNoTestSuiteFound
	}
	return
}

func (s *server) getSuiteFromTestTask(task *TestTask) (suite *testing.TestSuite, err error) {
	switch task.Kind {
	case "suite":
		suite, err = parseSuiteWithItems([]byte(task.Data))
	case "testcase":
		var testCase *testing.TestCase
		if testCase, err = testing.ParseTestCaseFromData([]byte(task.Data)); err != nil {
			return
		}
		suite = &testing.TestSuite{
			Items: []testing.TestCase{*testCase},
		}
	case "testcaseInSuite":
		suite, err = parseSuiteWithItems([]byte(task.Data))
		if err != nil {
			return
		}

		var targetTestcase *testing.TestCase
		for _, item := range suite.Items {
			if item.Name == task.CaseName {
				targetTestcase = &item
				break
			}
		}

		if targetTestcase != nil {
			parentCases := findParentTestCases(targetTestcase, suite)
			fmt.Printf("find %d parent cases\n", len(parentCases))
			suite.Items = append(parentCases, *targetTestcase)
		} else {
			err = fmt.Errorf("cannot found testcase %s", task.CaseName)
		}
	default:
		err = fmt.Errorf("not support '%s'", task.Kind)
	}
	return
}

func resetEnv(oldEnv map[string]string) {
	for key, val := range oldEnv {
		os.Setenv(key, val)
	}
}

func (s *server) getLoader(ctx context.Context) (loader testing.Writer) {
	var ok bool
	loader = s.loader

	var mdd metadata.MD
	if mdd, ok = metadata.FromIncomingContext(ctx); ok {
		storeNameMeta := mdd.Get(HeaderKeyStoreName)
		if len(storeNameMeta) > 0 {
			storeName := storeNameMeta[0]
			if storeName == "local" || storeName == "" {
				return
			}

			loader, _ = s.getLoaderByStoreName(storeName)
		}
	}
	return
}

// Run start to run the test task
func (s *server) Run(ctx context.Context, task *TestTask) (reply *TestResult, err error) {
	task.Level = withDefaultValue(task.Level, "info").(string)
	task.Env = withDefaultValue(task.Env, map[string]string{}).(map[string]string)

	var suite *testing.TestSuite

	// TODO may not safe in multiple threads
	oldEnv := map[string]string{}
	for key, val := range task.Env {
		oldEnv[key] = os.Getenv(key)
		os.Setenv(key, val)
	}

	defer func() {
		resetEnv(oldEnv)
	}()

	if suite, err = s.getSuiteFromTestTask(task); err != nil {
		return
	}

	fmt.Printf("prepare to run: %s, with level: %s\n", suite.Name, task.Level)
	fmt.Printf("task kind: %s, %d to run\n", task.Kind, len(suite.Items))
	dataContext := map[string]interface{}{}

	if err = suite.Render(dataContext); err != nil {
		reply.Error = err.Error()
		err = nil
		return
	}

	buf := new(bytes.Buffer)
	reply = &TestResult{}

	for _, testCase := range suite.Items {
		suiteRunner := runner.GetTestSuiteRunner(suite)
		suiteRunner.WithOutputWriter(buf)
		suiteRunner.WithWriteLevel(task.Level)

		// reuse the API prefix
		testCase.Request.RenderAPI(suite.API)

		output, testErr := suiteRunner.RunTestCase(&testCase, dataContext, ctx)
		if getter, ok := suiteRunner.(runner.HTTPResponseRecord); ok {
			resp := getter.GetResponseRecord()
			reply.TestCaseResult = append(reply.TestCaseResult, &TestCaseResult{
				StatusCode: int32(resp.StatusCode),
				Body:       resp.Body,
				Header:     mapToPair(resp.Header),
				Id:         testCase.ID,
				Output:     buf.String(),
			})
		}

		if testErr == nil {
			dataContext[testCase.Name] = output
		} else {
			reply.Error = testErr.Error()
			break
		}
	}

	if reply.Error != "" {
		fmt.Fprintln(buf, reply.Error)
	}
	reply.Message = buf.String()
	return
}

// GetVersion returns the version
func (s *server) GetVersion(ctx context.Context, in *Empty) (reply *HelloReply, err error) {
	reply = &HelloReply{Message: version.GetVersion()}
	return
}

func (s *server) GetSuites(ctx context.Context, in *Empty) (reply *Suites, err error) {
	loader := s.getLoader(ctx)
	reply = &Suites{
		Data: make(map[string]*Items),
	}

	var suites []testing.TestSuite
	if suites, err = loader.ListTestSuite(); err == nil && suites != nil {
		for _, suite := range suites {
			items := &Items{}
			for _, item := range suite.Items {
				items.Data = append(items.Data, item.Name)
			}
			reply.Data[suite.Name] = items
		}
	}

	return
}

func (s *server) CreateTestSuite(ctx context.Context, in *TestSuiteIdentity) (reply *HelloReply, err error) {
	loader := s.getLoader(ctx)
	err = loader.CreateSuite(in.Name, in.Api)
	return
}

func (s *server) GetTestSuite(ctx context.Context, in *TestSuiteIdentity) (result *TestSuite, err error) {
	loader := s.getLoader(ctx)
	var suite *testing.TestSuite
	if suite, _, err = loader.GetSuite(in.Name); err == nil && suite != nil {
		result = &TestSuite{
			Name:  suite.Name,
			Api:   suite.API,
			Param: mapToPair(suite.Param),
			Spec: &APISpec{
				Kind: suite.Spec.Kind,
				Url:  suite.Spec.URL,
			},
		}
	}
	return
}

func convertToTestingTestSuite(in *TestSuite) (suite *testing.TestSuite) {
	suite = &testing.TestSuite{
		Name:  in.Name,
		API:   in.Api,
		Param: pairToMap(in.Param),
	}
	if in.Spec != nil {
		suite.Spec = testing.APISpec{
			Kind: in.Spec.Kind,
			URL:  in.Spec.Url,
		}
	}
	return
}

func (s *server) UpdateTestSuite(ctx context.Context, in *TestSuite) (reply *HelloReply, err error) {
	reply = &HelloReply{}
	loader := s.getLoader(ctx)
	err = loader.UpdateSuite(*convertToTestingTestSuite(in))
	return
}

func (s *server) DeleteTestSuite(ctx context.Context, in *TestSuiteIdentity) (reply *HelloReply, err error) {
	reply = &HelloReply{}
	loader := s.getLoader(ctx)
	err = loader.DeleteSuite(in.Name)
	return
}

func (s *server) ListTestCase(ctx context.Context, in *TestSuiteIdentity) (result *Suite, err error) {
	var items []testing.TestCase
	loader := s.getLoader(ctx)
	if items, err = loader.ListTestCase(in.Name); err == nil {
		result = &Suite{}
		for _, item := range items {
			result.Items = append(result.Items, convertToGRPCTestCase(item))
		}
	}
	return
}

func (s *server) GetTestCase(ctx context.Context, in *TestCaseIdentity) (reply *TestCase, err error) {
	var result testing.TestCase
	loader := s.getLoader(ctx)
	if result, err = loader.GetTestCase(in.Suite, in.Testcase); err == nil {
		reply = convertToGRPCTestCase(result)
	}
	return
}

func convertToGRPCTestCase(testCase testing.TestCase) (result *TestCase) {
	req := &Request{
		Api:    testCase.Request.API,
		Method: testCase.Request.Method,
		Query:  mapToPair(testCase.Request.Query),
		Header: mapToPair(testCase.Request.Header),
		Form:   mapToPair(testCase.Request.Form),
		Body:   testCase.Request.Body,
	}

	resp := &Response{
		StatusCode:       int32(testCase.Expect.StatusCode),
		Body:             testCase.Expect.Body,
		Header:           mapToPair(testCase.Expect.Header),
		BodyFieldsExpect: mapInterToPair(testCase.Expect.BodyFieldsExpect),
		Verify:           testCase.Expect.Verify,
		Schema:           testCase.Expect.Schema,
	}

	result = &TestCase{
		Name:     testCase.Name,
		Request:  req,
		Response: resp,
	}
	return
}

func (s *server) RunTestCase(ctx context.Context, in *TestCaseIdentity) (result *TestCaseResult, err error) {
	var targetTestSuite testing.TestSuite

	loader := s.getLoader(ctx)
	targetTestSuite, err = loader.GetTestSuite(in.Suite, true)
	if err != nil {
		err = nil
		result = &TestCaseResult{
			Error: fmt.Sprintf("not found suite: %s", in.Suite),
		}
		return
	}

	var data []byte
	if data, err = yaml.Marshal(targetTestSuite); err == nil {
		task := &TestTask{
			Kind:     "testcaseInSuite",
			Data:     string(data),
			CaseName: in.Testcase,
			Level:    "debug",
		}

		var reply *TestResult
		if reply, err = s.Run(ctx, task); err == nil && len(reply.TestCaseResult) > 0 {
			lastIndex := len(reply.TestCaseResult) - 1
			lastItem := reply.TestCaseResult[lastIndex]

			result = &TestCaseResult{
				Output:     reply.Message,
				Error:      reply.Error,
				Body:       lastItem.Body,
				Header:     lastItem.Header,
				StatusCode: lastItem.StatusCode,
			}
		} else if err != nil {
			result = &TestCaseResult{
				Error: err.Error(),
			}
		} else {
			result = &TestCaseResult{
				Output: reply.Message,
				Error:  reply.Error,
			}
		}
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

func pairToInterMap(pairs []*Pair) (data map[string]interface{}) {
	data = make(map[string]interface{})
	for _, pair := range pairs {
		data[pair.Key] = pair.Value
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

func convertToTestingTestCase(in *TestCase) (result testing.TestCase) {
	result = testing.TestCase{
		Name: in.Name,
	}
	req := in.Request
	resp := in.Response

	if req != nil {
		result.Request.API = req.Api
		result.Request.Method = req.Method
		result.Request.Body = req.Body
		result.Request.Header = pairToMap(req.Header)
		result.Request.Form = pairToMap(req.Form)
		result.Request.Query = pairToMap(req.Query)
	}

	if resp != nil {
		result.Expect.Body = resp.Body
		result.Expect.Schema = resp.Schema
		result.Expect.StatusCode = int(resp.StatusCode)
		result.Expect.Verify = resp.Verify
		result.Expect.BodyFieldsExpect = pairToInterMap(resp.BodyFieldsExpect)
		result.Expect.Header = pairToMap(resp.Header)
	}
	return
}

func (s *server) CreateTestCase(ctx context.Context, in *TestCaseWithSuite) (reply *HelloReply, err error) {
	reply = &HelloReply{}
	loader := s.getLoader(ctx)
	err = loader.CreateTestCase(in.SuiteName, convertToTestingTestCase(in.Data))
	return
}

func (s *server) UpdateTestCase(ctx context.Context, in *TestCaseWithSuite) (reply *HelloReply, err error) {
	reply = &HelloReply{}
	if in.Data == nil {
		err = errors.New("data is required")
		return
	}
	loader := s.getLoader(ctx)
	err = loader.UpdateTestCase(in.SuiteName, convertToTestingTestCase(in.Data))
	return
}

func (s *server) DeleteTestCase(ctx context.Context, in *TestCaseIdentity) (reply *HelloReply, err error) {
	loader := s.getLoader(ctx)
	err = loader.DeleteTestCase(in.Suite, in.Testcase)
	return
}

// code generator
func (s *server) ListCodeGenerator(ctx context.Context, in *Empty) (reply *SimpleList, err error) {
	reply = &SimpleList{}

	generators := generator.GetCodeGenerators()
	for name := range generators {
		reply.Data = append(reply.Data, &Pair{
			Key: name,
		})
	}
	return
}

func (s *server) GenerateCode(ctx context.Context, in *CodeGenerateRequest) (reply *CommonResult, err error) {
	reply = &CommonResult{}

	instance := generator.GetCodeGenerator(in.Generator)
	if instance == nil {
		reply.Success = false
		reply.Message = fmt.Sprintf("generator '%s' not found", in.Generator)
	} else {
		var result testing.TestCase
		loader := s.getLoader(ctx)
		if result, err = loader.GetTestCase(in.TestSuite, in.TestCase); err == nil {
			output, genErr := instance.Generate(&result)
			reply.Success = genErr == nil
			reply.Message = util.OrErrorMessage(genErr, output)
		}
	}
	return
}

// converter
func (s *server) ListConverter(ctx context.Context, in *Empty) (reply *SimpleList, err error) {
	reply = &SimpleList{}
	converters := generator.GetTestSuiteConverters()
	for name := range converters {
		reply.Data = append(reply.Data, &Pair{
			Key: name,
		})
	}
	return
}

func (s *server) ConvertTestSuite(ctx context.Context, in *CodeGenerateRequest) (reply *CommonResult, err error) {
	reply = &CommonResult{}

	instance := generator.GetTestSuiteConverter(in.Generator)
	if instance == nil {
		reply.Success = false
		reply.Message = fmt.Sprintf("converter '%s' not found", in.Generator)
	} else {
		var result testing.TestSuite
		loader := s.getLoader(ctx)
		if result, err = loader.GetTestSuite(in.TestSuite, true); err == nil {
			output, genErr := instance.Convert(&result)
			reply.Success = genErr == nil
			reply.Message = util.OrErrorMessage(genErr, output)
		}
	}
	return
}

// Sample returns a sample of the test task
func (s *server) Sample(ctx context.Context, in *Empty) (reply *HelloReply, err error) {
	reply = &HelloReply{Message: sample.TestSuiteGitLab}
	return
}

// PopularHeaders returns a list of popular headers
func (s *server) PopularHeaders(ctx context.Context, in *Empty) (pairs *Pairs, err error) {
	pairs = &Pairs{
		Data: []*Pair{},
	}

	err = yaml.Unmarshal([]byte(popularHeaders), &pairs.Data)
	return
}

// GetSuggestedAPIs returns a list of suggested APIs
func (s *server) GetSuggestedAPIs(ctx context.Context, in *TestSuiteIdentity) (reply *TestCases, err error) {
	reply = &TestCases{}

	var suite *testing.TestSuite
	loader := s.getLoader(ctx)
	if suite, _, err = loader.GetSuite(in.Name); err != nil || suite == nil {
		return
	}

	if suite.Spec.URL == "" {
		return
	}

	fmt.Println("Finding APIs from", in.Name, "with loader", reflect.TypeOf(loader))
	var swaggerAPI *apispec.Swagger
	if swaggerAPI, err = apispec.ParseURLToSwagger(suite.Spec.URL); err == nil && swaggerAPI != nil {
		for api, item := range swaggerAPI.Paths {
			for method, oper := range item {
				reply.Data = append(reply.Data, &TestCase{
					Name: oper.OperationId,
					Request: &Request{
						Api:    api,
						Method: strings.ToUpper(method),
					},
				})
			}
		}
	}
	return
}

// FunctionsQuery returns a list of functions
func (s *server) FunctionsQuery(ctx context.Context, in *SimpleQuery) (reply *Pairs, err error) {
	reply = &Pairs{}
	in.Name = strings.ToLower(in.Name)

	for name, fn := range render.FuncMap() {
		lowerCaseName := strings.ToLower(name)
		if in.Name == "" || strings.Contains(lowerCaseName, in.Name) {
			reply.Data = append(reply.Data, &Pair{
				Key:   name,
				Value: fmt.Sprintf("%v", reflect.TypeOf(fn)),
			})
		}
	}
	return
}

// FunctionsQueryStream works like FunctionsQuery but is implemented in bidirectional streaming
func (s *server) FunctionsQueryStream(srv Runner_FunctionsQueryStreamServer) error {
	ctx := srv.Context()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			in, err := srv.Recv()
			if err != nil {
				if err == io.EOF {
					return nil
				}
				return err
			}
			reply := &Pairs{}
			in.Name = strings.ToLower(in.Name)

			for name, fn := range render.FuncMap() {
				lowerCaseName := strings.ToLower(name)
				if in.Name == "" || strings.Contains(lowerCaseName, in.Name) {
					reply.Data = append(reply.Data, &Pair{
						Key:   name,
						Value: fmt.Sprintf("%v", reflect.TypeOf(fn)),
					})
				}
			}
			if err := srv.Send(reply); err != nil {
				return err
			}
		}
	}
}

func (s *server) GetStores(ctx context.Context, in *Empty) (reply *Stores, err error) {
	storeFactory := testing.NewStoreFactory(s.configDir)
	var stores []testing.Store
	if stores, err = storeFactory.GetStores(); err == nil {
		reply = &Stores{
			Data: make([]*Store, 0),
		}
		for _, item := range stores {
			grpcStore := ToGRPCStore(item)

			storeStatus, sErr := s.VerifyStore(ctx, &SimpleQuery{Name: item.Name})
			grpcStore.Ready = sErr == nil && storeStatus.Success
			grpcStore.Password = "******" // return a placeholder instead of the actual value for the security reason

			reply.Data = append(reply.Data, grpcStore)
		}
		reply.Data = append(reply.Data, &Store{
			Name:  "local",
			Kind:  &StoreKind{},
			Ready: true,
		})
	}
	return
}
func (s *server) CreateStore(ctx context.Context, in *Store) (reply *Store, err error) {
	reply = &Store{}
	storeFactory := testing.NewStoreFactory(s.configDir)
	err = storeFactory.CreateStore(ToNormalStore(in))
	return
}
func (s *server) UpdateStore(ctx context.Context, in *Store) (reply *Store, err error) {
	reply = &Store{}
	storeFactory := testing.NewStoreFactory(s.configDir)
	err = storeFactory.UpdateStore(ToNormalStore(in))
	return
}
func (s *server) DeleteStore(ctx context.Context, in *Store) (reply *Store, err error) {
	reply = &Store{}
	storeFactory := testing.NewStoreFactory(s.configDir)
	err = storeFactory.DeleteStore(in.Name)
	return
}
func (s *server) VerifyStore(ctx context.Context, in *SimpleQuery) (reply *CommonResult, err error) {
	// TODO need to implement
	reply = &CommonResult{}
	var loader testing.Writer
	if loader, err = s.getLoaderByStoreName(in.Name); err == nil && loader != nil {
		verifyErr := loader.Verify()
		reply.Success = verifyErr == nil
		reply.Message = util.OKOrErrorMessage(verifyErr)
	}
	return
}

// secret related interfaces
func (s *server) GetSecrets(ctx context.Context, in *Empty) (reply *Secrets, err error) {
	return s.secretServer.GetSecrets(ctx, in)
}
func (s *server) CreateSecret(ctx context.Context, in *Secret) (reply *CommonResult, err error) {
	return s.secretServer.CreateSecret(ctx, in)
}
func (s *server) DeleteSecret(ctx context.Context, in *Secret) (reply *CommonResult, err error) {
	return s.secretServer.DeleteSecret(ctx, in)
}
func (s *server) UpdateSecret(ctx context.Context, in *Secret) (reply *CommonResult, err error) {
	return s.secretServer.UpdateSecret(ctx, in)
}

func (s *server) getLoaderByStoreName(storeName string) (loader testing.Writer, err error) {
	var store *testing.Store
	store, err = testing.NewStoreFactory(s.configDir).GetStore(storeName)
	if err == nil && store != nil {
		loader, err = s.storeWriterFactory.NewInstance(*store)
		if err != nil {
			err = fmt.Errorf("failed to new grpc loader from store %s, err: %v", store.Name, err)
		}
	} else {
		err = fmt.Errorf("failed to get store %s, err: %v", storeName, err)
	}
	return
}

//go:embed data/headers.yaml
var popularHeaders string

func findParentTestCases(testcase *testing.TestCase, suite *testing.TestSuite) (testcases []testing.TestCase) {
	reg, matchErr := regexp.Compile(`(.*?\{\{.*\.\w*.*?\}\})`)
	targetReg, targetErr := regexp.Compile(`\.\w*`)

	expectNames := new(UniqueSlice[string])
	if matchErr == nil && targetErr == nil {
		var expectName string
		for _, val := range testcase.Request.Header {
			if matched := reg.MatchString(val); matched {
				expectName = targetReg.FindString(val)
				expectName = strings.TrimPrefix(expectName, ".")
				expectNames.Push(expectName)
			}
		}

		findExpectNames(testcase.Request.API, expectNames)
		findExpectNames(testcase.Request.Body, expectNames)

		fmt.Println("expect test case names", expectNames.GetAll())
		for _, item := range suite.Items {
			if expectNames.Exist(item.Name) {
				testcases = append(testcases, item)
			}
		}
	}
	return
}

func findExpectNames(target string, expectNames *UniqueSlice[string]) {
	reg, _ := regexp.Compile(`(.*?\{\{.*\.\w*.*?\}\})`)
	targetReg, _ := regexp.Compile(`\.\w*`)

	for _, sub := range reg.FindStringSubmatch(target) {
		// remove {{ and }}
		if left, leftErr := regexp.Compile(`.*\{\{`); leftErr == nil {
			body := left.ReplaceAllString(sub, "")

			expectName := targetReg.FindString(body)
			expectName = strings.TrimPrefix(expectName, ".")
			expectNames.Push(expectName)
		}
	}
}

// UniqueSlice represents an unique slice
type UniqueSlice[T comparable] struct {
	data []T
}

// Push pushes an item if it's not exist
func (s *UniqueSlice[T]) Push(item T) *UniqueSlice[T] {
	if s.data == nil {
		s.data = []T{item}
	} else {
		for _, it := range s.data {
			if it == item {
				return s
			}
		}
		s.data = append(s.data, item)
	}
	return s
}

// Exist checks if the item exist, return true it exists
func (s *UniqueSlice[T]) Exist(item T) bool {
	if s.data != nil {
		for _, it := range s.data {
			if it == item {
				return true
			}
		}
	}
	return false
}

// GetAll returns all the items
func (s *UniqueSlice[T]) GetAll() []T {
	return s.data
}

var errNoTestSuiteFound = errors.New("no test suite found")
