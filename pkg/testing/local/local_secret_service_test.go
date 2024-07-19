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

package local

import (
	"context"
	"os"
	"testing"

	"github.com/linuxsuren/api-testing/pkg/server"
	"github.com/stretchr/testify/assert"
)

func TestLocalSecretService(t *testing.T) {
	dataDir, err := os.MkdirTemp(os.TempDir(), "local_secret_service_test")
	assert.NoError(t, err)
	defer os.RemoveAll(dataDir)
	ctx := context.Background()

	service := NewLocalSecretService(dataDir)

	t.Run("create secret", func(t *testing.T) {
		_, err := service.CreateSecret(ctx, &server.Secret{
			Name:  "test",
			Value: "test",
		})
		assert.NoError(t, err)
	})

	t.Run("update secret", func(t *testing.T) {
		_, err := service.UpdateSecret(ctx, &server.Secret{
			Name:  "test",
			Value: "test1",
		})
		assert.NoError(t, err)

		var secret *server.Secret
		secret, err = service.GetSecret(ctx, &server.Secret{
			Name: "test",
		})
		assert.NoError(t, err)
		assert.Equal(t, "test1", secret.Value)

		_, err = service.DeleteSecret(ctx, &server.Secret{
			Name: "test",
		})
		assert.NoError(t, err)

		var secrets *server.Secrets
		secrets, err = service.GetSecrets(ctx, &server.Empty{})
		assert.NoError(t, err)
		assert.Len(t, secrets.Data, 0)
	})
}
