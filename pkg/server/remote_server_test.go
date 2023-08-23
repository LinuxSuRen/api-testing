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

package server

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"sync"
	"sync/atomic"
	"testing"

	_ "embed"

	"github.com/h2non/gock"
	atesting "github.com/linuxsuren/api-testing/pkg/testing"
	"github.com/linuxsuren/api-testing/sample"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

const (
	BearerToken = "Bearer {{.login.data.access_token}}"
)

func TestRemoteServer(t *testing.T) {
	ctx := context.Background()

	loader := atesting.NewFileWriter("")
	loader.Put("testdata/simple.yaml")
	server := NewRemoteServer(loader, nil, nil, "")
	_, err := server.Run(ctx, &TestTask{
		Kind: "fake",
	})
	assert.NotNil(t, err)

	gock.New(urlFoo).Get("/").Reply(http.StatusOK).JSON(&server)
	gock.New(urlFoo).Get("/").Reply(http.StatusOK).JSON(&server)
	_, err = server.Run(ctx, &TestTask{
		Kind: "suite",
		Data: simpleSuite,
	})
	assert.Nil(t, err)

	gock.New(urlFoo).Get("/").Reply(http.StatusOK).JSON(&server)
	_, err = server.Run(ctx, &TestTask{
		Kind: "testcase",
		Data: simpleTestCase,
	})
	assert.Nil(t, err)

	gock.New(urlFoo).Get("/").Reply(http.StatusOK).JSON(&server)
	_, err = server.Run(ctx, &TestTask{
		Kind:     "testcaseInSuite",
		Data:     simpleSuite,
		CaseName: "get",
	})
	assert.Nil(t, err)

	gock.New(urlFoo).Get("/").Reply(http.StatusOK).JSON(&server)
	_, err = server.Run(ctx, &TestTask{
		Kind:     "testcaseInSuite",
		Data:     simpleSuite,
		CaseName: "fake",
		Env: map[string]string{
			"SERVER": "http://localhost:9090",
		},
	})
	assert.NotNil(t, err)

	var ver *HelloReply
	ver, err = server.GetVersion(ctx, &Empty{})
	assert.Empty(t, ver.Message)
	assert.Nil(t, err)

	ver, err = server.Sample(ctx, &Empty{})
	assert.Nil(t, err)
	assert.Equal(t, sample.TestSuiteGitLab, ver.Message)

	var suites *Suites
	suites, err = server.GetSuites(ctx, &Empty{})
	assert.Nil(t, err)
	assert.Equal(t, suites, &Suites{Data: map[string]*Items{
		"simple": {
			Data: []string{"get", "query"},
		},
	}})

	var testCase *TestCase
	testCase, err = server.GetTestCase(ctx, &TestCaseIdentity{
		Suite:    "simple",
		Testcase: "get",
	})
	assert.Nil(t, err)
	assert.Equal(t, "get", testCase.Name)
	assert.Equal(t, urlFoo, testCase.Request.Api)

	// secret functions
	_, err = server.GetSecrets(ctx, &Empty{})
	assert.Error(t, err)

	_, err = server.CreateSecret(ctx, &Secret{})
	assert.Error(t, err)

	_, err = server.DeleteSecret(ctx, &Secret{})
	assert.Error(t, err)

	_, err = server.UpdateSecret(ctx, &Secret{})
	assert.Error(t, err)
}

func TestRunTestCase(t *testing.T) {
	loader := atesting.NewFileWriter("")
	loader.Put("testdata/simple.yaml")
	server := NewRemoteServer(loader, nil, nil, "")

	defer gock.Clean()
	gock.New(urlFoo).Get("/").MatchHeader("key", "value").
		Reply(http.StatusOK).
		BodyString(`{"message": "hello"}`)

	result, err := server.RunTestCase(context.TODO(), &TestCaseIdentity{
		Suite:    "simple",
		Testcase: "get",
	})
	assert.NoError(t, err)
	assert.Contains(t, result.Output, "start to run: 'get'\nstart to send request to http://foo\nresponse body:")
}

