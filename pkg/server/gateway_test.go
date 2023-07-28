package server

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetadataStoreFunc(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(HeaderKeyStoreName, "test")
	if !assert.NoError(t, err) {
		return
	}

	md := MetadataStoreFunc(context.TODO(), req)
	assert.Equal(t, "test", md.Get(HeaderKeyStoreName)[0])
}
