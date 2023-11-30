There are {{ .Total }} test cases:
 
{{- if gt .Total 6 }}
{{- if gt .Error 0 }}

| API | Average | Max | Min | Count | Error |
|---|---|---|---|---|---|
{{- range $val := .Items}}
{{- if gt $val.Error 0 }}
| {{$val.API}} | {{$val.Average}} | {{$val.Max}} | {{$val.Min}} | {{$val.Count}} | {{$val.Error}} |
{{- end }}
{{- end }}
{{- end }}

<details>
  <summary><b>See all test records</b></summary>

| API | Average | Max | Min | Count | Error |
|---|---|---|---|---|---|
{{- range $val := .Items}}
| {{$val.API}} | {{$val.Average}} | {{$val.Max}} | {{$val.Min}} | {{$val.Count}} | {{$val.Error}} |
{{- end }}
</details>
{{- else }}

| API | Average | Max | Min | Count | Error |
|---|---|---|---|---|---|
{{- range $val := .Items}}
| {{$val.API}} | {{$val.Average}} | {{$val.Max}} | {{$val.Min}} | {{$val.Count}} | {{$val.Error}} |
{{- end }}
{{- end }}
