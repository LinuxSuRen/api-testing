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
package util

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
)

func TestClient(t *testing.T) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()

	addr := listener.Addr().String()
	port := addr[strings.LastIndex(addr, ":")+1:]

	visitCount := 0
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(ContentType, Plain)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello"))
		visitCount++
	}))

	go http.Serve(listener, mux)

	client := GetDefaultCachedHTTPClient()
	for i := 0; i < 10; i++ {
		resp, err := client.Get(fmt.Sprintf("http://localhost:%s", port))
		if err != nil {
			t.Fatal(err)
		}

		if val := resp.Header.Get(ContentType); val != Plain {
			t.Fatalf("invalid content type, got %q, expect %q", val, Plain)
		}

		if code := resp.Status; code != "200 OK" {
			t.Fatalf("invalid status code, got %q, expect %q", code, "200 OK")
		}

		var data []byte
		if data, err = io.ReadAll(resp.Body); err != nil {
			t.Fatal(err)
		} else if string(data) != "hello" {
			t.Fatal("invalid response")
		}
	}

	if visitCount != 1 {
		t.Fatal("invalid visit count")
	}
}
