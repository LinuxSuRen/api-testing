/*
Copyright 2023 API Testing Authors.

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
package oauth

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestAuthInterceptor(t *testing.T) {
	t.Run("GetVersion", func(t *testing.T) {
		inter := &authInter{}
		_, err := inter.authInterceptor(context.TODO(), nil, &grpc.UnaryServerInfo{
			FullMethod: "/server.Runner/GetVersion",
		}, func(ctx context.Context, req interface{}) (interface{}, error) {
			return nil, nil
		})
		assert.NoError(t, err)
	})

	t.Run("without auth", func(t *testing.T) {
		inter := &authInter{}
		resp, err := inter.authInterceptor(context.TODO(), nil, &grpc.UnaryServerInfo{
			FullMethod: "/server.Runner/GetSuites",
		}, func(ctx context.Context, req interface{}) (interface{}, error) {
			return nil, nil
		})
		assert.Error(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("normal", func(t *testing.T) {
		accessToken["fake"] = &UserInfo{}
		ctx := metadata.NewIncomingContext(context.TODO(), metadata.Pairs("auth", "fake"))

		inter := &authInter{}
		_, err := inter.authInterceptor(ctx, nil, &grpc.UnaryServerInfo{
			FullMethod: "/server.Runner/GetSuites",
		}, func(ctx context.Context, req interface{}) (interface{}, error) {
			return nil, nil
		})
		assert.NoError(t, err)
	})

	assert.NotNil(t, NewAuthInterceptor(nil)) // should have a better way to test it
}
