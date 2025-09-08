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
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/linuxsuren/api-testing/pkg/logging"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var pluginBridgeLogger = logging.DefaultLogger(logging.LogLevelInfo).WithName("ai_plugin_bridge")

// AIPluginBridge manages communication between the server and AI plugins
type AIPluginBridge struct {
	mu          sync.RWMutex
	pluginClients map[string]AIPluginClient
	defaultClient AIPluginClient
}

// PluginInfo contains basic information about a plugin
type PluginInfo struct {
	ID       string
	Name     string
	Address  string
	Priority int
	Healthy  bool
}

// NewAIPluginBridge creates a new AI plugin bridge
func NewAIPluginBridge() *AIPluginBridge {
	return &AIPluginBridge{
		pluginClients: make(map[string]AIPluginClient),
		defaultClient: NewMockAIPluginClient(),
	}
}

// RegisterPlugin registers an AI plugin client
func (b *AIPluginBridge) RegisterPlugin(id string, client AIPluginClient) {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	b.pluginClients[id] = client
	pluginBridgeLogger.Info("Registered AI plugin", "id", id)
}

// UnregisterPlugin removes an AI plugin client
func (b *AIPluginBridge) UnregisterPlugin(id string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	delete(b.pluginClients, id)
	pluginBridgeLogger.Info("Unregistered AI plugin", "id", id)
}

// GetPlugin returns the best available plugin for an operation
func (b *AIPluginBridge) GetPlugin(operation string) AIPluginClient {
	b.mu.RLock()
	defer b.mu.RUnlock()
	
	// For now, return the first available plugin or the default mock
	for _, client := range b.pluginClients {
		if client.IsHealthy(context.Background()) {
			return client
		}
	}
	
	return b.defaultClient
}

// GetAllPlugins returns information about all registered plugins
func (b *AIPluginBridge) GetAllPlugins() []PluginInfo {
	b.mu.RLock()
	defer b.mu.RUnlock()
	
	plugins := make([]PluginInfo, 0, len(b.pluginClients))
	for id, client := range b.pluginClients {
		plugins = append(plugins, PluginInfo{
			ID:      id,
			Name:    id, // For simplicity, use ID as name
			Healthy: client.IsHealthy(context.Background()),
		})
	}
	
	return plugins
}

// GenerateSQL routes SQL generation requests to appropriate plugins
func (b *AIPluginBridge) GenerateSQL(ctx context.Context, req *GenerateSQLRequest) (*GenerateSQLResponse, error) {
	client := b.GetPlugin("generate_sql")
	
	pluginBridgeLogger.Info("Routing GenerateSQL request", "natural_language", req.NaturalLanguage)
	
	return client.GenerateSQL(ctx, req)
}

// ValidateSQL routes SQL validation requests to appropriate plugins
func (b *AIPluginBridge) ValidateSQL(ctx context.Context, req *ValidateSQLRequest) (*ValidateSQLResponse, error) {
	client := b.GetPlugin("validate_sql")
	
	pluginBridgeLogger.Info("Routing ValidateSQL request", "sql_length", len(req.Sql))
	
	return client.ValidateSQL(ctx, req)
}

// GetAICapabilities returns consolidated capabilities from all healthy plugins
func (b *AIPluginBridge) GetAICapabilities(ctx context.Context) (*AICapabilitiesResponse, error) {
	client := b.GetPlugin("capabilities")
	
	return client.GetCapabilities(ctx)
}

// MessageTransformer handles transformation between API and plugin message formats
type MessageTransformer struct{}

// TransformDataQueryToGenerateSQL converts DataQuery to GenerateSQLRequest
func (t *MessageTransformer) TransformDataQueryToGenerateSQL(query *DataQuery) *GenerateSQLRequest {
	req := &GenerateSQLRequest{
		NaturalLanguage: query.NaturalLanguage,
		Context:         query.AiContext,
	}
	
	if query.DatabaseType != "" {
		req.DatabaseTarget = &DatabaseTarget{
			Type: query.DatabaseType,
		}
	}
	
	if query.ExplainQuery {
		req.Options = &GenerationOptions{
			IncludeExplanation: true,
			FormatOutput:       true,
		}
	}
	
	return req
}

// TransformDataQueryToValidateSQL converts DataQuery to ValidateSQLRequest
func (t *MessageTransformer) TransformDataQueryToValidateSQL(query *DataQuery) *ValidateSQLRequest {
	return &ValidateSQLRequest{
		Sql:          query.Sql,
		DatabaseType: query.DatabaseType,
		Context:      query.AiContext,
	}
}

