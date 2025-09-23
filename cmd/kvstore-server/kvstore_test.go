package main

import (
	"context"
	"testing"

	"censys-kvstore/proto"
)

func TestKVStore_Set(t *testing.T) {
	store := NewKVStore()
	ctx := context.Background()

	tests := []struct {
		name     string
		key      string
		value    string
		expected bool
	}{
		{
			name:     "Set valid key-value pair",
			key:      "test-key",
			value:    "test-value",
			expected: true,
		},
		{
			name:     "Set empty key",
			key:      "",
			value:    "test-value",
			expected: true,
		},
		{
			name:     "Set empty value",
			key:      "test-key",
			value:    "",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &proto.SetRequest{
				Key:   tt.key,
				Value: tt.value,
			}
			resp, err := store.Set(ctx, req)
			if err != nil {
				t.Fatalf("Set() error = %v", err)
			}
			if resp.Success != tt.expected {
				t.Errorf("Set() success = %v, expected %v", resp.Success, tt.expected)
			}
		})
	}
}

func TestKVStore_Get(t *testing.T) {
	store := NewKVStore()
	ctx := context.Background()

	// Set up test data
	store.Set(ctx, &proto.SetRequest{Key: "existing-key", Value: "existing-value"})

	tests := []struct {
		name           string
		key            string
		expectedValue  string
		expectedSuccess bool
	}{
		{
			name:           "Get existing key",
			key:            "existing-key",
			expectedValue:  "existing-value",
			expectedSuccess: true,
		},
		{
			name:           "Get non-existing key",
			key:            "non-existing-key",
			expectedValue:  "",
			expectedSuccess: false,
		},
		{
			name:           "Get empty key",
			key:            "",
			expectedValue:  "",
			expectedSuccess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &proto.GetRequest{Key: tt.key}
			resp, err := store.Get(ctx, req)
			if err != nil {
				t.Fatalf("Get() error = %v", err)
			}
			if resp.Success != tt.expectedSuccess {
				t.Errorf("Get() success = %v, expected %v", resp.Success, tt.expectedSuccess)
			}
			if resp.Value != tt.expectedValue {
				t.Errorf("Get() value = %v, expected %v", resp.Value, tt.expectedValue)
			}
		})
	}
}

func TestKVStore_Delete(t *testing.T) {
	store := NewKVStore()
	ctx := context.Background()

	// Set up test data
	store.Set(ctx, &proto.SetRequest{Key: "existing-key", Value: "existing-value"})

	tests := []struct {
		name           string
		key            string
		expectedSuccess bool
	}{
		{
			name:           "Delete existing key",
			key:            "existing-key",
			expectedSuccess: true,
		},
		{
			name:           "Delete non-existing key",
			key:            "non-existing-key",
			expectedSuccess: false,
		},
		{
			name:           "Delete empty key",
			key:            "",
			expectedSuccess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &proto.DeleteRequest{Key: tt.key}
			resp, err := store.Delete(ctx, req)
			if err != nil {
				t.Fatalf("Delete() error = %v", err)
			}
			if resp.Success != tt.expectedSuccess {
				t.Errorf("Delete() success = %v, expected %v", resp.Success, tt.expectedSuccess)
			}
		})
	}
}

func TestKVStore_ConcurrentAccess(t *testing.T) {
	store := NewKVStore()
	ctx := context.Background()

	// Test concurrent writes
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(i int) {
			key := "key" + string(rune(i))
			value := "value" + string(rune(i))
			_, err := store.Set(ctx, &proto.SetRequest{Key: key, Value: value})
			if err != nil {
				t.Errorf("Concurrent Set() error = %v", err)
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify all values were set correctly
	for i := 0; i < 10; i++ {
		key := "key" + string(rune(i))
		expectedValue := "value" + string(rune(i))
		resp, err := store.Get(ctx, &proto.GetRequest{Key: key})
		if err != nil {
			t.Fatalf("Get() error = %v", err)
		}
		if !resp.Success || resp.Value != expectedValue {
			t.Errorf("Get() for key %s = %v, expected success=true, value=%s", key, resp, expectedValue)
		}
	}
}
