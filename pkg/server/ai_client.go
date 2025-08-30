package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// AIClient wraps the gRPC client for AI plugin communication
type AIClient struct {
	client RunnerExtensionClient
	conn   *grpc.ClientConn
}

// Note: AIRequest and AIResponse types are now defined in server.pb.go
// We use the protobuf-generated types instead of custom definitions

// NewAIClient creates a new AI plugin client
func NewAIClient(address string) (*AIClient, error) {
	// Create gRPC connection with insecure credentials for local communication
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to AI plugin at %s: %w", address, err)
	}

	// Create the gRPC client
	client := NewRunnerExtensionClient(conn)

	return &AIClient{
		client: client,
		conn:   conn,
	}, nil
}

// Close closes the gRPC connection
func (c *AIClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// ConvertNLToSQL converts natural language to SQL using the AI plugin
func (c *AIClient) ConvertNLToSQL(ctx context.Context, naturalLanguage string, schema map[string]interface{}) (*AIResponse, error) {
	// Serialize schema to JSON string for protobuf Context field
	schemaJSON, err := json.Marshal(map[string]interface{}{
		"type":   "nl_to_sql",
		"schema": schema,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal schema: %w", err)
	}

	// Prepare AI request using protobuf types
	request := &AIRequest{
		Input:   naturalLanguage,
		Context: string(schemaJSON),
	}

	return c.executeAIRequest(ctx, request)
}

// GenerateTestCase generates test cases using the AI plugin
func (c *AIClient) GenerateTestCase(ctx context.Context, apiSpec string, requirements string) (*AIResponse, error) {
	// Serialize context to JSON string for protobuf Context field
	contextJSON, err := json.Marshal(map[string]interface{}{
		"type":     "generate_test_case",
		"api_spec": apiSpec,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal context: %w", err)
	}

	// Prepare AI request using protobuf types
	request := &AIRequest{
		Input:   requirements,
		Context: string(contextJSON),
	}

	return c.executeAIRequest(ctx, request)
}

// OptimizeQuery optimizes SQL queries using the AI plugin
func (c *AIClient) OptimizeQuery(ctx context.Context, sqlQuery string, performance map[string]interface{}) (*AIResponse, error) {
	// Serialize context to JSON string for protobuf Context field
	contextJSON, err := json.Marshal(map[string]interface{}{
		"type":             "optimize_query",
		"performance_data": performance,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal context: %w", err)
	}

	// Prepare AI request using protobuf types
	request := &AIRequest{
		Input:   sqlQuery,
		Context: string(contextJSON),
	}

	return c.executeAIRequest(ctx, request)
}

// executeAIRequest executes an AI request through the gRPC client
func (c *AIClient) executeAIRequest(ctx context.Context, request *AIRequest) (*AIResponse, error) {
	// Serialize the AI request to JSON
	requestData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal AI request: %w", err)
	}

	// Parse the context to get the request type
	var contextData map[string]interface{}
	if err := json.Unmarshal([]byte(request.Context), &contextData); err != nil {
		return nil, fmt.Errorf("failed to parse request context: %w", err)
	}
	requestType, _ := contextData["type"].(string)
	if requestType == "" {
		requestType = "ai-request"
	}

	// Create the TestSuiteWithCase structure for gRPC communication
	// Using the existing protobuf structure to carry JSON data
	testSuite := &TestSuiteWithCase{
		Suite: &TestSuite{
			Name: "ai-plugin-request",
			Spec: &APISpec{
				Kind: "ai",
				Url:  string(requestData), // Embed JSON data in URL field
			},
		},
		Case: &TestCase{
			Name: requestType,
		},
	}

	// Set timeout for the request
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Execute the gRPC call
	result, err := c.client.Run(ctxWithTimeout, testSuite)
	if err != nil {
		return nil, fmt.Errorf("AI plugin gRPC call failed: %w", err)
	}

	// Parse the response
	if result == nil {
		return &AIResponse{
			Success: false,
			Error:   "AI plugin returned nil result",
		}, nil
	}

	// Try to parse the message as JSON (AI response)
	var aiResponse AIResponse
	if err := json.Unmarshal([]byte(result.Message), &aiResponse); err != nil {
		// If JSON parsing fails, treat as plain text response
		log.Printf("Failed to parse AI response as JSON, treating as plain text: %v", err)
		aiResponse = AIResponse{
			Success: result.Success,
			Result:  result.Message,
			Error:   "",
		}
		if !result.Success {
			aiResponse.Error = result.Message
			aiResponse.Result = ""
		}
	}

	return &aiResponse, nil
}

// HealthCheck checks if the AI plugin is healthy and responsive
func (c *AIClient) HealthCheck(ctx context.Context) error {
	// Create a simple health check request
	contextJSON, err := json.Marshal(map[string]interface{}{
		"type": "health_check",
	})
	if err != nil {
		return fmt.Errorf("failed to marshal health check context: %w", err)
	}

	request := &AIRequest{
		Input:   "ping",
		Context: string(contextJSON),
	}

	response, err := c.executeAIRequest(ctx, request)
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}

	if !response.Success {
		return fmt.Errorf("AI plugin health check failed: %s", response.Error)
	}

	return nil
}
