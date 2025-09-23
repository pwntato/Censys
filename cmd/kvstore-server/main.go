package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"censys-kvstore/proto"

	"google.golang.org/grpc"
)

// In-memory key-value store implementation
type kvStore struct {
	proto.UnimplementedKeyValueStoreServer
	mu   sync.RWMutex
	data map[string]string
}

// NewKVStore creates a new key-value store instance
func NewKVStore() *kvStore {
	return &kvStore{
		data: make(map[string]string),
	}
}

// Set stores a value at the given key
func (k *kvStore) Set(ctx context.Context, req *proto.SetRequest) (*proto.SetResponse, error) {
	k.mu.Lock()
	defer k.mu.Unlock()
	
	k.data[req.Key] = req.Value
	return &proto.SetResponse{
		Success: true,
		Message: fmt.Sprintf("Key '%s' set successfully", req.Key),
	}, nil
}

// Get retrieves the value for the given key
func (k *kvStore) Get(ctx context.Context, req *proto.GetRequest) (*proto.GetResponse, error) {
	k.mu.RLock()
	defer k.mu.RUnlock()
	
	value, exists := k.data[req.Key]
	if !exists {
		return &proto.GetResponse{
			Success: false,
			Value:   "",
			Message: fmt.Sprintf("Key '%s' not found", req.Key),
		}, nil
	}
	
	return &proto.GetResponse{
		Success: true,
		Value:   value,
		Message: fmt.Sprintf("Key '%s' retrieved successfully", req.Key),
	}, nil
}

// Delete removes the given key
func (k *kvStore) Delete(ctx context.Context, req *proto.DeleteRequest) (*proto.DeleteResponse, error) {
	k.mu.Lock()
	defer k.mu.Unlock()
	
	_, exists := k.data[req.Key]
	if !exists {
		return &proto.DeleteResponse{
			Success: false,
			Message: fmt.Sprintf("Key '%s' not found", req.Key),
		}, nil
	}
	
	delete(k.data, req.Key)
	return &proto.DeleteResponse{
		Success: true,
		Message: fmt.Sprintf("Key '%s' deleted successfully", req.Key),
	}, nil
}

func main() {
	// Create the key-value store instance
	store := NewKVStore()
	
	// Create gRPC server
	grpcServer := grpc.NewServer()
	proto.RegisterKeyValueStoreServer(grpcServer, store)
	
	// Start listening on port 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	
	log.Println("Key-Value Store gRPC server starting on :50051")
	
	// Start the server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
