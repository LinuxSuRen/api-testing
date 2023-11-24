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
