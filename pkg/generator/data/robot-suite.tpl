*** Settings ***
Library         RequestsLibrary

*** Test Cases ***
{{- range $item := .Items}}
{{$item.Name}}
    {{- if $item.Request.Header}}
    ${headers}=    Create Dictionary        {{- range $key, $val := $item.Request.Header}}   {{$key}}    {{$val}}{{- end}}
    {{- end}}
    ${response}=    {{$item.Request.Method}}  {{$item.Request.API}}{{- if .Request.Header}}   headers=${headers}{{end}}
{{- end}}
