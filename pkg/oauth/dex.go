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

type dexOAuthProvider struct {
	server string
}

// AllScopes returns all the supported scopes
func (p *dexOAuthProvider) AllScopes() []string {
	return []string{"openid", "email", "groups", "profile", "offline_access"}
}
func (p *dexOAuthProvider) MinimalScopes() []string {
	return p.AllScopes()
}
func (p *dexOAuthProvider) GetName() string {
	return "dex"
}
func (p *dexOAuthProvider) GetServer() string {
	return p.server
}
func (p *dexOAuthProvider) SetServer(server string) {
	p.server = server
}
func (p *dexOAuthProvider) GetTokenURL() string {
	return "/api/dex/token"
}
func (p *dexOAuthProvider) GetAuthURL() string {
	return "/api/dex/auth"
}
func (p *dexOAuthProvider) GetUserInfoURL() string {
	return "/api/dex/userinfo"
}

func init() {
	RegisterOAuthProvider(&dexOAuthProvider{})
}
