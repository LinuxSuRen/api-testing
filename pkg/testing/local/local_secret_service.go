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
	"path/filepath"

	"github.com/linuxsuren/api-testing/pkg/server"
	"github.com/linuxsuren/api-testing/pkg/testing/remote"
	"gopkg.in/yaml.v2"
)

type localSecretService struct {
	dataDir string
	remote.UnimplementedSecretServiceServer
}

func NewLocalSecretService(dataDir string) remote.SecretServiceServer {
	return &localSecretService{
		dataDir: dataDir,
	}
}

func (s *localSecretService) GetSecret(ctx context.Context, in *server.Secret) (reply *server.Secret, err error) {
	var secrets *server.Secrets
	if secrets, err = s.GetSecrets(ctx, &server.Empty{}); err == nil {
		for _, secret := range secrets.Data {
			if secret.Name == in.Name {
				reply = secret
				break
			}
		}
	}
	return
}

func (s *localSecretService) GetSecrets(ctx context.Context, in *server.Empty) (reply *server.Secrets, err error) {
	secretFile := s.getDataFilePath()
	var data []byte
	if data, err = os.ReadFile(secretFile); err == nil {
		secretData := make(map[string]string)

		if err = yaml.Unmarshal(data, &secretData); err == nil {
			reply = &server.Secrets{
				Data: make([]*server.Secret, 0),
			}

			for key, value := range secretData {
				reply.Data = append(reply.Data, &server.Secret{
					Name:  key,
					Value: value,
				})
			}
		}
	}
	return
}

func (s *localSecretService) CreateSecret(ctx context.Context, in *server.Secret) (reply *server.CommonResult, err error) {
	secretFile := s.getDataFilePath()

	if _, fErr := os.Stat(secretFile); fErr != nil {
		var file *os.File
		if file, err = os.Create(secretFile); err != nil {
			return
		}
		file.Close()
	}

	var data []byte
	if data, err = os.ReadFile(secretFile); err == nil {
		secretData := make(map[string]string)

		if err = yaml.Unmarshal(data, &secretData); err == nil {
			secretData[in.Name] = in.Value

			if data, err = yaml.Marshal(secretData); err == nil {
				err = os.WriteFile(secretFile, data, 0644)
			}
		}
	}
	return
}

func (s *localSecretService) DeleteSecret(ctx context.Context, in *server.Secret) (reply *server.CommonResult, err error) {
	secretFile := s.getDataFilePath()

	var data []byte
	if data, err = os.ReadFile(secretFile); err == nil {
		secretData := make(map[string]string)

		if err = yaml.Unmarshal(data, &secretData); err == nil {
			delete(secretData, in.Name)

			if data, err = yaml.Marshal(secretData); err == nil {
				err = os.WriteFile(secretFile, data, 0644)
			}
		}
	}
	return
}

func (s *localSecretService) UpdateSecret(ctx context.Context, in *server.Secret) (reply *server.CommonResult, err error) {
	reply, err = s.CreateSecret(ctx, in)
	return
}

func (s *localSecretService) getDataFilePath() string {
	return filepath.Join(s.dataDir, "secret.yaml")
}