// TransformGenerateSQLToDataQueryResult converts GenerateSQLResponse to DataQueryResult
func (t *MessageTransformer) TransformGenerateSQLToDataQueryResult(resp *GenerateSQLResponse) *DataQueryResult {
	result := &DataQueryResult{
		Meta: &DataMeta{
			Duration: "AI SQL Generation",
		},
	}
	
	if resp.Error != nil {
		result.Data = []*Pair{
			{Key: "error", Value: resp.Error.Message, Description: resp.Error.Details},
			{Key: "error_code", Value: resp.Error.Code.String()},
		}
	} else {
		result.Data = []*Pair{
			{Key: "generated_sql", Value: resp.GeneratedSql, Description: "AI-generated SQL query"},
			{Key: "confidence_score", Value: fmt.Sprintf("%.2f", resp.ConfidenceScore), Description: "AI confidence level (0.0-1.0)"},
			{Key: "explanation", Value: resp.Explanation, Description: "AI explanation of the generated query"},
		}
		
		// Add suggestions as items
		for i, suggestion := range resp.Suggestions {
			result.Items = append(result.Items, &Pairs{
				Data: []*Pair{
					{Key: "suggestion", Value: suggestion, Description: fmt.Sprintf("AI suggestion #%d", i+1)},
				},
			})
		}
		
		// Add processing metadata
		if resp.Metadata != nil {
			result.AiInfo = &AIProcessingInfo{
				RequestId:        resp.Metadata.RequestId,
				ProcessingTimeMs: resp.Metadata.ProcessingTimeMs,
				ModelUsed:        resp.Metadata.ModelUsed,
				ConfidenceScore:  resp.ConfidenceScore,
			}
		}
	}
	
	return result
}

// TransformValidateSQLToDataQueryResult converts ValidateSQLResponse to DataQueryResult
func (t *MessageTransformer) TransformValidateSQLToDataQueryResult(resp *ValidateSQLResponse) *DataQueryResult {
	result := &DataQueryResult{
		Meta: &DataMeta{
			Duration: "AI SQL Validation",
		},
	}
	
	result.Data = []*Pair{
		{Key: "is_valid", Value: fmt.Sprintf("%t", resp.IsValid), Description: "SQL validation result"},
		{Key: "formatted_sql", Value: resp.FormattedSql, Description: "AI-formatted SQL query"},
	}
	
	// Add validation errors as items
	for i, validationError := range resp.Errors {
		result.Items = append(result.Items, &Pairs{
			Data: []*Pair{
				{Key: "error_type", Value: validationError.Type.String(), Description: "Type of validation error"},
				{Key: "error_message", Value: validationError.Message, Description: "Validation error description"},
				{Key: "line", Value: fmt.Sprintf("%d", validationError.Line), Description: "Error line number"},
				{Key: "column", Value: fmt.Sprintf("%d", validationError.Column), Description: "Error column number"},
			},
		})
		pluginBridgeLogger.Info("Validation error added", "index", i, "message", validationError.Message)
	}
	
	// Add warnings as items
	for i, warning := range resp.Warnings {
		result.Items = append(result.Items, &Pairs{
			Data: []*Pair{
				{Key: "warning", Value: warning, Description: fmt.Sprintf("Validation warning #%d", i+1)},
			},
		})
	}
	
	// Add processing metadata
	if resp.Metadata != nil {
		result.AiInfo = &AIProcessingInfo{
			RequestId:        fmt.Sprintf("validation-%d", time.Now().Unix()),
			ProcessingTimeMs: resp.Metadata.ValidationTimeMs,
			ModelUsed:        resp.Metadata.ValidatorVersion,
			ConfidenceScore:  1.0, // Validation is deterministic
		}
	}
	
	return result
}

// ValidateAIQuery validates AI query parameters
func (t *MessageTransformer) ValidateAIQuery(query *DataQuery) error {
	if query.NaturalLanguage == "" && query.Sql == "" {
		return fmt.Errorf("AI query must have either natural_language or sql field")
	}
	
	// Validate database type if provided
	if query.DatabaseType != "" {
		validTypes := []string{"mysql", "postgresql", "sqlite", "mongodb", "oracle", "mssql"}
		isValid := false
		for _, validType := range validTypes {
			if strings.EqualFold(query.DatabaseType, validType) {
				isValid = true
				break
			}
		}
		if !isValid {
			return fmt.Errorf("unsupported database type: %s", query.DatabaseType)
		}
	}
	
	return nil
}

// CreateErrorResult creates a DataQueryResult with error information
func (t *MessageTransformer) CreateErrorResult(err error) *DataQueryResult {
	return &DataQueryResult{
		Data: []*Pair{
			{Key: "error", Value: err.Error(), Description: "AI query processing error"},
			{Key: "timestamp", Value: time.Now().Format(time.RFC3339), Description: "Error timestamp"},
		},
		Meta: &DataMeta{
			Duration: "Error occurred",
		},
	}
}

