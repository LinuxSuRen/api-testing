/*
Copyright 2023-2024 API Testing Authors.

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
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/spf13/cobra"
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

func TestCreateRunner(t *testing.T) {

	t.Run("invalid port", func(t *testing.T) {
		extMgr := NewExtension("git", "store", 75530)
		extMgr.Port = 75530
		assert.NotNil(t, extMgr)
		assert.Error(t, CreateRunner(extMgr, nil, nil))
		assert.Error(t, CreateMonitor(extMgr, nil, nil))
		assert.Error(t, CreateExtensionRunner(extMgr, nil, nil))
	})

	t.Run("random port", func(t *testing.T) {
		extMgr := NewExtension("git", "store", -1)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		command := &cobra.Command{}
		command.SetContext(ctx)
		assert.Error(t, CreateRunner(extMgr, command, nil))
	})

	t.Run("random port, CreateMonitor", func(t *testing.T) {
		extMgr := NewExtension("git", "store", -1)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		command := &cobra.Command{}
		command.SetContext(ctx)
		assert.Error(t, CreateMonitor(extMgr, command, nil))
	})

	t.Run("random port, CreateExtensionRunner", func(t *testing.T) {
		extMgr := NewExtension("git", "store", -1)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		command := &cobra.Command{}
		command.SetContext(ctx)
		assert.Error(t, CreateExtensionRunner(extMgr, command, nil))
	})

	t.Run("socket", func(t *testing.T) {
		extMgr := NewExtension("git", "store", -1)
		extMgr.Socket = filepath.Join(os.TempDir(), time.Microsecond.String())

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		command := &cobra.Command{}
		command.SetContext(ctx)
		assert.Error(t, CreateRunner(extMgr, command, nil))
	})

	t.Run("socket, CreateMonitor", func(t *testing.T) {
		extMgr := NewExtension("git", "store", -1)
		extMgr.Socket = filepath.Join(os.TempDir(), time.Microsecond.String())

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		command := &cobra.Command{}
		command.SetContext(ctx)
		assert.Error(t, CreateMonitor(extMgr, command, nil))
	})

	t.Run("socket, CreateExtensionRunner", func(t *testing.T) {
		extMgr := NewExtension("git", "store", -1)
		extMgr.Socket = filepath.Join(os.TempDir(), time.Microsecond.String())

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		command := &cobra.Command{}
		command.SetContext(ctx)
		assert.Error(t, CreateExtensionRunner(extMgr, command, nil))
	})
}
