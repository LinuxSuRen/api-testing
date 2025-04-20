/*
Copyright 2023-2025 API Testing Authors.

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

package remote

import (
	context "context"
	"testing"

	atest "github.com/linuxsuren/api-testing/pkg/testing"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

func TestWithStoreContext(t *testing.T) {
	ctx := WithStoreContext(context.Background(), sampleStore)
	md, ok := metadata.FromOutgoingContext(ctx)

	assert.True(t, ok)
	assert.Equal(t, sampleStore, MDToStore(md))
}

func TestGetStoreFromContext(t *testing.T) {
	parentCtx := context.Background()
	ctx := WithIncomingStoreContext(parentCtx, sampleStore)

	assert.Equal(t, sampleStore, GetStoreFromContext(ctx))
	assert.Equal(t, make(map[string]string), GetDataFromContext(parentCtx))
	assert.Equal(t, sampleStore.ToMap(), GetDataFromContext(ctx))
}

var sampleStore = &atest.Store{
	Properties: make(map[string]string),
}
