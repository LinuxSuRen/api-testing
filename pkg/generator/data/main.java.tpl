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
        String body = "{{.Request.Body.String}}";
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
        conn.getOutputStream().write(postDataBytes);

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
