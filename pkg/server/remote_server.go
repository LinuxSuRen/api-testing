/**
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

	"log"

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
	storeExtMgr        ExtManager

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
func NewRemoteServer(loader testing.Writer, storeWriterFactory testing.StoreWriterFactory, secretServer SecretServiceServer, storeExtMgr ExtManager, configDir string) RunnerServer {
	if secretServer == nil {
		secretServer = &fakeSecretServer{}
	}

	return &server{
		loader:             loader,
		storeWriterFactory: storeWriterFactory,
		configDir:          configDir,
		secretServer:       secretServer,
		storeExtMgr:        storeExtMgr,
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
			log.Printf("find %d parent cases\n", len(parentCases))
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

	log.Printf("prepare to run: %s, with level: %s\n", suite.Name, task.Level)
	log.Printf("task kind: %s, %d to run\n", task.Kind, len(suite.Items))
	dataContext := map[string]interface{}{}

	if err = suite.Render(dataContext); err != nil {
		reply.Error = err.Error()
		err = nil
		return
	}
	// inject the parameters from input
	if len(task.Parameters) > 0 {
		dataContext[testing.ContextKeyGlobalParam] = pairToMap(task.Parameters)
	}

	buf := new(bytes.Buffer)
	reply = &TestResult{}

	for _, testCase := range suite.Items {
		suiteRunner := runner.GetTestSuiteRunner(suite)
		suiteRunner.WithOutputWriter(buf)
		suiteRunner.WithWriteLevel(task.Level)
		suiteRunner.WithSecure(suite.Spec.Secure)

		// reuse the API prefix
		testCase.Request.RenderAPI(suite.API)

		output, testErr := suiteRunner.RunTestCase(&testCase, dataContext, ctx)
		if getter, ok := suiteRunner.(runner.ResponseRecord); ok {
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
			items.Kind = suite.Spec.Kind
			reply.Data[suite.Name] = items
		}
	}

	return
}

func (s *server) CreateTestSuite(ctx context.Context, in *TestSuiteIdentity) (reply *HelloReply, err error) {
	reply = &HelloReply{}
	loader := s.getLoader(ctx)
	if loader == nil {
		reply.Error = "no loader found"
	} else {
		if err = loader.CreateSuite(in.Name, in.Api); err == nil {
			toUpdate := testing.TestSuite{
				Name: in.Name,
				API:  in.Api,
				Spec: testing.APISpec{
					Kind: in.Kind,
				},
			}

			switch strings.ToLower(in.Kind) {
			case "grpc", "trpc":
				toUpdate.Spec.RPC = &testing.RPCDesc{}
			}

			err = loader.UpdateSuite(toUpdate)
		}
	}
	return
}

func (s *server) ImportTestSuite(ctx context.Context, in *TestSuiteSource) (result *CommonResult, err error) {
	result = &CommonResult{}
	if in.Kind != "postman" && in.Kind != "" {
		result.Success = false
		result.Message = fmt.Sprintf("not support kind: %s", in.Kind)
		return
	}

	var suite *testing.TestSuite
	importer := generator.NewPostmanImporter()
	if in.Url != "" {
		suite, err = importer.ConvertFromURL(in.Url)
	} else if in.Data != "" {
		suite, err = importer.Convert([]byte(in.Data))
	} else {
		err = errors.New("url or data is required")
	}

	if err != nil {
		result.Success = false
		result.Message = err.Error()
		return
	}

	loader := s.getLoader(ctx)

	if err = loader.CreateSuite(suite.Name, suite.API); err != nil {
		return
	}

	for _, item := range suite.Items {
		if err = loader.CreateTestCase(suite.Name, item); err != nil {
			break
		}
	}
	result.Success = true
	return
}

func (s *server) GetTestSuite(ctx context.Context, in *TestSuiteIdentity) (result *TestSuite, err error) {
	loader := s.getLoader(ctx)
	var suite *testing.TestSuite
	if suite, _, err = loader.GetSuite(in.Name); err == nil && suite != nil {
		result = ToGRPCSuite(suite)
	}
	return
}

func (s *server) UpdateTestSuite(ctx context.Context, in *TestSuite) (reply *HelloReply, err error) {
	reply = &HelloReply{}
	loader := s.getLoader(ctx)
	err = loader.UpdateSuite(*ToNormalSuite(in))
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
			result.Items = append(result.Items, ToGRPCTestCase(item))
		}
	}
	return
}

func (s *server) GetTestCase(ctx context.Context, in *TestCaseIdentity) (reply *TestCase, err error) {
	var result testing.TestCase
	loader := s.getLoader(ctx)
	if result, err = loader.GetTestCase(in.Suite, in.Testcase); err == nil {
		reply = ToGRPCTestCase(result)
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
			Kind:       "testcaseInSuite",
			Data:       string(data),
			CaseName:   in.Testcase,
			Level:      "debug",
			Parameters: in.Parameters,
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
		if pair.Key == "" {
			continue
		}
		data[pair.Key] = pair.Value
	}
	return
}

func pairToMap(pairs []*Pair) (data map[string]string) {
	data = make(map[string]string)
	for _, pair := range pairs {
		if pair.Key == "" {
			continue
		}
		data[pair.Key] = pair.Value
	}
	return
}

func convertConditionalVerify(verify []*ConditionalVerify) (result []testing.ConditionalVerify) {
	if verify != nil {
		result = make([]testing.ConditionalVerify, 0)

		for _, item := range verify {
			result = append(result, testing.ConditionalVerify{
				Condition: item.Condition,
				Verify:    item.Verify,
			})
		}
	}
	return
}

func (s *server) CreateTestCase(ctx context.Context, in *TestCaseWithSuite) (reply *HelloReply, err error) {
	reply = &HelloReply{}
	loader := s.getLoader(ctx)
	err = loader.CreateTestCase(in.SuiteName, ToNormalTestCase(in.Data))
	return
}

func (s *server) UpdateTestCase(ctx context.Context, in *TestCaseWithSuite) (reply *HelloReply, err error) {
	reply = &HelloReply{}
	if in.Data == nil {
		err = errors.New("data is required")
		return
	}
	loader := s.getLoader(ctx)
	err = loader.UpdateTestCase(in.SuiteName, ToNormalTestCase(in.Data))
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
		var suite testing.TestSuite

		loader := s.getLoader(ctx)
		if suite, err = loader.GetTestSuite(in.TestSuite, true); err != nil {
			return
		}

		dataContext := map[string]interface{}{}
		if err = suite.Render(dataContext); err != nil {
			return
		}

		if result, err = loader.GetTestCase(in.TestSuite, in.TestCase); err == nil {
			result.Request.RenderAPI(suite.API)

			output, genErr := instance.Generate(&suite, &result)
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

	log.Println("Finding APIs from", in.Name, "with loader", reflect.TypeOf(loader))
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

func (s *server) GetStoreKinds(context.Context, *Empty) (kinds *StoreKinds, err error) {
	storeFactory := testing.NewStoreFactory(s.configDir)
	var stores []testing.StoreKind
	if stores, err = storeFactory.GetStoreKinds(); err == nil {
		kinds = &StoreKinds{}
		for _, store := range stores {
			kinds.Data = append(kinds.Data, &StoreKind{
				Name:    store.Name,
				Enabled: store.Enabled,
				Url:     store.URL,
			})
		}
	}
	return
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
			grpcStore.Ready = sErr == nil && storeStatus.Ready
			grpcStore.ReadOnly = storeStatus.ReadOnly
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
	store := ToNormalStore(in)

	if store.Kind.URL == "" {
		store.Kind.URL = fmt.Sprintf("unix://%s", os.ExpandEnv(fmt.Sprintf("$HOME/.config/atest/%s.sock", store.Kind.Name)))
	}

	if err = storeFactory.CreateStore(store); err == nil && s.storeExtMgr != nil {
		err = s.storeExtMgr.Start(store.Kind.Name, store.Kind.URL)
	}
	return
}
func (s *server) UpdateStore(ctx context.Context, in *Store) (reply *Store, err error) {
	reply = &Store{}
	storeFactory := testing.NewStoreFactory(s.configDir)
	store := ToNormalStore(in)
	if err = storeFactory.UpdateStore(store); err == nil && s.storeExtMgr != nil {
		// TODO need to restart extension if config was changed
		err = s.storeExtMgr.Start(store.Kind.Name, store.Kind.URL)
	}
	return
}
func (s *server) DeleteStore(ctx context.Context, in *Store) (reply *Store, err error) {
	reply = &Store{}
	storeFactory := testing.NewStoreFactory(s.configDir)
	err = storeFactory.DeleteStore(in.Name)
	return
}
func (s *server) VerifyStore(ctx context.Context, in *SimpleQuery) (reply *ExtensionStatus, err error) {
	reply = &ExtensionStatus{}
	var loader testing.Writer
	if loader, err = s.getLoaderByStoreName(in.Name); err == nil && loader != nil {
		readOnly, verifyErr := loader.Verify()
		reply.Ready = verifyErr == nil
		reply.ReadOnly = readOnly
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

		log.Println("expect test case names", expectNames.GetAll())
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
