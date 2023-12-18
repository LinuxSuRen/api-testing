There are {{ .Total }} test cases, failed count {{ .Error }}:
 
{{- if gt .Total 6 }}
{{- if gt .Error 0 }}

| Name | Average | Max | Min | Count | Error |
|---|---|---|---|---|---|
{{- range $val := .Items}}
{{- if gt $val.Error 0 }}
| {{$val.Name}} | {{$val.Average}} | {{$val.Max}} | {{$val.Min}} | {{$val.Count}} | {{$val.Error}} |
{{- end }}
{{- end }}
{{- end }}

<details>
  <summary><b>See all test records</b></summary>

| Name | Average | Max | Min | Count | Error |
|---|---|---|---|---|---|
{{- range $val := .Items}}
| {{$val.Name}} | {{$val.Average}} | {{$val.Max}} | {{$val.Min}} | {{$val.Count}} | {{$val.Error}} |
{{- end }}
</details>
{{- else }}

| Name | Average | Max | Min | Count | Error |
|---|---|---|---|---|---|
{{- range $val := .Items}}
| {{$val.Name}} | {{$val.Average}} | {{$val.Max}} | {{$val.Min}} | {{$val.Count}} | {{$val.Error}} |
{{- end }}
{{- end }}

{{- if gt .LastResourceUsage.Memory 0 }}

Resource usage:
* CPU: {{ .LastResourceUsage.CPU }}
* Memory: {{ .LastResourceUsage.Memory }}
{{- end }}

{{- if .Errors }}

<details>
  <summary><b>See the error message</b></summary>
{{- range $val := .Errors}}
* {{ $val }}
{{- end }}
</details>
{{- end }}

{{- if gt .Converage.Total 0 }}

API Coverage: {{ .Converage.Covered }}/{{ .Converage.Total }}
{{- end }}