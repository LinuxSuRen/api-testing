| API | Average | Max | Min | Count | Error |
|---|---|---|---|---|---|
{{- range $val := .}}
| {{$val.API}} | {{$val.Average}} | {{$val.Max}} | {{$val.Min}} | {{$val.Count}} | {{$val.Error}} |
{{- end}}
