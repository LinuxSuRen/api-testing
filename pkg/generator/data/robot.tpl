*** Settings ***
Library         RequestsLibrary

*** Test Cases ***
{{.Name}}
    {{- if .Request.Header}}
    ${headers}=    Create Dictionary        {{- range $key, $val := .Request.Header}}   {{$key}}    {{$val}}{{- end}}
    {{- end}}
    ${response}=    {{.Request.Method}}  {{.Request.API}}{{- if .Request.Header}}   headers=${headers}{{end}}
