
# API Specification (REST)

Base URL: `/api/v1`

Authentication: Bearer token (JWT) required for all endpoints except `/register` and `/login`

**Headers:**
```
Authorization: Bearer <jwt_token>
Content-Type: application/json
```

## Authentication

### Register
**POST** `/api/v1/register`

**Request:**
```json
{
  "email": "user@example.com",
  "password": "secure_password"
}
```

**Response:**
```json
{
  "token": "eyJhbGc...",
  "user": {
    "id": "uuid",
    "email": "user@example.com"
  }
}
```

### Login
**POST** `/api/v1/login`

**Request:**
```json
{
  "email": "user@example.com",
  "password": "secure_password"
}
```

**Response:**
```json
{
  "token": "eyJhbGc..."
}
```

## Dashboard

### Get Dashboard KPIs
**GET** `/api/v1/dashboard`

**Response:**
```json
{
  "totalCustomers": 15,
  "totalInventoryItems": 45,
  "itemsExpiringSoon": 8,
  "totalSales": 125000,
  "outstandingBalance": 45000
}
```

### Get Recent Activity
**GET** `/api/v1/dashboard/activity`

**Response:**
```json
[
  {
    "type": "sale",
    "description": "Sale to Customer X",
    "amount": 5000,
    "timestamp": "2025-11-16T10:30:00Z"
  }
]
```

## Inventory

### List Inventory
**GET** `/api/v1/inventory`

**Query Parameters:**
- `status` - Filter by status: `fresh`, `expiring_soon`, `expired`
- `sort` - Sort by: `expiry`, `name`, `quantity`

**Example:**
```
GET /api/v1/inventory?sort=expiry&status=expiring_soon
```

**Response:**
```json
[
  {
    "id": "uuid",
    "name": "Tomatoes",
    "variant": "Cherry",
    "lotNumber": "LOT-2025-001",
    "quantity": 50,
    "unit": "kg",
    "purchaseDate": "2025-11-10",
    "expiryDate": "2025-11-20",
    "category": "Vegetables",
    "subCategory": "Fresh Produce",
    "costPrice": 20,
    "sellingPrice": 30,
    "marginPercentage": 50,
    "supplierId": "uuid",
    "supplierName": "Farm Fresh",
    "purchaseInvoice": "INV-001",
    "minStockLevel": 10,
    "reorderPoint": 20,
    "shelfLifeDays": 7,
    "storageLocation": "Cold Storage A1",
    "barcode": "1234567890",
    "sku": "TOM-CHER-001",
    "hsnCode": "07020000",
    "gstRate": 5,
    "status": "fresh",
    "weightPerUnit": 0.5,
    "packagingType": "Box",
    "imageUrl": "https://...",
    "notes": "Organic certified",
    "totalSold": 100,
    "totalWasted": 5,
    "lastRestockDate": "2025-11-10",
    "createdAt": "2025-11-10T08:00:00Z",
    "updatedAt": "2025-11-15T10:00:00Z"
  }
]
```

### Create Inventory Item
**POST** `/api/v1/inventory`

**Request:** (All fields from list response, required: name, quantity, unit, purchaseDate, expiryDate)

### Get Inventory Item
**GET** `/api/v1/inventory/{id}`

### Update Inventory Item
**PUT** `/api/v1/inventory/{id}`

**Request:** Same as create (partial updates supported)

### Delete Inventory Item (Soft Delete)
**DELETE** `/api/v1/inventory/{id}`

**Request:**
```json
{
  "reason": "Spoiled before expiry due to storage issue",
  "attestation": "I CONFIRM DELETE"
}
```

**Response:**
```json
{
  "message": "Inventory item deleted successfully"
}
```

**Note:** This performs a soft delete. The record is marked with `deleted_at`, `deleted_by`, and `deletion_reason` but not physically removed.

## Customers

### List Customers
**GET** `/api/v1/customers`

