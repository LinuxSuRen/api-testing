/*
Copyright 2023-2025 API Testing Authors.

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
package server

import (
	"errors"
	"testing"
	"time"

	fakeruntime "github.com/linuxsuren/go-fake-runtime"
	"github.com/stretchr/testify/assert"
)

func TestStoreExtManager(t *testing.T) {
	t.Run("not found", func(t *testing.T) {
		mgr := NewStoreExtManagerInstance(&fakeruntime.FakeExecer{
			ExpectLookPathError: errors.New("not found"),
		})
		err := mgr.Start("fake", "")
		assert.Error(t, err)
	})

	t.Run("exist executable file", func(t *testing.T) {
		mgr := NewStoreExtManagerInstance(&fakeruntime.FakeExecer{
			ExpectLookPath: "/usr/local/bin/go",
		})
		err := mgr.Start("go", "")
		assert.NoError(t, err, err)

		time.Sleep(time.Microsecond * 100)
		err = mgr.Start("go", "")
		assert.NoError(t, err)

		time.Sleep(time.Microsecond * 100)
		err = mgr.StopAll()
		assert.NoError(t, err)
	})
}

func TestAIPluginManagement(t *testing.T) {
	mgr := NewStoreExtManagerInstance(&fakeruntime.FakeExecer{
		ExpectLookPath: "/usr/local/bin/go",
	})

	t.Run("register and discover AI plugins", func(t *testing.T) {
		// Test plugin registration
		pluginInfo := AIPluginInfo{
			Name:         "test-ai-plugin",
			Version:      "1.0.0",
			Description:  "Test AI plugin for unit testing",
			Capabilities: []string{"sql-generation", "code-analysis"},
			SocketPath:   "unix:///tmp/test-ai-plugin.sock",
			Metadata: map[string]string{
				"author": "test-team",
				"type":   "ai",
			},
		}

		err := mgr.RegisterAIPlugin(pluginInfo)
		assert.NoError(t, err)

		// Test plugin discovery
		plugins, err := mgr.DiscoverAIPlugins()
		assert.NoError(t, err)
		assert.Len(t, plugins, 1)
		assert.Equal(t, "test-ai-plugin", plugins[0].Name)
		assert.Equal(t, "1.0.0", plugins[0].Version)
		assert.Contains(t, plugins[0].Capabilities, "sql-generation")
	})

	t.Run("check AI plugin health", func(t *testing.T) {
		// Check individual plugin health
		health, err := mgr.CheckAIPluginHealth("test-ai-plugin")
		assert.NoError(t, err)
		assert.NotNil(t, health)
		assert.Equal(t, "test-ai-plugin", health.Name)
		// Since socket doesn't exist, status should be offline
		assert.Equal(t, "offline", health.Status)
		assert.Contains(t, health.ErrorMessage, "Plugin socket not found")

		// Check all plugins health
		allHealth, err := mgr.GetAllAIPluginHealth()
		assert.NoError(t, err)
		assert.Len(t, allHealth, 1)
		assert.Contains(t, allHealth, "test-ai-plugin")
	})

	t.Run("unregister AI plugin", func(t *testing.T) {
		// Unregister plugin
		err := mgr.UnregisterAIPlugin("test-ai-plugin")
		assert.NoError(t, err)

		// Verify plugin is removed
		plugins, err := mgr.DiscoverAIPlugins()
		assert.NoError(t, err)
		assert.Len(t, plugins, 0)

		// Try to unregister non-existent plugin
		err = mgr.UnregisterAIPlugin("non-existent")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("register plugin with invalid data", func(t *testing.T) {
		// Test empty name
		err := mgr.RegisterAIPlugin(AIPluginInfo{
			Name:       "",
			SocketPath: "unix:///tmp/test.sock",
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "name cannot be empty")

		// Test empty socket path
		err = mgr.RegisterAIPlugin(AIPluginInfo{
			Name:       "test-plugin",
			SocketPath: "",
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "socket path cannot be empty")
	})

	t.Run("check non-existent plugin health", func(t *testing.T) {
		health, err := mgr.CheckAIPluginHealth("non-existent")
		assert.Error(t, err)
		assert.Nil(t, health)
		assert.Contains(t, err.Error(), "not found")
	})

	// Cleanup
	err := mgr.StopAll()
	assert.NoError(t, err)
}
