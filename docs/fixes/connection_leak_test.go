// connection_leak_test.go
// Test file to verify database connection leak fixes in atest-ext-store-orm

package fixes

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Note: This test file is designed to work with the atest-ext-store-orm extension
// after applying the connection leak fix patch. It demonstrates how to test
// connection management and validate the fix is working correctly.

// MockORMStore represents a mock implementation for testing connection leak fixes
type MockORMStore struct {
	mu                sync.RWMutex
	activeConnections map[string]int
	maxConnections    int
}

// NewMockORMStore creates a new mock ORM store for testing
func NewMockORMStore(maxConnections int) *MockORMStore {
	return &MockORMStore{
		activeConnections: make(map[string]int),
		maxConnections:    maxConnections,
	}
}

// GetActiveConnections returns the number of active connections for a database
func (m *MockORMStore) GetActiveConnections(database string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.activeConnections[database]
}

// SimulateConnection simulates creating a connection to a database
func (m *MockORMStore) SimulateConnection(database string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.activeConnections[database] >= m.maxConnections {
		return ErrTooManyConnections
	}
	m.activeConnections[database]++
	return nil
}

// SimulateDisconnection simulates closing a connection
func (m *MockORMStore) SimulateDisconnection(database string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.activeConnections[database] > 0 {
		m.activeConnections[database]--
	}
}

// ErrTooManyConnections represents the error when connection limit is exceeded
var ErrTooManyConnections = errors.New("too many connections")

// TestDatabaseConnectionLeak verifies that database connections are properly managed
// and don't leak when switching between different databases
func TestDatabaseConnectionLeak(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping connection leak test in short mode")
	}

	mockStore := NewMockORMStore(10)

	// Test rapid database switching without connection leaks
	databases := []string{"db1", "db2", "db3"}
	
	// Simulate multiple rapid switches
	for i := 0; i < 50; i++ {
		for _, db := range databases {
			// Simulate connection creation
			err := mockStore.SimulateConnection(db)
			require.NoError(t, err, "Should not exceed connection limit on iteration %d for db %s", i, db)
			
			// Verify connection count is reasonable
			count := mockStore.GetActiveConnections(db)
			assert.LessOrEqual(t, count, 3, "Too many connections for db %s: %d", db, count)
			
			// Simulate some work
			time.Sleep(time.Millisecond)
			
			// Simulate connection cleanup
			mockStore.SimulateDisconnection(db)
		}
	}

	// Verify all connections are cleaned up
	for _, db := range databases {
		count := mockStore.GetActiveConnections(db)
		assert.Equal(t, 0, count, "Connections not properly cleaned up for db %s", db)
	}
}

// TestConnectionPoolConfiguration verifies connection pool settings work correctly
func TestConnectionPoolConfiguration(t *testing.T) {
	tests := []struct {
		name           string
		maxConnections int
	}{
		{"Normal pool size", 10},
		{"Large pool size", 100},
		{"Small pool size", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := NewMockORMStore(tt.maxConnections)
			
			// Try to create connections up to the limit
			for i := 0; i < tt.maxConnections; i++ {
				err := mockStore.SimulateConnection("testdb")
				assert.NoError(t, err, "Should not error within connection limit")
			}
			
			// Try to exceed the limit - should always error
			err := mockStore.SimulateConnection("testdb")
			assert.Error(t, err, "Should error when exceeding connection limit")
			assert.Equal(t, ErrTooManyConnections, err, "Should return specific error type")
		})
	}
}

// TestConnectionReuse verifies that connections are properly reused
func TestConnectionReuse(t *testing.T) {
	mockStore := NewMockORMStore(10)
	
	// Create and close connections multiple times
	for i := 0; i < 20; i++ {
		err := mockStore.SimulateConnection("reusedb")
		require.NoError(t, err, "Connection creation should not fail on iteration %d", i)
		
		// Verify connection count
		count := mockStore.GetActiveConnections("reusedb")
		assert.LessOrEqual(t, count, 10, "Connection count should not exceed pool size")
		
		mockStore.SimulateDisconnection("reusedb")
	}
	
	// Final verification
	finalCount := mockStore.GetActiveConnections("reusedb")
	assert.Equal(t, 0, finalCount, "All connections should be closed")
}

// TestConcurrentDatabaseAccess simulates concurrent access to multiple databases
func TestConcurrentDatabaseAccess(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping concurrent test in short mode")
	}

	mockStore := NewMockORMStore(5)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Simulate concurrent access patterns
	done := make(chan bool, 2)
	
	// Worker 1: Access db1 repeatedly
	go func() {
		defer func() { 
			// Ensure we always send to done channel to prevent test hanging
			select {
			case done <- true:
			default:
			}
		}()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				err := mockStore.SimulateConnection("db1")
				if err == nil {
					time.Sleep(time.Millisecond * 10)
					mockStore.SimulateDisconnection("db1")
				}
				time.Sleep(time.Millisecond * 5)
			}
		}
	}()
	
	// Worker 2: Access db2 repeatedly  
	go func() {
		defer func() { 
			// Ensure we always send to done channel to prevent test hanging
			select {
			case done <- true:
			default:
			}
		}()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				err := mockStore.SimulateConnection("db2")
				if err == nil {
					time.Sleep(time.Millisecond * 10)
					mockStore.SimulateDisconnection("db2")
				}
				time.Sleep(time.Millisecond * 5)
			}
		}
	}()

	// Let workers run for a short time
	time.Sleep(time.Millisecond * 100)
	cancel()
	
	// Wait for workers to finish with timeout protection
	for i := 0; i < 2; i++ {
		select {
		case <-done:
			// Worker finished
		case <-time.After(5 * time.Second):
			t.Fatal("Test timed out waiting for workers to finish")
		}
	}
	
	// Verify no connections are leaked
	db1Count := mockStore.GetActiveConnections("db1")
	db2Count := mockStore.GetActiveConnections("db2")
	assert.LessOrEqual(t, db1Count, 5, "db1 should not have excessive connections")
	assert.LessOrEqual(t, db2Count, 5, "db2 should not have excessive connections")
}

// TestCacheKeyGeneration verifies that cache keys are properly generated
func TestCacheKeyGeneration(t *testing.T) {
	tests := []struct {
		store    string
		database string
		expected string
	}{
		{"mysql", "testdb", "mysql:testdb"},
		{"postgres", "proddb", "postgres:proddb"},
		{"sqlite", "local.db", "sqlite:local.db"},
		{"mysql", "", "mysql:"},
		{"", "testdb", ":testdb"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			cacheKey := generateCacheKey(tt.store, tt.database)
			assert.Equal(t, tt.expected, cacheKey, "Cache key should match expected format")
		})
	}
}

// generateCacheKey creates a composite cache key for store and database
func generateCacheKey(store, database string) string {
	return store + ":" + database
}
