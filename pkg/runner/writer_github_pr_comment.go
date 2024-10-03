/*
Copyright 2023 API Testing Authors.

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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/linuxsuren/api-testing/pkg/logging"

	"github.com/linuxsuren/api-testing/pkg/apispec"
	"github.com/linuxsuren/api-testing/pkg/util"
)

var (
	githubLogger = logging.DefaultLogger(logging.LogLevelInfo).WithName("github")
)

type githubPRCommentWriter struct {
	*GithubPRCommentOption
}

func NewGithubPRCommentWriter(opt *GithubPRCommentOption) (ReportResultWriter, error) {
	var err error

	opt.Token = util.EmptyThenDefault(opt.Token, os.Getenv("GITHUB_TOKEN"))
	if opt.Repo == "" || opt.Identity == "" || opt.Token == "" {
		err = fmt.Errorf("GitHub report parameters are not enough")
	}
	return &githubPRCommentWriter{
		GithubPRCommentOption: opt,
	}, err
}

func (w *githubPRCommentWriter) loadExistData(newData []ReportResult) (result []ReportResult, err error) {
	result = newData
	if w.ReportFile == "" {
		return
	}

	var data []byte
	if data, err = os.ReadFile(w.ReportFile); err == nil {
		existData := make([]ReportResult, 0)

		if len(data) > 0 {
			if err = json.Unmarshal(data, &existData); err != nil {
				return
			}
		}

		for i := range result {
			for _, item := range existData {
				if result[i].API == item.API {
					result[i].Count += item.Count
					result[i].Error += item.Error
				}
			}
		}
	}

	// write data back to the file
	if data, err = json.Marshal(result); err == nil {
		err = os.WriteFile(w.ReportFile, data, 0644)
	}
	return
}

func (w *githubPRCommentWriter) Output(result []ReportResult) (err error) {
	if w.PR <= 0 {
		githubLogger.Info("skip reporting to GitHub due to without a valid PR number")
		return
	}

	if result, err = w.loadExistData(result); err != nil {
		err = fmt.Errorf("failed to load exist data: %v", err)
		return
	}

	var existCommentId int
	if existCommentId, err = w.exist(); err != nil {
		return
	}

	buf := new(bytes.Buffer)
	mdWriter := NewMarkdownResultWriter(buf)
	if err = mdWriter.Output(result); err == nil {
		content := buf.String() + "\n\n" + w.Identity

		err = w.createOrUpdate(content, existCommentId)
	}
	return
}

func (w *githubPRCommentWriter) exist() (id int, err error) {
	var req *http.Request
	api := fmt.Sprintf("https://api.github.com/repos/%s/issues/%d/comments", w.Repo, w.PR)
	req, err = http.NewRequest(http.MethodGet, api, nil)
	if err != nil {
		return
	}

	var resp *http.Response
	if resp, err = w.sendRequest(req); err == nil {
		comments := make([]comment, 0)

		if err = unmarshalResponseBody(resp, http.StatusOK, &comments); err == nil {
			for _, item := range comments {
				if strings.HasSuffix(item.Body, w.Identity) {
					id = item.ID
					return
				}
			}
		}
	}
	return
}

func (w *githubPRCommentWriter) createOrUpdate(content string, id int) (err error) {
	var req *http.Request
	var api string
	var requestMethod string
	var expectedCode int

	if id > 0 {
		requestMethod = http.MethodPatch
		expectedCode = http.StatusOK
		api = fmt.Sprintf("https://api.github.com/repos/%s/issues/comments/%d", w.Repo, id)
	} else {
		requestMethod = http.MethodPost
		expectedCode = http.StatusCreated
		api = fmt.Sprintf("https://api.github.com/repos/%s/issues/%d/comments", w.Repo, w.PR)
	}

	co := comment{
		Body: content,
	}

	githubLogger.Info("comment id", "id", id)
	var data []byte
	if data, err = json.Marshal(co); err != nil {
		err = fmt.Errorf("failed to marshal body when createOrupdate comment: %v", err)
		return
	}

	req, err = http.NewRequest(requestMethod, api, io.NopCloser(bytes.NewBuffer(data)))
	if err != nil {
		return
	}

	var resp *http.Response
	if resp, err = w.sendRequest(req); err == nil && resp.StatusCode != expectedCode {
		err = fmt.Errorf("failed to update or create comment, status code: %d", resp.StatusCode)
	}
	return
}

func (w *githubPRCommentWriter) sendRequest(req *http.Request) (resp *http.Response, err error) {
	w.setHeader(req)
	resp, err = http.DefaultClient.Do(req)
	return
}

func (w *githubPRCommentWriter) setHeader(req *http.Request) {
	req.Header.Set(util.Authorization, fmt.Sprintf("Bearer %s", w.Token))
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
}

func (w *githubPRCommentWriter) WithAPIConverage(apiConverage apispec.APIConverage) (r ReportResultWriter) {
	// not have this feature
	return
}

func (w *githubPRCommentWriter) WithResourceUsage([]ResourceUsage) ReportResultWriter {
	return w
}

func unmarshalResponseBody(resp *http.Response, expectedCode int, obj interface{}) (err error) {
	if resp.StatusCode != expectedCode {
		err = fmt.Errorf("expect status code: %d, but %d", expectedCode, resp.StatusCode)
		return
	}

	var data []byte
	if data, err = io.ReadAll(resp.Body); err == nil {
		err = json.Unmarshal(data, obj)
	}
	return
}

type GithubPRCommentOption struct {
	Repo       string
	PR         int
	Identity   string
	Token      string
	ReportFile string
}

type comment struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}
