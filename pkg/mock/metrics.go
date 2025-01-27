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
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

var _ RequestMetrics = &NoopMetrics{}

// NoopMetrics implements RequestMetrics but does nothing
type NoopMetrics struct{}

// NewNoopMetrics creates a new NoopMetrics instance
func NewNoopMetrics() *NoopMetrics {
	return &NoopMetrics{}
}

// RecordRequest implements RequestMetrics but does nothing
func (m *NoopMetrics) RecordRequest(path string) {}

// GetMetrics implements RequestMetrics but returns empty map
func (m *NoopMetrics) GetMetrics() MetricData {
	return MetricData{}
}

// AddMetricsHandler implements RequestMetrics but does nothing
func (m *NoopMetrics) AddMetricsHandler(mux MetricsHandler) {}

type MetricsHandler interface {
	HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) *mux.Route
}

type MetricData struct {
	FirstRequestTime time.Time
	LastRequestTime  time.Time
	Requests         map[string]int
}

// RequestMetrics represents an interface for collecting request metrics
type RequestMetrics interface {
	RecordRequest(path string)
	GetMetrics() MetricData
	AddMetricsHandler(MetricsHandler)
}

var _ RequestMetrics = &InMemoryMetrics{}

// InMemoryMetrics implements RequestMetrics with in-memory storage
type InMemoryMetrics struct {
	MetricData
	mu sync.RWMutex
}

// NewInMemoryMetrics creates a new InMemoryMetrics instance
func NewInMemoryMetrics() *InMemoryMetrics {
	return &InMemoryMetrics{
		MetricData: MetricData{
			Requests: make(map[string]int),
		},
	}
}

// RecordRequest records a request for the given path
func (m *InMemoryMetrics) RecordRequest(path string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Requests[path]++
	if m.FirstRequestTime.IsZero() {
		m.FirstRequestTime = time.Now()
	}
	m.LastRequestTime = time.Now()
}

// GetMetrics returns a copy of the current metrics
func (m *InMemoryMetrics) GetMetrics() MetricData {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.MetricData
}

func (m *InMemoryMetrics) AddMetricsHandler(mux MetricsHandler) {
	// Add metrics endpoint
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		metrics := m.GetMetrics()
		_ = json.NewEncoder(w).Encode(metrics)
	})
}
