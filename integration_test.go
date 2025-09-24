package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/pwntato/Censys/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Integration test helper functions
func startTestServers() error {
	// This would typically start the servers in separate goroutines
	// For integration tests, we assume servers are running
	return nil
}

func TestIntegration_SetAndGet(t *testing.T) {
	// Test gRPC connection
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Skipf("Skipping integration test: gRPC server not available: %v", err)
		return
	}
	defer conn.Close()

	client := proto.NewKeyValueStoreClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test Set
	setReq := &proto.SetRequest{
		Key:   "integration-test-key",
		Value: "integration-test-value",
	}
	setResp, err := client.Set(ctx, setReq)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	if !setResp.Success {
		t.Fatalf("Set failed: %s", setResp.Message)
	}

	// Test Get
	getReq := &proto.GetRequest{
		Key: "integration-test-key",
	}
	getResp, err := client.Get(ctx, getReq)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if !getResp.Success {
		t.Fatalf("Get failed: %s", getResp.Message)
	}
	if getResp.Value != "integration-test-value" {
		t.Fatalf("Expected value 'integration-test-value', got '%s'", getResp.Value)
	}

	// Test Delete
	deleteReq := &proto.DeleteRequest{
		Key: "integration-test-key",
	}
	deleteResp, err := client.Delete(ctx, deleteReq)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	if !deleteResp.Success {
		t.Fatalf("Delete failed: %s", deleteResp.Message)
	}

	// Verify deletion
	getResp2, err := client.Get(ctx, getReq)
	if err != nil {
		t.Fatalf("Get after delete failed: %v", err)
	}
	if getResp2.Success {
		t.Fatalf("Expected key to be deleted, but it still exists")
	}
}

func TestIntegration_HTTPAPI(t *testing.T) {
	baseURL := "http://localhost:8080"

	// Test health endpoint
	resp, err := http.Get(baseURL + "/health")
	if err != nil {
		t.Skipf("Skipping HTTP integration test: API server not available: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Health check failed with status: %d", resp.StatusCode)
	}

	// Test Set via HTTP
	setData := map[string]string{
		"key":   "http-test-key",
		"value": "http-test-value",
	}
	jsonData, _ := json.Marshal(setData)

	resp, err = http.Post(baseURL+"/kv/set", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("HTTP Set failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("HTTP Set failed with status: %d", resp.StatusCode)
	}

	// Test Get via HTTP
	resp, err = http.Get(baseURL + "/kv/get/http-test-key")
	if err != nil {
		t.Fatalf("HTTP Get failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("HTTP Get failed with status: %d", resp.StatusCode)
	}

	var getResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&getResp); err != nil {
		t.Fatalf("Failed to decode Get response: %v", err)
	}

	if getResp["value"] != "http-test-value" {
		t.Fatalf("Expected value 'http-test-value', got '%v'", getResp["value"])
	}

	// Test Delete via HTTP
	req, _ := http.NewRequest("DELETE", baseURL+"/kv/delete/http-test-key", nil)
	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("HTTP Delete failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("HTTP Delete failed with status: %d", resp.StatusCode)
	}

	// Verify deletion
	resp, err = http.Get(baseURL + "/kv/get/http-test-key")
	if err != nil {
		t.Fatalf("HTTP Get after delete failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("Expected 404 after delete, got status: %d", resp.StatusCode)
	}
}

func TestIntegration_EndToEnd(t *testing.T) {
	// This test requires both services to be running
	// It tests the complete flow: HTTP API -> gRPC -> Storage

	baseURL := "http://localhost:8080"

	// Test data
	testCases := []struct {
		key   string
		value string
	}{
		{"key1", "value1"},
		{"key2", "value2"},
		{"key3", "value3"},
	}

	// Set multiple values
	for _, tc := range testCases {
		setData := map[string]string{
			"key":   tc.key,
			"value": tc.value,
		}
		jsonData, _ := json.Marshal(setData)

		resp, err := http.Post(baseURL+"/kv/set", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatalf("Failed to set %s: %v", tc.key, err)
		}
		resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Failed to set %s, status: %d", tc.key, resp.StatusCode)
		}
	}

	// Get all values
	for _, tc := range testCases {
		resp, err := http.Get(baseURL + "/kv/get/" + tc.key)
		if err != nil {
			t.Fatalf("Failed to get %s: %v", tc.key, err)
		}
		resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Failed to get %s, status: %d", tc.key, resp.StatusCode)
		}
	}

	// Delete all values
	for _, tc := range testCases {
		req, _ := http.NewRequest("DELETE", baseURL+"/kv/delete/"+tc.key, nil)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to delete %s: %v", tc.key, err)
		}
		resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Failed to delete %s, status: %d", tc.key, resp.StatusCode)
		}
	}

	// Verify all values are deleted
	for _, tc := range testCases {
		resp, err := http.Get(baseURL + "/kv/get/" + tc.key)
		if err != nil {
			t.Fatalf("Failed to verify deletion of %s: %v", tc.key, err)
		}
		resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Fatalf("Expected %s to be deleted, but it still exists", tc.key)
		}
	}
}
