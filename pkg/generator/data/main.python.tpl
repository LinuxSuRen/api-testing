'''
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
'''
import requests
from urllib.parse import urlencode

def main():
    url = "{{.Request.API}}"
    method = "{{.Request.Method}}"

    headers = {
        {{- range $key, $val := .Request.Header}}
        "{{$key}}": "{{$val}}",
        {{- end}}
    }

    cookies = {
        {{- range $key, $val := .Request.Cookie}}
        "{{$key}}": "{{$val}}",
        {{- end}}
    }

    data = None
    if method in ["POST", "PUT", "PATCH"]:
        {{- if gt (len .Request.Form) 0 }}
        data = {
            {{- range $key, $val := .Request.Form}}
            "{{$key}}": "{{$val}}",
            {{- end}}
        }
        data = urlencode(data)
        headers["Content-Type"] = "application/x-www-form-urlencoded"
        {{- else}}
        data = "{{.Request.Body.String}}"
        {{- end}}

    try:
        response = requests.request(method, url, headers=headers, cookies=cookies, data=data)
        response.raise_for_status()
    except requests.RequestException as e:
        print(f"An error occurred: {e}")
        return

    print(response.text)

if __name__ == "__main__":
    main()
