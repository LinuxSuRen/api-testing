/*
Copyright 2024 API Testing Authors.

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
	"time"

	_ "embed"

	"github.com/h2non/gock"
	atest "github.com/linuxsuren/api-testing/pkg/testing"
	"github.com/linuxsuren/api-testing/sample"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

const (
	BearerToken = "Bearer {{.login.data.access_token}}"
)

func TestRemoteServer(t *testing.T) {
	ctx := context.Background()
	defer gock.Off()

	loader := atest.NewFileWriter("")
	loader.Put("testdata/simple.yaml")
	server := NewRemoteServer(loader, nil, nil, nil, "")
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
	assert.Equal(t, "unknown", ver.Message)
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
	loader := atest.NewFileWriter("")
	loader.Put("testdata/simple.yaml")
	server := NewRemoteServer(loader, nil, nil, nil, "")

	defer gock.Clean()
	gock.New(urlFoo).Get("/").MatchHeader("key", "value").
		Reply(http.StatusOK).
		BodyString(`{"message": "hello"}`)

	result, err := server.RunTestCase(context.TODO(), &TestCaseIdentity{
		Suite:    "simple",
		Testcase: "get",
	})
	assert.NoError(t, err)
	assert.Contains(t, "start to run: 'get'\nstart to send request to http://foo\nstatus code: 200\n", result.Output)
}

func TestFindParentTestCases(t *testing.T) {
	tests := []struct {
		name     string
		testcase *atest.TestCase
		suite    *atest.TestSuite
		expect   []atest.TestCase
	}{{
		name: "normal",
		testcase: &atest.TestCase{
			Request: atest.Request{
				Header: map[string]string{
					"Authorization": BearerToken,
				},
			},
		},
		suite: &atest.TestSuite{
			Items: []atest.TestCase{{
				Name: "login",
			}},
		},
		expect: []atest.TestCase{{
			Name: "login",
		}},
	}, {
		name: "body",
		testcase: &atest.TestCase{
			Request: atest.Request{
				Body: atest.NewRequestBody(`{{.login.data}}`),
			},
		},
		suite: &atest.TestSuite{
			Items: []atest.TestCase{{
				Name: "login",
			}, {
				Name: "user",
				Request: atest.Request{
					Body: atest.NewRequestBody(`{{.login.data}}`),
				},
			}},
		},
		expect: []atest.TestCase{{
			Name: "login",
		}},
	}, {
		name:     "empty cases",
		testcase: &atest.TestCase{},
		suite:    &atest.TestSuite{},
	}, {
		name: "complex",
		testcase: &atest.TestCase{
			Name: "user",
			Request: atest.Request{
				API: "/users/{{(index .login 0).name}}",
			},
		},
		suite: &atest.TestSuite{
			Items: []atest.TestCase{{
				Name: "login",
			}, {
				Name: "user",
				Request: atest.Request{
					API: "/users/{{(index .login 0).name}}",
				},
			}},
		},
		expect: []atest.TestCase{{
			Name: "login",
		}},
	}, {
		name: "nest dep",
		testcase: &atest.TestCase{
			Name: "user",
			Request: atest.Request{
				API: "/users/{{(index .users 0).name}}{{randomKubernetesName}}",
				Header: map[string]string{
					"Authorization": BearerToken,
				},
			},
		},
		suite: &atest.TestSuite{
			Items: []atest.TestCase{{
				Name: "login",
			}, {
				Name: "users",
				Request: atest.Request{
					API: "/users",
				},
			}, {
				Name: "user",
				Request: atest.Request{
					API: "/users/{{(index .users 0).name}}",
					Header: map[string]string{
						"Authorization": BearerToken,
					},
				},
			}},
		},
		expect: []atest.TestCase{{
			Name: "login",
		}, {
			Name: "users",
			Request: atest.Request{
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

		writer := atest.NewFileWriter("")
		err = writer.Put(tmpFile.Name())
		assert.NoError(t, err)

		ctx := context.Background()
		server := NewRemoteServer(writer, nil, nil, nil, "")
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
	writer := atest.NewFileWriter(os.TempDir())
	writer.Put(tmpFile.Name())

	server := NewRemoteServer(writer, nil, nil, nil, "")
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
			assert.Equal(t, atest.TestCase{
				Name: "get",
				Request: atest.Request{
					API: urlFoo,
					Header: map[string]string{
						"key": "value",
					},
					Cookie: map[string]string{},
					Query:  map[string]interface{}{},
					Form:   map[string]string{},
				},
				Expect: atest.Response{
					Header:           map[string]string{},
					BodyFieldsExpect: map[string]interface{}{},
					Verify:           nil,
				},
			}, ToNormalTestCase(result))
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

	t.Run("convertConditionalVerify", func(t *testing.T) {
		assert.Equal(t, []atest.ConditionalVerify{{
			Condition: []string{"1 == 1"},
			Verify:    []string{"1 == 1"},
		}}, convertConditionalVerify([]*ConditionalVerify{{
			Condition: []string{"1 == 1"},
			Verify:    []string{"1 == 1"},
		}}))
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
		assert.Equal(t, 6, len(pairs.Data))
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
		gock.Off()
		name := fmt.Sprintf("fake-%d", time.Now().Second())
		_, err := server.CreateTestSuite(ctx, &TestSuiteIdentity{
			Name: name,
		})
		assert.NoError(t, err)

		_, err = server.UpdateTestSuite(ctx, &TestSuite{
			Name: name,
			Spec: &APISpec{
				Url: urlFoo + "/v2",
			},
		})
		assert.NoError(t, err)

		gock.New(urlFoo).Get("/v2").Reply(500)

		_, err = server.GetSuggestedAPIs(ctx, &TestSuiteIdentity{Name: name})
		assert.NoError(t, err)
	})

	t.Run("normal", func(t *testing.T) {
		gock.Off()
		_, err := server.CreateTestSuite(ctx, &TestSuiteIdentity{
			Name: "fake-1",
		})
		assert.NoError(t, err)

		_, err = server.UpdateTestSuite(ctx, &TestSuite{
			Name: "fake-1",
			Spec: &APISpec{
				Kind: "swagger",
				Url:  urlFoo + "/v1",
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
		assert.Equal(t, 6, len(generators.Data))
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
		assert.Error(t, err, err)
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

	writer := atest.NewFileWriter(dir)
	server = NewRemoteServer(writer, newLocalloaderFromStore(), nil, nil, dir)
	return
}

type fakeLocalLoaderFactory struct {
}

func newLocalloaderFromStore() atest.StoreWriterFactory {
	return &fakeLocalLoaderFactory{}
}

func (l *fakeLocalLoaderFactory) NewInstance(store atest.Store) (writer atest.Writer, err error) {
	writer = atest.NewFileWriter("")
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
