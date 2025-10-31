#!/bin/bash

echo "=========================================="
echo "Testing TCTSSF Server with Swagger"
echo "=========================================="

# Start the server in the background
echo "Starting server..."
./tctssf &
SERVER_PID=$!

# Wait for server to start
echo "Waiting for server to start..."
sleep 5

# Test if server is running
echo -e "\nTesting server health..."
if curl -s http://localhost:3000/api/login > /dev/null 2>&1; then
    echo "✓ Server is running"
else
    echo "✗ Server failed to start"
    kill $SERVER_PID 2>/dev/null
    exit 1
fi

# Test Swagger UI
echo -e "\nTesting Swagger UI..."
if curl -s http://localhost:3000/swagger/index.html | grep -q "swagger"; then
    echo "✓ Swagger UI is accessible"
else
    echo "✗ Swagger UI not found"
fi

# Test Swagger JSON
echo -e "\nTesting Swagger JSON..."
if curl -s http://localhost:3000/swagger/doc.json | grep -q "TCTSSF API"; then
    echo "✓ Swagger JSON is accessible"
else
    echo "✗ Swagger JSON not found"
fi

# Display access information
echo -e "\n=========================================="
echo "Server is running successfully!"
echo "=========================================="
echo "Frontend:     http://localhost:3000"
echo "API Docs:     http://localhost:3000/swagger/index.html"
echo "API Base:     http://localhost:3000/api"
echo "=========================================="
echo -e "\nPress Ctrl+C to stop the server"

# Wait for user interrupt
trap "echo 'Stopping server...'; kill $SERVER_PID 2>/dev/null; exit 0" INT
wait $SERVER_PID
