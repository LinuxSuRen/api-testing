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

	req, err := http.NewRequest("{{.Request.Method}}," "{{.Request.API}}", body)
	if err != nil {
		panic(err)
	}

	{{- range $key, $val := .Request.Header}}
	req.Header.Set("{{$key}}", "{{$val}}")
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
