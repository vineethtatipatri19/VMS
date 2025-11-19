# VMS Architecture Documentation

## System Overview

VMS (Vendor Management System) is a full-stack application for managing perishable goods inventory, customers, and transactions. Built with Go backend and React frontend, deployed via Docker containers.

## Technology Stack

### Backend
- **Language:** Go 1.21
- **Framework:** Gorilla Mux (HTTP router)
- **Database:** PostgreSQL 15
- **Authentication:** JWT (JSON Web Tokens)
- **AI Integration:** Google Gemini API
- **Migrations:** Custom SQL-based migration system

### Frontend
- **Framework:** React 18
- **State Management:** Context API + Hooks
- **HTTP Client:** Axios
- **Charts:** Chart.js (react-chartjs-2)
- **Routing:** React Router v6
- **Build Tool:** Create React App
- **Web Server:** Nginx (production)

### Infrastructure
- **Containerization:** Docker + Docker Compose
- **Database Volume:** Persistent Docker volume
- **Networking:** Bridge network (docker-compose)
- **Ports:**
  - Frontend: 3000 (dev), 80 (container)
  - Backend: 8080
  - Database: 5432

## Architecture Pattern

### Backend: Clean Architecture (Repository Pattern)

```
backend/
├── main.go                      # Entry point, routing, middleware
├── startup.go                   # Application initialization
├── migrate.go                   # Database migrations
├── auth.go                      # Authentication handlers
├── internal/
│   ├── domain/                  # Business entities
│   │   ├── entities.go
│   │   ├── transaction.go
│   │   ├── customer.go
│   │   └── inventory.go
│   ├── repository/              # Data access layer
│   │   └── postgres/
│   │       ├── transaction.go   # Transaction repository
│   │       ├── customer.go
│   │       └── inventory.go
│   ├── service/                 # Business logic
│   │   ├── transaction.go
│   │   └── inventory.go
│   └── handlers/                # HTTP handlers
│       ├── transaction_handler.go
│       ├── customer_handler.go
│       └── inventory_handler.go
```

**Layers:**
1. **Domain Layer** - Pure business entities and interfaces
2. **Repository Layer** - Database interactions, SQL queries
3. **Service Layer** - Business logic, validation, orchestration
4. **Handler Layer** - HTTP request/response, routing

### Frontend: Component-Based Architecture

```
frontend/src/
├── App.js                       # Root component, routing
├── pages/                       # Page components
│   ├── Dashboard.js             # KPIs, charts, activity
│   ├── Inventory.js             # CRUD + filtering
│   ├── Customers.js             # Customer management
│   ├── Transactions.js          # Transaction ledger
│   ├── Crates.js                # Crate tracking
│   ├── Wastage.js               # Wastage logs
│   ├── ExpiryAlerts.js          # Expiry management
│   ├── Forecasting.js           # AI predictions
│   ├── Reports.js               # Report generation
│   ├── Login.js                 # Authentication
│   └── Register.js              # User registration
├── components/                  # Reusable components
│   ├── Layout.js                # Main layout + sidebar
│   ├── DeleteConfirmationModal.js
│   └── ui/                      # UI library
│       ├── Badge.js
│       ├── Button.js
│       ├── Card.js
│       ├── Input.js
│       ├── Modal.js
│       ├── Select.js
│       └── Toast.js
├── services/
│   └── api.js                   # Axios client, API methods
├── context/
│   └── AuthContext.js           # Authentication state
└── utils/
    └── dataHelpers.js           # Data transformation
```

## Data Flow

### Request Flow (Backend)

```
HTTP Request
    ↓
Middleware (CORS, Auth)
    ↓
Router (Gorilla Mux)
    ↓
Handler (HTTP layer)
    ↓
Service (Business logic)
    ↓
Repository (Data access)
    ↓
PostgreSQL Database
    ↓
← Response flows back up
```

### Request Flow (Frontend)

```
User Action
    ↓
React Component
    ↓
Event Handler
    ↓
API Service (axios)
    ↓
HTTP Request → Backend
    ↓
Response
    ↓
State Update (useState/Context)
    ↓
Component Re-render
    ↓
UI Update
```

## Database Schema

### Core Tables

**users**
- Authentication and user management
- Fields: id, email, password_hash, created_at

**customers** (26+ fields)
- Customer information, KYC, credit management
- Fields: id, name, email, phone, address, gstin, credit_limit, current_balance, etc.

**inventory_items** (35+ fields)
- Product inventory with expiry tracking
- Fields: id, name, variant, lot_number, quantity, unit, purchase_date, expiry_date, cost_price, selling_price, supplier_name, storage_location, etc.

