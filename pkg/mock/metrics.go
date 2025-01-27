/*
Copyright 2025 API Testing Authors.

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
package mock

import (
	"net/http"
	"sync"
)

// NoopMetrics implements RequestMetrics but does nothing
type NoopMetrics struct{}

// NewNoopMetrics creates a new NoopMetrics instance
func NewNoopMetrics() *NoopMetrics {
	return &NoopMetrics{}
}

// RecordRequest implements RequestMetrics but does nothing
func (m *NoopMetrics) RecordRequest(path string) {}

// GetMetrics implements RequestMetrics but returns empty map
func (m *NoopMetrics) GetMetrics() map[string]int {
	return make(map[string]int)
}

// AddMetricsHandler implements RequestMetrics but does nothing
func (m *NoopMetrics) AddMetricsHandler(mux *http.ServeMux, prefix string) {}

// RequestMetrics represents an interface for collecting request metrics
type RequestMetrics interface {
	RecordRequest(path string)
	GetMetrics() map[string]int
	AddMetricsHandler(mux *http.ServeMux, prefix string)
}

// InMemoryMetrics implements RequestMetrics with in-memory storage
type InMemoryMetrics struct {
	requests map[string]int
	mu       sync.RWMutex
}

// NewInMemoryMetrics creates a new InMemoryMetrics instance
func NewInMemoryMetrics() *InMemoryMetrics {
	return &InMemoryMetrics{
		requests: make(map[string]int),
	}
}

// RecordRequest records a request for the given path
func (m *InMemoryMetrics) RecordRequest(path string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.requests[path]++
}

// GetMetrics returns a copy of the current metrics
func (m *InMemoryMetrics) GetMetrics() map[string]int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Return a copy to avoid map races
	result := make(map[string]int)
	for k, v := range m.requests {
		result[k] = v
	}
	return result
}

func (m *InMemoryMetrics) AddMetricsHandler(mux *http.ServeMux, prefix string) {
	// Add metrics endpoint
	mux.HandleFunc(prefix+"/metrics", func(w http.ResponseWriter, r *http.Request) {
		// metrics handling code
	})
}
