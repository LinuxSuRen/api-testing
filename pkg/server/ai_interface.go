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

// AIRequest represents a standard request to an AI plugin
// Following the simplicity principle: model name + prompt + optional config
type AIRequest struct {
	Model  string                 `json:"model"`  // Model identifier (e.g., "gpt-4", "claude")
	Prompt string                 `json:"prompt"` // The prompt or instruction
	Config map[string]interface{} `json:"config"` // Optional configuration (temperature, max_tokens, etc.)
}

// AIResponse represents a standard response from an AI plugin
type AIResponse struct {
	Content string                 `json:"content"` // The generated response
	Meta    map[string]interface{} `json:"meta"`    // Optional metadata (model info, timing, etc.)
}

// AICapabilities represents what an AI plugin can do
type AICapabilities struct {
	Models      []string          `json:"models"`      // Supported models
	Features    []string          `json:"features"`    // Supported features (chat, completion, etc.)
	Limits      map[string]int    `json:"limits"`      // Limits (max_tokens, rate_limit, etc.)
	Description string            `json:"description"` // Plugin description
	Version     string            `json:"version"`     // Plugin version
}

// Standard AI plugin communication methods
const (
	AIMethodGenerate      = "ai.generate"      // Generate content from prompt
	AIMethodCapabilities  = "ai.capabilities"  // Get plugin capabilities
)

// Standard plugin communication message format
type PluginRequest struct {
	Method  string      `json:"method"`  // Method name
	Payload interface{} `json:"payload"` // Request payload
}

type PluginResponse struct {
	Success bool        `json:"success"` // Whether request succeeded
	Data    interface{} `json:"data"`    // Response data
	Error   string      `json:"error"`   // Error message if failed
}