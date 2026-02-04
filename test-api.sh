#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== API Key Authentication Test ===${NC}\n"

# 1. Create API key
echo -e "${GREEN}1. Create API key${NC}"
RESPONSE=$(curl -s -X POST http://localhost:8080/admin/keys \
  -d '{"name":"Test Key"}')
echo "$RESPONSE" | python3 -m json.tool

# Extract the API key from response
API_KEY=$(echo "$RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['key'])" 2>/dev/null)

if [ -z "$API_KEY" ]; then
  echo -e "${RED}Failed to create API key${NC}"
  exit 1
fi

echo -e "\n${BLUE}API Key: $API_KEY${NC}\n"

# 2. List all keys (should show masked key)
echo -e "${GREEN}2. List all keys (masked)${NC}"
curl -s http://localhost:8080/admin/keys | python3 -m json.tool
echo ""

# 3. Try to create product WITHOUT API key (should fail)
echo -e "\n${GREEN}3. Create product WITHOUT key (should fail 401)${NC}"
curl -s -w "\nHTTP Status: %{http_code}\n" \
  -X POST http://localhost:8080/v1/products \
  -d '{"id":"1","name":"MacBook","price":2499}'
echo ""

# 4. Create product WITH API key (should work)
echo -e "\n${GREEN}4. Create product WITH key (should work)${NC}"
curl -s -X POST http://localhost:8080/v1/products \
  -H "X-API-Key: $API_KEY" \
  -d '{"id":"1","name":"MacBook Pro","price":2499.99}' | python3 -m json.tool
echo ""

# 5. Create another product
echo -e "${GREEN}5. Create another product${NC}"
curl -s -X POST http://localhost:8080/v1/products \
  -H "X-API-Key: $API_KEY" \
  -d '{"id":"2","name":"iPhone 16","price":999.99}' | python3 -m json.tool
echo ""

# 6. Get products (public, no key needed)
echo -e "\n${GREEN}6. List products (public, no key needed)${NC}"
curl -s http://localhost:8080/v1/products | python3 -m json.tool
echo ""

# 7. Check usage stats
echo -e "\n${GREEN}7. Check API key usage stats${NC}"
curl -s http://localhost:8080/admin/keys | python3 -m json.tool
echo ""

# 8. Try to delete WITHOUT key (should fail)
echo -e "\n${GREEN}8. Delete product WITHOUT key (should fail 401)${NC}"
curl -s -w "\nHTTP Status: %{http_code}\n" \
  -X DELETE http://localhost:8080/v1/products/1
echo ""

# 9. Delete WITH key (should work)
echo -e "\n${GREEN}9. Delete product WITH key (should work)${NC}"
curl -s -w "\nHTTP Status: %{http_code}\n" \
  -X DELETE http://localhost:8080/v1/products/1 \
  -H "X-API-Key: $API_KEY"
echo ""

# 10. Final usage stats
echo -e "\n${GREEN}10. Final usage stats (should show 3 requests)${NC}"
curl -s http://localhost:8080/admin/keys | python3 -m json.tool

echo -e "\n${BLUE}=== Test Complete ===${NC}"
