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
package main

import (
	"io"
	"net/http"
)

func main() {
	{{- if gt (len .Request.Form) 0 }}
	data := url.Values{}
	{{- range $key, $val := .Request.Form}}
	data.Set("{{$key}}", "{{$val}}")
	{{- end}}
	body := strings.NewReader(data.Encode())
	{{- else}}
	body := bytes.NewBufferString("{{.Request.Body.String}}")
	{{- end }}

	req, err := http.NewRequest("{{.Request.Method}}", "{{.Request.API}}", body)
	if err != nil {
		panic(err)
	}

	{{- range $key, $val := .Request.Header}}
	req.Header.Set("{{$key}}", "{{$val}}")
 	{{- end}}

	{{- if gt (len .Request.Cookie) 0 }}
	{{- range $key, $val := .Request.Cookie}}
	req.AddCookie(&http.Cookie{
		Name:  "{{$key}}",
		Value: "{{$val}}",
	})
	{{- end}}
	{{- end}}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != http.StatusOK {
		panic("status code is not 200")
	}

	data, err := io.ReadAll(resp.Body)
	println(string(data))
}
