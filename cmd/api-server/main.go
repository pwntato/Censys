package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"censys-kvstore/proto"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// APIServer handles HTTP requests and forwards them to the gRPC service
type APIServer struct {
	grpcClient proto.KeyValueStoreClient
}

// NewAPIServer creates a new API server instance
func NewAPIServer(grpcAddr string) (*APIServer, error) {
	// Connect to the gRPC server
	conn, err := grpc.Dial(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := proto.NewKeyValueStoreClient(conn)
	return &APIServer{grpcClient: client}, nil
}

// SetRequest represents the JSON request body for setting a key-value pair
type SetRequest struct {
	Key   string `json:"key" binding:"required"`
	Value string `json:"value" binding:"required"`
}

// SetResponse represents the JSON response for setting a key-value pair
type SetResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// GetResponse represents the JSON response for getting a value
type GetResponse struct {
	Success bool   `json:"success"`
	Value   string `json:"value,omitempty"`
	Message string `json:"message"`
}

// DeleteResponse represents the JSON response for deleting a key
type DeleteResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// Set handles POST /kv/set
func (s *APIServer) Set(c *gin.Context) {
	var req SetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if gRPC client is available (for testing)
	if s.grpcClient == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gRPC client not available"})
		return
	}

	// Call gRPC service
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	grpcResp, err := s.grpcClient.Set(ctx, &proto.SetRequest{
		Key:   req.Key,
		Value: req.Value,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	status := http.StatusOK
	if !grpcResp.Success {
		status = http.StatusBadRequest
	}

	c.JSON(status, SetResponse{
		Success: grpcResp.Success,
		Message: grpcResp.Message,
	})
}

// Get handles GET /kv/get/:key
func (s *APIServer) Get(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "key parameter is required"})
		return
	}

	// Check if gRPC client is available (for testing)
	if s.grpcClient == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gRPC client not available"})
		return
	}

	// Call gRPC service
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	grpcResp, err := s.grpcClient.Get(ctx, &proto.GetRequest{Key: key})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	status := http.StatusOK
	if !grpcResp.Success {
		status = http.StatusNotFound
	}

	c.JSON(status, GetResponse{
		Success: grpcResp.Success,
		Value:   grpcResp.Value,
		Message: grpcResp.Message,
	})
}

// Delete handles DELETE /kv/delete/:key
func (s *APIServer) Delete(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "key parameter is required"})
		return
	}

	// Check if gRPC client is available (for testing)
	if s.grpcClient == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gRPC client not available"})
		return
	}

	// Call gRPC service
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	grpcResp, err := s.grpcClient.Delete(ctx, &proto.DeleteRequest{Key: key})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	status := http.StatusOK
	if !grpcResp.Success {
		status = http.StatusNotFound
	}

	c.JSON(status, DeleteResponse{
		Success: grpcResp.Success,
		Message: grpcResp.Message,
	})
}

// Health handles GET /health
func (s *APIServer) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "status": "healthy"})
}

func main() {
	// Create API server
	apiServer, err := NewAPIServer("kvstore-server:50051")
	if err != nil {
		log.Fatalf("Failed to create API server: %v", err)
	}

	// Set Gin to release mode for production
	gin.SetMode(gin.ReleaseMode)

	// Create Gin router
	router := gin.Default()

	// Add CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Define routes
	router.GET("/health", apiServer.Health)
	router.POST("/kv/set", apiServer.Set)
	router.GET("/kv/get/:key", apiServer.Get)
	router.DELETE("/kv/delete/:key", apiServer.Delete)

	// Start server
	log.Println("API server starting on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start API server: %v", err)
	}
}
