package remote

import (
	"testing"

	atest "github.com/linuxsuren/api-testing/pkg/testing"
	"github.com/stretchr/testify/assert"
)

func TestConvert(t *testing.T) {
	t.Run("convertToNormalTestSuite, empty object", func(t *testing.T) {
		assert.Equal(t, &atest.TestSuite{
			Param: map[string]string{},
		}, convertToNormalTestSuite(&TestSuite{}))
	})

	t.Run("convertToNormalTestSuite, normal object", func(t *testing.T) {
		assert.Equal(t, &atest.TestSuite{
			Param: defaultMap,
			Spec: atest.APISpec{
				Kind: "http",
				URL:  "/v1",
			},
		}, convertToNormalTestSuite(&TestSuite{
			Param: defaultPairs,
			Spec: &APISpec{
				Url:  "/v1",
				Kind: "http",
			},
		}))
	})

	t.Run("convertToGRPCTestSuite, normal object", func(t *testing.T) {
		result := convertToGRPCTestSuite(&atest.TestSuite{
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
		}, convertToNormalTestCase(&TestCase{
			Request: &Request{
				Api:    "/v1",
				Header: defaultPairs,
			},
			Response: &Response{
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
var defaultPairs []*Pair = []*Pair{{Key: "foo", Value: "bar"}}
