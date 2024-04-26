/*
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
*/

async function makeRequest() {
    {{- if gt (len .Request.Header) 0 }}
    const headers = new Headers();
    {{- range $key, $val := .Request.Header}}
    headers.append("{{$key}}", "{{$val}}");
    {{- end}}
    {{- end}}

    {{- if gt (len .Request.Cookie) 0 }}
    const cookieHeader = [
        {{- range $key, $val := .Request.Cookie}}
        "{{$key}}={{$val}}",
        {{- end}}
    ].join("; ");
    headers.append("Cookie", cookieHeader);
    {{- end}}

    const requestOptions = {
        method: "{{.Request.Method}}",
        {{- if gt (len .Request.Header) 0 }}
        headers: headers,
        {{- end}}
        redirect: "follow"
    };

    {{- if ne .Request.Body.String "" }}
    console.log('the body is ignored because this is a GET request')
    /*
    {{- else if gt (len .Request.Form) 0 }}
    console.log('the body is ignored because this is a GET request')
    /*
    {{- end}}
    {{- if gt (len .Request.Form) 0 }}
    const formData = new URLSearchParams();
    {{- range $key, $val := .Request.Form}}
    formData.append("{{$key}}", "{{$val}}");
    {{- end}}
    let body = formData;
    requestOptions.body = body;
    {{- else if ne .Request.Body.String ""}}
    let body = `{{.Request.Body.String | safeString}}`;
    requestOptions.body = body;
    {{- end}}
    {{- if ne .Request.Body.String "" }}
    */
    {{- else if gt (len .Request.Form) 0 }}
    */
    {{- end}}

    try {
        const response = await fetch("{{.Request.API}}", requestOptions);
        if (response.ok) { // Check if HTTP status code is 200/OK
            const text = await response.text();
            console.log(text);
        } else {
            throw new Error('Network response was not ok. Status code: ' + response.status);
        }
    } catch (error) {
        console.error('Fetch error:', error);
    }
}

makeRequest();