func TestFindParentTestCases(t *testing.T) {
	tests := []struct {
		name     string
		testcase *atesting.TestCase
		suite    *atesting.TestSuite
		expect   []atesting.TestCase
	}{{
		name: "normal",
		testcase: &atesting.TestCase{
			Request: atesting.Request{
				Header: map[string]string{
					"Authorization": BearerToken,
				},
			},
		},
		suite: &atesting.TestSuite{
			Items: []atesting.TestCase{{
				Name: "login",
			}},
		},
		expect: []atesting.TestCase{{
			Name: "login",
		}},
	}, {
		name: "body",
		testcase: &atesting.TestCase{
			Request: atesting.Request{
				Body: `{{.login.data}}`,
			},
		},
		suite: &atesting.TestSuite{
			Items: []atesting.TestCase{{
				Name: "login",
			}, {
				Name: "user",
				Request: atesting.Request{
					Body: `{{.login.data}}`,
				},
			}},
		},
		expect: []atesting.TestCase{{
			Name: "login",
		}},
	}, {
		name:     "empty cases",
		testcase: &atesting.TestCase{},
		suite:    &atesting.TestSuite{},
	}, {
		name: "complex",
		testcase: &atesting.TestCase{
			Name: "user",
			Request: atesting.Request{
				API: "/users/{{(index .login 0).name}}",
			},
		},
		suite: &atesting.TestSuite{
			Items: []atesting.TestCase{{
				Name: "login",
			}, {
				Name: "user",
				Request: atesting.Request{
					API: "/users/{{(index .login 0).name}}",
				},
			}},
		},
		expect: []atesting.TestCase{{
			Name: "login",
		}},
	}, {
		name: "nest dep",
		testcase: &atesting.TestCase{
			Name: "user",
			Request: atesting.Request{
				API: "/users/{{(index .users 0).name}}{{randomKubernetesName}}",
				Header: map[string]string{
					"Authorization": BearerToken,
				},
			},
		},
		suite: &atesting.TestSuite{
			Items: []atesting.TestCase{{
				Name: "login",
			}, {
				Name: "users",
				Request: atesting.Request{
					API: "/users",
				},
			}, {
				Name: "user",
				Request: atesting.Request{
					API: "/users/{{(index .users 0).name}}",
					Header: map[string]string{
						"Authorization": BearerToken,
					},
				},
			}},
		},
		expect: []atesting.TestCase{{
			Name: "login",
		}, {
			Name: "users",
			Request: atesting.Request{
				API: "/users",
			},
		}},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := findParentTestCases(tt.testcase, tt.suite)
			assert.Equal(t, tt.expect, result)
		})
	}
}

func TestUniqueSlice(t *testing.T) {
	uniqueSlice := new(UniqueSlice[string])
	uniqueSlice.Push("a").Push("a").Push("b")
	assert.Equal(t, []string{"a", "b"}, uniqueSlice.GetAll())
}

func TestWithDefaultValue(t *testing.T) {
	assert.Equal(t, withDefaultValue("a", "b"), "a")
	assert.Equal(t, withDefaultValue("", "b"), "b")
	assert.Equal(t, withDefaultValue(nil, map[string]string{"key": "val"}), map[string]string{"key": "val"})
	assert.Equal(t, withDefaultValue(map[string]string{"key": "val"}, map[string]string{"key": "value"}), map[string]string{"key": "val"})
}

func TestMapInterToPair(t *testing.T) {
	assert.Equal(t, []*Pair{{Key: "key", Value: "val"}},
		mapInterToPair(map[string]interface{}{"key": "val"}))
}

