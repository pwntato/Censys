.PHONY: help build test clean run-docker run-local stop-docker proto

# Default target
help:
	@echo "Available targets:"
	@echo "  build        - Build both services"
	@echo "  test         - Run all tests"
	@echo "  test-unit    - Run unit tests only"
	@echo "  test-integration - Run integration tests (requires services running)"
	@echo "  proto        - Generate protobuf Go files"
	@echo "  run-docker   - Start services with Docker Compose"
	@echo "  stop-docker  - Stop Docker Compose services"
	@echo "  run-local    - Run services locally (requires Go)"
	@echo "  clean        - Clean build artifacts"
	@echo "  logs         - Show Docker Compose logs"

# Generate protobuf files
proto:
	@echo "Generating protobuf files..."
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/kvstore.proto

# Build both services
build: proto
	@echo "Building Key-Value Store service..."
	go build -o bin/kvstore-server ./cmd/kvstore-server
	@echo "Building API server..."
	go build -o bin/api-server ./cmd/api-server
	@echo "Build complete!"

# Run unit tests
test-unit:
	@echo "Running unit tests..."
	go test -v ./cmd/kvstore-server/...
	go test -v ./cmd/api-server/...

# Run integration tests (requires services to be running)
test-integration:
	@echo "Running integration tests..."
	go test -v ./integration_test.go

# Run all tests
test: test-unit
	@echo "Unit tests completed. Run 'make test-integration' after starting services with 'make run-docker'"

# Start services with Docker Compose
run-docker:
	@echo "Starting services with Docker Compose..."
	docker-compose up --build -d
	@echo "Services started! API available at http://localhost:8080"
	@echo "gRPC service available at localhost:50051"

# Stop Docker Compose services
stop-docker:
	@echo "Stopping Docker Compose services..."
	docker-compose down

# Run services locally (requires Go and protobuf compiler)
run-local: proto
	@echo "Starting Key-Value Store service..."
	go run cmd/kvstore-server/main.go &
	@echo "Waiting for gRPC service to start..."
	sleep 3
	@echo "Starting API server..."
	go run cmd/api-server/main.go

# Show Docker Compose logs
logs:
	docker-compose logs -f

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	go clean
	docker-compose down --volumes --remove-orphans

# Test the API endpoints
test-api:
	@echo "Testing API endpoints..."
	@echo "Health check:"
	curl -s http://localhost:8080/health | jq .
	@echo "\nSetting a key-value pair:"
	curl -s -X POST http://localhost:8080/kv/set \
		-H "Content-Type: application/json" \
		-d '{"key": "test-key", "value": "test-value"}' | jq .
	@echo "\nGetting the value:"
	curl -s http://localhost:8080/kv/get/test-key | jq .
	@echo "\nDeleting the key:"
	curl -s -X DELETE http://localhost:8080/kv/delete/test-key | jq .
	@echo "\nVerifying deletion:"
	curl -s http://localhost:8080/kv/get/test-key | jq .

# Run comprehensive API tests
test-api-comprehensive:
	@echo "Running comprehensive API tests..."
	./test_api.sh

# Full test suite
test-full: run-docker
	@echo "Waiting for services to be ready..."
	sleep 10
	@echo "Running full test suite..."
	make test-integration
	make test-api
	make stop-docker
