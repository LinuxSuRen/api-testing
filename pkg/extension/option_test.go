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
package extension

import (
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
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

func TestExtension(t *testing.T) {
	extMgr := NewExtension("git", "store", -1)
	assert.NotNil(t, extMgr)
	assert.Equal(t, "atest-store-git", extMgr.GetFullName())

	flags := &pflag.FlagSet{}
	extMgr.AddFlags(flags)

	assert.NotNil(t, flags.Lookup("port"))
	assert.NotNil(t, flags.Lookup("socket"))
}