**Response:**
```json
[
  {
    "id": "uuid",
    "name": "ABC Store",
    "email": "abc@store.com",
    "address": "123 Main St",
    "contactNumber": "+919876543210",
    "alternateContact": "+919876543211",
    "whatsappNumber": "+919876543210",
    "photoUrl": "https://...",
    "businessName": "ABC Retail",
    "gstin": "29ABCDE1234F1Z5",
    "customerType": "wholesale",
    "aadhaarVerified": true,
    "kycDocumentType": "aadhaar",
    "kycDocumentNumber": "1234-5678-9012",
    "creditLimit": 100000,
    "currentBalance": 25000,
    "paymentTermsDays": 30,
    "interestRate": 2.5,
    "status": "active",
    "notes": "Reliable customer",
    "tags": ["vip", "wholesale"],
    "lastTransactionDate": "2025-11-15",
    "totalPurchases": 500000,
    "loyaltyPoints": 150,
    "createdAt": "2025-01-01T00:00:00Z",
    "updatedAt": "2025-11-15T10:00:00Z"
  }
]
```

### Create Customer
**POST** `/api/v1/customers`

**Request:** (All fields from list response, required: name, contactNumber)

### Get Customer
**GET** `/api/v1/customers/{id}`

### Update Customer
**PUT** `/api/v1/customers/{id}`

### Delete Customer (Soft Delete)
**DELETE** `/api/v1/customers/{id}`

**Request:**
```json
{
  "reason": "Customer moved to different region",
  "attestation": "I CONFIRM DELETE"
}
```

## Transactions

### List Transactions
**GET** `/api/v1/transactions`

**Query Parameters:**
- `customerId` - Filter by customer
- `type` - Filter by type: `sale`, `payment`

**Response:**
```json
[
  {
    "id": "uuid",
    "customerId": "uuid",
    "customerName": "ABC Store",
    "date": "2025-11-16",
    "type": "sale",
    "paymentAmount": 5000,
    "totalAmount": 5000,
    "details": {...},
    "paymentMethod": "cash",
    "paymentReference": "REF-001",
    "dueDate": "2025-12-16",
    "isOverdue": false,
    "discountAmount": 500,
    "taxAmount": 250,
    "notes": "Bulk order",
    "invoiceNumber": "INV-2025-001",
    "receiptSent": true,
    "balanceAfter": 20000,
    "saleType": "credit",
    "deliveryStatus": "delivered",
    "deliveryDate": "2025-11-16",
    "deliveryAddress": "123 Main St",
    "createdAt": "2025-11-16T10:00:00Z",
    "updatedAt": "2025-11-16T11:00:00Z",
    "updatedBy": "user-uuid"
  }
]
```

### Create Transaction
**POST** `/api/v1/transactions`

**Request:**
```json
{
  "customerId": "uuid",
  "type": "sale",
  "paymentMethod": "cash",
  "paymentAmount": 5000,
  "items": [
    {
      "inventoryLotId": "uuid",
      "itemName": "Tomatoes",
      "quantity": 10,
      "unit": "kg",
      "pricePerUnit": 30
    }
  ]
}
```

### Get Transaction
**GET** `/api/v1/transactions/{id}`

### Update Transaction
**PUT** `/api/v1/transactions/{id}`

**Request:** Partial update with fields to change

### Delete Transaction (Soft Delete)
**DELETE** `/api/v1/transactions/{id}`

**Request:**
```json
{
  "reason": "Duplicate entry created by mistake",
  "attestation": "I CONFIRM DELETE"
}
```

## Crates

### List Crate Entries
**GET** `/api/v1/crates`

**Query Parameters:**
- `customerId` - Filter by customer

**Response:**
```json
[
  {
    "id": "uuid",
    "customerId": "uuid",
    "date": "2025-11-16",
    "cratesIssued": 20,
    "cratesReturned": 10,
    "balance": 10,
    "notes": "Delivered with morning order",
    "crateType": "standard",
    "crateValue": 50,
    "transactionId": "uuid"
  }
]
```

