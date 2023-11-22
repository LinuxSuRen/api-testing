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
	"fmt"
	"log"
	"net/http"

	"context"
	"crypto/tls"

	"github.com/linuxsuren/api-testing/pkg/util"
	"golang.org/x/oauth2"
)

var accessToken = make(map[string]*UserInfo)

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
	log.Println("get code", code)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: a.skipTlsVerify},
	}
	sslcli := &http.Client{Transport: tr}
	ctx := context.WithValue(r.Context(), oauth2.HTTPClient, sslcli)

	token, err := a.config.Exchange(ctx, code, oauth2.VerifierOption(a.verifier))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	accessToken[token.AccessToken] = nil
	// get userInfo, save it to session
	if userInfo, err := GetUserInfo(a.provider, token.AccessToken, a.skipTlsVerify); err == nil {
		accessToken[token.AccessToken] = userInfo
		log.Println(userInfo.Name, "has login")
	} else {
		log.Println("failed to get userinfo", err)
	}

	http.Redirect(w, r, "/?access_token="+token.AccessToken, http.StatusFound)
}

func (a *auth) RequestCode(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	ref := r.Header.Get("Referer")
	log.Println("callback host", r.Host)

	if ref == "" {
		a.config.RedirectURL = fmt.Sprintf("https://%s/oauth2/callback", r.Host)
	} else {
		a.config.RedirectURL = fmt.Sprintf("%soauth2/callback", ref)
	}

	u := a.config.AuthCodeURL(a.state, oauth2.S256ChallengeOption(a.verifier))
	http.Redirect(w, r, u, http.StatusFound)
}
