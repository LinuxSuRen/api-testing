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

import "log"

type githubOAuthProvider struct {
}

// AllScopes returns all the supported scopes
// See also https://docs.github.com/en/apps/oauth-apps/building-oauth-apps/scopes-for-oauth-apps
func (p *githubOAuthProvider) AllScopes() []string {
	return []string{"repo", "public_repo", "read:org", "user", "read:user", "user:email"}
}
func (p *githubOAuthProvider) MinimalScopes() []string {
	return []string{"read:user", "public_repo"}
}
func (p *githubOAuthProvider) GetName() string {
	return "github"
}
func (p *githubOAuthProvider) GetServer() string {
	return "https://github.com/login"
}
func (p *githubOAuthProvider) SetServer(_ string) {
	log.Println("not support")
}
func (p *githubOAuthProvider) GetTokenURL() string {
	return "/oauth/access_token"
}
func (p *githubOAuthProvider) GetAuthURL() string {
	return "/oauth/authorize"
}
func (p *githubOAuthProvider) GetUserInfoURL() string {
	return "https://api.github.com/user"
}

func init() {
	RegisterOAuthProvider(&githubOAuthProvider{})
}