func TestUpdateTestCase(t *testing.T) {
	t.Run("no suite found", func(t *testing.T) {
		server, clean := getRemoteServerInTempDir()
		defer clean()
		server.UpdateTestCase(context.TODO(), &TestCaseWithSuite{
			Data: &TestCase{},
		})
	})

	t.Run("no data", func(t *testing.T) {
		server, clean := getRemoteServerInTempDir()
		defer clean()
		_, err := server.UpdateTestCase(context.TODO(), &TestCaseWithSuite{})
		assert.Error(t, err)
	})

	t.Run("normal", func(t *testing.T) {
		tmpFile, err := os.CreateTemp(os.TempDir(), "test")
		assert.NoError(t, err)
		defer os.Remove(tmpFile.Name())

		fmt.Fprint(tmpFile, simpleSuite)

		writer := atesting.NewFileWriter("")
		err = writer.Put(tmpFile.Name())
		assert.NoError(t, err)

		ctx := context.Background()
		server := NewRemoteServer(writer, nil, nil, "")
		_, err = server.UpdateTestCase(ctx, &TestCaseWithSuite{
			SuiteName: "simple",
			Data: &TestCase{
				Name: "get",
				Request: &Request{
					Api: "http://foo.json",
				},
				Response: &Response{
					StatusCode: 200,
					BodyFieldsExpect: []*Pair{{
						Key:   "key",
						Value: "value",
					}},
				},
			},
		})
		assert.NoError(t, err)

		_, err = server.UpdateTestCase(ctx, &TestCaseWithSuite{
			SuiteName: "simple",
			Data: &TestCase{
				Name: "post",
				Request: &Request{
					Method: http.MethodPost,
					Api:    urlFoo,
					Header: []*Pair{{
						Key:   "key",
						Value: "value",
					}},
				},
			},
		})
		assert.NoError(t, err)

		var testCase *TestCase
		testCase, err = server.GetTestCase(ctx, &TestCaseIdentity{
			Suite: "simple", Testcase: "get",
		})
		assert.NoError(t, err)
		if assert.NotNil(t, testCase) {
			assert.Equal(t, "http://foo.json", testCase.Request.Api)
			assert.Equal(t, int32(200), testCase.Response.StatusCode)
		}

		_, err = server.CreateTestSuite(ctx, &TestSuiteIdentity{Name: "fake"})
		assert.NoError(t, err)

		var suites *Suites
		suites, err = server.GetSuites(ctx, &Empty{})
		if assert.NoError(t, err) {
			assert.Equal(t, 2, len(suites.Data))
		}

		_, err = server.DeleteTestCase(ctx, &TestCaseIdentity{Suite: "simple", Testcase: "get"})
		assert.NoError(t, err)

		testCase, err = server.GetTestCase(ctx, &TestCaseIdentity{Suite: "simple", Testcase: "get"})
		assert.Nil(t, testCase)
		assert.Error(t, err)
	})
}

