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

package util

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
)

func GetDefaultCachedHTTPClient() *http.Client {
	return &http.Client{
		Transport: defaultCachedClient,
	}
}

var defaultCachedClient = &cachedClient{
	cache:     make(map[string]*http.Response),
	cacheData: make(map[string][]byte),
}

type cachedClient struct {
	cache     map[string]*http.Response
	cacheData map[string][]byte
}

type cachedResponse struct {
	body   []byte
	header http.Header
	status int
}

func (c *cachedClient) RoundTrip(req *http.Request) (*http.Response, error) {
	key := req.URL.String()

	var cachedData *http.Response
	var ok bool
	if cachedData, ok = c.cache[key]; ok {
		cachedData.Body = ioutil.NopCloser(bytes.NewReader(c.cacheData[key]))
		return cachedData, nil
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return resp, err
	}
	var data []byte
	if data, err = io.ReadAll(resp.Body); err == nil {
		resp.Body = ioutil.NopCloser(bytes.NewReader(data))
	}
	c.cache[key] = resp
	c.cacheData[key] = data

	return resp, err
}
