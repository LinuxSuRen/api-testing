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
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	fakeruntime "github.com/linuxsuren/go-fake-runtime"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockAIPlugin represents a mock AI plugin for testing
type MockAIPlugin struct {
	socketPath    string
	responseDelay time.Duration
	shouldError   bool
	errorMessage  string
	status        string
}

func NewMockAIPlugin(socketPath string) *MockAIPlugin {
	return &MockAIPlugin{
		socketPath:    socketPath,
		responseDelay: 100 * time.Millisecond,
		shouldError:   false,
		status:        "online",
	}
}

func (m *MockAIPlugin) SetResponseDelay(delay time.Duration) {
	m.responseDelay = delay
}

func (m *MockAIPlugin) SetError(shouldError bool, message string) {
	m.shouldError = shouldError
	m.errorMessage = message
	if shouldError {
		m.status = "error"
	} else {
		m.status = "online"
	}
}

func (m *MockAIPlugin) SetStatus(status string) {
	m.status = status
}

func (m *MockAIPlugin) CreateSocketFile() error {
	// Create socket directory
	dir := filepath.Dir(strings.TrimPrefix(m.socketPath, "unix://"))
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	
	// Create mock socket file
	socketFile := strings.TrimPrefix(m.socketPath, "unix://")
	file, err := os.Create(socketFile)
	if err != nil {
		return err
	}
	file.Close()
	
	return nil
}

func (m *MockAIPlugin) RemoveSocketFile() error {
	socketFile := strings.TrimPrefix(m.socketPath, "unix://")
	return os.RemoveAll(socketFile)
}

