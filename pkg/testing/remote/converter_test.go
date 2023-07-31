package remote

import (
	"testing"

	server "github.com/linuxsuren/api-testing/pkg/server"
	atest "github.com/linuxsuren/api-testing/pkg/testing"
	"github.com/stretchr/testify/assert"
)

func TestConvert(t *testing.T) {
	t.Run("convertToNormalTestSuite, empty object", func(t *testing.T) {
		assert.Equal(t, &atest.TestSuite{
			Param: map[string]string{},
		}, ConvertToNormalTestSuite(&TestSuite{}))
	})

	t.Run("convertToNormalTestSuite, normal object", func(t *testing.T) {
		assert.Equal(t, &atest.TestSuite{
			Param: defaultMap,
			Spec: atest.APISpec{
				Kind: "http",
				URL:  "/v1",
			},
		}, ConvertToNormalTestSuite(&TestSuite{
			Param: defaultPairs,
			Spec: &server.APISpec{
				Url:  "/v1",
				Kind: "http",
			},
		}))
	})

	t.Run("convertToGRPCTestSuite, normal object", func(t *testing.T) {
		result := ConvertToGRPCTestSuite(&atest.TestSuite{
			API:   "v1",
			Param: defaultMap,
		})
		assert.Equal(t, "v1", result.Api)
		assert.Equal(t, defaultPairs, result.Param)
	})

	t.Run("convertToNormalTestCase", func(t *testing.T) {
		assert.Equal(t, atest.TestCase{
			Request: atest.Request{
				API:    "/v1",
				Header: defaultMap,
				Query:  map[string]string{},
				Form:   map[string]string{},
			},
			Expect: atest.Response{
				BodyFieldsExpect: defaultInterMap,
				Header:           map[string]string{},
			},
		}, convertToNormalTestCase(&server.TestCase{
			Request: &server.Request{
				Api:    "/v1",
				Header: defaultPairs,
			},
			Response: &server.Response{
				BodyFieldsExpect: defaultPairs,
			},
		}))
	})

	t.Run("convertToGRPCTestCase", func(t *testing.T) {
		result := convertToGRPCTestCase(atest.TestCase{
			Expect: atest.Response{
				BodyFieldsExpect: defaultInterMap,
				Header:           defaultMap,
			},
		})
		if !assert.NotNil(t, result) {
			return
		}
		assert.Equal(t, defaultPairs, result.Response.BodyFieldsExpect)
		assert.Equal(t, defaultPairs, result.Response.Header)
	})
}

var defaultInterMap = map[string]interface{}{"foo": "bar"}
var defaultMap map[string]string = map[string]string{"foo": "bar"}
var defaultPairs []*server.Pair = []*server.Pair{{Key: "foo", Value: "bar"}}
