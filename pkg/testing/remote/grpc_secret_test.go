/*
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

package remote

import (
	context "context"
	"errors"
	"testing"

	"github.com/linuxsuren/api-testing/pkg/secret"
	server "github.com/linuxsuren/api-testing/pkg/server"
	"github.com/stretchr/testify/assert"
)

func TestGetSecret(t *testing.T) {
	t.Run("no err", func(t *testing.T) {
		getter := NewGRPCSecretGetter(&fakeSecretGetter{
			secret: &server.Secret{
				Name:  "fake",
				Value: "value",
			},
		})
		result, err := getter.GetSecret("fake")
		assert.Nil(t, err)
		assert.Equal(t, secret.Secret{
			Name:  "fake",
			Value: "value",
		}, result)
	})

	t.Run("have err", func(t *testing.T) {
		getter := NewGRPCSecretGetter(&fakeSecretGetter{
			err: errors.New("fake"),
		})
		secret, err := getter.GetSecret("fake")
		assert.NotNil(t, err)
		assert.NotNil(t, secret)
	})

	secretServer, err := NewGRPCSecretFrom("fake")
	assert.Nil(t, err)
	assert.NotNil(t, secretServer)

	ctx := context.Background()
	_, err = secretServer.GetSecret(ctx, &server.Secret{Name: "fake"})
	assert.NotNil(t, err)

	_, err = secretServer.GetSecrets(ctx, &server.Empty{})
	assert.NotNil(t, err)

	_, err = secretServer.CreateSecret(ctx, &server.Secret{})
	assert.NotNil(t, err)

	_, err = secretServer.DeleteSecret(ctx, &server.Secret{})
	assert.NotNil(t, err)

	_, err = secretServer.UpdateSecret(ctx, &server.Secret{})
	assert.NotNil(t, err)
}

type fakeSecretGetter struct {
	secret *server.Secret
	err    error
}

func (f *fakeSecretGetter) GetSecret(context.Context, *server.Secret) (
	secret *server.Secret, err error) {
	secret = f.secret
	err = f.err
	return
}
