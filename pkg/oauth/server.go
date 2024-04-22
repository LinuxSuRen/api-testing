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
	"net/http"

	"context"

	"github.com/linuxsuren/api-testing/pkg/logging"
	"github.com/linuxsuren/api-testing/pkg/util"

	"golang.org/x/oauth2"
)

var (
	accessToken = make(map[string]*UserInfo)
	oauthLogger = logging.DefaultLogger(logging.LogLevelInfo).WithName("oauth")
)

func GetUser(token string) *UserInfo {
	return accessToken[token]
}

type auth struct {
	provider      OAuthProvider
	config        oauth2.Config
	verifier      string
	skipTlsVerify bool
	state         string
}

// NewAuth creates a new auth handler
func NewAuth(provider OAuthProvider, config oauth2.Config, skipTlsVerify bool) *auth {
	config.Scopes = provider.MinimalScopes()
	config.Endpoint.TokenURL = fmt.Sprintf("%s%s", provider.GetServer(), provider.GetTokenURL())
	config.Endpoint.AuthURL = fmt.Sprintf("%s%s", provider.GetServer(), provider.GetAuthURL())
	config.Endpoint.DeviceAuthURL = "https://github.com/login/device/code"
	return &auth{
		provider:      provider,
		config:        config,
		verifier:      oauth2.GenerateVerifier(),
		skipTlsVerify: skipTlsVerify,
		state:         util.String(6),
	}
}

func (a *auth) Callback(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	r.ParseForm()
	state := r.Form.Get("state")
	if state != a.state {
		http.Error(w, "State invalid", http.StatusBadRequest)
		return
	}
	code := r.Form.Get("code")
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}
	oauthLogger.Info("get code", "code", code)

	sslcli := util.TlsAwareHTTPClient(a.skipTlsVerify)
	ctx := context.WithValue(r.Context(), oauth2.HTTPClient, sslcli)

	token, err := a.config.Exchange(ctx, code, oauth2.VerifierOption(a.verifier))
	a.getUserInfo(w, r, token, err)
}

func (a *auth) getUserInfo(w http.ResponseWriter, r *http.Request, token *oauth2.Token, err error) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	accessToken[token.AccessToken] = nil
	// get userInfo, save it to session
	if userInfo, err := GetUserInfo(a.provider, token.AccessToken, a.skipTlsVerify); err == nil {
		accessToken[token.AccessToken] = userInfo
		oauthLogger.Info("has login", "username", userInfo.Name)
	} else {
		oauthLogger.Info("failed to get userinfo", "error", err)
	}

	http.Redirect(w, r, "/?access_token="+token.AccessToken, http.StatusFound)
}

func (a *auth) RequestLocalToken(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	deviceCode := r.FormValue("device_code")
	response, ok := deviceAuthResponseMap[deviceCode]
	if !ok {
		http.Error(w, "device code not found", http.StatusBadRequest)
		return
	}

	token, err := a.config.DeviceAccessToken(r.Context(), response)
	a.getUserInfo(w, r, token, err)
}

var deviceAuthResponseMap = map[string]*oauth2.DeviceAuthResponse{}

func (a *auth) RequestLocalCode(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	response, err := a.config.DeviceAuth(context.Background())
	if err != nil {
		oauthLogger.Info("failed to get device auth", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	deviceAuthResponseMap[response.DeviceCode] = response

	data, _ := json.Marshal(response)
	w.Write(data)
}

func (a *auth) RequestCode(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	ref := r.Header.Get("Referer")
	oauthLogger.Info("callback host", "host", r.Host)

	if ref == "" {
		a.config.RedirectURL = fmt.Sprintf("https://%s/oauth2/callback", r.Host)
	} else {
		a.config.RedirectURL = fmt.Sprintf("%soauth2/callback", ref)
	}

	u := a.config.AuthCodeURL(a.state, oauth2.S256ChallengeOption(a.verifier))
	http.Redirect(w, r, u, http.StatusFound)
}
