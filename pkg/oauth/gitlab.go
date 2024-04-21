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

type gitlabOAuthProvider struct {
}

// AllScopes returns all the supported scopes
// See also https://docs.gitlab.com/ee/integration/oauth_provider.html
func (p *gitlabOAuthProvider) AllScopes() []string {
	return []string{"api", "read_user", "read_api", "read_repository", "profile", "email"}
}
func (p *gitlabOAuthProvider) MinimalScopes() []string {
	return []string{"read_user", "read_api", "read_repository", "profile"}
}
func (p *gitlabOAuthProvider) GetName() string {
	return "gitlab"
}
func (p *gitlabOAuthProvider) GetServer() string {
	return "https://gitlab.com"
}
func (p *gitlabOAuthProvider) SetServer(_ string) {
	githubLogger.Info("not support")
}
func (p *gitlabOAuthProvider) GetTokenURL() string {
	return "/oauth/token"
}
func (p *gitlabOAuthProvider) GetAuthURL() string {
	return "/oauth/authorize"
}
func (p *gitlabOAuthProvider) GetUserInfoURL() string {
	return "/oauth/userinfo"
}

type privateGitlabOAuthProvider struct {
	*gitlabOAuthProvider
	server string
}

func (p *privateGitlabOAuthProvider) GetName() string {
	return "private-gitlab"
}
func (p *privateGitlabOAuthProvider) GetServer() string {
	return p.server
}
func (p *privateGitlabOAuthProvider) SetServer(server string) {
	p.server = server
}

func init() {
	RegisterOAuthProvider(&gitlabOAuthProvider{})
	RegisterOAuthProvider(&privateGitlabOAuthProvider{})
}
