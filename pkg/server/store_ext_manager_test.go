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

func TestUnifiedPluginManagement(t *testing.T) {
	mgr := NewStoreExtManagerInstance(&fakeruntime.FakeExecer{
		ExpectLookPath: "/usr/local/bin/go",
	})

	t.Run("discover AI plugins via category filter", func(t *testing.T) {
		// Test AI plugin discovery using the new unified method
		aiPlugins, err := mgr.GetPluginsByCategory("ai")
		assert.NoError(t, err)
		// Initially should be empty since no AI plugins are configured in test environment
		assert.NotNil(t, aiPlugins)
	})

	t.Run("plugin health check functionality", func(t *testing.T) {
		// Test plugin health check with standard plugin
		health, err := mgr.CheckPluginHealth("go")
		assert.NoError(t, err)
		assert.NotNil(t, health)
		assert.Equal(t, "go", health.Name)
		// Plugin should be offline since it's not actually started
		assert.Equal(t, "offline", health.Status)
	})

	t.Run("get all plugin health", func(t *testing.T) {
		// Test getting all plugin health status
		allHealth, err := mgr.GetAllPluginHealth()
		assert.NoError(t, err)
		assert.NotNil(t, allHealth)
		// Should return map of plugin health status
	})

	t.Run("plugin category filtering", func(t *testing.T) {
		// Test filtering by different categories
		dbPlugins, err := mgr.GetPluginsByCategory("database")
		assert.NoError(t, err)
		assert.NotNil(t, dbPlugins)
		
		webPlugins, err := mgr.GetPluginsByCategory("web")
		assert.NoError(t, err)
		assert.NotNil(t, webPlugins)
	})

	t.Run("health check for non-existent plugin", func(t *testing.T) {
		health, err := mgr.CheckPluginHealth("non-existent-plugin")
		assert.NoError(t, err) // Should not error for non-existent plugins
		assert.NotNil(t, health)
		assert.Equal(t, "non-existent-plugin", health.Name)
		assert.Equal(t, "offline", health.Status)
	})

	// Cleanup
	err := mgr.StopAll()
	assert.NoError(t, err)
}
