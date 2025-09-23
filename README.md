# Censys Key-Value Store

A decomposed Key-Value store implementation with two services communicating over gRPC.

## Architecture

This project consists of two microservices:

1. **Key-Value Store Service (Backend)**: A gRPC service that implements the core key-value storage functionality
2. **API Server (Frontend)**: A JSON REST API service that communicates with the backend via gRPC

## Features

- **Set**: Store a value at a given key
- **Get**: Retrieve the value for a given key
- **Delete**: Delete a given key
- **Health Check**: Monitor service health
- **Concurrent Access**: Thread-safe operations
- **Docker Support**: Containerized deployment
- **CORS Support**: Cross-origin resource sharing enabled
- **Comprehensive Testing**: Unit and integration tests

## API Endpoints

### REST API (Port 8080)

- `GET /health` - Health check endpoint
- `POST /kv/set` - Set a key-value pair
- `GET /kv/get/:key` - Get value by key
- `DELETE /kv/delete/:key` - Delete a key

### gRPC API (Port 50051)

- `Set(SetRequest) returns (SetResponse)` - Store a key-value pair
- `Get(GetRequest) returns (GetResponse)` - Retrieve a value by key
- `Delete(DeleteRequest) returns (DeleteResponse)` - Delete a key

## Quick Start

### Using Docker Compose (Recommended)

1. Clone the repository:

```bash
git clone <repository-url>
cd Censys
```

2. Build and start the services:

```bash
docker compose up --build
```

3. Test the API:

```bash
# Health check
curl http://localhost:8080/health

# Set a key-value pair
curl -X POST http://localhost:8080/kv/set \
  -H "Content-Type: application/json" \
  -d '{"key": "test-key", "value": "test-value"}'

# Get the value
curl http://localhost:8080/kv/get/test-key

# Delete the key
curl -X DELETE http://localhost:8080/kv/delete/test-key
```

### Manual Build and Run

1. Install dependencies (requires Go 1.23.0 or later):

```bash
go mod download
```

2. Generate protobuf files:

```bash
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/kvstore.proto
```

3. Start the Key-Value Store service:

```bash
go run cmd/kvstore-server/main.go
```

4. In another terminal, start the API server:

```bash
go run cmd/api-server/main.go
```

## Testing

### Unit Tests

Run unit tests for both services:

```bash
# Test the Key-Value Store service
go test ./cmd/kvstore-server/...

# Test the API server
go test ./cmd/api-server/...
```

### Integration Tests

Run integration tests (requires both services to be running):

```bash
# Start services first
docker compose up -d

# Run integration tests
go test -v ./integration_test.go

# Stop services
docker compose down
```

### Automated API Testing

Use the provided test script for comprehensive API testing:

```bash
# Make the script executable
chmod +x test_api.sh

# Run the test script (requires services to be running)
./test_api.sh
```

### Manual Testing

1. **Set a value**:

```bash
curl -X POST http://localhost:8080/kv/set \
  -H "Content-Type: application/json" \
  -d '{"key": "my-key", "value": "my-value"}'
```

Expected response:

```json
{
  "success": true,
  "message": "Key 'my-key' set successfully"
}
```

2. **Get a value**:

```bash
curl http://localhost:8080/kv/get/my-key
```

Expected response:

```json
{
  "success": true,
  "value": "my-value",
  "message": "Key 'my-key' retrieved successfully"
}
```

3. **Delete a value**:

```bash
curl -X DELETE http://localhost:8080/kv/delete/my-key
```

Expected response:

```json
{
  "success": true,
  "message": "Key 'my-key' deleted successfully"
}
```

4. **Get non-existent key**:

```bash
curl http://localhost:8080/kv/get/non-existent
```

Expected response:

```json
{
  "success": false,
  "value": "",
  "message": "Key 'non-existent' not found"
}
```

## Docker Images

### Build individual images:

```bash
# Build Key-Value Store service
docker build -f Dockerfile.kvstore-server -t censys-kvstore-server .

# Build API server
docker build -f Dockerfile.api-server -t censys-api-server .
```

### Run individual containers:

```bash
# Start Key-Value Store service
docker run -p 50051:50051 censys-kvstore-server

# Start API server (in another terminal)
docker run -p 8080:8080 --network host censys-api-server
```

## Development

### Project Structure

```
.
├── cmd/
│   ├── kvstore-server/     # gRPC Key-Value Store service
│   └── api-server/         # REST API service
├── proto/                  # Protocol Buffer definitions
├── Dockerfile.kvstore-server
├── Dockerfile.api-server
├── docker-compose.yml
├── go.mod
├── go.sum
└── README.md
```

### Adding New Features

The architecture is designed to be extensible:

1. **Adding new KV operations**: Extend the protobuf definition in `proto/kvstore.proto`
2. **Adding new REST endpoints**: Add new handlers in `cmd/api-server/main.go`
3. **Changing transport protocol**: The gRPC layer can be easily replaced with other protocols

### Code Quality

- **Thread Safety**: All operations use proper locking mechanisms
- **Error Handling**: Comprehensive error handling and meaningful error messages
- **Logging**: Structured logging for debugging and monitoring
- **Health Checks**: Built-in health check endpoints for both services
- **Security**: Non-root user execution in Docker containers

## Performance Considerations

- **In-Memory Storage**: Current implementation uses in-memory storage for simplicity
- **Concurrent Access**: Thread-safe operations support concurrent reads and writes
- **Connection Pooling**: gRPC connections are reused efficiently
- **Timeouts**: All operations have appropriate timeouts to prevent hanging

## Monitoring and Observability

- **Health Endpoints**: Both services expose health check endpoints
- **Docker Health Checks**: Container health monitoring
- **Structured Logging**: Easy to parse logs for monitoring systems
- **Metrics Ready**: Architecture supports easy addition of metrics collection

## Troubleshooting

### Common Issues

1. **Port conflicts**: Ensure ports 8080 and 50051 are available
2. **gRPC connection issues**: Verify the Key-Value Store service is running before starting the API server
3. **Docker build failures**: Ensure Docker is running and has sufficient resources

### Debugging

1. **Check service logs**:

```bash
docker compose logs kvstore-server
docker compose logs api-server
```

2. **Test gRPC connection directly**:

```bash
grpcurl -plaintext localhost:50051 list
```

3. **Verify API endpoints**:

```bash
curl -v http://localhost:8080/health
```

## License

This project is part of a technical assessment and is provided as-is for evaluation purposes.
