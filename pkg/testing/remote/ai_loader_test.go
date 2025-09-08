/*
Copyright 2024 API Testing Authors.

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
package remote_test

import (
	"context"
	"testing"

	"github.com/linuxsuren/api-testing/pkg/server"
	"github.com/linuxsuren/api-testing/pkg/testing/remote"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net"
	"time"
)

// MockAILoader implements the AI methods of remote.LoaderServer for testing
type MockAILoader struct {
	remote.UnimplementedLoaderServer
}

func (m *MockAILoader) GenerateSQL(ctx context.Context, req *server.GenerateSQLRequest) (*server.GenerateSQLResponse, error) {
	if req.NaturalLanguage == "" {
		return &server.GenerateSQLResponse{
			Error: &server.AIError{
				Code:    server.AIErrorCode_INVALID_INPUT,
				Message: "Natural language input is required",
				Details: "The natural_language field cannot be empty",
			},
		}, nil
	}

	return &server.GenerateSQLResponse{
		GeneratedSql:    "SELECT * FROM users WHERE created_at >= NOW() - INTERVAL '30 days'",
		ConfidenceScore: 0.95,
		Explanation:     "Generated SQL query based on natural language input",
		Suggestions:     []string{"Consider adding LIMIT clause", "Index on created_at recommended"},
		Metadata: &server.GenerationMetadata{
			RequestId:        "mock-req-123",
			ProcessingTimeMs: 100.0,
			ModelUsed:        "mock-ai-model",
			TokenCount:       25,
			Timestamp:        timestamppb.New(time.Now()),
		},
	}, nil
}

func (m *MockAILoader) ValidateSQL(ctx context.Context, req *server.ValidateSQLRequest) (*server.ValidateSQLResponse, error) {
	if req.Sql == "" {
		return &server.ValidateSQLResponse{
			IsValid: false,
			Errors: []*server.ValidationError{
				{
					Message: "SQL query is required",
					Line:    1,
					Column:  1,
					Type:    server.ValidationErrorType_SYNTAX_ERROR,
				},
			},
		}, nil
	}

	// Simple validation: check if it contains SELECT
	if req.Sql == "SELECT * FROM users" {
		return &server.ValidateSQLResponse{
			IsValid:      true,
			FormattedSql: "SELECT *\nFROM users",
			Metadata: &server.ValidationMetadata{
				ValidatorVersion:  "mock-validator-1.0",
				ValidationTimeMs: 10.0,
				Timestamp:        timestamppb.New(time.Now()),
			},
		}, nil
	}

	return &server.ValidateSQLResponse{
		IsValid: false,
		Errors: []*server.ValidationError{
			{
				Message: "Invalid SQL syntax",
				Line:    1,
				Column:  1,
				Type:    server.ValidationErrorType_SYNTAX_ERROR,
			},
		},
		Warnings: []string{"Consider using proper SQL syntax"},
	}, nil
}

func (m *MockAILoader) GetAICapabilities(ctx context.Context, req *server.Empty) (*server.AICapabilitiesResponse, error) {
	return &server.AICapabilitiesResponse{
		SupportedDatabases: []string{"mysql", "postgresql", "sqlite"},
		Features: []*server.AIFeature{
			{
				Name:        "sql_generation",
				Enabled:     true,
				Description: "Generate SQL from natural language",
				Parameters: map[string]string{
					"max_complexity": "high",
					"model":          "mock-ai-model",
				},
			},
		},
		Version: "mock-1.0.0",
		Status:  server.HealthStatus_HEALTHY,
		Limits: map[string]string{
			"max_requests_per_minute": "100",
			"max_query_length":        "2000",
		},
	}, nil
}

func setupMockAIServer(t *testing.T) (*grpc.ClientConn, func()) {
	buffer := 101024 * 1024
	listener := bufconn.Listen(buffer)

	server := grpc.NewServer()
	remote.RegisterLoaderServer(server, &MockAILoader{})

	go func() {
		if err := server.Serve(listener); err != nil {
			t.Logf("Server error: %v", err)
		}
	}()

	conn, err := grpc.DialContext(context.Background(), "",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return listener.Dial()
		}),
		grpc.WithInsecure())
	require.NoError(t, err)

	closeFunc := func() {
		err := listener.Close()
		if err != nil {
			t.Logf("Error closing listener: %v", err)
		}
		server.Stop()
	}

	return conn, closeFunc
}

func TestAILoader_GenerateSQL(t *testing.T) {
	conn, cleanup := setupMockAIServer(t)
	defer cleanup()

	client := remote.NewLoaderClient(conn)

	tests := []struct {
		name           string
		request        *server.GenerateSQLRequest
		expectError    bool
		expectedSQL    string
		expectedErrMsg string
	}{
		{
			name: "successful generation",
			request: &server.GenerateSQLRequest{
				NaturalLanguage: "Find all users created in the last 30 days",
				DatabaseTarget: &server.DatabaseTarget{
					Type:    "postgresql",
					Version: "13.0",
				},
			},
			expectError: false,
			expectedSQL: "SELECT * FROM users WHERE created_at >= NOW() - INTERVAL '30 days'",
		},
		{
			name: "empty natural language input",
			request: &server.GenerateSQLRequest{
				NaturalLanguage: "",
			},
			expectError:    false, // Server returns error in response, not as gRPC error
			expectedErrMsg: "Natural language input is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			resp, err := client.GenerateSQL(ctx, tt.request)

			if tt.expectError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)

			if tt.expectedErrMsg != "" {
				require.NotNil(t, resp.Error)
				assert.Equal(t, tt.expectedErrMsg, resp.Error.Message)
			} else {
				assert.Equal(t, tt.expectedSQL, resp.GeneratedSql)
				assert.Greater(t, resp.ConfidenceScore, float32(0))
				assert.NotEmpty(t, resp.Explanation)
				assert.NotNil(t, resp.Metadata)
			}
		})
	}
}

func TestAILoader_ValidateSQL(t *testing.T) {
	conn, cleanup := setupMockAIServer(t)
	defer cleanup()

	client := remote.NewLoaderClient(conn)

	tests := []struct {
		name        string
		request     *server.ValidateSQLRequest
		expectValid bool
		expectError bool
	}{
		{
			name: "valid SQL",
			request: &server.ValidateSQLRequest{
				Sql:          "SELECT * FROM users",
				DatabaseType: "postgresql",
			},
			expectValid: true,
			expectError: false,
		},
		{
			name: "invalid SQL",
			request: &server.ValidateSQLRequest{
				Sql:          "INVALID QUERY",
				DatabaseType: "mysql",
			},
			expectValid: false,
			expectError: false,
		},
		{
			name: "empty SQL",
			request: &server.ValidateSQLRequest{
				Sql: "",
			},
			expectValid: false,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			resp, err := client.ValidateSQL(ctx, tt.request)

			if tt.expectError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
			assert.Equal(t, tt.expectValid, resp.IsValid)

			if tt.expectValid {
				assert.NotEmpty(t, resp.FormattedSql)
				assert.NotNil(t, resp.Metadata)
			} else {
				assert.NotEmpty(t, resp.Errors)
			}
		})
	}
}

func TestAILoader_GetAICapabilities(t *testing.T) {
	conn, cleanup := setupMockAIServer(t)
	defer cleanup()

	client := remote.NewLoaderClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.GetAICapabilities(ctx, &server.Empty{})

	require.NoError(t, err)
	require.NotNil(t, resp)

	// Test response structure
	assert.NotEmpty(t, resp.SupportedDatabases)
	assert.Contains(t, resp.SupportedDatabases, "mysql")
	assert.Contains(t, resp.SupportedDatabases, "postgresql")

	assert.NotEmpty(t, resp.Features)
	assert.Equal(t, "sql_generation", resp.Features[0].Name)
	assert.True(t, resp.Features[0].Enabled)

	assert.Equal(t, "mock-1.0.0", resp.Version)
	assert.Equal(t, server.HealthStatus_HEALTHY, resp.Status)
	assert.NotEmpty(t, resp.Limits)
}

func TestAILoader_ErrorHandling(t *testing.T) {
	// Test connection error handling
	conn, err := grpc.Dial("invalid:address", grpc.WithInsecure())
	require.NoError(t, err)
	defer conn.Close()

	client := remote.NewLoaderClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, err = client.GenerateSQL(ctx, &server.GenerateSQLRequest{
		NaturalLanguage: "test query",
	})

	require.Error(t, err)
	
	// Check that it's a gRPC error
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.Unavailable, st.Code())
}

func TestAILoader_ContextCancellation(t *testing.T) {
	conn, cleanup := setupMockAIServer(t)
	defer cleanup()

	client := remote.NewLoaderClient(conn)

	// Create a context that's already cancelled
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, err := client.GenerateSQL(ctx, &server.GenerateSQLRequest{
		NaturalLanguage: "Find all users",
	})

	require.Error(t, err)
	
	// Check that it's a cancellation error
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.Canceled, st.Code())
}

func TestAILoader_Timeout(t *testing.T) {
	conn, cleanup := setupMockAIServer(t)
	defer cleanup()

	client := remote.NewLoaderClient(conn)

	// Use very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	// Wait a bit to ensure timeout
	time.Sleep(1 * time.Millisecond)

	_, err := client.GenerateSQL(ctx, &server.GenerateSQLRequest{
		NaturalLanguage: "Find all users",
	})

	require.Error(t, err)
	
	// Check that it's a deadline exceeded error
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.DeadlineExceeded, st.Code())
}