#!/bin/bash
set -e

BASE_URL="http://localhost:8080"

echo "=========================================="
echo "Setting up Demo User for API Testing"
echo "=========================================="
echo ""

# Demo user credentials
EMAIL="demo@vms.com"
PASSWORD="demo123"
NAME="Demo User"

echo "1. Creating demo user..."
REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/register" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\",\"name\":\"$NAME\"}" 2>&1)

if echo "$REGISTER_RESPONSE" | grep -q "id"; then
  echo "✓ User created successfully"
elif echo "$REGISTER_RESPONSE" | grep -q "duplicate key"; then
  echo "✓ User already exists (skipping creation)"
else
  echo "Note: $REGISTER_RESPONSE"
fi

echo ""
echo "2. Logging in to get JWT token..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/login" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}")

TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.data.token // .token')

if [ "$TOKEN" == "null" ] || [ -z "$TOKEN" ]; then
  echo "✗ Failed to get token"
  echo "Response: $LOGIN_RESPONSE"
  exit 1
fi

echo "✓ Token obtained successfully"
echo ""
echo "=========================================="
echo "Demo User Credentials:"
echo "=========================================="
echo "Email:    $EMAIL"
echo "Password: $PASSWORD"
echo ""
echo "JWT Token (valid for 24 hours):"
echo "$TOKEN"
echo ""
echo "=========================================="
echo "Saving token to .demo-token file..."
echo -n "$TOKEN" > .demo-token
echo "✓ Token saved to .demo-token"
echo ""
echo "You can now use this token in API requests:"
echo "  curl -H 'Authorization: Bearer $TOKEN' $BASE_URL/api/v1/dashboard"
echo ""
echo "Or run the updated test-api.sh script which will use this token automatically."
echo "=========================================="
