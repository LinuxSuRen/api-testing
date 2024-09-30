/*
Copyright 2024 API Testing Authors.

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

package generator

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/linuxsuren/api-testing/pkg/testing"
)

type nativeImporter struct {
}

type nativeData struct {
	Data string `json:"data"`
}

func NewNativeImporter() Importer {
	return &nativeImporter{}
}

func (p *nativeImporter) Convert(data []byte) (suite *testing.TestSuite, err error) {
	nativeData := nativeData{}
	if err = json.Unmarshal(data, &nativeData); err == nil {
		var data []byte
		if data, err = base64.StdEncoding.DecodeString(nativeData.Data); err == nil {
			suite, err = testing.Parse(data)
		}
	}
	return
}

func (p *nativeImporter) ConvertFromFile(dataFile string) (*testing.TestSuite, error) {
	return convertFromFile(dataFile, p)
}

func (p *nativeImporter) ConvertFromURL(dataURLStr string) (*testing.TestSuite, error) {
	return convertFromURL(dataURLStr, p)
}

func convertFromFile(dataFile string, dataImport DataImporter) (suite *testing.TestSuite, err error) {
	var data []byte
	if data, err = os.ReadFile(dataFile); err == nil {
		suite, err = dataImport.Convert(data)
	}
	return
}

func convertFromURL(dataURLStr string, dataImport DataImporter) (suite *testing.TestSuite, err error) {
	var req *http.Request
	var resp *http.Response
	var dataURL *url.URL

	if dataURL, err = url.Parse(dataURLStr); err == nil {
		req, err = http.NewRequest(http.MethodGet, dataURLStr, nil)
	}

	if err == nil {
		// put all query params as headers
		for k, v := range dataURL.Query() {
			req.Header.Add(k, v[0])
		}

		if resp, err = http.DefaultClient.Do(req); err == nil {
			var data []byte
			if data, err = io.ReadAll(resp.Body); err == nil {
				suite, err = dataImport.Convert(data)
			}
		}
	}
	return
}
