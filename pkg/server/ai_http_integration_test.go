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
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	fakeruntime "github.com/linuxsuren/go-fake-runtime"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// AIPluginAPIResponse represents the response structure for AI plugin API endpoints
type AIPluginAPIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func TestAIPluginHTTPEndpoints(t *testing.T) {
	// Setup test environment
	tempDir, err := os.MkdirTemp("", "ai_http_test_*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Initialize ExtManager
	mgr := NewStoreExtManagerInstance(&fakeruntime.FakeExecer{
		ExpectLookPath: "/usr/local/bin/go",
	})
	defer mgr.StopAll()

	t.Run("POST /api/v1/ai/plugins/register", func(t *testing.T) {
		// Create mock socket for testing
		socketPath := fmt.Sprintf("unix://%s/register-test-plugin.sock", tempDir)
		mockPlugin := NewMockAIPlugin(socketPath)
		err := mockPlugin.CreateSocketFile()
		require.NoError(t, err)
		defer mockPlugin.RemoveSocketFile()

		pluginInfo := AIPluginInfo{
			Name:         "register-test-plugin",
			Version:      "1.0.0",
			Description:  "Test plugin for HTTP endpoint",
			Capabilities: []string{"sql-generation", "code-analysis"},
			SocketPath:   socketPath,
			Metadata: map[string]string{
				"author": "test-team",
				"type":   "ai",
			},
		}

		// Create HTTP handler for AI plugin registration
		handler := func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}

			var receivedPlugin AIPluginInfo
			if err := json.NewDecoder(r.Body).Decode(&receivedPlugin); err != nil {
				response := AIPluginAPIResponse{
					Success: false,
					Error:   "Invalid JSON payload",
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(response)
				return
			}

			// Register plugin using ExtManager
			err := mgr.RegisterAIPlugin(receivedPlugin)
			if err != nil {
				response := AIPluginAPIResponse{
					Success: false,
					Error:   err.Error(),
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(response)
				return
			}

			response := AIPluginAPIResponse{
				Success: true,
				Message: "Plugin registered successfully",
				Data:    receivedPlugin,
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(response)
		}

		// Test the endpoint
		jsonData, err := json.Marshal(pluginInfo)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/ai/plugins/register", bytes.NewReader(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		
		var response AIPluginAPIResponse
		err = json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)
		assert.True(t, response.Success)
		assert.Contains(t, response.Message, "registered successfully")
	})

	t.Run("GET /api/v1/ai/plugins/discover", func(t *testing.T) {
		// Pre-register some test plugins
		testPlugins := []AIPluginInfo{
			{
				Name:         "discover-plugin-1",
				Version:      "1.0.0",
				Description:  "Discovery test plugin 1",
				Capabilities: []string{"text-analysis"},
				SocketPath:   fmt.Sprintf("unix://%s/discover-1.sock", tempDir),
			},
			{
				Name:         "discover-plugin-2",
				Version:      "1.1.0",
				Description:  "Discovery test plugin 2",
				Capabilities: []string{"image-analysis", "nlp"},
				SocketPath:   fmt.Sprintf("unix://%s/discover-2.sock", tempDir),
			},
		}

		for _, plugin := range testPlugins {
			err := mgr.RegisterAIPlugin(plugin)
			require.NoError(t, err)
		}

		// Create discovery endpoint handler
		handler := func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}

			plugins, err := mgr.DiscoverAIPlugins()
			if err != nil {
				response := AIPluginAPIResponse{
					Success: false,
					Error:   err.Error(),
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(response)
				return
			}

			response := AIPluginAPIResponse{
				Success: true,
				Data:    plugins,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}

		req := httptest.NewRequest(http.MethodGet, "/api/v1/ai/plugins/discover", nil)
		w := httptest.NewRecorder()

		handler(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var response AIPluginAPIResponse
		err = json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)
		assert.True(t, response.Success)

		// Verify plugin data
		pluginsData, ok := response.Data.([]interface{})
		require.True(t, ok)
		assert.GreaterOrEqual(t, len(pluginsData), 2)
	})

	t.Run("GET /api/v1/ai/plugins/{name}/health", func(t *testing.T) {
		// Setup test plugin with socket
		socketPath := fmt.Sprintf("unix://%s/health-test-plugin.sock", tempDir)
		mockPlugin := NewMockAIPlugin(socketPath)
		err := mockPlugin.CreateSocketFile()
		require.NoError(t, err)
		defer mockPlugin.RemoveSocketFile()

		pluginInfo := AIPluginInfo{
			Name:       "health-test-plugin",
			Version:    "1.0.0",
			SocketPath: socketPath,
		}
		err = mgr.RegisterAIPlugin(pluginInfo)
		require.NoError(t, err)

		// Create health check endpoint handler
		handler := func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}

			// Extract plugin name from URL path (simulated)
			pluginName := "health-test-plugin"

			health, err := mgr.CheckAIPluginHealth(pluginName)
			if err != nil {
				response := AIPluginAPIResponse{
					Success: false,
					Error:   err.Error(),
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(response)
				return
			}

			response := AIPluginAPIResponse{
				Success: true,
				Data:    health,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}

		req := httptest.NewRequest(http.MethodGet, "/api/v1/ai/plugins/health-test-plugin/health", nil)
		w := httptest.NewRecorder()

		handler(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var response AIPluginAPIResponse
		err = json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)
		assert.True(t, response.Success)
		assert.NotNil(t, response.Data)
	})

	t.Run("GET /api/v1/ai/plugins/health", func(t *testing.T) {
		// Create bulk health check endpoint handler
		handler := func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}

			healthMap, err := mgr.GetAllAIPluginHealth()
			if err != nil {
				response := AIPluginAPIResponse{
					Success: false,
					Error:   err.Error(),
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(response)
				return
			}

			response := AIPluginAPIResponse{
				Success: true,
				Data:    healthMap,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}

		req := httptest.NewRequest(http.MethodGet, "/api/v1/ai/plugins/health", nil)
		w := httptest.NewRecorder()

		handler(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var response AIPluginAPIResponse
		err = json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)
		assert.True(t, response.Success)
		assert.NotNil(t, response.Data)
	})

	t.Run("DELETE /api/v1/ai/plugins/{name}", func(t *testing.T) {
		// Register a plugin to be deleted
		pluginInfo := AIPluginInfo{
			Name:       "delete-test-plugin",
			Version:    "1.0.0",
			SocketPath: fmt.Sprintf("unix://%s/delete-test.sock", tempDir),
		}
		err := mgr.RegisterAIPlugin(pluginInfo)
		require.NoError(t, err)

		// Create delete endpoint handler
		handler := func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodDelete {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}

			// Extract plugin name from URL path (simulated)
			pluginName := "delete-test-plugin"

			err := mgr.UnregisterAIPlugin(pluginName)
			if err != nil {
				response := AIPluginAPIResponse{
					Success: false,
					Error:   err.Error(),
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(response)
				return
			}

			response := AIPluginAPIResponse{
				Success: true,
				Message: "Plugin unregistered successfully",
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/ai/plugins/delete-test-plugin", nil)
		w := httptest.NewRecorder()

		handler(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var response AIPluginAPIResponse
		err = json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)
		assert.True(t, response.Success)
		assert.Contains(t, response.Message, "unregistered successfully")

		// Verify plugin is actually removed
		plugins, err := mgr.DiscoverAIPlugins()
		require.NoError(t, err)
		
		pluginExists := false
		for _, p := range plugins {
			if p.Name == "delete-test-plugin" {
				pluginExists = true
				break
			}
		}
		assert.False(t, pluginExists)
	})

	t.Run("error handling and validation", func(t *testing.T) {
		// Test invalid JSON payload
		handler := func(w http.ResponseWriter, r *http.Request) {
			var pluginInfo AIPluginInfo
			if err := json.NewDecoder(r.Body).Decode(&pluginInfo); err != nil {
				response := AIPluginAPIResponse{
					Success: false,
					Error:   "Invalid JSON payload",
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(response)
				return
			}

			// Test validation
			if pluginInfo.Name == "" {
				response := AIPluginAPIResponse{
					Success: false,
					Error:   "Plugin name is required",
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(response)
				return
			}

			w.WriteHeader(http.StatusOK)
		}

		// Test malformed JSON
		req := httptest.NewRequest(http.MethodPost, "/api/v1/ai/plugins/register", 
			bytes.NewReader([]byte(`{"invalid": json}`)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Test empty plugin name
		emptyPlugin := AIPluginInfo{Name: ""}
		jsonData, _ := json.Marshal(emptyPlugin)
		req = httptest.NewRequest(http.MethodPost, "/api/v1/ai/plugins/register", 
			bytes.NewReader(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()

		handler(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestAIPluginAPIPerformance(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "ai_api_perf_*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	mgr := NewStoreExtManagerInstance(&fakeruntime.FakeExecer{
		ExpectLookPath: "/usr/local/bin/go",
	})
	defer mgr.StopAll()

	t.Run("response time benchmarks", func(t *testing.T) {
		// Test discovery endpoint response time
		handler := func(w http.ResponseWriter, r *http.Request) {
			plugins, _ := mgr.DiscoverAIPlugins()
			response := AIPluginAPIResponse{
				Success: true,
				Data:    plugins,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}

		// Measure response time
		start := time.Now()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/ai/plugins/discover", nil)
		w := httptest.NewRecorder()
		handler(w, req)
		elapsed := time.Since(start)

		// Response should be under 100ms for discovery
		assert.Less(t, elapsed.Milliseconds(), int64(100), 
			"Discovery endpoint response time exceeded 100ms")

		// Test health check endpoint response time
		healthHandler := func(w http.ResponseWriter, r *http.Request) {
			healthMap, _ := mgr.GetAllAIPluginHealth()
			response := AIPluginAPIResponse{
				Success: true,
				Data:    healthMap,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}

		start = time.Now()
		req = httptest.NewRequest(http.MethodGet, "/api/v1/ai/plugins/health", nil)
		w = httptest.NewRecorder()
		healthHandler(w, req)
		elapsed = time.Since(start)

		// Health check response should be under 500ms
		assert.Less(t, elapsed.Milliseconds(), int64(500), 
			"Health check endpoint response time exceeded 500ms")
	})
}