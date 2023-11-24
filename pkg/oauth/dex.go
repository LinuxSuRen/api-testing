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
