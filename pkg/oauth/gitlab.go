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

import "log"

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
	log.Println("not support")
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
