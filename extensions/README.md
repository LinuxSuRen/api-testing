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

The AI extension provides intelligent testing capabilities through gRPC interface on port 50051. It offers a unified AI interface that can handle various content generation tasks.

### Available Services

#### GenerateContent (Recommended)
A general-purpose AI interface that can handle multiple content types:
- **SQL Generation**: Convert natural language to SQL queries
- **Test Case Generation**: Automatically generate test cases based on API specifications
- **Mock Service Creation**: Generate mock services and responses
- **Test Data Generation**: Create realistic test data using AI models
- **Result Analysis**: Intelligent analysis of test results and failure patterns

**Request Parameters:**
- `prompt`: Natural language description of what to generate
- `contentType`: Type of content ("sql", "testcase", "mock", "analysis")
- `context`: Additional context information (schemas, API specs, etc.)
- `sessionId`: Optional session ID for conversational context
- `parameters`: AI model configuration (model_provider, model_name, api_key, etc.)

#### GenerateSQLFromNaturalLanguage (Legacy)
Specific interface for SQL generation, maintained for backward compatibility.

### Configuration Parameters

These parameters are configured in the store extension settings and passed to the AI service:

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
```

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
