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
	"testing"
	"time"

	"github.com/linuxsuren/api-testing/pkg/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestGenerateSQLRequest_Serialization(t *testing.T) {
	tests := []struct {
		name string
		req  *server.GenerateSQLRequest
	}{
		{
			name: "complete request",
			req: &server.GenerateSQLRequest{
				NaturalLanguage: "Find all users created in the last 30 days",
				DatabaseTarget: &server.DatabaseTarget{
					Type:    "postgresql",
					Version: "13.0",
					Schemas: []string{"public", "users"},
					Metadata: map[string]string{
						"host": "localhost",
						"port": "5432",
					},
				},
				Options: &server.GenerationOptions{
					IncludeExplanation:   true,
					FormatOutput:         true,
					MaxSuggestions:       3,
					ConfidenceThreshold:  0.8,
					EnableOptimization:   true,
				},
				Context: map[string]string{
					"table":  "users",
					"schema": "public",
				},
			},
		},
		{
			name: "minimal request",
			req: &server.GenerateSQLRequest{
				NaturalLanguage: "SELECT * FROM users",
			},
		},
		{
			name: "empty request",
			req:  &server.GenerateSQLRequest{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test serialization
			data, err := proto.Marshal(tt.req)
			require.NoError(t, err)
			// Empty messages may serialize to empty bytes - this is expected in protobuf
			if tt.name != "empty request" {
				require.NotEmpty(t, data)
			}

			// Test deserialization
			unmarshaled := &server.GenerateSQLRequest{}
			err = proto.Unmarshal(data, unmarshaled)
			require.NoError(t, err)

			// Test equality
			assert.True(t, proto.Equal(tt.req, unmarshaled))
		})
	}
}

