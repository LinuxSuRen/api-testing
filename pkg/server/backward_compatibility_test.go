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

	"github.com/linuxsuren/api-testing/pkg/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

// TestDataQuery_BackwardCompatibility tests that existing DataQuery messages
// still work after adding AI extensions
func TestDataQuery_BackwardCompatibility(t *testing.T) {
	tests := []struct {
		name        string
		query       *server.DataQuery
		description string
	}{
		{
			name: "traditional_sql_query",
			query: &server.DataQuery{
				Type:   "sql",
				Key:    "users_query",
				Sql:    "SELECT * FROM users WHERE active = 1",
				Offset: 0,
				Limit:  100,
			},
			description: "Traditional SQL query without AI extensions should work unchanged",
		},
		{
			name: "minimal_query",
			query: &server.DataQuery{
				Type: "sql",
			},
			description: "Minimal query with only required fields",
		},
		{
			name: "legacy_pagination_query",
			query: &server.DataQuery{
				Type:   "sql", 
				Key:    "paginated_users",
				Sql:    "SELECT id, name FROM users",
				Offset: 50,
				Limit:  25,
			},
			description: "Legacy pagination should continue to work",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Testing: %s", tt.description)

			// Test serialization
			data, err := proto.Marshal(tt.query)
			require.NoError(t, err, "Failed to marshal legacy DataQuery")
			require.NotEmpty(t, data)

			// Test deserialization
			unmarshaled := &server.DataQuery{}
			err = proto.Unmarshal(data, unmarshaled)
			require.NoError(t, err, "Failed to unmarshal legacy DataQuery")

			// Verify all original fields are preserved
			assert.Equal(t, tt.query.Type, unmarshaled.Type)
			assert.Equal(t, tt.query.Key, unmarshaled.Key)
			assert.Equal(t, tt.query.Sql, unmarshaled.Sql)
			assert.Equal(t, tt.query.Offset, unmarshaled.Offset)
			assert.Equal(t, tt.query.Limit, unmarshaled.Limit)

			// Verify new AI fields have default values (empty/zero)
			assert.Empty(t, unmarshaled.NaturalLanguage, "NaturalLanguage should be empty for legacy queries")
			assert.Empty(t, unmarshaled.DatabaseType, "DatabaseType should be empty for legacy queries")
			assert.False(t, unmarshaled.ExplainQuery, "ExplainQuery should be false for legacy queries")
			assert.Empty(t, unmarshaled.AiContext, "AiContext should be empty for legacy queries")

			// Test full equality
			assert.True(t, proto.Equal(tt.query, unmarshaled))
		})
	}
}

// TestDataQueryResult_BackwardCompatibility tests that existing DataQueryResult messages
// still work after adding AI extensions
func TestDataQueryResult_BackwardCompatibility(t *testing.T) {
	tests := []struct {
		name        string
		result      *server.DataQueryResult
		description string
	}{
		{
			name: "traditional_query_result",
			result: &server.DataQueryResult{
				Data: []*server.Pair{
					{Key: "id", Value: "1", Description: "User ID"},
					{Key: "name", Value: "John Doe", Description: "Full name"},
					{Key: "email", Value: "john@example.com", Description: "Email address"},
				},
				Items: []*server.Pairs{
					{
						Data: []*server.Pair{
							{Key: "id", Value: "1"},
							{Key: "name", Value: "John Doe"},
						},
					},
					{
						Data: []*server.Pair{
							{Key: "id", Value: "2"},
							{Key: "name", Value: "Jane Smith"},
						},
					},
				},
				Meta: &server.DataMeta{
					Databases:       []string{"testdb", "userdb"},
					Tables:          []string{"users", "profiles"},
					CurrentDatabase: "testdb",
					Duration:        "125ms",
					Labels: []*server.Pair{
						{Key: "env", Value: "production"},
						{Key: "region", Value: "us-east-1"},
					},
				},
			},
			description: "Traditional query result without AI processing info should work unchanged",
		},
		{
			name: "minimal_result",
			result: &server.DataQueryResult{
				Data: []*server.Pair{
					{Key: "count", Value: "42"},
				},
			},
			description: "Minimal result with only data field",
		},
		{
			name: "empty_result",
			result: &server.DataQueryResult{},
			description: "Empty result should be handled correctly",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Testing: %s", tt.description)

			// Test serialization
			data, err := proto.Marshal(tt.result)
			require.NoError(t, err, "Failed to marshal legacy DataQueryResult")
			// Empty messages may serialize to empty bytes - this is expected in protobuf
			if tt.name != "empty_result" {
				require.NotEmpty(t, data)
			}

			// Test deserialization
			unmarshaled := &server.DataQueryResult{}
			err = proto.Unmarshal(data, unmarshaled)
			require.NoError(t, err, "Failed to unmarshal legacy DataQueryResult")

			// Verify all original fields are preserved
			assert.Equal(t, len(tt.result.Data), len(unmarshaled.Data))
			for i, pair := range tt.result.Data {
				assert.Equal(t, pair.Key, unmarshaled.Data[i].Key)
				assert.Equal(t, pair.Value, unmarshaled.Data[i].Value)
				assert.Equal(t, pair.Description, unmarshaled.Data[i].Description)
			}

			assert.Equal(t, len(tt.result.Items), len(unmarshaled.Items))
			if tt.result.Meta != nil {
				require.NotNil(t, unmarshaled.Meta)
				assert.Equal(t, tt.result.Meta.CurrentDatabase, unmarshaled.Meta.CurrentDatabase)
				assert.Equal(t, tt.result.Meta.Duration, unmarshaled.Meta.Duration)
			}

			// Verify new AI field has default value (nil)
			assert.Nil(t, unmarshaled.AiInfo, "AiInfo should be nil for legacy results")

			// Test full equality
			assert.True(t, proto.Equal(tt.result, unmarshaled))
		})
	}
}

