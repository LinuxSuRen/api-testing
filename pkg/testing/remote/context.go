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
	"github.com/linuxsuren/api-testing/pkg/testing"
	"google.golang.org/grpc/metadata"
)

func WithStoreContext(ctx context.Context, store *testing.Store) context.Context {
	return metadata.NewOutgoingContext(ctx, metadata.New(store.ToMap()))
}

func GetStoreFromContext(ctx context.Context) (store *testing.Store) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		store = MDToStore(md)
	}
	return
}

func GetDataFromContext(ctx context.Context) (data map[string]string) {
	data = make(map[string]string)
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		for key, val := range md {
			data[key] = val[0]
		}
	}
	return
}

func WithIncomingStoreContext(ctx context.Context, store *testing.Store) context.Context {
	return metadata.NewIncomingContext(ctx, metadata.New(store.ToMap()))
}

func MDToStore(md metadata.MD) *testing.Store {
	data := make(map[string]string)
	for key, val := range md {
		data[key] = val[0]
	}
	store := testing.MapToStore(data)
	return &store
}
