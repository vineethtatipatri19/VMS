# New API Endpoints

## Wastage Tracking

### List Wastage Logs
```
GET /api/v1/wastage?reasonCode=expired
Authorization: Bearer <token>

Response: Array of WastageLog objects
[
  {
    "id": "uuid",
    "inventoryItemId": "uuid",
    "itemName": "Tomatoes",
    "lotNumber": "LOT-001",
    "quantity": 5.5,
    "unit": "kg",
    "reasonCode": "expired",
    "reasonDetails": "Past expiry date",
    "estimatedValue": 275.0,
    "reportedBy": "John Doe",
    "createdAt": "2025-11-16T15:00:00Z"
  }
]
```

### Create Wastage Entry
```
POST /api/v1/wastage
Authorization: Bearer <token>
Content-Type: application/json

{
  "inventoryItemId": "uuid",
  "itemName": "Tomatoes",
  "lotNumber": "LOT-001",
  "quantity": 5.5,
  "unit": "kg",
  "reasonCode": "expired",
  "reasonDetails": "Past expiry date",
  "estimatedValue": 275.0,
  "reportedBy": "John Doe"
}

Response: 201 Created with WastageLog object
```

**Reason Codes**: `expired`, `damaged`, `spoiled`, `recalled`, `other`

## Expiry Alerts

### List Expiry Alerts
```
GET /api/v1/expiry-alerts?status=pending
Authorization: Bearer <token>

Response: Array of ExpiryAlert objects
[
  {
    "id": "uuid",
    "inventoryItemId": "uuid",
    "itemName": "Milk",
    "lotNumber": "LOT-002",
    "expiryDate": "2025-11-19",
    "daysToExpiry": 3,
    "quantity": 20.0,
    "unit": "kg",
    "alertStatus": "pending",
    "createdAt": "2025-11-16T15:00:00Z"
  }
]
```

**Alert Status**: `pending`, `acknowledged`, `resolved`

### Update Expiry Alert Status
```
PUT /api/v1/expiry-alerts/{id}
Authorization: Bearer <token>
Content-Type: application/json

{
  "status": "acknowledged"
}

Response: 204 No Content
```

## Payment Schedules

### List Payment Schedules
```
GET /api/v1/payment-schedules?customerId=uuid&status=pending
Authorization: Bearer <token>

Response: Array of PaymentSchedule objects
[
  {
    "id": "uuid",
    "transactionId": "uuid",
    "customerId": "uuid",
    "installmentNum": 1,
    "dueDate": "2025-12-01",
    "amount": 5000.0,
    "paidAmount": 0.0,
    "status": "pending",
    "paymentDate": null,
    "paymentMethod": null,
    "paymentRef": null,
    "notes": "First installment",
    "createdAt": "2025-11-16T15:00:00Z",
    "updatedAt": "2025-11-16T15:00:00Z"
  }
]
```

**Payment Status**: `pending`, `partial`, `paid`, `overdue`

## Reports & Analytics

### Overdue Payments Report
```
GET /api/v1/reports/overdue-payments
Authorization: Bearer <token>

Response: Array of overdue payment records
[
  {
    "customerId": "uuid",
    "customerName": "ABC Store",
    "transactionId": "uuid",
    "invoiceNumber": "INV-20251115-0001",
    "totalAmount": 15000.0,
    "dueDate": "2025-11-10",
    "daysOverdue": 6,
    "balanceAfter": 15000.0
  }
]
```

### Wastage Summary Report
```
GET /api/v1/reports/wastage-summary
Authorization: Bearer <token>

Response: Array of wastage summary by item
[
  {
    "itemName": "Tomatoes",
    "totalQuantityWasted": 125.5,
    "totalValueLost": 6275.0,
    "mostCommonReason": "expired",
    "wasteCount": 15
  }
]
```

## Enhanced Existing Endpoints

### Inventory - Now Returns 35 Fields
```
GET /api/v1/inventory?status=available
Authorization: Bearer <token>

Response includes new fields:
- category, subCategory
- costPrice, sellingPrice, marginPercentage
- supplierId, supplierName, purchaseInvoice
- minStockLevel, reorderPoint, status
- barcode, sku, hsnCode, gstRate
- storageLocation, shelfLifeDays
- totalSold, totalWasted
```

### Transactions - Now Returns 21 Fields
```
GET /api/v1/transactions?customerId=uuid
Authorization: Bearer <token>

Response includes new fields:
- paymentMethod, paymentReference
- dueDate, isOverdue
- discountAmount, taxAmount
- invoiceNumber (auto-generated)
- receiptSent, balanceAfter
- saleType, deliveryStatus, deliveryDate, deliveryAddress
- notes
```

### Customers - Now Returns 26 Fields
```
GET /api/v1/customers
Authorization: Bearer <token>

Response includes new fields:
- email, gstin, customerType
- creditLimit, currentBalance
- paymentTermsDays, interestRate
- status, businessName
- kycDocumentType, kycDocumentNumber
- whatsappNumber, alternateContact
- lastTransactionDate, totalPurchases, loyaltyPoints
- notes, tags
```

### Crates - Now Returns 9 Fields
```
GET /api/v1/crates?customerId=uuid
Authorization: Bearer <token>

Response includes new fields:
- crateType (plastic/wooden/metal)
- crateValue (deposit amount)
- transactionId (linked transaction)
```

## Query Parameters

### Common Filters
- `customerId` - Filter by customer
- `status` - Filter by status
- `sort` - Sort results (expiry, name, quantity, date)
- `reasonCode` - Filter wastage by reason

## Response Codes
- `200 OK` - Success
- `201 Created` - Resource created
- `204 No Content` - Update successful
- `400 Bad Request` - Invalid input
- `401 Unauthorized` - Missing/invalid token
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error

## Authentication
All endpoints require JWT authentication:
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

Obtain token via:
```
POST /api/v1/login
Content-Type: application/json

{
  "username": "user@example.com",
  "password": "password"
}
```
