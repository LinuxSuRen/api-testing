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

package pkg

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/linuxsuren/api-testing/pkg/server"
	atest "github.com/linuxsuren/api-testing/pkg/testing"
	"github.com/linuxsuren/api-testing/pkg/testing/remote"
	"github.com/stretchr/testify/assert"
)

func newRemoteServer(t *testing.T) remote.LoaderServer {
	remoteServer := NewRemoteServer(&fakeS3{data: map[*string][]byte{
		aws.String("invalid"): []byte("invalid"),
	}})
	assert.NotNil(t, remoteServer)
	return remoteServer
}

func TestNewRemoteServer(t *testing.T) {
	emptyCtx := context.Background()
	defaultCtx := remote.WithIncomingStoreContext(emptyCtx, &atest.Store{})

	t.Run("ListTestSuite, no required info in context", func(t *testing.T) {
		_, err := newRemoteServer(t).ListTestSuite(emptyCtx, nil)
		assert.Error(t, err)
	})

	t.Run("ListTestSuite", func(t *testing.T) {
		_, err := newRemoteServer(t).ListTestSuite(defaultCtx, nil)
		assert.NoError(t, err)

		var result *server.ExtensionStatus
		result, err = newRemoteServer(t).Verify(defaultCtx, &server.Empty{})
		assert.NoError(t, err)
		assert.True(t, result.Ready)
	})

	t.Run("CreateTestSuite", func(t *testing.T) {
		server := newRemoteServer(t)
		_, err := server.CreateTestSuite(defaultCtx, &remote.TestSuite{
			Name: "fake",
		})
		assert.NoError(t, err)

		var suites *remote.TestSuites
		suites, err = server.ListTestSuite(defaultCtx, nil)
		if assert.NoError(t, err) {
			assert.Equal(t, "fake", suites.Data[0].Name)
		}

		var suite *remote.TestSuite
		suite, err = server.GetTestSuite(defaultCtx, &remote.TestSuite{Name: "fake"})
		if assert.NoError(t, err) {
			assert.Equal(t, "fake", suite.Name)
		}
	})

	t.Run("GetTestSuite", func(t *testing.T) {
		_, err := newRemoteServer(t).GetTestSuite(defaultCtx, &remote.TestSuite{
			Name: "fake",
		})
		assert.NoError(t, err)
	})

	t.Run("UpdateTestSuite", func(t *testing.T) {
		_, err := newRemoteServer(t).UpdateTestSuite(defaultCtx, &remote.TestSuite{
			Name: "fake",
		})
		assert.NoError(t, err)
	})

	t.Run("DeleteTestSuite", func(t *testing.T) {
		server := newRemoteServer(t)
		_, err := server.CreateTestSuite(defaultCtx, &remote.TestSuite{
			Name: "fake",
		})
		assert.NoError(t, err)

		_, err = server.DeleteTestSuite(defaultCtx, &remote.TestSuite{
			Name: "fake",
		})
		assert.NoError(t, err)
	})

	t.Run("ListTestCases", func(t *testing.T) {
		_, err := newRemoteServer(t).ListTestCases(defaultCtx, &remote.TestSuite{
			Name: "fake",
		})
		assert.NoError(t, err)
	})

	t.Run("CreateTestCase", func(t *testing.T) {
		_, err := newRemoteServer(t).CreateTestCase(defaultCtx, &server.TestCase{
			Name:      "fake",
			SuiteName: "fake",
		})
		assert.NoError(t, err)
	})

	t.Run("GetTestCase", func(t *testing.T) {
		_, err := newRemoteServer(t).GetTestCase(defaultCtx, &server.TestCase{
			Name:      "fake",
			SuiteName: "fake",
		})
		assert.NoError(t, err)
	})

	t.Run("UpdateTestCase", func(t *testing.T) {
		_, err := newRemoteServer(t).UpdateTestCase(defaultCtx, &server.TestCase{
			Name:      "fake",
			SuiteName: "fake",
		})
		assert.NoError(t, err)
	})

	t.Run("DeleteTestCase", func(t *testing.T) {
		_, err := newRemoteServer(t).DeleteTestCase(defaultCtx, &server.TestCase{
			Name:      "fake",
			SuiteName: "fake",
		})
		assert.NoError(t, err)
	})
}

func TestCommonFuns(t *testing.T) {
	t.Run("generateKey", func(t *testing.T) {
		assert.Equal(t, "test.yaml", *generateKey("test"))
	})

	t.Run("mapToS3Options", func(t *testing.T) {
		assert.Equal(t, s3Options{
			AccessKeyID:     "id",
			SecretAccessKey: "secret",
			SessionToken:    "token",
			Region:          "region",
			DisableSSL:      true,
			ForcePathStyle:  true,
			Bucket:          "bucket",
		}, mapToS3Options(map[string]string{
			"accesskeyid":     "id",
			"secretaccesskey": "secret",
			"sessiontoken":    "token",
			"region":          "region",
			"disablessl":      "true",
			"forcepathstyle":  "true",
			"bucket":          "bucket",
		}))
	})

	t.Run("removeTestCaseByName, an empty TestSuite", func(t *testing.T) {
		suite := &remote.TestSuite{
			Items: []*server.TestCase{{
				Name: "fake",
			}},
		}

		assert.Equal(t, suite, removeTestCaseByName(suite, "test"))
	})

	t.Run("removeTestCaseByName, a normal TestSuite", func(t *testing.T) {
		suite := &remote.TestSuite{
			Items: []*server.TestCase{{
				Name: "fake",
			}},
		}

		assert.Empty(t, removeTestCaseByName(suite, "fake").Items)
	})

	t.Run("updateTestCase", func(t *testing.T) {
		suite := &remote.TestSuite{
			Items: []*server.TestCase{{
				Name: "fake",
				Request: &server.Request{
					Method: "GET",
				},
			}},
		}

		suite = updateTestCase(suite, &server.TestCase{
			Name: "fake",
			Request: &server.Request{
				Method: "POST",
			},
		})
		assert.Equal(t, "POST", suite.Items[0].Request.Method)
	})

	t.Run("getTestCaseByName", func(t *testing.T) {
		testCase := &server.TestCase{
			Name: "fake",
			Request: &server.Request{
				Api: "http://fake.com",
			},
		}
		sampleTestSuite := &remote.TestSuite{
			Items: []*server.TestCase{testCase},
		}

		testcase := getTestCaseByName(sampleTestSuite, "fake")
		assert.Equal(t, testCase, testcase)
	})
}