func TestAIIntegrationEndToEnd(t *testing.T) {
	// Setup test environment
	tempDir, err := os.MkdirTemp("", "ai_test_*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Initialize ExtManager with fake execer
	mgr := NewStoreExtManagerInstance(&fakeruntime.FakeExecer{
		ExpectLookPath: "/usr/local/bin/go",
	})

	t.Run("complete AI plugin lifecycle", func(t *testing.T) {
		// Setup mock AI plugin
		socketPath := fmt.Sprintf("unix://%s/test-ai-plugin.sock", tempDir)
		mockPlugin := NewMockAIPlugin(socketPath)
		
		// Create socket file to simulate online plugin
		err := mockPlugin.CreateSocketFile()
		require.NoError(t, err)
		defer mockPlugin.RemoveSocketFile()

		// Test plugin registration
		pluginInfo := AIPluginInfo{
			Name:         "integration-test-plugin",
			Version:      "1.0.0",
			Description:  "Integration test AI plugin",
			Capabilities: []string{"sql-generation", "code-analysis", "natural-language"},
			SocketPath:   socketPath,
			Metadata: map[string]string{
				"author":      "integration-test",
				"type":        "ai",
				"environment": "test",
			},
		}

		// Register plugin
		err = mgr.RegisterAIPlugin(pluginInfo)
		assert.NoError(t, err)

		// Test plugin discovery
		plugins, err := mgr.DiscoverAIPlugins()
		assert.NoError(t, err)
		assert.Len(t, plugins, 1)
		assert.Equal(t, "integration-test-plugin", plugins[0].Name)
		assert.Contains(t, plugins[0].Capabilities, "sql-generation")

		// Test health check
		health, err := mgr.CheckAIPluginHealth("integration-test-plugin")
		assert.NoError(t, err)
		assert.NotNil(t, health)
		assert.Equal(t, "integration-test-plugin", health.Name)
		assert.Equal(t, "online", health.Status)
		assert.Empty(t, health.ErrorMessage)

		// Test bulk health check
		allHealth, err := mgr.GetAllAIPluginHealth()
		assert.NoError(t, err)
		assert.Len(t, allHealth, 1)
		assert.Contains(t, allHealth, "integration-test-plugin")

		// Test plugin unregistration
		err = mgr.UnregisterAIPlugin("integration-test-plugin")
		assert.NoError(t, err)

		// Verify plugin is removed
		plugins, err = mgr.DiscoverAIPlugins()
		assert.NoError(t, err)
		assert.Len(t, plugins, 0)
	})

	t.Run("error scenarios and recovery", func(t *testing.T) {
		// Test offline plugin scenario
		socketPath := fmt.Sprintf("unix://%s/offline-plugin.sock", tempDir)
		pluginInfo := AIPluginInfo{
			Name:       "offline-test-plugin",
			Version:    "1.0.0",
			SocketPath: socketPath,
		}

		// Register plugin without creating socket (offline state)
		err := mgr.RegisterAIPlugin(pluginInfo)
		assert.NoError(t, err)

		// Check health should detect offline status
		health, err := mgr.CheckAIPluginHealth("offline-test-plugin")
		assert.NoError(t, err)
		assert.Equal(t, "offline", health.Status)
		assert.Contains(t, health.ErrorMessage, "Plugin socket not found")

		// Test plugin comes online
		mockPlugin := NewMockAIPlugin(socketPath)
		err = mockPlugin.CreateSocketFile()
		require.NoError(t, err)
		defer mockPlugin.RemoveSocketFile()

		// Health check should now show online
		health, err = mgr.CheckAIPluginHealth("offline-test-plugin")
		assert.NoError(t, err)
		assert.Equal(t, "online", health.Status)
		assert.Empty(t, health.ErrorMessage)

		// Cleanup
		err = mgr.UnregisterAIPlugin("offline-test-plugin")
		assert.NoError(t, err)
	})

	t.Run("performance and resource usage", func(t *testing.T) {
		// Create multiple plugins to test system performance
		const numPlugins = 10
		var plugins []*MockAIPlugin
		var pluginInfos []AIPluginInfo

		// Setup multiple mock plugins
		for i := 0; i < numPlugins; i++ {
			socketPath := fmt.Sprintf("unix://%s/perf-plugin-%d.sock", tempDir, i)
			mockPlugin := NewMockAIPlugin(socketPath)
			
			err := mockPlugin.CreateSocketFile()
			require.NoError(t, err)
			plugins = append(plugins, mockPlugin)

			pluginInfo := AIPluginInfo{
				Name:         fmt.Sprintf("perf-test-plugin-%d", i),
				Version:      "1.0.0",
				Description:  fmt.Sprintf("Performance test plugin %d", i),
				Capabilities: []string{"performance-test"},
				SocketPath:   socketPath,
				Metadata: map[string]string{
					"test_id": fmt.Sprintf("%d", i),
				},
			}
			pluginInfos = append(pluginInfos, pluginInfo)
		}

		defer func() {
			for _, plugin := range plugins {
				plugin.RemoveSocketFile()
			}
		}()

		// Measure registration performance
		startTime := time.Now()
		for _, info := range pluginInfos {
			err := mgr.RegisterAIPlugin(info)
			assert.NoError(t, err)
		}
		registrationTime := time.Since(startTime)

		// Registration should complete within reasonable time
		assert.Less(t, registrationTime.Milliseconds(), int64(1000), "Plugin registration took too long")

		// Test bulk health check performance
		startTime = time.Now()
		allHealth, err := mgr.GetAllAIPluginHealth()
		assert.NoError(t, err)
		healthCheckTime := time.Since(startTime)

		// Health check should be fast
		assert.Less(t, healthCheckTime.Milliseconds(), int64(500), "Bulk health check took too long")
		assert.Len(t, allHealth, numPlugins)

		// Test discovery performance
		startTime = time.Now()
		discoveredPlugins, err := mgr.DiscoverAIPlugins()
		assert.NoError(t, err)
		discoveryTime := time.Since(startTime)

		assert.Less(t, discoveryTime.Milliseconds(), int64(200), "Plugin discovery took too long")
		assert.Len(t, discoveredPlugins, numPlugins)

		// Cleanup all plugins
		for i := 0; i < numPlugins; i++ {
			err := mgr.UnregisterAIPlugin(fmt.Sprintf("perf-test-plugin-%d", i))
			assert.NoError(t, err)
		}
	})

	t.Run("concurrent operations", func(t *testing.T) {
		// Test concurrent plugin operations
		const numConcurrent = 5
		doneChan := make(chan bool, numConcurrent)
		errorChan := make(chan error, numConcurrent)

		// Run concurrent plugin registrations
		for i := 0; i < numConcurrent; i++ {
			go func(id int) {
				socketPath := fmt.Sprintf("unix://%s/concurrent-plugin-%d.sock", tempDir, id)
				mockPlugin := NewMockAIPlugin(socketPath)
				
				err := mockPlugin.CreateSocketFile()
				if err != nil {
					errorChan <- err
					return
				}
				defer mockPlugin.RemoveSocketFile()

				pluginInfo := AIPluginInfo{
					Name:       fmt.Sprintf("concurrent-plugin-%d", id),
					Version:    "1.0.0",
					SocketPath: socketPath,
				}

				// Register plugin
				err = mgr.RegisterAIPlugin(pluginInfo)
				if err != nil {
					errorChan <- err
					return
				}

				// Check health
				_, err = mgr.CheckAIPluginHealth(fmt.Sprintf("concurrent-plugin-%d", id))
				if err != nil {
					errorChan <- err
					return
				}

				// Unregister plugin
				err = mgr.UnregisterAIPlugin(fmt.Sprintf("concurrent-plugin-%d", id))
				if err != nil {
					errorChan <- err
					return
				}

				doneChan <- true
			}(i)
		}

		// Wait for all operations to complete
		timeout := time.After(10 * time.Second)
		completed := 0
		for completed < numConcurrent {
			select {
			case <-doneChan:
				completed++
			case err := <-errorChan:
				t.Errorf("Concurrent operation failed: %v", err)
			case <-timeout:
				t.Error("Concurrent operations timed out")
				return
			}
		}
	})

	// Cleanup ExtManager
	err = mgr.StopAll()
	assert.NoError(t, err)
}

func TestAIPluginHealthMonitoring(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "ai_health_test_*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	mgr := NewStoreExtManagerInstance(&fakeruntime.FakeExecer{
		ExpectLookPath: "/usr/local/bin/go",
	})

	t.Run("automatic health monitoring", func(t *testing.T) {
		socketPath := fmt.Sprintf("unix://%s/health-monitor-plugin.sock", tempDir)
		mockPlugin := NewMockAIPlugin(socketPath)

		pluginInfo := AIPluginInfo{
			Name:       "health-monitor-plugin",
			Version:    "1.0.0",
			SocketPath: socketPath,
		}

		// Register plugin without socket (offline)
		err := mgr.RegisterAIPlugin(pluginInfo)
		require.NoError(t, err)

		// Initial health should be offline
		health, err := mgr.CheckAIPluginHealth("health-monitor-plugin")
		assert.NoError(t, err)
		assert.Equal(t, "offline", health.Status)

		// Create socket to simulate plugin coming online
		err = mockPlugin.CreateSocketFile()
		require.NoError(t, err)
		defer mockPlugin.RemoveSocketFile()

		// Wait for health monitoring to detect the change
		// Note: In a real scenario, the health monitoring runs every 30 seconds
		// For testing, we trigger manual health checks
		time.Sleep(100 * time.Millisecond)
		
		health, err = mgr.CheckAIPluginHealth("health-monitor-plugin")
		assert.NoError(t, err)
		assert.Equal(t, "online", health.Status)

		// Remove socket to simulate plugin going offline
		err = mockPlugin.RemoveSocketFile()
		require.NoError(t, err)

		health, err = mgr.CheckAIPluginHealth("health-monitor-plugin")
		assert.NoError(t, err)
		assert.Equal(t, "offline", health.Status)

		// Cleanup
		err = mgr.UnregisterAIPlugin("health-monitor-plugin")
		assert.NoError(t, err)
	})

	err = mgr.StopAll()
	assert.NoError(t, err)
}

