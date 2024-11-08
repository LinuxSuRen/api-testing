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
package util

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
)

func TlsAwareHTTPClient(insecure bool) *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
	}
	return &http.Client{Transport: tr}
}

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

func (c *cachedClient) RoundTrip(req *http.Request) (*http.Response, error) {
	key := req.URL.String()

	var cachedData *http.Response
	var ok bool
	if cachedData, ok = c.cache[key]; ok {
		cachedData.Body = io.NopCloser(bytes.NewReader(c.cacheData[key]))
		return cachedData, nil
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return resp, err
	}
	var data []byte
	if data, err = io.ReadAll(resp.Body); err == nil {
		resp.Body = io.NopCloser(bytes.NewReader(data))
	}
	c.cache[key] = resp
	c.cacheData[key] = data

	return resp, err
}

func WriteAsJSON(obj interface{}, w http.ResponseWriter) (n int, err error) {
	w.Header().Set(ContentType, JSON)
	var data []byte
	if data, err = json.Marshal(obj); err == nil {
		n, err = w.Write(data)
	}
	return
}
