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
	"encoding/json"
	"fmt"
	"github.com/linuxsuren/api-testing/pkg/logging"
	"io"
	"net/http"
	"strings"

	"github.com/linuxsuren/api-testing/pkg/util"
)

var (
	userLogger = logging.DefaultLogger(logging.LogLevelInfo).WithName("user")
)

type OAuthProvider interface {
	AllScopes() []string
	MinimalScopes() []string
	GetName() string
	GetServer() string
	SetServer(string)
	GetTokenURL() string
	GetAuthURL() string
	GetUserInfoURL() string
}

var allOAuthProviders = make(map[string]OAuthProvider)

func RegisterOAuthProvider(provier OAuthProvider) {
	name := provier.GetName()
	_, ok := allOAuthProviders[name]
	if !ok {
		allOAuthProviders[name] = provier
	} else {
		panic(fmt.Sprintf("duplicated oauth provider: %q", name))
	}
}

func GetOAuthProvider(name string) OAuthProvider {
	return allOAuthProviders[name]
}

type UserInfo struct {
	Sub               string   `json:"sub"`
	Name              string   `json:"name"`
	PreferredUsername string   `json:"preferred_username"`
	Email             string   `json:"email"`
	Picture           string   `json:"picture"`
	Groups            []string `json:"groups"`
}

func GetUserInfo(server OAuthProvider, token string, skipTlsVerify bool) (userInfo *UserInfo, err error) {
	api := server.GetUserInfoURL()
	if !strings.HasPrefix(api, "http://") && !strings.HasPrefix(api, "https://") {
		api = fmt.Sprintf("%s%s", server.GetServer(), server.GetUserInfoURL())
	}
	req, err := http.NewRequest(http.MethodGet, api, nil)
	if err != nil {
		return
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := util.TlsAwareHTTPClient(skipTlsVerify)
	var resp *http.Response
	if resp, err = client.Do(req); err != nil {
		return
	}
	defer resp.Body.Close()

	userLogger.Info("getting userinfo from", server.GetName())
	if resp.StatusCode == http.StatusOK {
		var data []byte
		if data, err = io.ReadAll(resp.Body); err != nil {
			return
		}

		userInfo = &UserInfo{}
		err = json.Unmarshal(data, userInfo)
	}
	return
}
