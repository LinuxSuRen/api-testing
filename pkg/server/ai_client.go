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

// AIRequest represents the request structure for AI operations
type AIRequest struct {
	Type    string                 `json:"type"`    // "nl_to_sql", "generate_test_case", etc.
	Input   string                 `json:"input"`   // Natural language input
	Context map[string]interface{} `json:"context"` // Additional context data
}

// AIResponse represents the response structure from AI operations
type AIResponse struct {
	Success bool                   `json:"success"`
	Result  string                 `json:"result"` // Generated SQL, test case, etc.
	Error   string                 `json:"error"`  // Error message if any
	Meta    map[string]interface{} `json:"meta"`   // Additional metadata
}

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
	// Prepare AI request
	request := AIRequest{
		Type:  "nl_to_sql",
		Input: naturalLanguage,
		Context: map[string]interface{}{
			"schema": schema,
		},
	}

	return c.executeAIRequest(ctx, request)
}

// GenerateTestCase generates test cases using the AI plugin
func (c *AIClient) GenerateTestCase(ctx context.Context, apiSpec string, requirements string) (*AIResponse, error) {
	// Prepare AI request
	request := AIRequest{
		Type:  "generate_test_case",
		Input: requirements,
		Context: map[string]interface{}{
			"api_spec": apiSpec,
		},
	}

	return c.executeAIRequest(ctx, request)
}

// OptimizeQuery optimizes SQL queries using the AI plugin
func (c *AIClient) OptimizeQuery(ctx context.Context, sqlQuery string, performance map[string]interface{}) (*AIResponse, error) {
	// Prepare AI request
	request := AIRequest{
		Type:  "optimize_query",
		Input: sqlQuery,
		Context: map[string]interface{}{
			"performance_data": performance,
		},
	}

	return c.executeAIRequest(ctx, request)
}

// executeAIRequest executes an AI request through the gRPC client
func (c *AIClient) executeAIRequest(ctx context.Context, request AIRequest) (*AIResponse, error) {
	// Serialize the AI request to JSON
	requestData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal AI request: %w", err)
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
			Name: request.Type,
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
	request := AIRequest{
		Type:    "health_check",
		Input:   "ping",
		Context: map[string]interface{}{},
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
