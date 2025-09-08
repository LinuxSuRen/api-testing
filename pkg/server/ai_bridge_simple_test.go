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
package server

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// formatSQL provides basic SQL formatting for testing
func formatSQL(sql string) string {
	formatted := strings.ReplaceAll(sql, " FROM ", "\nFROM ")
	formatted = strings.ReplaceAll(formatted, " WHERE ", "\nWHERE ")
	formatted = strings.ReplaceAll(formatted, " ORDER BY ", "\nORDER BY ")
	formatted = strings.ReplaceAll(formatted, " GROUP BY ", "\nGROUP BY ")
	formatted = strings.ReplaceAll(formatted, " HAVING ", "\nHAVING ")
	return strings.TrimSpace(formatted)
}

func TestAIPluginBridge_BasicFunctionality(t *testing.T) {
	bridge := NewAIPluginBridge()
	
	// Test default client works
	client := bridge.GetPlugin("generate_sql")
	assert.NotNil(t, client)
	
	// Test plugin registration
	mockClient := NewExtendedMockAIPluginClient(true)
	bridge.RegisterPlugin("test-plugin", mockClient)
	
	plugins := bridge.GetAllPlugins()
	assert.Len(t, plugins, 1)
	assert.Equal(t, "test-plugin", plugins[0].ID)
}

func TestMessageTransformer_DataQueryTransformations(t *testing.T) {
	transformer := &MessageTransformer{}
	
	t.Run("Transform to GenerateSQLRequest", func(t *testing.T) {
		query := &DataQuery{
			Type:            "ai",
			NaturalLanguage: "Find all users",
			DatabaseType:    "postgresql",
			ExplainQuery:    true,
			AiContext: map[string]string{
				"table": "users",
			},
		}
		
		req := transformer.TransformDataQueryToGenerateSQL(query)
		
		assert.Equal(t, "Find all users", req.NaturalLanguage)
		assert.Equal(t, "postgresql", req.DatabaseTarget.Type)
		assert.True(t, req.Options.IncludeExplanation)
		assert.True(t, req.Options.FormatOutput)
		assert.Equal(t, "users", req.Context["table"])
	})
	
	t.Run("Transform to ValidateSQLRequest", func(t *testing.T) {
		query := &DataQuery{
			Type:         "ai",
			Sql:          "SELECT * FROM users",
			DatabaseType: "mysql",
			AiContext: map[string]string{
				"version": "8.0",
			},
		}
		
		req := transformer.TransformDataQueryToValidateSQL(query)
		
		assert.Equal(t, "SELECT * FROM users", req.Sql)
		assert.Equal(t, "mysql", req.DatabaseType)
		assert.Equal(t, "8.0", req.Context["version"])
	})
}

func TestMessageTransformer_ResponseTransformations(t *testing.T) {
	transformer := &MessageTransformer{}
	
	t.Run("Transform GenerateSQLResponse", func(t *testing.T) {
		resp := &GenerateSQLResponse{
			GeneratedSql:    "SELECT * FROM users WHERE active = 1",
			ConfidenceScore: 0.95,
			Explanation:     "Query to find active users",
			Suggestions:     []string{"Add LIMIT", "Add INDEX"},
		}
		
		result := transformer.TransformGenerateSQLToDataQueryResult(resp)
		
		require.NotNil(t, result)
		require.NotEmpty(t, result.Data)
		
		dataMap := make(map[string]string)
		for _, pair := range result.Data {
			dataMap[pair.Key] = pair.Value
		}
		
		assert.Equal(t, "SELECT * FROM users WHERE active = 1", dataMap["generated_sql"])
		assert.Equal(t, "0.95", dataMap["confidence_score"])
		assert.Equal(t, "Query to find active users", dataMap["explanation"])
		
		// Check suggestions are in items
		assert.Len(t, result.Items, 2)
	})
	
	t.Run("Transform ValidateSQLResponse", func(t *testing.T) {
		resp := &ValidateSQLResponse{
			IsValid:      true,
			FormattedSql: "SELECT *\nFROM users",
			Warnings:     []string{"Consider adding LIMIT"},
		}
		
		result := transformer.TransformValidateSQLToDataQueryResult(resp)
		
		require.NotNil(t, result)
		require.NotEmpty(t, result.Data)
		
		dataMap := make(map[string]string)
		for _, pair := range result.Data {
			dataMap[pair.Key] = pair.Value
		}
		
		assert.Equal(t, "true", dataMap["is_valid"])
		assert.Equal(t, "SELECT *\nFROM users", dataMap["formatted_sql"])
		
		// Check warning is in items
		assert.Len(t, result.Items, 1)
	})
}

