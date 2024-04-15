[
{{range $index, $result := .}}
    {
        "Name": "{{$result.Name}}",
        "API": "{{$result.API}}",
        "Count": {{$result.Count}},
        "Average": "{{$result.Average}}",
        "Max": "{{$result.Max}}",
        "Min": "{{$result.Min}}",
        "QPS": {{$result.QPS}},
        "Error": {{$result.Error}},
        "LastErrorMessage": "{{$result.LastErrorMessage}}"
    }
{{end}}
]
