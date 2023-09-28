curl -X {{.Request.Method}} '{{.Request.API}}'
{{- range $key, $val := .Request.Header}}
  -H '{{$key}}: {{$val}}'
{{- end}}
{{- if .Request.Body }}
  --data-raw '{{.Request.Body}}'
{{- end}}
