package pkg_test

import (
	"testing"

	"github.com/linuxsuren/api-testing/extensions/store-orm/pkg"
	"github.com/linuxsuren/api-testing/pkg/testing/remote"
	"github.com/stretchr/testify/assert"
)

func TestConvertToRemoteTestCase(t *testing.T) {
	result := pkg.ConvertToRemoteTestCase(&pkg.TestCase{
		Name:   "name",
		API:    "api",
		Method: "get",
		Body:   "body",
		Query:  sampleJSONMap,
		Header: sampleJSONMap,
		Form:   sampleJSONMap,

		ExpectStatusCode: 1,
		ExpectBody:       "expect body",
		ExpectSchema:     "schema",
		ExpectHeader:     sampleJSONMap,
		ExpectBodyFields: sampleJSONMap,
		ExpectVerify:     `["one"]`,
	})
	assert.Equal(t, &remote.TestCase{
		Name: "name",
		Request: &remote.Request{
			Api:    "api",
			Method: "get",
			Body:   "body",
			Query:  samplePairs,
			Header: samplePairs,
			Form:   samplePairs,
		},
		Response: &remote.Response{
			StatusCode:       1,
			BodyFieldsExpect: samplePairs,
			Verify:           []string{"one"},
			Schema:           "schema",
			Body:             "expect body",
			Header:           samplePairs,
		},
	}, result)
}

func TestConverToDBTestCase(t *testing.T) {
	t.Run("without request and response", func(t *testing.T) {
		result := pkg.ConverToDBTestCase(&remote.TestCase{})
		assert.Equal(t, pkg.TestCase{}, result)
	})

	t.Run("only have request", func(t *testing.T) {
		result := pkg.ConverToDBTestCase(&remote.TestCase{
			Request: &remote.Request{
				Api:    "api",
				Method: "get",
				Body:   "body",
				Header: samplePairs,
				Query:  samplePairs,
				Form:   samplePairs,
			},
		})
		assert.Equal(t, pkg.TestCase{
			API:    "api",
			Method: "get",
			Body:   "body",
			Header: sampleJSONMap,
			Query:  sampleJSONMap,
			Form:   sampleJSONMap,
		}, result)
	})

	t.Run("only have response", func(t *testing.T) {
		result := pkg.ConverToDBTestCase(&remote.TestCase{
			Response: &remote.Response{
				StatusCode:       1,
				Body:             "body",
				Schema:           "schema",
				Header:           samplePairs,
				BodyFieldsExpect: samplePairs,
				Verify:           []string{"one"},
			},
		})
		assert.Equal(t, pkg.TestCase{
			ExpectBody:       "body",
			ExpectStatusCode: 1,
			ExpectSchema:     "schema",
			ExpectVerify:     `["one"]`,
			ExpectHeader:     sampleJSONMap,
			ExpectBodyFields: sampleJSONMap,
		}, result)
	})
}

const sampleJSONMap = `{"key":"value"}`

var samplePairs []*remote.Pair = []*remote.Pair{{
	Key:   "key",
	Value: "value",
}}