**transactions** (22 fields)
- Sales and payment records
- Fields: id, customer_id, date, type, total_amount, payment_amount, invoice_number, sale_type, delivery_status, discount_amount, tax_amount, notes, etc.

**sale_items**
- Line items for each transaction
- Fields: id, transaction_id, inventory_lot_id, item_name, quantity, unit, price_per_unit, total, etc.
- **Linked to transactions** - Enables real product tracking in transactions

**crate_ledger**
- Track returnable crates per customer
- Fields: id, customer_id, date, type (issue/return), quantity, balance_after, notes

**wastage_log**
- Track wasted/damaged inventory
- Fields: id, inventory_item_id, date, quantity, reason, cost_impact, recorded_by, notes, photo_url

**expiry_alerts**
- Automatic alerts for expiring items
- Fields: id, inventory_item_id, alert_date, expiry_date, days_until_expiry, acknowledged, acknowledged_at

### Audit Trail

All entities support soft deletes:
- `deleted_at` - Timestamp of deletion
- `deleted_by` - User who deleted
- `deletion_reason` - Why it was deleted

All entities track updates:
- `created_at` - Creation timestamp
- `updated_at` - Last update timestamp
- `updated_by` - User who last updated

## Key Features Implementation

### 1. Dashboard Charts (Real Data)

**Sales Trend Chart:**
```javascript
// Frontend: Dashboard.js
const getSalesTrendData = () => {
  // Calculate last 7 days
  const last7Days = generateDateArray(7);
  
  // Aggregate sales by date from transactions
  transactions.forEach(tx => {
    if (tx.Type === 'sale') {
      // Sum TotalAmount per date
    }
  });
  
  return { labels, data };
};
```

**Top Products Chart:**
```javascript
// Frontend: Dashboard.js
const getTopProductsData = () => {
  const productSales = {};
  
  // Extract items from transaction.Details.items
  transactions.forEach(tx => {
    tx.Details.items.forEach(item => {
      // Aggregate quantity by item_name
      productSales[item.item_name] += item.quantity;
    });
  });
  
  // Sort and return top 5
};
```

**Backend Support:**
```go
// backend/internal/repository/postgres/transaction.go
func (r *transactionRepository) populateSaleItems(ctx, transactions) {
  // Query sale_items for all transaction IDs
  // Build map of items by transaction
  // Populate Details field with items array
}
```

### 2. Transaction with Sale Items

**Creating a Sale:**
```
1. Frontend submits transaction with items array
2. Backend creates transaction record
3. Backend creates sale_items records for each item
4. Backend updates inventory quantities (ACID transaction)
5. Backend returns transaction with items populated
```

**Querying Transactions:**
```
1. Repository fetches transactions from DB
2. Repository calls populateSaleItems()
3. Queries sale_items table for all transaction IDs
4. Groups items by transaction_id
5. Populates Details.items for each transaction
6. Returns complete transactions with items
```

### 3. Soft Delete with Audit

**Delete Flow:**
```
Frontend:
1. User clicks delete button
2. DeleteConfirmationModal opens
3. User must type "I CONFIRM DELETE"
4. User enters deletion reason
5. Sends DELETE request

Backend:
1. Validates deletion reason exists
2. Updates: deleted_at = NOW(), deleted_by = userID, deletion_reason = reason
3. Does NOT actually delete row
4. Returns success
```

**Recovery:**
```sql
-- Restore deleted record
UPDATE transactions 
SET deleted_at = NULL, deleted_by = NULL, deletion_reason = NULL
WHERE id = 'uuid';
```

### 4. JWT Authentication

**Registration:**
```go
// Hash password with bcrypt
passwordHash := bcrypt.GenerateFromPassword(password)
// Store user in database
// Generate JWT token with user ID
token := jwt.NewWithClaims(HS256, claims)
```

**Protected Endpoints:**
```go
// Middleware extracts and validates JWT
func authMiddleware(next http.Handler) http.Handler {
  // Extract Bearer token from header
  // Parse and validate JWT
  // Store user ID in request context
  // Call next handler
}
```

### 5. FEFO Inventory Sorting

First Expired First Out - Items sorted by expiry date:

```sql
SELECT * FROM inventory_items 
WHERE deleted_at IS NULL 
ORDER BY expiry_date ASC
```

Frontend applies color coding:
- **Green (Fresh):** > 7 days until expiry
- **Orange (Expiring Soon):** 1-7 days
- **Red (Expired):** <= 0 days

## API Response Format

### Success Response
```json
{
  "success": true,
  "data": { ... } or [ ... ]
}
```

