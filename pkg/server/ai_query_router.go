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
	"time"

	"github.com/linuxsuren/api-testing/pkg/logging"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var aiQueryLogger = logging.DefaultLogger(logging.LogLevelInfo).WithName("ai_query_router")

// ServerInterface defines minimal server interface needed by router
type ServerInterface interface {
	// Add any server methods needed by the router in the future
}

// AIQueryRouter handles routing and processing of AI-type queries
type AIQueryRouter struct {
	server       ServerInterface
	pluginBridge *AIPluginBridge
	transformer  *MessageTransformer
}

// AIPluginClient defines the interface for communicating with AI plugins
type AIPluginClient interface {
	GenerateSQL(ctx context.Context, req *GenerateSQLRequest) (*GenerateSQLResponse, error)
	ValidateSQL(ctx context.Context, req *ValidateSQLRequest) (*ValidateSQLResponse, error)
	GetCapabilities(ctx context.Context) (*AICapabilitiesResponse, error)
	IsHealthy(ctx context.Context) bool
}

// NewAIQueryRouter creates a new AI query router
func NewAIQueryRouter(s ServerInterface) *AIQueryRouter {
	return &AIQueryRouter{
		server:       s,
		pluginBridge: NewAIPluginBridge(),
		transformer:  &MessageTransformer{},
	}
}

// NewAIQueryRouterWithBridge creates a new AI query router with a custom bridge
func NewAIQueryRouterWithBridge(s ServerInterface, bridge *AIPluginBridge) *AIQueryRouter {
	return &AIQueryRouter{
		server:       s,
		pluginBridge: bridge,
		transformer:  &MessageTransformer{},
	}
}

// IsAIQuery determines if a DataQuery is an AI-type query
func (r *AIQueryRouter) IsAIQuery(query *DataQuery) bool {
	if query == nil {
		return false
	}
	
	// Check if query type is explicitly "ai"
	if strings.ToLower(query.Type) == "ai" {
		return true
	}
	
	// Check if query has AI-specific fields
	if query.NaturalLanguage != "" || query.DatabaseType != "" || len(query.AiContext) > 0 {
		return true
	}
	
	return false
}

// RouteAIQuery processes AI queries and returns appropriate results
func (r *AIQueryRouter) RouteAIQuery(ctx context.Context, query *DataQuery) (*DataQueryResult, error) {
	aiQueryLogger.Info("Routing AI query", "type", query.Type, "natural_language", query.NaturalLanguage)
	
	// Validate AI query
	if err := r.transformer.ValidateAIQuery(query); err != nil {
		return r.transformer.CreateErrorResult(err), nil
	}
	
	// Determine AI operation type and route accordingly
	if query.NaturalLanguage != "" && (query.Sql == "" || query.ExplainQuery) {
		return r.routeGenerateSQL(ctx, query)
	} else if query.Sql != "" {
		return r.routeValidateSQL(ctx, query)
	}
	
	return r.transformer.CreateErrorResult(fmt.Errorf("invalid AI query: cannot determine operation type")), nil
}


// routeGenerateSQL handles SQL generation requests
func (r *AIQueryRouter) routeGenerateSQL(ctx context.Context, query *DataQuery) (*DataQueryResult, error) {
	// Transform DataQuery to GenerateSQLRequest
	req := r.transformer.TransformDataQueryToGenerateSQL(query)
	
	// Call AI plugin via bridge
	resp, err := r.pluginBridge.GenerateSQL(ctx, req)
	if err != nil {
		aiQueryLogger.Error(err, "AI plugin GenerateSQL failed")
		return r.transformer.CreateErrorResult(fmt.Errorf("AI service error: %w", err)), nil
	}
	
	// Transform response back to DataQueryResult
	return r.transformer.TransformGenerateSQLToDataQueryResult(resp), nil
}

// routeValidateSQL handles SQL validation requests
func (r *AIQueryRouter) routeValidateSQL(ctx context.Context, query *DataQuery) (*DataQueryResult, error) {
	// Transform DataQuery to ValidateSQLRequest
	req := r.transformer.TransformDataQueryToValidateSQL(query)
	
	// Call AI plugin via bridge
	resp, err := r.pluginBridge.ValidateSQL(ctx, req)
	if err != nil {
		aiQueryLogger.Error(err, "AI plugin ValidateSQL failed")
		return r.transformer.CreateErrorResult(fmt.Errorf("AI service error: %w", err)), nil
	}
	
	// Transform response back to DataQueryResult
	return r.transformer.TransformValidateSQLToDataQueryResult(resp), nil
}


// MockAIPluginClient provides a mock implementation for testing
type MockAIPluginClient struct{}

func NewMockAIPluginClient() *MockAIPluginClient {
	return &MockAIPluginClient{}
}

func (m *MockAIPluginClient) GenerateSQL(ctx context.Context, req *GenerateSQLRequest) (*GenerateSQLResponse, error) {
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
		Explanation:     fmt.Sprintf("Generated SQL query based on natural language input: %s", req.NaturalLanguage),
		Suggestions:     []string{"Consider adding LIMIT clause for large datasets", "Verify table and column names exist"},
		Metadata: &GenerationMetadata{
			RequestId:        fmt.Sprintf("mock-req-%d", time.Now().Unix()),
			ProcessingTimeMs: 100.0,
			ModelUsed:        "mock-ai-model",
			TokenCount:       25,
			Timestamp:        timestamppb.New(time.Now()),
		},
	}, nil
}

func (m *MockAIPluginClient) ValidateSQL(ctx context.Context, req *ValidateSQLRequest) (*ValidateSQLResponse, error) {
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
	
	// Basic validation
	isValid := strings.Contains(strings.ToUpper(req.Sql), "SELECT") ||
		strings.Contains(strings.ToUpper(req.Sql), "INSERT") ||
		strings.Contains(strings.ToUpper(req.Sql), "UPDATE") ||
		strings.Contains(strings.ToUpper(req.Sql), "DELETE")
	
	if isValid {
		return &ValidateSQLResponse{
			IsValid:      true,
			FormattedSql: formatSQL(req.Sql),
			Metadata: &ValidationMetadata{
				ValidatorVersion:  "mock-validator-1.0",
				ValidationTimeMs: 10.0,
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
		Warnings: []string{"Consider using standard SQL keywords"},
	}, nil
}

func (m *MockAIPluginClient) GetCapabilities(ctx context.Context) (*AICapabilitiesResponse, error) {
	return &AICapabilitiesResponse{
		SupportedDatabases: []string{"mysql", "postgresql", "sqlite", "mongodb"},
		Features: []*AIFeature{
			{
				Name:        "sql_generation",
				Enabled:     true,
				Description: "Generate SQL from natural language",
				Parameters: map[string]string{
					"max_complexity": "high",
					"max_tokens":     "1000",
				},
			},
			{
				Name:        "sql_validation",
				Enabled:     true,
				Description: "Validate and format SQL queries",
				Parameters: map[string]string{
					"syntax_check": "true",
					"format":       "true",
				},
			},
		},
		Version: "mock-1.0.0",
		Status:  HealthStatus_HEALTHY,
		Limits: map[string]string{
			"max_requests_per_minute": "100",
			"max_query_length":        "5000",
		},
	}, nil
}

func (m *MockAIPluginClient) IsHealthy(ctx context.Context) bool {
	return true
}