### Create Crate Entry
**POST** `/api/v1/crates`

### Update Crate Entry
**PUT** `/api/v1/crates/{id}`

### Delete Crate Entry (Soft Delete)
**DELETE** `/api/v1/crates/{id}`

### Get Crate Balance
**GET** `/api/v1/crates/balance/{customerId}`

**Response:**
```json
{
  "customerId": "uuid",
  "balance": 10
}
```

## Wastage

### List Wastage Entries
**GET** `/api/v1/wastage`

**Query Parameters:**
- `reason` - Filter by reason: `expired`, `damaged`, `contaminated`, `spillage`, `other`

**Response:**
```json
[
  {
    "id": "uuid",
    "inventoryItemId": "uuid",
    "itemName": "Tomatoes",
    "quantity": 5,
    "unit": "kg",
    "reason": "expired",
    "reasonDetails": "Expired 2 days past date",
    "costValue": 150,
    "loggedBy": "manager@store.com",
    "loggedAt": "2025-11-16T08:00:00Z",
    "photoUrl": "https://..."
  }
]
```

### Create Wastage Entry
**POST** `/api/v1/wastage`

### Update Wastage Entry
**PUT** `/api/v1/wastage/{id}`

### Delete Wastage Entry (Soft Delete)
**DELETE** `/api/v1/wastage/{id}`

## Expiry Alerts

### List Expiry Alerts
**GET** `/api/v1/expiry-alerts`

**Query Parameters:**
- `acknowledged` - Filter by status: `true`, `false`

**Response:**
```json
[
  {
    "id": "uuid",
    "inventoryItemId": "uuid",
    "itemName": "Tomatoes",
    "alertDate": "2025-11-14",
    "expiryDate": "2025-11-18",
    "daysUntilExpiry": 2,
    "acknowledged": false,
    "acknowledgedAt": null,
    "acknowledgedBy": null,
    "createdAt": "2025-11-14T00:00:00Z"
  }
]
```

### Acknowledge Alert
**PUT** `/api/v1/expiry-alerts/{id}/acknowledge`

### Delete Expiry Alert (Soft Delete)
**DELETE** `/api/v1/expiry-alerts/{id}`

## Forecasting

### Generate Forecast
**POST** `/api/v1/forecast`

**Request:**
```json
{
  "itemName": "Tomatoes",
  "historicalData": [
    {"date": "2025-11-01", "quantity": 100},
    {"date": "2025-11-08", "quantity": 120}
  ],
  "forecastPeriod": 7
}
```

**Response:**
```json
{
  "forecast": "Based on historical data...",
  "confidence": "high",
  "predictions": [...]
}
```

## Reports

### Generate Report
**POST** `/api/v1/reports/generate`

**Request:**
```json
{
  "type": "sales",
  "startDate": "2025-11-01",
  "endDate": "2025-11-30",
  "customerId": "uuid"
}
```

**Response:**
```json
{
  "report": {...},
  "summary": {...}
}
```

## Error Responses

All endpoints return consistent error responses:

```json
{
  "error": "Error message description"
}
```

**Common HTTP Status Codes:**
- `200` - Success
- `201` - Created
- `400` - Bad Request (validation error)
- `401` - Unauthorized (missing/invalid token)
- `404` - Not Found
- `500` - Internal Server Error

## Soft Delete Pattern

All DELETE endpoints require attestation:

**Request Body:**
```json
{
  "reason": "Detailed explanation of why this record is being deleted",
  "attestation": "I CONFIRM DELETE"
}
```

**Audit Fields Added:**
- `deleted_at` - Timestamp when deleted
- `deleted_by` - User ID who deleted
- `deletion_reason` - Reason + attestation text

**Restoration:**
Records can be restored by setting `deleted_at = NULL` in the database.
