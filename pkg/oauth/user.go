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
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
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

	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: skipTlsVerify},
		},
	}

	var resp *http.Response
	if resp, err = client.Do(req); err != nil {
		return
	}
	defer resp.Body.Close()

	log.Println("getting userinfo from", server.GetName())
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
