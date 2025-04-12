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
	"net/http"

	"google.golang.org/grpc/metadata"
)

// MetadataStoreFunc is a function that extracts metadata from a request.
func MetadataStoreFunc(ctx context.Context, r *http.Request) (md metadata.MD) {
	store := r.Header.Get(HeaderKeyStoreName)
	md = metadata.Pairs(HeaderKeyStoreName, store)

	if auth := r.Header.Get("X-Auth"); auth != "" {
		md.Set("Auth", auth)
	}
	return
}
