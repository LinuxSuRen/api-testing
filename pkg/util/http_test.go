/**
MIT License

Copyright (c) 2023 API Testing Authors.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
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
