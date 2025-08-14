Ports in extensions:

| Type | Name                                                                     | Port |
|------|--------------------------------------------------------------------------|------|
| Store | [orm](https://github.com/LinuxSuRen/atest-ext-store-orm)                 | 4071 |
| Store | [s3](https://github.com/LinuxSuRen/atest-ext-store-s3)                   | 4072 |
| Store | [etcd](https://github.com/LinuxSuRen/atest-ext-store-etcd)               | 4073 |
| Store | [git](https://github.com/LinuxSuRen/atest-ext-store-git)                 | 4074 |
| Store | [mongodb](https://github.com/LinuxSuRen/atest-ext-store-mongodb)         | 4075 |
| Store | [redis](https://github.com/LinuxSuRen/atest-ext-store-redis)             |  |
| Store | [iotdb](https://github.com/LinuxSuRen/atest-ext-store-iotdb) | |
| Store | [Cassandra](https://github.com/LinuxSuRen/atest-ext-store-cassandra) | |
| Store | [Elasticsearch](https://github.com/LinuxSuRen/atest-ext-store-elasticsearch) | |
| Monitor | [docker-monitor](https://github.com/LinuxSuRen/atest-ext-monitor-docker) |  |
| Agent | [collector](https://github.com/LinuxSuRen/atest-ext-collector)           |  |
| Secret | [Vault](https://github.com/LinuxSuRen/api-testing-vault-extension)       | |
| Data | [Swagger](https://github.com/LinuxSuRen/atest-ext-data-swagger) | |
| AI | [ai-extension](https://github.com/LinuxSuRen/atest-ext-ai) | 50051 |

## AI Extension Interface

This extension provides AI-powered content generation capabilities through a unified gRPC interface.

### Unified Interface Design

#### `GenerateContent`
The single, unified interface for all AI content generation tasks. This interface can handle various content types including SQL generation, test case writing, mock service creation, and more through a single endpoint.

**Request Parameters:**
- `prompt`: Natural language description of what you want to generate
- `contentType`: Type of content to generate (e.g., "sql", "testcase", "mock")
- `context`: Additional context information as key-value pairs
- `parameters`: Additional parameters specific to the content type

**Response:**
- `success.content`: Generated content
- `success.explanation`: Explanation of the generated content
- `success.confidenceScore`: Confidence score (0.0 to 1.0)
- `success.metadata`: Additional metadata about the generation
- `error`: Error information if generation fails

### Content Types and Task Identification

The `contentType` parameter serves as the task identifier, allowing the system to:

1. **Apply appropriate prompting strategies** for different content types
2. **Use specialized processing logic** for each task type
3. **Provide task-specific validation and post-processing**
4. **Generate relevant metadata** for each content type

#### Supported Content Types:

1. **SQL Generation** (`contentType: "sql"`)
   - Generates SQL queries from natural language
   - Supports multiple database types via parameters
   - Context can include schema information
   - Applies SQL-specific validation and dialect adaptation

2. **Test Case Generation** (`contentType: "testcase"`)
   - Generates comprehensive test cases
   - Context can include code specifications
   - Uses test-specific prompting strategies

3. **Mock Service Generation** (`contentType: "mock"`)
   - Generates mock services and API responses
   - Context can include API specifications
   - Applies service-specific formatting

4. **Generic Content** (any other `contentType`)
   - Handles any other content generation requests
   - Flexible prompt-based generation
   - Extensible for future content types

### Architecture Benefits

**Single Interface Advantages:**
- **Consistency**: All content generation follows the same request/response pattern
- **Extensibility**: New content types can be added without API changes
- **Simplicity**: Clients only need to integrate with one interface
- **Flexibility**: Rich context and parameter support for all content types

**Task Differentiation:**
While using a single interface, the system internally routes requests to specialized handlers based on `contentType`, ensuring:
- Appropriate AI model prompting for each task
- Task-specific validation and processing
- Relevant metadata and confidence scoring
- Optimized performance for different content types

### Configuration Parameters

The following parameters can be passed via `GenerateContentRequest.parameters`:

- `model_provider`: AI provider ("ollama", "openai") - used in GenerateContentRequest.parameters
- `model_name`: Specific model name (e.g., "llama3.2", "gpt-4") - passed via parameters["model_name"]
- `api_key`: Authentication key for external providers - used in parameters["api_key"]
- `base_url`: Service endpoint URL - specified in parameters["base_url"]
- `temperature`: Response randomness control (0.0-1.0) - set via parameters["temperature"]
- `max_tokens`: Maximum response length - controlled by parameters["max_tokens"]

### Usage Examples

```protobuf
// Generate SQL query
GenerateContentRequest {
  prompt: "Show me all users from California"
  contentType: "sql"
  context: {
    "schema": "CREATE TABLE users (id INT, name VARCHAR(100), state VARCHAR(50))"
  }
  parameters: {
    "model_provider": "ollama"
    "model_name": "llama3.2"
    "database_type": "mysql"
  }
}

// Generate test case
GenerateContentRequest {
  prompt: "Create a test case for user registration API"
  contentType: "testcase"
  context: {
    "api_spec": "POST /api/users {name, email, password}"
  }
}

// Generate mock service
GenerateContentRequest {
  prompt: "Create a mock REST API for user management"
  contentType: "mock"
  context: {
    "endpoints": "GET /users, POST /users, PUT /users/{id}, DELETE /users/{id}"
    "response_format": "JSON"
  }
}
```

### Migration Notes

This version removes the legacy `GenerateSQLFromNaturalLanguage` interface in favor of the unified `GenerateContent` approach. All SQL generation should now use:

```protobuf
GenerateContentRequest {
  prompt: "your natural language query"
  contentType: "sql"
  parameters: {"database_type": "mysql|postgresql|sqlite"}
  context: {"schemas": "table definitions"}
}
```

### Frequently Asked Questions

**Q1: Do I need different interfaces for different tasks (SQL generation, test case writing, mock service creation)?**

A: No, you only need the single `GenerateContent` interface. The `contentType` parameter tells the system what kind of content you want to generate, and the system internally routes your request to the appropriate specialized handler. This design provides:
- **Task identification**: The system knows exactly what you're trying to accomplish
- **Specialized processing**: Each content type gets optimized prompting and validation
- **Consistent API**: You always use the same interface regardless of the task
- **Future extensibility**: New content types can be added without changing the API

**Q2: Will removing the SQL-specific interface affect other files?**

A: The removal has been carefully implemented to ensure minimal impact:
- **Protobuf files**: All generated code (Go, Python, gRPC-gateway, OpenAPI) has been updated
- **Implementation files**: The Python gRPC server has been updated to use the new unified interface
- **Backward compatibility**: While the old interface is removed, the same functionality is available through `GenerateContent` with `contentType: "sql"`
- **No breaking changes**: The core functionality remains the same, just accessed through a cleaner, unified interface

The migration ensures that all SQL generation capabilities are preserved while providing a more maintainable and extensible architecture.

## Contribute a new extension

* First, create a repository. And please keep the same naming convertion.
* Second, implement the `Loader` gRPC service which defined by [this proto](../pkg/testing/remote/loader.proto).
* Finally, add the extension's name into function [SupportedExtensions](../console/atest-ui/src/views/store.ts).

## Naming conventions

Please follow the following conventions if you want to add a new store extension:

`store-xxx`

`xxx` should be a type of a backend storage.

## Test

First, build and copy the binary file into the system path. You can run the following
command in the root directory of this project:

```shell
make build-ext-etcd copy-ext
```
