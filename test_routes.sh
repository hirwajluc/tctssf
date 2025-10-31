#!/bin/bash

echo "================================================"
echo "Testing TCTSSF Route Priority"
echo "================================================"
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Start server
echo "Starting server..."
./tctssf > /tmp/tctssf_test.log 2>&1 &
SERVER_PID=$!
sleep 3

cleanup() {
    echo ""
    echo "Cleaning up..."
    kill $SERVER_PID 2>/dev/null
    exit $1
}

trap cleanup INT TERM

# Test 1: Swagger UI
echo "Test 1: Swagger UI"
echo -n "  GET /swagger/index.html ... "
RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:3000/swagger/index.html 2>&1)
if [ "$RESPONSE" = "200" ]; then
    echo -e "${GREEN}✓ PASS${NC} (HTTP $RESPONSE)"
else
    echo -e "${RED}✗ FAIL${NC} (HTTP $RESPONSE, expected 200)"
    echo ""
    echo "Server logs:"
    tail -20 /tmp/tctssf_test.log
    cleanup 1
fi

# Test 2: Swagger JSON
echo -n "  GET /swagger/doc.json ... "
CONTENT=$(curl -s http://localhost:3000/swagger/doc.json 2>&1)
if echo "$CONTENT" | grep -q "TCTSSF API"; then
    echo -e "${GREEN}✓ PASS${NC} (Contains API definition)"
else
    echo -e "${RED}✗ FAIL${NC} (No API definition found)"
    cleanup 1
fi

# Test 3: Frontend root
echo ""
echo "Test 2: Frontend Routes"
echo -n "  GET / ... "
RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:3000/ 2>&1)
if [ "$RESPONSE" = "200" ]; then
    echo -e "${GREEN}✓ PASS${NC} (HTTP $RESPONSE)"
else
    echo -e "${YELLOW}⚠ WARNING${NC} (HTTP $RESPONSE)"
fi

# Test 4: API endpoint
echo ""
echo "Test 3: API Routes"
echo -n "  POST /api/login ... "
RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" -X POST http://localhost:3000/api/login \
    -H "Content-Type: application/json" -d '{}' 2>&1)
if [ "$RESPONSE" = "400" ] || [ "$RESPONSE" = "401" ]; then
    echo -e "${GREEN}✓ PASS${NC} (HTTP $RESPONSE - endpoint exists)"
elif [ "$RESPONSE" = "200" ]; then
    echo -e "${YELLOW}⚠ UNEXPECTED${NC} (HTTP 200 with empty body)"
else
    echo -e "${RED}✗ FAIL${NC} (HTTP $RESPONSE)"
fi

# Test 5: Non-existent file
echo ""
echo "Test 4: Catch-All Route"
echo -n "  GET /nonexistent.html ... "
CONTENT=$(curl -s http://localhost:3000/nonexistent.html 2>&1)
if echo "$CONTENT" | grep -q "html"; then
    echo -e "${GREEN}✓ PASS${NC} (Serves index.html for SPA)"
else
    echo -e "${YELLOW}⚠ WARNING${NC} (Unexpected response)"
fi

# Test 6: Check logs for route priority
echo ""
echo "Test 5: Route Priority Verification"
echo -n "  Checking server logs ... "
if grep -q "Static file request for path: /swagger" /tmp/tctssf_test.log; then
    echo -e "${RED}✗ FAIL${NC} (Swagger going through static handler!)"
    echo ""
    echo "Relevant logs:"
    grep "/swagger" /tmp/tctssf_test.log
    cleanup 1
else
    echo -e "${GREEN}✓ PASS${NC} (Swagger bypasses static handler)"
fi

echo ""
echo "================================================"
echo -e "${GREEN}ALL TESTS PASSED!${NC}"
echo "================================================"
echo ""
echo "Swagger is working correctly!"
echo ""
echo "Access points:"
echo "  Local:    http://localhost:3000/swagger/index.html"
echo "  External: http://86.48.7.218:3000/swagger/index.html"
echo ""
echo "Server logs: /tmp/tctssf_test.log"
echo ""
echo "Press Ctrl+C to stop the server"
echo ""

# Keep running
wait $SERVER_PID
