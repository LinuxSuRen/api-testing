package server

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

// AIManager manages AI plugin connections and operations
type AIManager struct {
	client     *AIClient
	connected  bool
	mutex      sync.RWMutex
	address    string
	retryCount int
	maxRetries int
}

// NewAIManager creates a new AI manager instance
func NewAIManager() *AIManager {
	return &AIManager{
		address:    "localhost:50052", // Default AI plugin address
		maxRetries: 3,
		retryCount: 0,
	}
}

// Initialize initializes the AI manager and establishes connection to AI plugin
func (m *AIManager) Initialize() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Try to connect to AI plugin
	client, err := NewAIClient(m.address)
	if err != nil {
		log.Printf("Failed to connect to AI plugin at %s: %v", m.address, err)
		m.connected = false
		return err
	}

	m.client = client
	m.connected = true
	m.retryCount = 0
	log.Printf("Successfully connected to AI plugin at %s", m.address)
	return nil
}

// IsConnected returns whether the AI plugin is connected
func (m *AIManager) IsConnected() bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.connected
}

// ConvertNLToSQL converts natural language to SQL using the AI plugin
func (m *AIManager) ConvertNLToSQL(ctx context.Context, input string, context map[string]interface{}) (*AIResponse, error) {
	m.mutex.RLock()
	if !m.connected || m.client == nil {
		m.mutex.RUnlock()
		return m.getFallbackResponse("convert_nl_to_sql", input), nil
	}
	client := m.client
	m.mutex.RUnlock()

	return client.ConvertNLToSQL(ctx, input, context)
}

// GenerateTestCase generates test cases using the AI plugin
func (m *AIManager) GenerateTestCase(ctx context.Context, apiSpec string, requirements string) (*AIResponse, error) {
	m.mutex.RLock()
	if !m.connected || m.client == nil {
		m.mutex.RUnlock()
		return m.getFallbackResponse("generate_test_case", requirements), nil
	}
	client := m.client
	m.mutex.RUnlock()

	return client.GenerateTestCase(ctx, apiSpec, requirements)
}

// OptimizeQuery optimizes SQL queries using the AI plugin
func (m *AIManager) OptimizeQuery(ctx context.Context, sqlQuery string, performance map[string]interface{}) (*AIResponse, error) {
	m.mutex.RLock()
	if !m.connected || m.client == nil {
		m.mutex.RUnlock()
		return m.getFallbackResponse("optimize_query", sqlQuery), nil
	}
	client := m.client
	m.mutex.RUnlock()

	return client.OptimizeQuery(ctx, sqlQuery, performance)
}

// HealthCheck checks the health of the AI plugin
func (m *AIManager) HealthCheck(ctx context.Context) error {
	m.mutex.RLock()
	if !m.connected || m.client == nil {
		m.mutex.RUnlock()
		return fmt.Errorf("AI plugin not connected")
	}
	client := m.client
	m.mutex.RUnlock()

	return client.HealthCheck(ctx)
}

// Reconnect attempts to reconnect to the AI plugin
func (m *AIManager) Reconnect() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.retryCount >= m.maxRetries {
		return fmt.Errorf("max retry attempts (%d) exceeded", m.maxRetries)
	}

	// Close existing connection if any
	if m.client != nil {
		m.client.Close()
	}

	// Wait before retry
	time.Sleep(time.Duration(m.retryCount+1) * time.Second)

	// Try to reconnect
	client, err := NewAIClient(m.address)
	if err != nil {
		m.retryCount++
		m.connected = false
		return fmt.Errorf("reconnection attempt %d failed: %w", m.retryCount, err)
	}

	m.client = client
	m.connected = true
	m.retryCount = 0
	log.Printf("Successfully reconnected to AI plugin at %s", m.address)
	return nil
}

// Close closes the AI manager and its connections
func (m *AIManager) Close() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.client != nil {
		err := m.client.Close()
		m.client = nil
		m.connected = false
		return err
	}
	return nil
}

// SetAddress sets the AI plugin address
func (m *AIManager) SetAddress(address string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.address = address
}

// GetAddress returns the current AI plugin address
func (m *AIManager) GetAddress() string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.address
}

// getFallbackResponse returns a fallback response when AI plugin is not available
func (m *AIManager) getFallbackResponse(requestType, input string) *AIResponse {
	switch requestType {
	case "convert_nl_to_sql":
		return &AIResponse{
			Success: true,
			Result:  "SELECT * FROM users WHERE status = 'active' ORDER BY created_at DESC LIMIT 10; -- Fallback response",
			Error:   "",
			Meta: map[string]interface{}{
				"fallback": true,
				"reason":   "AI plugin not available",
				"input":    input,
			},
		}
	case "generate_test_case":
		return &AIResponse{
			Success: true,
			Result:  `{"name": "Test API Endpoint", "method": "GET", "url": "/api/test", "expected_status": 200}`,
			Error:   "",
			Meta: map[string]interface{}{
				"fallback": true,
				"reason":   "AI plugin not available",
				"input":    input,
			},
		}
	case "optimize_query":
		return &AIResponse{
			Success: true,
			Result:  input + " -- Query optimization not available (fallback mode)",
			Error:   "",
			Meta: map[string]interface{}{
				"fallback": true,
				"reason":   "AI plugin not available",
				"input":    input,
			},
		}
	default:
		return &AIResponse{
			Success: false,
			Result:  "",
			Error:   "Unknown request type",
			Meta: map[string]interface{}{
				"fallback": true,
				"reason":   "Unknown request type",
				"input":    input,
			},
		}
	}
}