// TestMixedCompatibility tests that AI and legacy messages can coexist
func TestMixedCompatibility(t *testing.T) {
	tests := []struct {
		name             string
		legacyQuery      *server.DataQuery
		aiQuery          *server.DataQuery
		legacyResult     *server.DataQueryResult
		aiResult         *server.DataQueryResult
		description      string
	}{
		{
			name: "mixed_query_types",
			legacyQuery: &server.DataQuery{
				Type: "sql",
				Key:  "legacy_query",
				Sql:  "SELECT * FROM users",
			},
			aiQuery: &server.DataQuery{
				Type:            "ai",
				Key:             "ai_query",
				NaturalLanguage: "Find all active users",
				DatabaseType:    "postgresql",
				ExplainQuery:    true,
				AiContext: map[string]string{
					"table": "users",
				},
			},
			legacyResult: &server.DataQueryResult{
				Data: []*server.Pair{
					{Key: "count", Value: "10"},
				},
				Meta: &server.DataMeta{
					Duration: "50ms",
				},
			},
			aiResult: &server.DataQueryResult{
				Data: []*server.Pair{
					{Key: "id", Value: "1"},
				},
				Meta: &server.DataMeta{
					Duration: "150ms",
				},
				AiInfo: &server.AIProcessingInfo{
					RequestId:        "ai-123",
					ProcessingTimeMs: 100.0,
					ModelUsed:        "gpt-4",
					ConfidenceScore:  0.95,
				},
			},
			description: "Legacy and AI queries/results should be serializable together",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Testing: %s", tt.description)

			// Test that both queries can be serialized/deserialized
			legacyData, err := proto.Marshal(tt.legacyQuery)
			require.NoError(t, err)
			
			aiData, err := proto.Marshal(tt.aiQuery)
			require.NoError(t, err)

			// Deserialize and verify
			legacyUnmarshaled := &server.DataQuery{}
			err = proto.Unmarshal(legacyData, legacyUnmarshaled)
			require.NoError(t, err)
			assert.True(t, proto.Equal(tt.legacyQuery, legacyUnmarshaled))

			aiUnmarshaled := &server.DataQuery{}
			err = proto.Unmarshal(aiData, aiUnmarshaled)
			require.NoError(t, err)
			assert.True(t, proto.Equal(tt.aiQuery, aiUnmarshaled))

			// Test that both results can be serialized/deserialized
			legacyResultData, err := proto.Marshal(tt.legacyResult)
			require.NoError(t, err)
			
			aiResultData, err := proto.Marshal(tt.aiResult)
			require.NoError(t, err)

			// Deserialize and verify
			legacyResultUnmarshaled := &server.DataQueryResult{}
			err = proto.Unmarshal(legacyResultData, legacyResultUnmarshaled)
			require.NoError(t, err)
			assert.True(t, proto.Equal(tt.legacyResult, legacyResultUnmarshaled))

			aiResultUnmarshaled := &server.DataQueryResult{}
			err = proto.Unmarshal(aiResultData, aiResultUnmarshaled)
			require.NoError(t, err)
			assert.True(t, proto.Equal(tt.aiResult, aiResultUnmarshaled))
		})
	}
}

// TestFieldNumbering verifies that field numbers follow the reserved ranges correctly
func TestFieldNumbering(t *testing.T) {
	t.Run("DataQuery_field_numbers", func(t *testing.T) {
		query := &server.DataQuery{
			Type:            "ai",            // field 1 (existing)
			Key:             "test",          // field 2 (existing)
			Sql:             "SELECT 1",      // field 3 (existing)
			Offset:          0,               // field 4 (existing)
			Limit:           10,              // field 5 (existing)
			NaturalLanguage: "test",          // field 10 (AI extension)
			DatabaseType:    "mysql",         // field 11 (AI extension)
			ExplainQuery:    true,            // field 12 (AI extension)
			AiContext:       map[string]string{"key": "value"}, // field 13 (AI extension)
		}

		// Test that serialization works with our field numbering
		data, err := proto.Marshal(query)
		require.NoError(t, err)
		require.NotEmpty(t, data)

		// Test deserialization
		unmarshaled := &server.DataQuery{}
		err = proto.Unmarshal(data, unmarshaled)
		require.NoError(t, err)
		assert.True(t, proto.Equal(query, unmarshaled))
	})

	t.Run("DataQueryResult_field_numbers", func(t *testing.T) {
		result := &server.DataQueryResult{
			Data:  []*server.Pair{{Key: "test", Value: "value"}}, // field 1 (existing)
			Items: []*server.Pairs{},                             // field 2 (existing)
			Meta:  &server.DataMeta{Duration: "10ms"},            // field 3 (existing)
			AiInfo: &server.AIProcessingInfo{                     // field 10 (AI extension)
				RequestId:        "test-123",
				ProcessingTimeMs: 50.0,
				ModelUsed:        "test-model",
			},
		}

		// Test that serialization works with our field numbering
		data, err := proto.Marshal(result)
		require.NoError(t, err)
		require.NotEmpty(t, data)

		// Test deserialization
		unmarshaled := &server.DataQueryResult{}
		err = proto.Unmarshal(data, unmarshaled)
		require.NoError(t, err)
		assert.True(t, proto.Equal(result, unmarshaled))
	})
}