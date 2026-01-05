#!/bin/bash
# Check if API server is running and healthy

API_URL="${1:-http://localhost:8080}"

echo "🏥 Checking API health..."
echo "URL: $API_URL/health"
echo ""

# Attempt health check
response=$(curl -s -w "\n%{http_code}" "$API_URL/health" 2>/dev/null)

# Check if curl succeeded
if [ $? -ne 0 ]; then
    echo "❌ Failed to connect to API"
    echo ""
    echo "Is the server running?"
    echo "  Start with: make run"
    echo "  Or: go run cmd/api/main.go"
    exit 1
fi

# Parse response
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | head -n-1)

# Check status code
if [ "$http_code" = "200" ]; then
    echo "✅ API is healthy"
    echo "Status: $http_code"
    echo "Response: $body"
    exit 0
else
    echo "❌ API returned unexpected status"
    echo "Status: $http_code"
    echo "Response: $body"
    exit 1
fi
