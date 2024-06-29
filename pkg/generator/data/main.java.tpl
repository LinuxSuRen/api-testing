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
import java.io.BufferedReader;
import java.io.InputStreamReader;
import java.net.HttpURLConnection;
import java.net.URL;
import java.net.URLEncoder;

public class Main {
    public static void main(String[] args) throws Exception {
        {{- if gt (len .Request.Form) 0 }}
        StringBuilder postData = new StringBuilder();
        {{- range $key, $val := .Request.Form}}
        postData.append(URLEncoder.encode("{{$key}}", "UTF-8"));
        postData.append("=");
        postData.append(URLEncoder.encode("{{$val}}", "UTF-8"));
        postData.append("&");
        {{- end}}
        byte[] postDataBytes = postData.toString().getBytes("UTF-8");
        {{- else}}
        String body = """
        {{.Request.Body.String}}""";
        byte[] postDataBytes = body.getBytes("UTF-8");
        {{- end }}

        URL url = new URL("{{.Request.API}}");
        HttpURLConnection conn = (HttpURLConnection) url.openConnection();
        conn.setRequestMethod("{{.Request.Method}}");

        {{- range $key, $val := .Request.Header}}
        conn.setRequestProperty("{{$key}}", "{{$val}}");
        {{- end}}

        {{- if gt (len .Request.Cookie) 0 }}
        {{- range $key, $val := .Request.Cookie}}
        conn.setRequestProperty("Cookie", "{{$key}}={{$val}}");
        {{- end}}
        {{- end}}

        conn.setDoOutput(true);
        {{- if ne .Request.Method "GET" }}
        conn.getOutputStream().write(postDataBytes);
        {{- end}}

        BufferedReader reader = new BufferedReader(new InputStreamReader(conn.getInputStream(), "UTF-8"));
        StringBuilder response = new StringBuilder();
        String line;
        while ((line = reader.readLine()) != null) {
            response.append(line);
        }
        reader.close();

        System.out.println(response);

        if (conn.getResponseCode() != HttpURLConnection.HTTP_OK) {
            throw new RuntimeException("status code is not 200");
        }
    }
}
