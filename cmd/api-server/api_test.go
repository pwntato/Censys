package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Mock API server for testing
	apiServer := &APIServer{}

	router.GET("/health", apiServer.Health)
	router.POST("/kv/set", apiServer.Set)
	router.GET("/kv/get/:key", apiServer.Get)
	router.DELETE("/kv/delete/:key", apiServer.Delete)

	return router
}

func TestHealthEndpoint(t *testing.T) {
	router := setupTestRouter()

	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["status"] != "healthy" {
		t.Errorf("Expected status 'healthy', got '%s'", response["status"])
	}

	if response["success"] != true {
		t.Errorf("Expected success true, got %v", response["success"])
	}
}

func TestSetEndpoint(t *testing.T) {
	router := setupTestRouter()

	tests := []struct {
		name           string
		requestBody    SetRequest
		expectedStatus int
	}{
		{
			name: "Valid set request",
			requestBody: SetRequest{
				Key:   "test-key",
				Value: "test-value",
			},
			expectedStatus: http.StatusInternalServerError, // Will fail due to no gRPC connection
		},
		{
			name: "Empty key",
			requestBody: SetRequest{
				Key:   "",
				Value: "test-value",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Empty value",
			requestBody: SetRequest{
				Key:   "test-key",
				Value: "",
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/kv/set", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestGetEndpoint(t *testing.T) {
	router := setupTestRouter()

	tests := []struct {
		name           string
		key            string
		expectedStatus int
	}{
		{
			name:           "Valid key",
			key:            "test-key",
			expectedStatus: http.StatusInternalServerError, // Will fail due to no gRPC connection
		},
		{
			name:           "Empty key",
			key:            "",
			expectedStatus: http.StatusNotFound, // No route match for /kv/get/
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Always append the key to the URL, even if empty
			// This ensures the route pattern matches
			url := "/kv/get/" + tt.key
			req, _ := http.NewRequest("GET", url, nil)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestDeleteEndpoint(t *testing.T) {
	router := setupTestRouter()

	tests := []struct {
		name           string
		key            string
		expectedStatus int
	}{
		{
			name:           "Valid key",
			key:            "test-key",
			expectedStatus: http.StatusInternalServerError, // Will fail due to no gRPC connection
		},
		{
			name:           "Empty key",
			key:            "",
			expectedStatus: http.StatusNotFound, // No route match for /kv/get/
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := "/kv/delete/" + tt.key
			req, _ := http.NewRequest("DELETE", url, nil)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestInvalidJSON(t *testing.T) {
	router := setupTestRouter()

	req, _ := http.NewRequest("POST", "/kv/set", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}
