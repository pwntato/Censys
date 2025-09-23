#!/bin/bash

# Test script for the Censys Key-Value Store API
# This script demonstrates all the API endpoints

echo "🚀 Testing Censys Key-Value Store API"
echo "======================================"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to test API endpoint
test_endpoint() {
    local method=$1
    local url=$2
    local data=$3
    local description=$4
    local expect_success=${5:-true}  # Default to expecting success
    
    echo -e "\n${YELLOW}Testing: $description${NC}"
    echo "Request: $method $url"
    
    if [ -n "$data" ]; then
        response=$(curl -s -X "$method" "$url" -H "Content-Type: application/json" -d "$data")
    else
        response=$(curl -s -X "$method" "$url")
    fi
    
    echo "Response: $response"
    
    # Check if response contains success based on expectation
    if [ "$expect_success" = "true" ]; then
        if echo "$response" | grep -q '"success":true'; then
            echo -e "${GREEN}✅ SUCCESS${NC}"
        else
            echo -e "${RED}❌ FAILED${NC}"
        fi
    else
        if echo "$response" | grep -q '"success":false' || echo "$response" | grep -q '"error"'; then
            echo -e "${GREEN}✅ SUCCESS${NC}"
        else
            echo -e "${RED}❌ FAILED${NC}"
        fi
    fi
}

# Wait for services to be ready
echo "⏳ Waiting for services to start..."
sleep 5

# Test health endpoint
test_endpoint "GET" "http://localhost:8080/health" "" "Health Check"

# Test setting a key-value pair
test_endpoint "POST" "http://localhost:8080/kv/set" '{"key": "demo-key", "value": "demo-value"}' "Set Key-Value Pair"

# Test getting the value
test_endpoint "GET" "http://localhost:8080/kv/get/demo-key" "" "Get Value by Key"

# Test setting another key-value pair
test_endpoint "POST" "http://localhost:8080/kv/set" '{"key": "another-key", "value": "another-value"}' "Set Another Key-Value Pair"

# Test getting the second value
test_endpoint "GET" "http://localhost:8080/kv/get/another-key" "" "Get Second Value"

# Test getting non-existent key
test_endpoint "GET" "http://localhost:8080/kv/get/non-existent" "" "Get Non-existent Key (should fail)" "false"

# Test deleting a key
test_endpoint "DELETE" "http://localhost:8080/kv/delete/demo-key" "" "Delete Key"

# Test getting deleted key (should fail)
test_endpoint "GET" "http://localhost:8080/kv/get/demo-key" "" "Get Deleted Key (should fail)" "false"

# Test deleting non-existent key (should fail)
test_endpoint "DELETE" "http://localhost:8080/kv/delete/non-existent" "" "Delete Non-existent Key (should fail)" "false"

# Test edge cases
echo -e "\n${YELLOW}Testing Edge Cases:${NC}"

# Test empty key
test_endpoint "POST" "http://localhost:8080/kv/set" '{"key": "", "value": "empty-key-value"}' "Set Empty Key" "false"

# Test empty value
test_endpoint "POST" "http://localhost:8080/kv/set" '{"key": "empty-value-key", "value": ""}' "Set Empty Value" "false"

# Test invalid JSON
echo -e "\n${YELLOW}Testing Invalid JSON:${NC}"
echo "Request: POST http://localhost:8080/kv/set with invalid JSON"
response=$(curl -s -X POST "http://localhost:8080/kv/set" -H "Content-Type: application/json" -d '{"key": "invalid", "value":}')
echo "Response: $response"
if echo "$response" | grep -q '"error"'; then
    echo -e "${GREEN}✅ SUCCESS (properly handled invalid JSON)${NC}"
else
    echo -e "${RED}❌ FAILED (should handle invalid JSON)${NC}"
fi

echo -e "\n${GREEN}🎉 API Testing Complete!${NC}"
echo "======================================"
