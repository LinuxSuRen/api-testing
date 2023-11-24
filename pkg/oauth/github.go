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