func TestListTestCase(t *testing.T) {
	tmpFile, err := os.CreateTemp(os.TempDir(), "test")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	fmt.Fprint(tmpFile, simpleSuite)
	writer := atesting.NewFileWriter(os.TempDir())
	writer.Put(tmpFile.Name())

	server := NewRemoteServer(writer, nil, nil, "")
	ctx := context.Background()

	t.Run("get two testcases", func(t *testing.T) {
		suite, err := server.ListTestCase(ctx, &TestSuiteIdentity{Name: "simple"})
		assert.NoError(t, err)
		if assert.NotNil(t, suite) {
			assert.Equal(t, 2, len(suite.Items))
		}
	})

	t.Run("get one testcase", func(t *testing.T) {
		result, err := server.GetTestCase(ctx, &TestCaseIdentity{Suite: "simple", Testcase: "get"})
		assert.NoError(t, err)
		if assert.NotNil(t, result) {
			assert.Equal(t, atesting.TestCase{
				Name: "get",
				Request: atesting.Request{
					API: urlFoo,
					Header: map[string]string{
						"key": "value",
					},
					Query: map[string]string{},
					Form:  map[string]string{},
				},
				Expect: atesting.Response{
					Header:           map[string]string{},
					BodyFieldsExpect: map[string]interface{}{},
					Verify:           nil,
				},
			}, convertToTestingTestCase(result))
		}
	})

	t.Run("create testcase", func(t *testing.T) {
		reply, err := server.CreateTestCase(ctx, &TestCaseWithSuite{
			SuiteName: "simple",
			Data: &TestCase{
				Name: "put",
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, reply)
	})
}

func TestRemoteServerSuite(t *testing.T) {
	t.Run("Get suite not found", func(t *testing.T) {
		ctx := context.Background()
		server, clean := getRemoteServerInTempDir()
		defer clean()

		suite, err := server.GetTestSuite(ctx, &TestSuiteIdentity{Name: "fake"})
		assert.NoError(t, err)
		assert.Nil(t, suite)
	})

	t.Run("Get existing suite", func(t *testing.T) {
		ctx := context.Background()
		server, clean := getRemoteServerInTempDir()
		defer clean()

		// create a new suite
		_, err := server.CreateTestSuite(ctx, &TestSuiteIdentity{Name: "fake"})
		assert.NoError(t, err)

		suite, err := server.GetTestSuite(ctx, &TestSuiteIdentity{Name: "fake"})
		assert.NoError(t, err)
		if assert.NotNil(t, suite) {
			assert.Equal(t, "fake", suite.Name)
		}

		_, err = server.UpdateTestSuite(ctx, &TestSuite{Name: "fake", Api: "http://foo"})
		assert.NoError(t, err)

		// check if the value was updated successfully
		suite, err = server.GetTestSuite(ctx, &TestSuiteIdentity{Name: "fake"})
		assert.NoError(t, err)
		if assert.NotNil(t, suite) {
			assert.Equal(t, "http://foo", suite.Api)
		}
	})

	t.Run("Delete non-exist suite", func(t *testing.T) {
		ctx := context.Background()
		server, clean := getRemoteServerInTempDir()
		defer clean()

		_, err := server.DeleteTestSuite(ctx, &TestSuiteIdentity{Name: "fake"})
		assert.Error(t, err)
	})
}

func TestPopularHeaders(t *testing.T) {
	ctx := context.Background()
	server, clean := getRemoteServerInTempDir()
	defer clean()

	pairs, err := server.PopularHeaders(ctx, &Empty{})
	if assert.NoError(t, err) {
		assert.Equal(t, 5, len(pairs.Data))
	}
}

func TestGetSuggestedAPIs(t *testing.T) {
	ctx := context.Background()
	server, clean := getRemoteServerInTempDir()
	defer clean()

	t.Run("not found TestSuite", func(t *testing.T) {
		reply, err := server.GetSuggestedAPIs(ctx, &TestSuiteIdentity{Name: "fake"})
		if assert.NoError(t, err) {
			assert.Equal(t, 0, len(reply.Data))
		}
	})

	t.Run("no swagger URL", func(t *testing.T) {
		_, err := server.CreateTestSuite(ctx, &TestSuiteIdentity{Name: "fake"})
		assert.NoError(t, err)

		reply, err := server.GetSuggestedAPIs(ctx, &TestSuiteIdentity{Name: "fake"})
		if assert.NoError(t, err) {
			assert.Equal(t, 0, len(reply.Data))
		}
	})

	t.Run("with swagger URL, not accessed", func(t *testing.T) {
		_, err := server.CreateTestSuite(ctx, &TestSuiteIdentity{
			Name: "fake",
		})

		_, err = server.UpdateTestSuite(ctx, &TestSuite{
			Name: "fake",
			Spec: &APISpec{
				Url: urlFoo + "/v2",
			},
		})
		assert.NoError(t, err)

		gock.New(urlFoo).Get("/v2").Reply(500)

		_, err = server.GetSuggestedAPIs(ctx, &TestSuiteIdentity{Name: "fake"})
		assert.NoError(t, err)
	})

	t.Run("normal", func(t *testing.T) {
		_, err := server.CreateTestSuite(ctx, &TestSuiteIdentity{
			Name: "fake-1",
		})
		assert.NoError(t, err)

		_, err = server.UpdateTestSuite(ctx, &TestSuite{
			Name: "fake-1",
			Spec: &APISpec{
				Url: urlFoo + "/v1",
			},
		})
		assert.NoError(t, err)

		gock.New(urlFoo).Get("/v1").Reply(200).File("testdata/swagger.json")

		var testcases *TestCases
		testcases, err = server.GetSuggestedAPIs(ctx, &TestSuiteIdentity{Name: "fake-1"})
		assert.NoError(t, err)
		if assert.NotNil(t, testcases) {
			assert.Equal(t, 5, len(testcases.Data))
		}
	})
}

func TestFunctionsQuery(t *testing.T) {
	ctx := context.Background()
	server, clean := getRemoteServerInTempDir()
	defer clean()

	t.Run("match exactly", func(t *testing.T) {
		reply, err := server.FunctionsQuery(ctx, &SimpleQuery{Name: "randNumeric"})
		if assert.NoError(t, err) {
			assert.Equal(t, 1, len(reply.Data))
			assert.Equal(t, "randNumeric", reply.Data[0].Key)
			assert.Equal(t, "func(int) string", reply.Data[0].Value)
		}
	})

	t.Run("ignore letter case", func(t *testing.T) {
		reply, err := server.FunctionsQuery(ctx, &SimpleQuery{Name: "randnumer"})
		if assert.NoError(t, err) {
			assert.Equal(t, 1, len(reply.Data))
		}
	})
}

func TestCodeGenerator(t *testing.T) {
	ctx := context.Background()
	server, clean := getRemoteServerInTempDir()
	defer clean()

	t.Run("ListCodeGenerator", func(t *testing.T) {
		generators, err := server.ListCodeGenerator(ctx, &Empty{})
		assert.NoError(t, err)
		assert.Equal(t, 1, len(generators.Data))
	})

	t.Run("GenerateCode, no generator found", func(t *testing.T) {
		result, err := server.GenerateCode(ctx, &CodeGenerateRequest{
			Generator: "fake",
		})
		assert.NoError(t, err)
		assert.False(t, result.Success)
	})

	t.Run("GenerateCode, no TestCase found", func(t *testing.T) {
		result, err := server.GenerateCode(ctx, &CodeGenerateRequest{
			Generator: "golang",
			TestSuite: "fake",
			TestCase:  "fake",
		})
		assert.Error(t, err)
		assert.NotNil(t, result)
	})

	t.Run("GenerateCode, normal", func(t *testing.T) {
		// create a new suite
		_, err := server.CreateTestSuite(ctx, &TestSuiteIdentity{Name: "fake"})
		assert.NoError(t, err)

		_, err = server.CreateTestCase(ctx, &TestCaseWithSuite{
			SuiteName: "fake",
			Data: &TestCase{
				Name: "fake",
			},
		})
		assert.NoError(t, err)

		result, err := server.GenerateCode(ctx, &CodeGenerateRequest{
			Generator: "golang",
			TestSuite: "fake",
			TestCase:  "fake",
		})
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.Success)
		assert.NotEmpty(t, result.Message)
	})

	t.Run("ListConverter", func(t *testing.T) {
		list, err := server.ListConverter(ctx, &Empty{})
		assert.NoError(t, err)
		assert.Equal(t, 2, len(list.Data))
	})

	t.Run("ConvertTestSuite no converter given", func(t *testing.T) {
		reply, err := server.ConvertTestSuite(ctx, &CodeGenerateRequest{})
		assert.NoError(t, err)
		if assert.NotNil(t, reply) {
			assert.False(t, reply.Success)
		}
	})

	t.Run("ConvertTestSuite no suite found", func(t *testing.T) {
		reply, err := server.ConvertTestSuite(ctx, &CodeGenerateRequest{Generator: "jmeter"})
		assert.NoError(t, err)
		if assert.NotNil(t, reply) {
			assert.True(t, reply.Success)
		}
	})

	t.Run("ImportTestSuite, url or data is required", func(t *testing.T) {
		result, err := server.ImportTestSuite(ctx, &TestSuiteSource{})
		assert.Error(t, err)
		assert.False(t, result.Success)
		assert.Equal(t, "url or data is required", result.Message)
	})

	t.Run("ImportTestSuite, invalid kind", func(t *testing.T) {
		result, err := server.ImportTestSuite(ctx, &TestSuiteSource{Kind: "fake"})
		assert.NoError(t, err)
		assert.False(t, result.Success)
		assert.Equal(t, "not support kind: fake", result.Message)
	})

	t.Run("ImportTestSuite, import from string", func(t *testing.T) {
		result, err := server.ImportTestSuite(ctx, &TestSuiteSource{
			Kind: "postman",
			Data: simplePostman,
		})
		assert.NoError(t, err)
		assert.True(t, result.Success)
	})

	t.Run("ImportTestSuite, import from URL", func(t *testing.T) {
		defer gock.Off()
		gock.New(urlFoo).Get("/").Reply(http.StatusOK).BodyString(simplePostman)

		// already exist
		result, err := server.ImportTestSuite(ctx, &TestSuiteSource{
			Kind: "postman",
			Url:  urlFoo,
		})
		assert.Error(t, err)
		assert.False(t, result.Success)
	})
}

func TestFunctionsQueryStream(t *testing.T) {
	ctx := context.Background()
	server, clean := getRemoteServerInTempDir()
	defer clean()

	fakess := &fakeServerStream{
		p:      new(uint32),
		lock:   sync.Mutex{},
		Ctx:    ctx,
		Inputs: []any{&SimpleQuery{Name: "randNumeric"}, &SimpleQuery{Name: "randnumer"}},
		Outpus: []any{},
	}
	err := server.FunctionsQueryStream(&runnerFunctionsQueryStreamServer{fakess})
	t.Run("match outputs length", func(t *testing.T) {
		if assert.NoError(t, err) {
			assert.Equal(t, 2, len(fakess.Outpus))
		}
	})
	t.Run("match exactly", func(t *testing.T) {
		if assert.NoError(t, err) {
			reply := fakess.Outpus[0]
			assert.IsType(t, &Pairs{}, reply)
			assert.Equal(t, 1, len(reply.(*Pairs).Data))
			assert.Equal(t, "randNumeric", reply.(*Pairs).Data[0].Key)
			assert.Equal(t, "func(int) string", reply.(*Pairs).Data[0].Value)
		}
	})
	t.Run("ignore letter case", func(t *testing.T) {
		if assert.NoError(t, err) {
			reply := fakess.Outpus[1]
			assert.IsType(t, &Pairs{}, reply)
			assert.Equal(t, 1, len(reply.(*Pairs).Data))
		}
	})
}

func TestStoreManager(t *testing.T) {
	ctx := context.Background()

	// always have a local store
	t.Run("GetStores, no external stores", func(t *testing.T) {
		server, clean := getRemoteServerInTempDir()
		defer clean()

		reply, err := server.GetStores(ctx, &Empty{})
		assert.NoError(t, err)
		if assert.Equal(t, 1, len(reply.Data)) {
			assert.Equal(t, "local", reply.Data[0].Name)
		}
	})

	t.Run("CreateStore", func(t *testing.T) {
		server, clean := getRemoteServerInTempDir()
		defer clean()
		reply, err := server.CreateStore(ctx, &Store{
			Name: "fake",
		})
		assert.NoError(t, err)
		assert.NotNil(t, reply)

		var stores *Stores
		stores, err = server.GetStores(ctx, &Empty{})
		assert.NoError(t, err)
		assert.Equal(t, 2, len(stores.Data))
	})

	t.Run("DeleteStore", func(t *testing.T) {
		server, clean := getRemoteServerInTempDir()
		defer clean()
		reply, err := server.DeleteStore(ctx, &Store{})
		assert.NoError(t, err)
		assert.NotNil(t, reply)
	})

	t.Run("VerifyStore", func(t *testing.T) {
		server, clean := getRemoteServerInTempDir()
		defer clean()

		reply, err := server.VerifyStore(ctx, &SimpleQuery{})
		assert.Error(t, err)
		assert.NotNil(t, reply)
	})

	t.Run("UpdateStore", func(t *testing.T) {
		server, clean := getRemoteServerInTempDir()
		defer clean()

		_, err := server.UpdateStore(ctx, &Store{})
		assert.Error(t, err)
	})
}

func TestFakeSecretServer(t *testing.T) {
	fakeSecret := &fakeSecretServer{}
	ctx := context.Background()

	_, err := fakeSecret.GetSecrets(ctx, &Empty{})
	assert.Error(t, err)

	_, err = fakeSecret.CreateSecret(ctx, &Secret{})
	assert.Error(t, err)

	_, err = fakeSecret.DeleteSecret(ctx, &Secret{})
	assert.Error(t, err)

	_, err = fakeSecret.UpdateSecret(ctx, &Secret{})
	assert.Error(t, err)
}

func getRemoteServerInTempDir() (server RunnerServer, call func()) {
	dir, _ := os.MkdirTemp(os.TempDir(), "remote-server-test")
	call = func() { os.RemoveAll(dir) }

	writer := atesting.NewFileWriter(dir)
	server = NewRemoteServer(writer, newLocalloaderFromStore(), nil, dir)
	return
}

type fakeLocalLoaderFactory struct {
}

func newLocalloaderFromStore() atesting.StoreWriterFactory {
	return &fakeLocalLoaderFactory{}
}

func (l *fakeLocalLoaderFactory) NewInstance(store atesting.Store) (writer atesting.Writer, err error) {
	writer = atesting.NewFileWriter("")
	return
}

//go:embed testdata/simple.yaml
var simpleSuite string

//go:embed testdata/simple_testcase.yaml
var simpleTestCase string

//go:embed testdata/postman.json
var simplePostman string

const urlFoo = "http://foo"

type fakeServerStream struct {
	p      *uint32
	lock   sync.Mutex
	Ctx    context.Context
	Inputs []any
	Outpus []any
}

func (s *fakeServerStream) SetInputs(in []any) {
	s.Inputs = in
	s.Outpus = make([]any, len(in))
}

func (s *fakeServerStream) SetHeader(metadata.MD) error { return nil }

func (s *fakeServerStream) SendHeader(metadata.MD) error { return nil }

func (s *fakeServerStream) SetTrailer(metadata.MD) {}

func (s *fakeServerStream) Context() context.Context { return s.Ctx }

func (s *fakeServerStream) SendMsg(m interface{}) error {
	s.lock.Lock()
	s.Outpus = append(s.Outpus, m)
	s.lock.Unlock()
	return nil
}

func (s *fakeServerStream) RecvMsg(m interface{}) error {
	defer atomic.AddUint32(s.p, 1)
	index := atomic.LoadUint32(s.p)
	if index == uint32(len(s.Inputs)) {
		return io.EOF
	}

	mv := reflect.ValueOf(m)
	if mv.Kind() == reflect.Pointer && !mv.IsNil() {
		mv = mv.Elem()
	} else {
		return fmt.Errorf("cannot receive to %v", m)
	}

	iv := reflect.ValueOf(s.Inputs[index])
	if iv.Kind() == reflect.Pointer && !iv.IsNil() {
		iv = iv.Elem()
	} else {
		return fmt.Errorf("invalid fake input at index %v", index)
	}

	if mv.CanSet() && mv.Type() == iv.Type() {
		mv.Set(iv)
	} else {
		return fmt.Errorf("cannot set fake input to %v at index %v", m, index)
	}

	return nil
}
