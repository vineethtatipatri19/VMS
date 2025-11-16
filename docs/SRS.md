
# PGVMS - Software Requirements Specification (SRS)

## Overview

PGVMS (Perishable Goods Vendor Management System) is an enterprise-grade inventory and customer management system designed for wholesale/retail businesses dealing with perishable goods. The system includes comprehensive audit trail capabilities, soft delete functionality, and data compliance features.

---

## System Components

### Backend
- **Technology**: Go (Golang)
- **Database**: PostgreSQL 15+
- **API**: RESTful endpoints (see `docs/API.md`)
- **Authentication**: JWT-based with bcrypt password hashing
- **Deployment**: Docker containers on Google Cloud Run

### Frontend
- **Technology**: React (Web UI)
- **Styling**: Custom CSS with design system
- **State Management**: React Context API
- **UI Components**: Custom component library (Button, Card, Modal, Input, etc.)

### Mobile (Future)
- React Native mobile client (skeleton provided)
- Consumes same REST API as web frontend

### Database
- **PostgreSQL** with comprehensive schema
- **Migrations**: SQL-based migrations in `infra/migrations/`
- **Current Schema**: 4 migrations applied (001-004)

---

## Functional Requirements

### 1. Authentication & Authorization

**FR-1.1**: User Registration
- Users can register with email and password
- Passwords hashed using bcrypt
- Email validation required

**FR-1.2**: User Login
- JWT token-based authentication
- Tokens include user ID and email claims
- Token expiration handling

**FR-1.3**: Session Management
- Persistent sessions with localStorage
- Auto-logout on token expiration
- Secure token storage

---

### 2. Inventory Management

**FR-2.1**: Inventory CRUD Operations
- Create new inventory items with full metadata
- Read/list inventory with filtering and search
- Update inventory item details
- Soft delete inventory items with attestation

**FR-2.2**: Inventory Fields
- Basic: name, category, subcategory, variant
- Quantities: quantity, unit, weight_per_unit
- Pricing: cost_price, selling_price, margin_percentage (auto-calculated)
- Lot Management: lot_number, purchase_date, expiry_date
- Supplier: supplier_id, supplier_name, purchase_invoice
- Status: Automatically calculated (available, low_stock, out_of_stock, expired, damaged)
- Compliance: HSN code, GST rate
- Tracking: total_sold, total_wasted

**FR-2.3**: FEFO (First Expiry First Out) Sorting
- Inventory automatically sorted by expiry date
- Near-expiry items highlighted
- Expired items flagged

**FR-2.4**: Stock Status Indicators
- Visual badges for stock status
- Color-coded alerts (green=available, yellow=low stock, red=out of stock, black=expired)

---

### 3. Customer Management

**FR-3.1**: Customer CRUD Operations
- Create customers with comprehensive details
- View customer profiles and transaction history
- Update customer information
- Soft delete customers with attestation

**FR-3.2**: Customer Fields
- Basic: name, email, business_name
- Contact: contact_number, alternate_contact, whatsapp_number, address
- KYC: photo_url, aadhaar_verified, kyc_document_type, kyc_document_number, gstin
- Classification: customer_type (b2b, b2c, retail, wholesale), status (active, inactive, blocked)
- Credit: credit_limit, current_balance, payment_terms_days, interest_rate
- Analytics: last_transaction_date, total_purchases, loyalty_points
- Customization: tags (array), notes

**FR-3.3**: Customer Balance Tracking
- Running balance calculation
- Current balance display
- Outstanding payments tracking

---

### 4. Transaction Management