func TestMessageTransformer_Validation(t *testing.T) {
	transformer := &MessageTransformer{}
	
	tests := []struct {
		name          string
		query         *DataQuery
		expectError   bool
		errorContains string
	}{
		{
			name: "valid AI query with natural language",
			query: &DataQuery{
				Type:            "ai",
				NaturalLanguage: "Find users",
				DatabaseType:    "postgresql",
			},
			expectError: false,
		},
		{
			name: "valid AI query with SQL",
			query: &DataQuery{
				Type: "ai",
				Sql:  "SELECT * FROM users",
			},
			expectError: false,
		},
		{
			name: "invalid - empty query",
			query: &DataQuery{
				Type: "ai",
			},
			expectError:   true,
			errorContains: "must have either natural_language or sql field",
		},
		{
			name: "invalid database type",
			query: &DataQuery{
				Type:            "ai",
				NaturalLanguage: "Find users",
				DatabaseType:    "unsupported_db",
			},
			expectError:   true,
			errorContains: "unsupported database type",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := transformer.ValidateAIQuery(tt.query)
			
			if tt.expectError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestExtendedMockAIPluginClient_Operations(t *testing.T) {
	client := NewExtendedMockAIPluginClient(true)
	ctx := context.Background()
	
	t.Run("GenerateSQL", func(t *testing.T) {
		req := &GenerateSQLRequest{
			NaturalLanguage: "Find all users",
		}
		
		resp, err := client.GenerateSQL(ctx, req)
		
		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.NotEmpty(t, resp.GeneratedSql)
		assert.Greater(t, resp.ConfidenceScore, float32(0))
		assert.NotEmpty(t, resp.Explanation)
	})
	
	t.Run("ValidateSQL", func(t *testing.T) {
		req := &ValidateSQLRequest{
			Sql: "SELECT * FROM users",
		}
		
		resp, err := client.ValidateSQL(ctx, req)
		
		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.True(t, resp.IsValid)
		assert.NotEmpty(t, resp.FormattedSql)
	})
	
	t.Run("GetCapabilities", func(t *testing.T) {
		resp, err := client.GetCapabilities(ctx)
		
		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.NotEmpty(t, resp.SupportedDatabases)
		assert.NotEmpty(t, resp.Features)
		assert.Equal(t, HealthStatus_HEALTHY, resp.Status)
	})
	
	t.Run("IsHealthy", func(t *testing.T) {
		healthy := client.IsHealthy(ctx)
		assert.True(t, healthy)
	})
}

func TestAIPluginBridge_ErrorSimulation(t *testing.T) {
	bridge := NewAIPluginBridge()
	client := NewExtendedMockAIPluginClient(true)
	
	// Test error simulation
	client.SetSimulateErrors(true, false)
	bridge.RegisterPlugin("error-test", client)
	
	ctx := context.Background()
	req := &GenerateSQLRequest{
		NaturalLanguage: "Find users",
	}
	
	_, err := bridge.GenerateSQL(ctx, req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "simulated generation error")
}

func TestAIQueryRouter_IsAIQuery_Logic(t *testing.T) {
	// Create a minimal server-like struct for testing
	router := NewAIQueryRouter(nil)
	
	tests := []struct {
		name     string
		query    *DataQuery
		expected bool
	}{
		{
			name: "AI type query",
			query: &DataQuery{
				Type: "ai",
				Sql:  "SELECT * FROM users",
			},
			expected: true,
		},
		{
			name: "Natural language query",
			query: &DataQuery{
				Type:            "sql",
				NaturalLanguage: "Find users",
			},
			expected: true,
		},
		{
			name: "Database type specified",
			query: &DataQuery{
				Type:         "sql",
				Sql:          "SELECT * FROM users",
				DatabaseType: "postgresql",
			},
			expected: true,
		},
		{
			name: "Traditional SQL query",
			query: &DataQuery{
				Type: "sql",
				Sql:  "SELECT * FROM users",
			},
			expected: false,
		},
		{
			name:     "Nil query",
			query:    nil,
			expected: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := router.IsAIQuery(tt.query)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAIPluginBridge_Integration(t *testing.T) {
	// Test full integration flow
	bridge := NewAIPluginBridge()
	transformer := &MessageTransformer{}
	
	// Test SQL generation flow
	query := &DataQuery{
		Type:            "ai",
		NaturalLanguage: "Count all users",
		DatabaseType:    "postgresql",
		ExplainQuery:    true,
	}
	
	// Transform query
	req := transformer.TransformDataQueryToGenerateSQL(query)
	assert.NotNil(t, req)
	
	// Call bridge
	ctx := context.Background()
	resp, err := bridge.GenerateSQL(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	
	// Transform response
	result := transformer.TransformGenerateSQLToDataQueryResult(resp)
	require.NotNil(t, result)
	require.NotEmpty(t, result.Data)
	
	// Verify result structure
	dataMap := make(map[string]string)
	for _, pair := range result.Data {
		dataMap[pair.Key] = pair.Value
	}
	
	assert.Contains(t, dataMap, "generated_sql")
	assert.Contains(t, dataMap, "confidence_score")
	assert.Contains(t, dataMap, "explanation")
}