func TestGenerateSQLResponse_Serialization(t *testing.T) {
	tests := []struct {
		name string
		resp *server.GenerateSQLResponse
	}{
		{
			name: "successful response",
			resp: &server.GenerateSQLResponse{
				GeneratedSql:     "SELECT * FROM users WHERE created_at >= NOW() - INTERVAL '30 days'",
				ConfidenceScore:  0.95,
				Explanation:      "This query finds all users created in the last 30 days",
				Suggestions:      []string{"Add LIMIT clause", "Consider indexing created_at"},
				Metadata: &server.GenerationMetadata{
					RequestId:        "req-123",
					ProcessingTimeMs: 150.5,
					ModelUsed:        "gpt-4",
					TokenCount:       45,
					Timestamp:        timestamppb.New(time.Now()),
				},
			},
		},
		{
			name: "error response",
			resp: &server.GenerateSQLResponse{
				Error: &server.AIError{
					Code:    server.AIErrorCode_INVALID_INPUT,
					Message: "Natural language input is required",
					Details: "The natural_language field cannot be empty",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test serialization
			data, err := proto.Marshal(tt.resp)
			require.NoError(t, err)
			require.NotEmpty(t, data)

			// Test deserialization
			unmarshaled := &server.GenerateSQLResponse{}
			err = proto.Unmarshal(data, unmarshaled)
			require.NoError(t, err)

			// Test equality
			assert.True(t, proto.Equal(tt.resp, unmarshaled))
		})
	}
}

func TestValidateSQLRequest_Serialization(t *testing.T) {
	tests := []struct {
		name string
		req  *server.ValidateSQLRequest
	}{
		{
			name: "complete request",
			req: &server.ValidateSQLRequest{
				Sql:          "SELECT * FROM users WHERE id = ?",
				DatabaseType: "mysql",
				Context: map[string]string{
					"version": "8.0",
					"schema":  "main",
				},
			},
		},
		{
			name: "minimal request",
			req: &server.ValidateSQLRequest{
				Sql: "SELECT 1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test serialization
			data, err := proto.Marshal(tt.req)
			require.NoError(t, err)
			require.NotEmpty(t, data)

			// Test deserialization
			unmarshaled := &server.ValidateSQLRequest{}
			err = proto.Unmarshal(data, unmarshaled)
			require.NoError(t, err)

			// Test equality
			assert.True(t, proto.Equal(tt.req, unmarshaled))
		})
	}
}

func TestValidateSQLResponse_Serialization(t *testing.T) {
	tests := []struct {
		name string
		resp *server.ValidateSQLResponse
	}{
		{
			name: "valid SQL response",
			resp: &server.ValidateSQLResponse{
				IsValid:      true,
				FormattedSql: "SELECT *\nFROM users\nWHERE id = ?",
				Metadata: &server.ValidationMetadata{
					ValidatorVersion:  "1.0.0",
					ValidationTimeMs: 25.0,
					Timestamp:        timestamppb.New(time.Now()),
				},
			},
		},
		{
			name: "invalid SQL response",
			resp: &server.ValidateSQLResponse{
				IsValid: false,
				Errors: []*server.ValidationError{
					{
						Message: "Syntax error near 'FROM'",
						Line:    1,
						Column:  15,
						Type:    server.ValidationErrorType_SYNTAX_ERROR,
					},
				},
				Warnings: []string{"Missing table alias"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test serialization
			data, err := proto.Marshal(tt.resp)
			require.NoError(t, err)
			require.NotEmpty(t, data)

			// Test deserialization
			unmarshaled := &server.ValidateSQLResponse{}
			err = proto.Unmarshal(data, unmarshaled)
			require.NoError(t, err)

			// Test equality
			assert.True(t, proto.Equal(tt.resp, unmarshaled))
		})
	}
}

func TestAICapabilitiesResponse_Serialization(t *testing.T) {
	resp := &server.AICapabilitiesResponse{
		SupportedDatabases: []string{"mysql", "postgresql", "sqlite"},
		Features: []*server.AIFeature{
			{
				Name:        "sql_generation",
				Enabled:     true,
				Description: "Generate SQL from natural language",
				Parameters: map[string]string{
					"max_complexity": "high",
					"supported_joins": "true",
				},
			},
			{
				Name:        "query_optimization",
				Enabled:     false,
				Description: "Optimize existing SQL queries",
			},
		},
		Version: "1.0.0",
		Status:  server.HealthStatus_HEALTHY,
		Limits: map[string]string{
			"max_requests_per_minute": "60",
			"max_query_length":        "1000",
		},
	}

	// Test serialization
	data, err := proto.Marshal(resp)
	require.NoError(t, err)
	require.NotEmpty(t, data)

	// Test deserialization
	unmarshaled := &server.AICapabilitiesResponse{}
	err = proto.Unmarshal(data, unmarshaled)
	require.NoError(t, err)

	// Test equality
	assert.True(t, proto.Equal(resp, unmarshaled))
}

func TestDataQuery_AIExtensions(t *testing.T) {
	tests := []struct {
		name string
		req  *server.DataQuery
	}{
		{
			name: "AI query with extensions",
			req: &server.DataQuery{
				Type:            "ai",
				Key:             "test-key",
				Sql:             "",
				Offset:          0,
				Limit:           10,
				NaturalLanguage: "Find all active users",
				DatabaseType:    "postgresql",
				ExplainQuery:    true,
				AiContext: map[string]string{
					"table":  "users",
					"schema": "public",
				},
			},
		},
		{
			name: "traditional query (backward compatibility)",
			req: &server.DataQuery{
				Type:   "sql",
				Key:    "test-key",
				Sql:    "SELECT * FROM users",
				Offset: 0,
				Limit:  10,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test serialization
			data, err := proto.Marshal(tt.req)
			require.NoError(t, err)
			require.NotEmpty(t, data)

			// Test deserialization
			unmarshaled := &server.DataQuery{}
			err = proto.Unmarshal(data, unmarshaled)
			require.NoError(t, err)

			// Test equality
			assert.True(t, proto.Equal(tt.req, unmarshaled))
		})
	}
}

func TestDataQueryResult_AIExtensions(t *testing.T) {
	tests := []struct {
		name   string
		result *server.DataQueryResult
	}{
		{
			name: "AI query result with processing info",
			result: &server.DataQueryResult{
				Data: []*server.Pair{
					{Key: "id", Value: "1"},
					{Key: "name", Value: "John"},
				},
				Meta: &server.DataMeta{
					Databases:       []string{"testdb"},
					Tables:          []string{"users"},
					CurrentDatabase: "testdb",
					Duration:        "150ms",
				},
				AiInfo: &server.AIProcessingInfo{
					RequestId:        "ai-req-123",
					ProcessingTimeMs: 150.5,
					ModelUsed:        "gpt-4",
					ConfidenceScore:  0.92,
					DebugInfo:        []string{"Used table schema", "Applied query optimization"},
				},
			},
		},
		{
			name: "traditional query result (backward compatibility)",
			result: &server.DataQueryResult{
				Data: []*server.Pair{
					{Key: "count", Value: "42"},
				},
				Meta: &server.DataMeta{
					Duration: "25ms",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test serialization
			data, err := proto.Marshal(tt.result)
			require.NoError(t, err)
			require.NotEmpty(t, data)

			// Test deserialization
			unmarshaled := &server.DataQueryResult{}
			err = proto.Unmarshal(data, unmarshaled)
			require.NoError(t, err)

			// Test equality
			assert.True(t, proto.Equal(tt.result, unmarshaled))
		})
	}
}

func TestAIErrorCode_EnumValues(t *testing.T) {
	tests := []struct {
		name     string
		code     server.AIErrorCode
		expected string
	}{
		{
			name:     "unspecified",
			code:     server.AIErrorCode_AI_ERROR_CODE_UNSPECIFIED,
			expected: "AI_ERROR_CODE_UNSPECIFIED",
		},
		{
			name:     "invalid input",
			code:     server.AIErrorCode_INVALID_INPUT,
			expected: "INVALID_INPUT",
		},
		{
			name:     "model unavailable",
			code:     server.AIErrorCode_MODEL_UNAVAILABLE,
			expected: "MODEL_UNAVAILABLE",
		},
		{
			name:     "rate limited",
			code:     server.AIErrorCode_RATE_LIMITED,
			expected: "RATE_LIMITED",
		},
		{
			name:     "processing error",
			code:     server.AIErrorCode_PROCESSING_ERROR,
			expected: "PROCESSING_ERROR",
		},
		{
			name:     "configuration error",
			code:     server.AIErrorCode_CONFIGURATION_ERROR,
			expected: "CONFIGURATION_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.code.String())
		})
	}
}

func TestValidationErrorType_EnumValues(t *testing.T) {
	tests := []struct {
		name     string
		errType  server.ValidationErrorType
		expected string
	}{
		{
			name:     "unspecified",
			errType:  server.ValidationErrorType_VALIDATION_ERROR_TYPE_UNSPECIFIED,
			expected: "VALIDATION_ERROR_TYPE_UNSPECIFIED",
		},
		{
			name:     "syntax error",
			errType:  server.ValidationErrorType_SYNTAX_ERROR,
			expected: "SYNTAX_ERROR",
		},
		{
			name:     "semantic error",
			errType:  server.ValidationErrorType_SEMANTIC_ERROR,
			expected: "SEMANTIC_ERROR",
		},
		{
			name:     "performance warning",
			errType:  server.ValidationErrorType_PERFORMANCE_WARNING,
			expected: "PERFORMANCE_WARNING",
		},
		{
			name:     "security warning",
			errType:  server.ValidationErrorType_SECURITY_WARNING,
			expected: "SECURITY_WARNING",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.errType.String())
		})
	}
}

func TestHealthStatus_EnumValues(t *testing.T) {
	tests := []struct {
		name     string
		status   server.HealthStatus
		expected string
	}{
		{
			name:     "unspecified",
			status:   server.HealthStatus_HEALTH_STATUS_UNSPECIFIED,
			expected: "HEALTH_STATUS_UNSPECIFIED",
		},
		{
			name:     "healthy",
			status:   server.HealthStatus_HEALTHY,
			expected: "HEALTHY",
		},
		{
			name:     "degraded",
			status:   server.HealthStatus_DEGRADED,
			expected: "DEGRADED",
		},
		{
			name:     "unhealthy",
			status:   server.HealthStatus_UNHEALTHY,
			expected: "UNHEALTHY",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.status.String())
		})
	}
}