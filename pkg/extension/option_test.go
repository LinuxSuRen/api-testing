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

package extension

import (
	"testing"
)

func TestGetListenAddress(t *testing.T) {
	opt := &Extension{}
	opt.Socket = "test"
	opt.Port = 8080
	protocol, address := opt.GetListenAddress()
	if protocol != "unix" {
		t.Errorf("Expected unix, but got %s", protocol)
	}
	if address != "test" {
		t.Errorf("Expected test, but got %s", address)
	}
	opt.Socket = ""
	protocol, address = opt.GetListenAddress()
	if protocol != "tcp" {
		t.Errorf("Expected tcp, but got %s", protocol)
	}
	if address != ":8080" {
		t.Errorf("Expected :8080, but got %s", address)
	}
}
