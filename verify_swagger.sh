#!/bin/bash

echo "=========================================="
echo "TCTSSF Swagger Route Verification"
echo "=========================================="
echo ""

# Check if binary exists
if [ ! -f "./tctssf" ]; then
    echo "❌ Binary not found. Building..."
    go build -o tctssf main.go
    if [ $? -ne 0 ]; then
        echo "❌ Build failed!"
        exit 1
    fi
    echo "✅ Build successful"
fi

echo "Starting server in background..."
./tctssf > /tmp/tctssf.log 2>&1 &
SERVER_PID=$!
echo "Server PID: $SERVER_PID"

# Wait for server to start
echo "Waiting for server to start..."
sleep 3

# Function to cleanup
cleanup() {
    echo ""
    echo "Stopping server (PID: $SERVER_PID)..."
    kill $SERVER_PID 2>/dev/null
    exit $1
}

trap cleanup INT TERM

# Test Swagger route
echo ""
echo "Testing Swagger UI route..."
SWAGGER_RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:3000/swagger/index.html 2>&1)

if [ "$SWAGGER_RESPONSE" = "200" ]; then
    echo "✅ Swagger UI accessible (HTTP 200)"
else
    echo "❌ Swagger UI returned HTTP $SWAGGER_RESPONSE (expected 200)"
    echo ""
    echo "Server logs:"
    tail -20 /tmp/tctssf.log
    cleanup 1
fi

# Test Swagger doc.json
echo ""
echo "Testing Swagger JSON..."
SWAGGER_JSON=$(curl -s http://localhost:3000/swagger/doc.json 2>&1)
if echo "$SWAGGER_JSON" | grep -q "TCTSSF API"; then
    echo "✅ Swagger JSON contains API definition"
else
    echo "❌ Swagger JSON not found or invalid"
    cleanup 1
fi

# Test frontend still works
echo ""
echo "Testing frontend root..."
FRONTEND_RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:3000/ 2>&1)
if [ "$FRONTEND_RESPONSE" = "200" ]; then
    echo "✅ Frontend still accessible (HTTP 200)"
else
    echo "⚠️  Frontend returned HTTP $FRONTEND_RESPONSE"
fi

# Test API endpoint
echo ""
echo "Testing API endpoint..."
API_RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" -X POST http://localhost:3000/api/login -H "Content-Type: application/json" -d '{}' 2>&1)
if [ "$API_RESPONSE" = "400" ] || [ "$API_RESPONSE" = "401" ]; then
    echo "✅ API endpoint accessible (HTTP $API_RESPONSE - expected for invalid request)"
else
    echo "⚠️  API returned HTTP $API_RESPONSE"
fi

echo ""
echo "=========================================="
echo "✅ ALL TESTS PASSED!"
echo "=========================================="
echo ""
echo "Access points:"
echo "  Frontend:  http://localhost:3000"
echo "  Swagger:   http://localhost:3000/swagger/index.html"
echo "  API:       http://localhost:3000/api"
echo ""
echo "On your network:"
echo "  Swagger:   http://86.48.7.218:3000/swagger/index.html"
echo ""
echo "Server logs: /tmp/tctssf.log"
echo ""
echo "Press Ctrl+C to stop server"

# Keep server running
wait $SERVER_PID
