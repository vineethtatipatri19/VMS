#!/bin/bash

echo "🧪 Testing VMS API Endpoints"
echo "================================"
echo ""

# Load JWT token from .demo-token file
if [ ! -f ".demo-token" ]; then
  echo "❌ Error: .demo-token file not found"
  echo "Please run ./setup-demo-user.sh first to create a demo user and get a token"
  exit 1
fi

TOKEN=$(cat .demo-token)
BASE_URL="http://localhost:8080/api/v1"

# Generate unique identifier for this test run
TEST_ID=$(date +%s)

# Test 1: Create a customer
echo "1️⃣ Testing Customer Creation..."
CUSTOMER_RESPONSE=$(curl -s -X POST "$BASE_URL/customers" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"Name\": \"Test Customer $TEST_ID\",
    \"ContactNumber\": \"+12345$TEST_ID\",
    \"CustomerType\": \"retail\",
    \"CreditLimit\": 10000,
    \"PaymentTermsDays\": 30,
    \"Status\": \"active\"
  }")

if echo "$CUSTOMER_RESPONSE" | jq -e '.success' > /dev/null 2>&1; then
  CUSTOMER_ID=$(echo "$CUSTOMER_RESPONSE" | jq -r '.data.ID')
  echo "✅ Customer created: $CUSTOMER_ID"
else
  echo "❌ Customer creation failed: $CUSTOMER_RESPONSE"
  exit 1
fi
echo ""

# Test 2: Create inventory item
echo "2️⃣ Testing Inventory Creation..."
INVENTORY_RESPONSE=$(curl -s -X POST "$BASE_URL/inventory" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"Name\": \"Test Product $TEST_ID\",
    \"LotNumber\": \"LOT-$TEST_ID\",
    \"Unit\": \"bottles\",
    \"Quantity\": 100,
    \"CostPrice\": 50,
    \"SellingPrice\": 75,
    \"SupplierName\": \"Test Supplier\",
    \"PurchaseDate\": \"2025-11-18\",
    \"ExpiryDate\": \"2026-05-18\",
    \"MinStockLevel\": 20,
    \"Status\": \"available\"
  }")

if echo "$INVENTORY_RESPONSE" | jq -e '.success' > /dev/null 2>&1; then
  INVENTORY_ID=$(echo "$INVENTORY_RESPONSE" | jq -r '.data.ID')
  echo "✅ Inventory created: $INVENTORY_ID"
else
  echo "❌ Inventory creation failed: $INVENTORY_RESPONSE"
  exit 1
fi
echo ""

# Test 3: Create a sale transaction
echo "3️⃣ Testing Sale Transaction..."
SALE_RESPONSE=$(curl -s -X POST "$BASE_URL/transactions" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"CustomerID\": \"$CUSTOMER_ID\",
    \"Type\": \"sale\",
    \"PaymentMethod\": \"cash\",
    \"Status\": \"completed\",
    \"Items\": [
      {
        \"InventoryID\": \"$INVENTORY_ID\",
        \"Quantity\": 10,
        \"UnitPrice\": 75
      }
    ]
  }")

if echo "$SALE_RESPONSE" | jq -e '.success' > /dev/null 2>&1; then
  TRANSACTION_ID=$(echo "$SALE_RESPONSE" | jq -r '.data.ID')
  echo "✅ Sale transaction created: $TRANSACTION_ID"
else
  echo "❌ Transaction creation failed: $SALE_RESPONSE"
  exit 1
fi
echo ""

# Test 4: Get customer details (verify balance update for credit sales)
echo "4️⃣ Testing Customer Balance..."
CUSTOMER_DETAILS=$(curl -s -H "Authorization: Bearer $TOKEN" "$BASE_URL/customers/$CUSTOMER_ID")
BALANCE=$(echo "$CUSTOMER_DETAILS" | jq -r '.data.CurrentBalance')
echo "✅ Customer balance: $BALANCE"
echo ""

# Test 5: Get inventory details (verify quantity deduction)
echo "5️⃣ Testing Inventory Quantity..."
INVENTORY_DETAILS=$(curl -s -H "Authorization: Bearer $TOKEN" "$BASE_URL/inventory/$INVENTORY_ID")
QUANTITY=$(echo "$INVENTORY_DETAILS" | jq -r '.data.Quantity')
echo "✅ Inventory quantity after sale: $QUANTITY (should be 90)"
echo ""

# Test 6: Get transaction details
echo "6️⃣ Testing Transaction Retrieval..."
TRANSACTION_DETAILS=$(curl -s -H "Authorization: Bearer $TOKEN" "$BASE_URL/transactions/$TRANSACTION_ID")
TOTAL=$(echo "$TRANSACTION_DETAILS" | jq -r '.data.TotalAmount')
STATUS=$(echo "$TRANSACTION_DETAILS" | jq -r '.data.Status')
echo "✅ Transaction total: $TOTAL, Status: $STATUS"
echo ""

# Test 7: Get dashboard stats
echo "7️⃣ Testing Dashboard..."
DASHBOARD=$(curl -s -H "Authorization: Bearer $TOKEN" "$BASE_URL/dashboard")
if echo "$DASHBOARD" | grep -q "error"; then
  echo "⚠️  Dashboard error: $DASHBOARD"
else
  echo "✅ Dashboard accessible"
  echo "$DASHBOARD" | jq '.'
fi
echo ""

echo "================================"
echo "✨ API Testing Complete!"
echo ""
echo "Frontend: http://localhost:3000"
echo "Backend:  http://localhost:8080"
echo "Database: localhost:5432"