**FR-4.1**: Sale Transactions
- Create multi-item sales transactions
- Link transactions to customers
- Auto-generate invoice numbers (format: INV-YYYYMMDD-####)
- Support for discounts and tax

**FR-4.2**: Payment Transactions
- Record payments against customer accounts
- Multiple payment methods (cash, UPI, card, bank transfer, cheque, credit)
- Payment reference tracking (UPI ref, cheque number, etc.)

**FR-4.3**: Transaction Fields
- Basic: customer_id, date, type (sale/payment), invoice_number
- Financial: total_amount, discount_amount, tax_amount, payment_amount, balance_after
- Payment: payment_method, payment_reference, due_date, is_overdue (auto-calculated)
- Delivery: delivery_status (pending, packed, dispatched, delivered, cancelled), delivery_date, delivery_address
- Classification: sale_type (regular, wholesale, credit, return)
- Details: JSONB field with line items, notes, receipt_sent flag

**FR-4.4**: Transaction Edit & Delete
- Edit transactions with audit trail
- Soft delete with attestation requirement
- Track who modified what and when

---

### 5. Crate/Asset Management

**FR-5.1**: Crate Ledger Tracking
- Issue crates to customers
- Record crate returns
- Maintain running balance per customer

**FR-5.2**: Crate Fields
- customer_id, date, transaction_id
- crates_issued, crates_returned, balance
- crate_type, crate_value
- notes, updated_at
- Audit fields: deleted_at, deleted_by, deletion_reason

---

### 6. Wastage Tracking (Module 8)

**FR-6.1**: Wastage Log Entry
- Log wasted inventory with quantity and reason
- Categorize by reason: expired, damaged, spoiled, pest, other
- Calculate cost impact
- Attach photo evidence

**FR-6.2**: Wastage Fields
- inventory_item_id, item_name
- quantity, unit
- reason, reason_details
- cost_value
- logged_by, logged_at
- photo_url
- Audit fields: deleted_at, deleted_by, deletion_reason

**FR-6.3**: Wastage Reports
- List all wastage entries
- Filter by date, reason, item
- Calculate total wastage cost
- View wastage trends

---

### 7. Expiry Alert System (Module 9)

**FR-7.1**: Automatic Alert Generation
- Generate alerts for items expiring within 3 days
- Daily background process (planned)

**FR-7.2**: Alert Management
- View all pending alerts
- Acknowledge alerts
- Track acknowledgment history

**FR-7.3**: Expiry Alert Fields
- inventory_item_id
- alert_date, expiry_date
- days_until_expiry (auto-calculated)
- acknowledged, acknowledged_at, acknowledged_by
- created_at
- Audit fields: deleted_at, deleted_by, deletion_reason

---

### 8. Audit Trail & Compliance System

**FR-8.1**: Soft Delete System
- **No permanent deletion** of any records
- All entities support soft delete
- Records marked as deleted remain in database
- Deleted records excluded from default queries

**FR-8.2**: Audit Fields (All Entities)
- `deleted_at` - Timestamp when soft-deleted (NULL = active)
- `deleted_by` - User ID who performed deletion
- `deletion_reason` - Required reason for deletion

**FR-8.3**: Attestation Requirement
- All delete operations require user attestation
- Must type exact phrase: "I CONFIRM DELETE"
- Attestation enforced at API level
- Failed attestation returns error

**FR-8.4**: Delete Request Payload
```json
{
  "reason": "Reason for deletion (required)",
  "attestation": "I CONFIRM DELETE"
}
```

**FR-8.5**: Audit Coverage
- Customers: Full audit trail
- Inventory Items: Full audit trail
- Transactions: Full audit trail
- Sale Items: Full audit trail
- Crate Ledger: Full audit trail
- Wastage Log: Full audit trail
- Expiry Alerts: Full audit trail

---

### 9. Dashboard & Reporting

**FR-9.1**: Dashboard KPIs
- Total customers (active only)
- Total inventory value
- Today's sales amount
- Outstanding receivables
- Low stock items count
- Items expiring soon

**FR-9.2**: Reports
- Sales reports with date filtering
- Inventory reports with status breakdown
- Customer financial statements
- Wastage cost analysis (planned)
- Profit margin reports (planned)

**FR-9.3**: Export Functionality
- Export reports to Excel (planned)
- Print-friendly layouts
- PDF generation (planned)

---

### 10. AI Forecasting (Module 7)

**FR-10.1**: Demand Prediction
- Google Gemini API integration (stub mode)
- Forecast future demand based on historical data
- Suggest reorder quantities

**FR-10.2**: Forecasting UI
- Input: item name
- Output: Predicted demand, suggested reorder quantity, confidence level
- Fallback to stub mode if API unavailable

---

## Non-Functional Requirements

### NFR-1: Performance
- API response time < 200ms for 95% of requests
- Support for 1000+ concurrent users
- Database query optimization with indexes
- Pagination for large datasets

### NFR-2: Security
- **Authentication**: JWT tokens with secure secret
- **Password Security**: bcrypt hashing with salt
- **HTTPS**: Enforced in production (Cloud Run default)
- **CORS**: Configured for trusted origins
- **SQL Injection Protection**: Parameterized queries
- **XSS Protection**: Input sanitization
- **Audit Trail**: Complete logging of all deletions
- **Attestation**: Required for all destructive operations

### NFR-3: Data Integrity
- **Soft Delete**: No data loss, full recovery capability
- **Foreign Key Constraints**: Referential integrity enforced
- **Transaction Atomicity**: ACID compliance
- **Backup Strategy**: Automated daily backups (planned)
- **Data Retention**: Indefinite retention of all records

### NFR-4: Scalability
- **Horizontal Scaling**: Cloud Run auto-scaling
- **Database**: Cloud SQL with read replicas (planned)
- **Caching**: Redis for session management (planned)
- **CDN**: Static asset delivery (planned)

### NFR-5: Availability
- **Uptime Target**: 99.9%
- **Error Handling**: Graceful degradation
- **Monitoring**: Cloud Monitoring integration
- **Alerting**: Uptime checks and notifications

### NFR-6: Compliance
- **Data Retention**: All deleted records retained indefinitely
- **Audit Trail**: Who, what, when, why for all deletions
- **GST Compliance**: HSN codes, tax rates ready
- **KYC Support**: Document type and number fields
- **Data Privacy**: Secure storage of customer data

### NFR-7: Usability
- **Responsive Design**: Mobile-first UI
- **Accessibility**: ARIA labels, keyboard navigation
- **Error Messages**: Clear, actionable feedback
- **Loading States**: Visual feedback for async operations
- **Toast Notifications**: Success/error alerts

### NFR-8: Maintainability
- **Code Quality**: Consistent formatting, linting
- **Documentation**: Comprehensive API docs, code comments
- **Testing**: Unit tests, integration tests
- **Migrations**: Versioned database schema changes
- **CI/CD**: Automated builds and deployments

---

## Technical Constraints

### TC-1: Technology Stack
- Backend: Go 1.21+
- Frontend: React 18+
- Database: PostgreSQL 15+
- Deployment: Docker + Cloud Run

### TC-2: External Dependencies
- Google Gemini API (optional)
- WhatsApp Business API (planned)
- SMS Gateway (planned)
- Payment Gateway (planned)

### TC-3: Browser Support
- Chrome 90+
- Firefox 88+
- Safari 14+
- Edge 90+

---

## Data Model

See `infra/migrations/` for complete schema:
- **001_init.sql**: Core tables (customers, inventory_items, transactions, crates)
- **002_users.sql**: Authentication tables
- **003_add_indexes.sql**: Performance indexes
- **004_enhance_entities.sql**: Enhanced fields, wastage, expiry alerts, audit fields

See `docs/ER.sql` for entity-relationship diagrams and `docs/ENHANCED_ENTITIES.md` for detailed field documentation.

---

## API Endpoints

See `docs/API.md` for complete API documentation including:
- Authentication endpoints
- CRUD operations for all entities
- Query parameters and filtering
- Request/response examples
- Error codes and handling
- Attestation requirements for delete operations

---

## Deployment Architecture

**Development**:
- Docker Compose with 3 services (backend, frontend, database)
- Local PostgreSQL on port 5432
- Backend on port 8080
- Frontend on port 3000

**Production**:
- Backend: Cloud Run (auto-scaling)
- Frontend: Cloud Run or Cloud Storage + CDN
- Database: Cloud SQL PostgreSQL
- Monitoring: Cloud Monitoring & Logging
- CI/CD: Cloud Build

See `docs/DEPLOYMENT.md` for detailed deployment instructions.

---

## Future Enhancements

See `docs/NEXT_STEPS.md` for detailed roadmap:
- Multi-tenant/multi-business support
- Credit management enforcement
- GST billing and tax compliance
- WhatsApp integration
- Supplier management
- Payment plans and installments
- Mobile-first POS interface
- AI/ML demand forecasting enhancement
- Customer segmentation and loyalty

---

**Version**: 2.0  
**Last Updated**: November 16, 2025  
**Status**: Production-ready with audit compliance