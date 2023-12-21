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