// Benchmark tests for performance validation
func BenchmarkAIPluginOperations(b *testing.B) {
	tempDir, err := os.MkdirTemp("", "ai_bench_*")
	require.NoError(b, err)
	defer os.RemoveAll(tempDir)

	mgr := NewStoreExtManagerInstance(&fakeruntime.FakeExecer{
		ExpectLookPath: "/usr/local/bin/go",
	})

	b.Run("RegisterAIPlugin", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			pluginInfo := AIPluginInfo{
				Name:       fmt.Sprintf("bench-plugin-%d", i),
				Version:    "1.0.0",
				SocketPath: fmt.Sprintf("unix://%s/bench-%d.sock", tempDir, i),
			}
			
			mgr.RegisterAIPlugin(pluginInfo)
		}
	})

	b.Run("DiscoverAIPlugins", func(b *testing.B) {
		// Pre-register some plugins
		for i := 0; i < 100; i++ {
			pluginInfo := AIPluginInfo{
				Name:       fmt.Sprintf("discover-bench-plugin-%d", i),
				Version:    "1.0.0",
				SocketPath: fmt.Sprintf("unix://%s/discover-bench-%d.sock", tempDir, i),
			}
			mgr.RegisterAIPlugin(pluginInfo)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			mgr.DiscoverAIPlugins()
		}
	})

	b.Run("GetAllAIPluginHealth", func(b *testing.B) {
		// Pre-register plugins with sockets
		for i := 0; i < 50; i++ {
			socketPath := fmt.Sprintf("unix://%s/health-bench-%d.sock", tempDir, i)
			mockPlugin := NewMockAIPlugin(socketPath)
			mockPlugin.CreateSocketFile()
			
			pluginInfo := AIPluginInfo{
				Name:       fmt.Sprintf("health-bench-plugin-%d", i),
				Version:    "1.0.0",
				SocketPath: socketPath,
			}
			mgr.RegisterAIPlugin(pluginInfo)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			mgr.GetAllAIPluginHealth()
		}
	})

	mgr.StopAll()
}