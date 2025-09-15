/*
Copyright 2025 API Testing Authors.

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

// AI Plugin Communication Interface Standards
// AI plugins use the existing testing.Loader.Query(map[string]string) interface

// Standard AI plugin communication methods
const (
	AIMethodGenerate     = "ai.generate"     // Generate content from prompt
	AIMethodCapabilities = "ai.capabilities" // Get plugin capabilities
)

// AI Plugin Query Parameter Standards
// AI plugins are called using loader.Query(query map[string]string) with these parameters:

// For ai.generate:
//   - "method": "ai.generate"
//   - "model":  model identifier (e.g., "gpt-4", "claude")
//   - "prompt": the prompt or instruction
//   - "config": optional JSON configuration string (e.g., `{"temperature": 0.7, "max_tokens": 1000}`)

// For ai.capabilities:
//   - "method": "ai.capabilities"

// AI Plugin Response Standards
// AI plugins return response through testing.DataResult.Pairs with these keys:

// For successful ai.generate:
//   - "content": the generated content
//   - "meta":    optional JSON metadata string (model info, timing, etc.)
//   - "success": "true"

// For successful ai.capabilities:
//   - "capabilities": JSON string containing plugin capabilities
//   - "models":       JSON array of supported models (fallback if capabilities not available)
//   - "features":     JSON array of supported features (fallback)
//   - "description":  plugin description (fallback)
//   - "version":      plugin version (fallback)
//   - "success":      "true"

// For errors:
//   - "error":   error message
//   - "success": "false"

// Plugin Discovery
// AI plugins are identified by having "ai" in their categories field:
//   categories: ["ai"]

// Usage Examples:
//
// Get AI plugins:
//   stores, err := server.GetStores(ctx, &SimpleQuery{Kind: "ai"})
//
// Call AI plugin:
//   loader, err := server.getLoaderByStoreName("my-ai-plugin")
//   result, err := loader.Query(map[string]string{
//       "method": "ai.generate",
//       "model":  "gpt-4",
//       "prompt": "Hello world",
//       "config": `{"temperature": 0.7}`,
//   })
//   content := result.Pairs["content"]

// Documentation structures (for reference only, actual types are generated from proto)
// See server.proto for the actual message definitions:
//
// AIRequest fields:
//   - plugin_name: AI plugin name
//   - model: Model identifier (e.g., "gpt-4", "claude")
//   - prompt: The prompt or instruction
//   - config: JSON configuration string (optional)
//
// AIResponse fields:
//   - content: Generated content
//   - meta: JSON metadata string (optional)
//   - success: Whether the call succeeded
//   - error: Error message if failed
//
// AICapabilitiesResponse fields:
//   - models: Supported models
//   - features: Supported features
//   - description: Plugin description
//   - version: Plugin version
//   - success: Whether the call succeeded
//   - error: Error message if failed
