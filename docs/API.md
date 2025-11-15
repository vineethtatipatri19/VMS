
# API Specification (REST)

Base URL: /api/v1

Authentication: Bearer token (JWT) - not implemented in starter (see NFRs)

Resources:
- /inventory [GET, POST]
- /inventory/{id} [GET, PUT, DELETE]
- /customers [GET, POST]
- /customers/{id} [GET, PUT, DELETE]
- /transactions [GET, POST]
- /transactions/{id} [GET, PUT, DELETE]
- /crates [GET, POST]
- /crates/{id} [GET, PUT, DELETE]
- /forecast [POST]  (accepts conditions + returns stub prediction)

Query params: filtering, pagination (limit, offset)

Examples:
GET /api/v1/inventory?sort=expiry&status=expiring_soon
POST /api/v1/transactions
{
  "customerId":"...",
  "type":"sale",
  "date":"2025-11-15T08:00:00Z",
  "items":[{"inventoryLotId":"...","itemName":"Tomato","quantity":10,"unit":"kg","pricePerUnit":20}]
}
