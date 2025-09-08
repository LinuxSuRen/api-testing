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
package server_test

import (
	"context"
	"testing"

	"github.com/linuxsuren/api-testing/pkg/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServer_GenerateSQL(t *testing.T) {
	s := server.NewInMemoryServer("", nil)

	tests := []struct {
		name           string
		request        *server.GenerateSQLRequest
		expectError    bool
		expectedErrMsg string
		expectSuccess  bool
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
			expectSuccess: true,
		},
		{
			name: "empty natural language input",
			request: &server.GenerateSQLRequest{
				NaturalLanguage: "",
			},
			expectError:    true,
			expectedErrMsg: "Natural language input is required",
		},
		{
			name: "with database context",
			request: &server.GenerateSQLRequest{
				NaturalLanguage: "Count active users",
				DatabaseTarget: &server.DatabaseTarget{
					Type:    "mysql",
					Version: "8.0",
					Schemas: []string{"main", "analytics"},
				},
				Options: &server.GenerationOptions{
					IncludeExplanation:  true,
					FormatOutput:        true,
					MaxSuggestions:      3,
					ConfidenceThreshold: 0.8,
				},
				Context: map[string]string{
					"table": "users",
				},
			},
			expectSuccess: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			resp, err := s.GenerateSQL(ctx, tt.request)

			require.NoError(t, err, "GenerateSQL should not return gRPC error")
			require.NotNil(t, resp)

			if tt.expectError {
				require.NotNil(t, resp.Error, "Response should contain error")
				assert.Equal(t, tt.expectedErrMsg, resp.Error.Message)
				assert.Equal(t, server.AIErrorCode_INVALID_INPUT, resp.Error.Code)
			} else if tt.expectSuccess {
				assert.Nil(t, resp.Error, "Response should not contain error")
				assert.NotEmpty(t, resp.GeneratedSql, "Generated SQL should not be empty")
				assert.Greater(t, resp.ConfidenceScore, float32(0), "Confidence score should be positive")
				assert.NotEmpty(t, resp.Explanation, "Explanation should not be empty")
				assert.NotNil(t, resp.Metadata, "Metadata should not be nil")
				assert.NotEmpty(t, resp.Metadata.RequestId, "Request ID should not be empty")
				assert.Greater(t, resp.Metadata.ProcessingTimeMs, float64(0), "Processing time should be positive")
			}
		})
	}
}

func TestServer_ValidateSQL(t *testing.T) {
	s := server.NewInMemoryServer("", nil)

	tests := []struct {
		name        string
		request     *server.ValidateSQLRequest
		expectValid bool
		expectError bool
	}{
		{
			name: "valid SELECT query",
			request: &server.ValidateSQLRequest{
				Sql:          "SELECT * FROM users WHERE active = 1",
				DatabaseType: "postgresql",
			},
			expectValid: true,
		},
		{
			name: "valid INSERT query",
			request: &server.ValidateSQLRequest{
				Sql:          "INSERT INTO users (name, email) VALUES ('John', 'john@example.com')",
				DatabaseType: "mysql",
			},
			expectValid: true,
		},
		{
			name: "valid UPDATE query",
			request: &server.ValidateSQLRequest{
				Sql:          "UPDATE users SET active = 0 WHERE id = 1",
				DatabaseType: "sqlite",
			},
			expectValid: true,
		},
		{
			name: "valid DELETE query",
			request: &server.ValidateSQLRequest{
				Sql:          "DELETE FROM users WHERE inactive_date < NOW() - INTERVAL 1 YEAR",
				DatabaseType: "postgresql",
			},
			expectValid: true,
		},
		{
			name: "empty SQL query",
			request: &server.ValidateSQLRequest{
				Sql: "",
			},
			expectValid: false,
			expectError: true,
		},
		{
			name: "invalid SQL syntax",
			request: &server.ValidateSQLRequest{
				Sql:          "INVALID QUERY SYNTAX",
				DatabaseType: "mysql",
			},
			expectValid: false,
		},
		{
			name: "complex valid query with context",
			request: &server.ValidateSQLRequest{
				Sql:          "SELECT u.name, p.title FROM users u JOIN posts p ON u.id = p.user_id",
				DatabaseType: "postgresql",
				Context: map[string]string{
					"schema": "public",
					"tables": "users,posts",
				},
			},
			expectValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			resp, err := s.ValidateSQL(ctx, tt.request)

			require.NoError(t, err, "ValidateSQL should not return gRPC error")
			require.NotNil(t, resp)

			assert.Equal(t, tt.expectValid, resp.IsValid)

			if tt.expectError {
				assert.NotEmpty(t, resp.Errors, "Should have validation errors")
				assert.Equal(t, "SQL query is required", resp.Errors[0].Message)
				assert.Equal(t, server.ValidationErrorType_SYNTAX_ERROR, resp.Errors[0].Type)
			} else if tt.expectValid {
				assert.Empty(t, resp.Errors, "Valid SQL should have no errors")
				assert.NotEmpty(t, resp.FormattedSql, "Should have formatted SQL")
				assert.NotNil(t, resp.Metadata, "Should have validation metadata")
			} else {
				assert.NotEmpty(t, resp.Errors, "Invalid SQL should have errors")
			}
		})
	}
}

