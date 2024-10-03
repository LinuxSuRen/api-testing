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

package runner

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/linuxsuren/api-testing/pkg/apispec"
	"github.com/linuxsuren/api-testing/pkg/util"
)

type httpResultWriter struct {
	requestMethod string
	targetUrl     string
	parameters    map[string]string
	templateFile  *TemplateOption
}

type TemplateOption struct {
	filename string
	fileType string
}

// NewHTTPResultWriter creates a new httpResultWriter
func NewHTTPResultWriter(requestType string, url string, parameters map[string]string, templateFile *TemplateOption) ReportResultWriter {
	return &httpResultWriter{
		requestMethod: requestType,
		targetUrl:     url,
		parameters:    parameters,
		templateFile:  templateFile,
	}
}

func NewTemplateOption(filename string, fileType string) *TemplateOption {
	return &TemplateOption{
		filename: filename,
		fileType: fileType,
	}
}

// Output writes the JSON base report to target writer
func (w *httpResultWriter) Output(result []ReportResult) (err error) {
	url := w.targetUrl
	for key, value := range w.parameters {
		if url == w.targetUrl {
			url = fmt.Sprintf("%s?%s=%s", url, key, value)
		} else {
			url = fmt.Sprintf("%s&%s=%s", url, key, value)
		}
	}
	log.Println("will send report to:" + url)

	var tmpl *template.Template
	if w.templateFile == nil {
		// use the default template file to serialize the data to JSON format
		tmpl, err = template.New("HTTP report template").Parse(defaultTemplate)
		if err != nil {
			log.Fatalf("Failed to parse template: %v", err)
		}
	} else {
		content, err := os.ReadFile(w.templateFile.filename)
		if err != nil {
			log.Println("Error reading file:", err)
			return err
		}

		tmpl, err = template.New("HTTP report template").Parse(string(content))
		if err != nil {
			log.Println("Error parsing template:", err)
			return err
		}
	}

	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, result)
	if err != nil {
		log.Printf("Failed to render template: %v", err)
		return
	}
	req, err := http.NewRequest(w.requestMethod, url, buf)
	if err != nil {
		log.Println("Error creating request:", err)
		return
	}

	var contentType string
	if w.templateFile != nil {
		switch w.templateFile.fileType {
		case "html":
			contentType = "text/html"
		case "yaml":
			contentType = "application/yaml"
		case "xml":
			contentType = "application/xml"
		default:
			contentType = "application/json"
		}
	} else {
		contentType = "application/json"
	}
	req.Header.Set(util.ContentType, contentType)

	var resp *http.Response
	if resp, err = http.DefaultClient.Do(req); err != nil {
		log.Println("error when client do", err)
		return
	}
	if resp.StatusCode == http.StatusOK {
		var data []byte
		if data, err = io.ReadAll(resp.Body); err != nil {
			log.Println("error when ReadAll", err)
			return
		}
		log.Println("getting response back", data)
	}
	return
}

//go:embed writer_templates/example.tpl
var defaultTemplate string

// WithAPIConverage sets the api coverage
func (w *httpResultWriter) WithAPIConverage(apiConverage apispec.APIConverage) ReportResultWriter {
	return w
}

func (w *httpResultWriter) WithResourceUsage([]ResourceUsage) ReportResultWriter {
	return w
}