// ExtendedMockAIPluginClient provides enhanced mock functionality for testing
type ExtendedMockAIPluginClient struct {
	healthy                bool
	simulateGenerationError bool
	simulateValidationError bool
}

// NewExtendedMockAIPluginClient creates a new enhanced mock client
func NewExtendedMockAIPluginClient(healthy bool) *ExtendedMockAIPluginClient {
	return &ExtendedMockAIPluginClient{
		healthy: healthy,
	}
}

func (m *ExtendedMockAIPluginClient) GenerateSQL(ctx context.Context, req *GenerateSQLRequest) (*GenerateSQLResponse, error) {
	if m.simulateGenerationError {
		return nil, fmt.Errorf("simulated generation error")
	}
	
	if req.NaturalLanguage == "" {
		return &GenerateSQLResponse{
			Error: &AIError{
				Code:    AIErrorCode_INVALID_INPUT,
				Message: "Natural language input is required",
				Details: "The natural_language field cannot be empty",
			},
		}, nil
	}
	
	return &GenerateSQLResponse{
		GeneratedSql:    fmt.Sprintf("-- Generated from: %s\nSELECT * FROM table_name WHERE condition;", req.NaturalLanguage),
		ConfidenceScore: 0.85,
		Explanation:     fmt.Sprintf("AI-generated SQL query for: %s", req.NaturalLanguage),
		Suggestions:     []string{"Consider adding LIMIT clause", "Verify table schema"},
		Metadata: &GenerationMetadata{
			RequestId:        fmt.Sprintf("mock-%d", time.Now().Unix()),
			ProcessingTimeMs: 100.0,
			ModelUsed:        "mock-enhanced-ai",
			TokenCount:       25,
			Timestamp:        timestamppb.New(time.Now()),
		},
	}, nil
}

func (m *ExtendedMockAIPluginClient) ValidateSQL(ctx context.Context, req *ValidateSQLRequest) (*ValidateSQLResponse, error) {
	if m.simulateValidationError {
		return nil, fmt.Errorf("simulated validation error")
	}
	
	if req.Sql == "" {
		return &ValidateSQLResponse{
			IsValid: false,
			Errors: []*ValidationError{
				{
					Message: "SQL query is required",
					Line:    1,
					Column:  1,
					Type:    ValidationErrorType_SYNTAX_ERROR,
				},
			},
		}, nil
	}
	
	// Enhanced validation logic
	sqlUpper := strings.ToUpper(req.Sql)
	isValid := strings.Contains(sqlUpper, "SELECT") ||
		strings.Contains(sqlUpper, "INSERT") ||
		strings.Contains(sqlUpper, "UPDATE") ||
		strings.Contains(sqlUpper, "DELETE")
	
	if isValid {
		return &ValidateSQLResponse{
			IsValid:      true,
			FormattedSql: formatSQL(req.Sql),
			Metadata: &ValidationMetadata{
				ValidatorVersion:  "mock-enhanced-validator-1.0",
				ValidationTimeMs: 15.0,
				Timestamp:        timestamppb.New(time.Now()),
			},
		}, nil
	}
	
	return &ValidateSQLResponse{
		IsValid: false,
		Errors: []*ValidationError{
			{
				Message: "Invalid SQL syntax detected",
				Line:    1,
				Column:  1,
				Type:    ValidationErrorType_SYNTAX_ERROR,
			},
		},
		Warnings: []string{"Consider using standard SQL syntax"},
	}, nil
}

func (m *ExtendedMockAIPluginClient) GetCapabilities(ctx context.Context) (*AICapabilitiesResponse, error) {
	return &AICapabilitiesResponse{
		SupportedDatabases: []string{"mysql", "postgresql", "sqlite", "mongodb"},
		Features: []*AIFeature{
			{
				Name:        "sql_generation",
				Enabled:     true,
				Description: "Enhanced SQL generation from natural language",
				Parameters: map[string]string{
					"max_complexity": "high",
					"max_tokens":     "1000",
					"model_version":  "enhanced",
				},
			},
			{
				Name:        "sql_validation",
				Enabled:     true,
				Description: "Enhanced SQL syntax and semantic validation",
				Parameters: map[string]string{
					"syntax_check":   "true",
					"semantic_check": "true",
					"format":         "true",
				},
			},
		},
		Version: "enhanced-mock-1.0.0",
		Status:  HealthStatus_HEALTHY,
		Limits: map[string]string{
			"max_requests_per_minute": "150",
			"max_query_length":        "10000",
		},
	}, nil
}

func (m *ExtendedMockAIPluginClient) IsHealthy(ctx context.Context) bool {
	return m.healthy
}

// SetSimulateErrors configures error simulation for testing
func (m *ExtendedMockAIPluginClient) SetSimulateErrors(generation, validation bool) {
	m.simulateGenerationError = generation
	m.simulateValidationError = validation
}