func TestServer_GetAICapabilities(t *testing.T) {
	s := server.NewInMemoryServer("", nil)

	ctx := context.Background()
	resp, err := s.GetAICapabilities(ctx, &server.Empty{})

	require.NoError(t, err)
	require.NotNil(t, resp)

	// Test response structure
	assert.NotEmpty(t, resp.SupportedDatabases, "Should support multiple databases")
	assert.Contains(t, resp.SupportedDatabases, "mysql")
	assert.Contains(t, resp.SupportedDatabases, "postgresql")
	assert.Contains(t, resp.SupportedDatabases, "sqlite")

	assert.NotEmpty(t, resp.Features, "Should have AI features")
	
	// Check for SQL generation feature
	var sqlGenFeature *server.AIFeature
	for _, feature := range resp.Features {
		if feature.Name == "sql_generation" {
			sqlGenFeature = feature
			break
		}
	}
	require.NotNil(t, sqlGenFeature, "Should have sql_generation feature")
	assert.True(t, sqlGenFeature.Enabled)
	assert.NotEmpty(t, sqlGenFeature.Description)
	assert.NotEmpty(t, sqlGenFeature.Parameters)

	// Check for SQL validation feature
	var sqlValFeature *server.AIFeature
	for _, feature := range resp.Features {
		if feature.Name == "sql_validation" {
			sqlValFeature = feature
			break
		}
	}
	require.NotNil(t, sqlValFeature, "Should have sql_validation feature")
	assert.True(t, sqlValFeature.Enabled)

	assert.NotEmpty(t, resp.Version, "Should have version")
	assert.NotEqual(t, server.HealthStatus_HEALTH_STATUS_UNSPECIFIED, resp.Status)
	assert.NotEmpty(t, resp.Limits, "Should have limits")
}

func TestServer_AIErrorHandling(t *testing.T) {
	s := server.NewInMemoryServer("", nil)

	tests := []struct {
		name          string
		testFunc      func(context.Context) error
		expectedError string
	}{
		{
			name: "GenerateSQL with empty input",
			testFunc: func(ctx context.Context) error {
				resp, err := s.GenerateSQL(ctx, &server.GenerateSQLRequest{})
				if err != nil {
					return err
				}
				if resp.Error != nil && resp.Error.Code == server.AIErrorCode_INVALID_INPUT {
					return nil // Expected error
				}
				return assert.AnError // Unexpected response
			},
		},
		{
			name: "ValidateSQL with empty input",
			testFunc: func(ctx context.Context) error {
				resp, err := s.ValidateSQL(ctx, &server.ValidateSQLRequest{})
				if err != nil {
					return err
				}
				if !resp.IsValid && len(resp.Errors) > 0 {
					return nil // Expected validation failure
				}
				return assert.AnError // Unexpected response
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			err := tt.testFunc(ctx)
			assert.NoError(t, err, "Error handling test should pass")
		})
	}
}

func TestServer_AIMethodsIntegration(t *testing.T) {
	s := server.NewInMemoryServer("", nil)
	ctx := context.Background()

	// First, check AI capabilities
	capResp, err := s.GetAICapabilities(ctx, &server.Empty{})
	require.NoError(t, err)
	require.NotNil(t, capResp)

	// Test SQL generation
	genResp, err := s.GenerateSQL(ctx, &server.GenerateSQLRequest{
		NaturalLanguage: "Find all users",
	})
	require.NoError(t, err)
	require.NotNil(t, genResp)
	require.Nil(t, genResp.Error, "Generation should succeed")

	// Test SQL validation with generated SQL
	valResp, err := s.ValidateSQL(ctx, &server.ValidateSQLRequest{
		Sql: genResp.GeneratedSql,
	})
	require.NoError(t, err)
	require.NotNil(t, valResp)

	// The generated SQL should be valid (contains SELECT keyword)
	assert.True(t, valResp.IsValid, "Generated SQL should be valid")
}

func TestFormatSQL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple SELECT",
			input:    "SELECT * FROM users",
			expected: "SELECT *\nFROM users",
		},
		{
			name:     "SELECT with WHERE",
			input:    "SELECT id, name FROM users WHERE active = 1",
			expected: "SELECT id, name\nFROM users\nWHERE active = 1",
		},
		{
			name:     "complex query with multiple clauses",
			input:    "SELECT u.name, COUNT(*) FROM users u WHERE u.active = 1 GROUP BY u.name ORDER BY COUNT(*) DESC",
			expected: "SELECT u.name, COUNT(*)\nFROM users u\nWHERE u.active = 1\nGROUP BY u.name\nORDER BY COUNT(*) DESC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This tests the formatSQL function indirectly through ValidateSQL
			s := server.NewInMemoryServer("", nil)
			ctx := context.Background()

			resp, err := s.ValidateSQL(ctx, &server.ValidateSQLRequest{
				Sql: tt.input,
			})

			require.NoError(t, err)
			require.NotNil(t, resp)

			if resp.IsValid {
				assert.Equal(t, tt.expected, resp.FormattedSql)
			}
		})
	}
}