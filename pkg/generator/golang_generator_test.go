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
limitations under the License.*/

package generator

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	atest "github.com/linuxsuren/api-testing/pkg/testing"
)

func TestGolangGenerator_Generate(t *testing.T) {
	generator := NewGolangGenerator()

	testSuite := &atest.TestSuite{
		Name: "Test Suite Example",
	}
	testcase := &atest.TestCase{
		Name: "Test Case Example",
		Request: atest.Request{
			Method: http.MethodGet,
			API:    urlTest,
			Header: map[string]string{
				"Content-Type": "application/json",
			},
			Body: atest.RequestBody{Value: ""},
		},
	}

	// Create a fake HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check the request method and path
		if r.Method != "POST" || r.URL.Path != urlTest {
			http.Error(w, "Unexpected request", http.StatusBadRequest)
			return
		}
		// Write a response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	// Replace the API URL with the fake server's URL
	testcase.Request.API = server.URL + testcase.Request.API

	result, err := generator.Generate(testSuite, testcase)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedSubstring := "http.NewRequest"
	if !strings.Contains(result, expectedSubstring) {
		t.Errorf("Expected result to contain %q, got %q", expectedSubstring, result)
	}

	expectedImport := "\n\t\"bytes\""
	if testcase.Request.Method == http.MethodPost && !strings.Contains(result, expectedImport) {
		t.Errorf("Expected result to contain %q when method is POST, got %q", expectedImport, result)
	}
}

const urlTest = "http://foo"
