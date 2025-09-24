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

1. Clone and start services:

```bash
git clone <repository-url>
cd Censys
make run-docker
```

2. Test the API:

```bash
make test-api
```

3. Stop services:

```bash
make stop-docker
```

## Available Commands

| Command            | Description                        |
| ------------------ | ---------------------------------- |
| `make help`        | Show all commands                  |
| `make run-docker`  | Start services with Docker Compose |
| `make stop-docker` | Stop Docker Compose services       |
| `make test-all`    | Run complete test suite            |
| `make test-unit`   | Run unit tests                     |
| `make test-api`    | Test API endpoints                 |
| `make run-local`   | Run services locally (requires Go) |
| `make build`       | Build both services                |
| `make logs`        | Show service logs                  |
| `make clean`       | Clean build artifacts              |

### Environment Variables

The services can be configured using the following environment variables:

| Variable              | Default                | Description                                                 |
| --------------------- | ---------------------- | ----------------------------------------------------------- |
| `KVSTORE_PORT`        | `50051`                | Port for the Key-Value Store gRPC service                   |
| `API_PORT`            | `8080`                 | Port for the API Server HTTP service                        |
| `GRPC_SERVER_ADDRESS` | `kvstore-server:50051` | Address of the gRPC server for the API server to connect to |

You can set these variables in your environment or create a `.env` file in the project root:

```bash
# Example .env file
KVSTORE_PORT=50051
API_PORT=8080
GRPC_SERVER_ADDRESS=kvstore-server:50051
```

## Testing

```bash
# Run all tests
make test-all

# Individual test types
make test-unit        # Unit tests
make test-integration # Integration tests (requires services running)
make test-api         # API endpoint tests
```

## Docker

```bash
make run-docker  # Start services
make stop-docker # Stop services
make logs        # View logs
make clean       # Clean up
```

### Project Structure

```
cmd/
├── kvstore-server/  # gRPC service
└── api-server/      # REST API
proto/               # Protocol definitions
```

## Troubleshooting

```bash
make logs  # View service logs
curl http://localhost:8080/health  # Test API
```

## License

This project is part of a technical assessment and is provided as-is for evaluation purposes.
