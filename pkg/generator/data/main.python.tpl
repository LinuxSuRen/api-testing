import io
import requests
from urllib.parse import urlencode

def main():
    {{- if gt (len .Request.Form) 0 }}
    data = {}
    {{- range $key, $val := .Request.Form}}
    data["{{$key}}"] = "{{$val}}"
    encoded_data = urlencode(data)
    {{- end}}
    body = io.BytesIO(encoded_data.encode("utf-8"))
    {{- else}}
    body = io.BytesIO(b"{{.Request.Body.String}}")
    {{- end}}
    {{- range $key, $val := .Request.Header}}
    headers = {"{{$key}}": "{{$val}}"}
    {{- end}}
    {{- if gt (len .Request.Cookie) 0 }}
    {{- range $key, $val := .Request.Cookie}}
    cookies = {"{{$key}}": "{{$val}}"}
    {{- end}}
    {{- end}}
    {{- if gt (len .Request.Cookie) 0 }}
    try:
        req = requests.Request("{{.Request.Method}}", "{{.Request.API}}", headers=headers, cookies=cookies, data=body)
    except requests.RequestException as e:
        raise e
    {{- else}}
    try:
        req = requests.Request("{{.Request.Method}}", "{{.Request.API}}", headers=headers, data=body)
    except requests.RequestException as e:
        raise e
    {{- end}}

    resp = requests.Session().send(req.prepare())
    if resp.status_code != 200:
        raise Exception("status code is not 200")

    data = resp.content
    print(data.decode("utf-8"))

if __name__ == "__main__":
    main()
