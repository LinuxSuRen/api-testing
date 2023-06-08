package pkg_test

import (
	"bytes"
	"net/http"
	"testing"

	_ "embed"

	"github.com/linuxsuren/api-testing/extensions/collector/pkg"
	"github.com/stretchr/testify/assert"
)

func TestSampleExporter(t *testing.T) {
	exporter := pkg.NewSampleExporter(true)
	assert.Equal(t, "sample", exporter.TestSuite.Name)

	request, err := newRequest()
	assert.NoError(t, err)
	exporter.Add(&pkg.RequestAndResponse{Request: request})

	request, err = newRequest()
	exporter.Add(&pkg.RequestAndResponse{
		Request: request,
		Response: &pkg.SimpleResponse{
			Body:       "hello",
			StatusCode: http.StatusOK,
		},
	})

	var result string
	result, err = exporter.Export()
	assert.NoError(t, err)
	assert.Equal(t, sampleSuite, result)
}

func newRequest() (request *http.Request, err error) {
	request, err = http.NewRequest(http.MethodGet, "http://foo/api/v1",
		bytes.NewBuffer([]byte("hello")))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer token")
	return
}

//go:embed testdata/sample_suite.yaml
var sampleSuite string