### Error Response
```json
{
  "success": false,
  "error": "Error message"
}
```

**Frontend Access:**
```javascript
// Axios wraps response in .data
const response = await axios.get('/api/v1/customers');

// Backend returns: {success: true, data: [...]}
// Access: response.data.data (nested)
// Or fallback: response.data
const customers = response.data.data || response.data || [];
```

## Security

### Backend Security
- ✅ Password hashing with bcrypt (cost 10)
- ✅ JWT token expiration (24 hours)
- ✅ CORS middleware with proper origins
- ✅ SQL injection protection (parameterized queries)
- ✅ Input validation on all endpoints
- ✅ Error messages don't leak sensitive data

### Frontend Security
- ✅ JWT stored in localStorage
- ✅ Token sent in Authorization header
- ✅ Protected routes (redirect to login if no token)
- ✅ Logout clears token
- ✅ HTTPS enforced in production

### Database Security
- ✅ User credentials in environment variables
- ✅ No default/weak passwords in production
- ✅ Connection pooling with limits
- ✅ Audit trail for all deletions

## Performance Optimizations

### Backend
- Database connection pooling
- Index on frequently queried columns (customer_id, date, type)
- Batch queries for sale_items (single query for all transactions)
- Prepared statements for repeated queries

### Frontend
- Code splitting (React.lazy for routes)
- Memoization with useMemo/useCallback
- Debounced search inputs
- Pagination for large lists
- Chart data calculated once per transaction load

### Database
- Indexes on foreign keys
- Indexes on date fields for filtering
- GIN index on JSONB fields (if used)
- Materialized views for complex reports (future)

## Deployment

### Docker Compose (Local)
```yaml
services:
  db:        # PostgreSQL with persistent volume
  backend:   # Go app with auto migrations
  frontend:  # React app served by Nginx
```

### Google Cloud Run (Production)
- Backend: Containerized Go app
- Frontend: Static files in Cloud Storage + CDN
- Database: Cloud SQL PostgreSQL
- Secrets: Secret Manager
- CI/CD: Cloud Build

## Testing Strategy

### Backend Tests
- Unit tests for service layer
- Integration tests for repository layer
- E2E tests for full API flows
- Located in: `backend/tests/`

### Frontend Tests
- Component tests with React Testing Library
- Integration tests for page flows
- Manual testing for UI/UX

## Monitoring & Logging

### Backend
- Structured logging to stdout
- HTTP request/response logging
- Error tracking with stack traces
- Health check endpoint: `/api/v1/health`

### Database
- Query performance monitoring
- Connection pool metrics
- Slow query log

## Future Enhancements

### Planned Features
- [ ] Real-time notifications (WebSockets)
- [ ] Mobile app (React Native)
- [ ] Barcode scanning
- [ ] Multi-tenant support
- [ ] Advanced analytics dashboard
- [ ] Email notifications
- [ ] PDF invoice generation
- [ ] WhatsApp integration for alerts

### Scalability
- Horizontal scaling: Multiple backend instances behind load balancer
- Database: Read replicas for reporting queries
- Caching: Redis for frequently accessed data
- Queue: RabbitMQ for async processing (reports, notifications)

## Development Guidelines

### Backend
- Follow Go standard project layout
- Use interfaces for dependencies
- Write tests for business logic
- Keep handlers thin (delegate to services)
- Use context for request-scoped data

### Frontend
- Keep components small and focused
- Extract reusable logic to custom hooks
- Use TypeScript for type safety (future)
- Follow React best practices
- Maintain consistent naming conventions

### Database
- Always use migrations for schema changes
- Never modify data directly in production
- Test migrations in staging first
- Keep backups before major changes

## Troubleshooting

### Common Issues

**Backend won't start:**
- Check DATABASE_URL is correct
- Ensure PostgreSQL is running
- Verify migrations completed: `SELECT * FROM schema_migrations;`

**Frontend API errors:**
- Check REACT_APP_API_URL points to correct backend
- Verify backend is running: `curl http://localhost:8080/api/v1/health`
- Check browser console for CORS errors

**Database connection errors:**
- Verify credentials in docker-compose.yml
- Check if port 5432 is available
- Ensure database container is healthy: `docker ps`

**Charts not showing data:**
- Verify transactions have sale_items: `SELECT * FROM sale_items;`
- Check transaction repository populates Details
- Verify frontend processes Details.items correctly

## Contact & Support

- **GitHub:** https://github.com/vineethtatipatri19/VMS
- **Issues:** https://github.com/vineethtatipatri19/VMS/issues
- **Documentation:** See `/docs` folder for detailed guides
