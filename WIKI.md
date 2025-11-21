# PGVMS Wiki - Complete Project Guide

> **Perishable Goods Vendor Management System**  
> A production-ready, enterprise-grade inventory and customer management system with AI-powered forecasting

---

## 📚 Table of Contents

1. [Overview](#overview)
2. [Quick Start](#quick-start)
3. [System Architecture](#system-architecture)
4. [Features Guide](#features-guide)
5. [API Reference](#api-reference)
6. [Database Schema](#database-schema)
7. [Development Guide](#development-guide)
8. [Deployment](#deployment)
9. [Troubleshooting](#troubleshooting)
10. [Contributing](#contributing)

---

## Overview

### What is PGVMS?

PGVMS is a comprehensive management system designed for businesses dealing with perishable goods. It provides complete inventory tracking, customer management, transaction processing, and AI-powered demand forecasting—all with enterprise-grade audit trails and security.

### Key Statistics

- **8 Complete Modules**: Dashboard, Inventory, Customers, Transactions, Crates, Wastage, Alerts, AI Forecasting
- **45+ API Endpoints**: Full REST API with JWT authentication
- **35+ Fields per Inventory Item**: Comprehensive tracking
- **26+ Fields per Customer**: Complete KYC and credit management
- **Zero Data Loss**: Soft deletes with complete audit trail
- **Production Ready**: Docker containerized, Cloud Run deployable

### Technology Stack

| Component | Technology | Version |
|-----------|-----------|---------|
| Backend | Go (Golang) | 1.21+ |
| Frontend | React | 18.x |
| Database | PostgreSQL | 15+ |
| Authentication | JWT | - |
| AI Engine | Google Gemini | 1.5 Pro |
| Charts | Chart.js | 4.x |
| Containerization | Docker | - |
| Orchestration | Docker Compose | - |

### Project Links

- **GitHub**: https://github.com/vineethtatipatri19/VMS
- **Issues**: https://github.com/vineethtatipatri19/VMS/issues
- **Demo**: Access via setup (demo@vms.com / demo123)

---

## Quick Start

### Prerequisites

**All Platforms:**
- Docker Desktop (Windows/Mac) or Docker Engine (Linux)
- Git
- 4GB RAM available
- Ports 3000, 5432, 8080 available

**Windows-Specific:**
- WSL 2 enabled (`wsl --install`)
- Docker Desktop with WSL 2 backend
- Git Bash (recommended) or PowerShell

**macOS:**
- Homebrew (optional, for easy installation)
- Docker Desktop

**Linux:**
- Docker Engine and Docker Compose plugin
- Build essentials

### Installation (One Command)

```bash
# Clone repository
git clone https://github.com/vineethtatipatri19/VMS.git
cd VMS

# Run unified setup script (works everywhere)
bash setup.sh
```

**What happens:**
1. ✅ Detects environment (local/Codespaces)
2. ✅ Configures correct URLs
3. ✅ Builds all Docker images
4. ✅ Starts all services (frontend, backend, database)
5. ✅ Runs 7 database migrations
6. ✅ Loads comprehensive demo data
7. ✅ Creates demo user account
8. ✅ Verifies everything works

### Access Your Application

**Local Development:**
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080/api/v1
- Database: localhost:5432

**GitHub Codespaces:**
- URLs shown in setup output
- Format: `https://<codespace>-3000.app.github.dev`

**Login Credentials:**
- Email: `demo@vms.com`
- Password: `demo123`

### Demo Data Included

After setup, you'll have:
- 15 customers (mix of B2B, wholesale, retail)
- 45 inventory items (vegetables, fruits with expiry dates)
- 12 transactions with actual sale items
- 5 crate ledger entries
- 7 wastage logs
- 12 expiry alerts at different urgency levels
- Working dashboard with real charts

### Platform-Specific Notes

**Windows:**
- Use Git Bash (recommended) or PowerShell
- Ensure Docker Desktop is running (check system tray)
- WSL 2 must be enabled
- See [SETUP.md Windows section](SETUP.md#windows-1011-complete-walkthrough)

**macOS:**
- Start Docker Desktop from Applications
- May need to allow Docker in Security & Privacy settings
- Works with both Intel and Apple Silicon

**Linux:**
- Add user to docker group: `sudo usermod -aG docker $USER`
- Logout and login after adding to docker group
- Use `docker compose` (without hyphen) on newer versions

---

## System Architecture

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         Client Browser                           │
│                    (React Web Application)                       │
└────────────────────────┬────────────────────────────────────────┘
                         │ HTTPS/REST API
                         │ JWT Authentication
┌────────────────────────▼────────────────────────────────────────┐
│                      API Gateway                                 │
│                   (Go Chi Router)                                │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  Middleware: Auth, CORS, Logging, Recovery               │  │
│  └──────────────────────────────────────────────────────────┘  │
└────────────────────────┬────────────────────────────────────────┘
                         │
        ┌────────────────┼────────────────┐
        │                │                │
┌───────▼───────┐ ┌─────▼─────┐ ┌───────▼────────┐
│   Handlers    │ │  Services │ │  Repositories  │
│   (HTTP)      │ │ (Business)│ │ (Data Access)  │
└───────────────┘ └───────────┘ └────────┬───────┘
                                          │
                                ┌─────────▼─────────┐
                                │   PostgreSQL      │
                                │   Database        │
                                └───────────────────┘
```

### Clean Architecture Layers

**Domain Layer** (`internal/domain/`)
- Business entities (Customer, Inventory, Transaction, etc.)
- Entity validation rules
- Domain errors
- **No external dependencies**

**Repository Layer** (`internal/repository/`)
- Data access interfaces
- PostgreSQL implementations
- Mock implementations for testing
- CRUD operations with soft delete support

**Service Layer** (`internal/service/`)
- Business logic implementation
- Orchestrates multiple repositories
- Transaction management
- Complex validations

**Handler Layer** (`internal/handlers/`)
- HTTP request/response handling
- Input validation
- Error mapping
- Response formatting

**Middleware** (`internal/middleware/`)
- JWT authentication
- CORS configuration
- Request logging
- Panic recovery

### Frontend Architecture

```
src/
├── pages/              # Full page components (Dashboard, Inventory, etc.)
├── components/         # Reusable UI components
│   ├── Button.js
│   ├── Card.js
│   ├── Modal.js
│   ├── Table.js
│   └── DeleteConfirmationModal.js
├── services/
│   └── api.js         # Centralized API client (axios)
├── context/
│   └── AuthContext.js # Authentication state management
└── utils/             # Helper functions
```

### Database Schema Overview

**Core Tables:**
- `customers` - Customer information with KYC and credit tracking
- `inventory_items` - Inventory with FEFO tracking and expiry dates
- `transactions` - Sales transactions with payment tracking
- `sale_items` - Line items for each transaction
- `users` - User accounts for authentication

**Supporting Tables:**
- `crate_ledger` - Returnable crate tracking
- `wastage_log` - Damaged/expired item logging
- `expiry_alerts` - Automatic expiry notifications
- `payment_schedules` - Installment payment tracking
- `pricing_tiers` - Volume-based pricing

**Audit Features:**
All tables include:
- `deleted_at` - Soft delete timestamp
- `deleted_by` - Who deleted the record
- `deletion_reason` - Why it was deleted
- Standard timestamps (`created_at`, `updated_at`)

---

## Features Guide

### 1. Dashboard Module

**Purpose:** Real-time business overview with key performance indicators

**Features:**
- **Quick Stats**: Total customers, expiring items (next 7 days), today's sales, outstanding balances
- **Sales Trend Chart**: Last 7 days of actual sales data (line chart)
- **Top Products Chart**: Best-selling products by quantity (bar chart)
- **Recent Activity Feed**: Latest transactions and updates
- **Quick Actions**: Create transaction, add inventory, register customer

**API Endpoints:**
```
GET /api/v1/dashboard          # Main dashboard data
GET /api/v1/dashboard/activity # Recent activity feed
```

**Demo Data:**
- 15 customers with varying balances
- 21 items expiring in next 7 days
- ₹189,480.50 total inventory value
- Real sales trend from 4 transactions

### 2. Inventory Management Module

**Purpose:** Complete perishable goods tracking with FEFO sorting

**Features:**
- **FEFO Sorting**: Automatically sorts by expiry date (oldest first)
- **Status Badges**: Fresh (green), Expiring Soon (yellow), Expired (red), Damaged (gray)
- **35+ Fields per Item**:
  - Basic: name, category, subcategory, variant
  - Quantities: quantity, unit, weight_per_unit, min_stock_level
  - Pricing: cost_price, selling_price, margin_percentage (auto-calculated)
  - Lot: lot_number, batch_number, purchase_date, expiry_date
  - Supplier: supplier_id, supplier_name, purchase_invoice
  - Storage: storage_location, warehouse_section, rack_number
  - Compliance: hsn_code, gst_rate
  - Quality: quality_grade, perishability_index
  - Analytics: total_sold, total_wasted, turnover_rate

**CRUD Operations:**
- **Create**: Add new inventory with full metadata
- **Read**: View with filters (category, status, expiry range)
- **Update**: Modify any field except lot_number
- **Delete**: Soft delete with attestation ("I CONFIRM DELETE") and reason

**API Endpoints:**
```
GET    /api/v1/inventory           # List all inventory
POST   /api/v1/inventory           # Create new item
GET    /api/v1/inventory/:id       # Get specific item
PUT    /api/v1/inventory/:id       # Update item
DELETE /api/v1/inventory/:id       # Soft delete item
```

**Demo Data:** 45 items including:
- Tomatoes (3 variants: Cherry, Roma, Beefsteak)
- Onions (Red, Yellow, Shallots)
- Potatoes (Russet, Red, Fingerling)
- Leafy Greens (Spinach, Lettuce, Kale)
- Fruits (Apples, Bananas, Mangoes)

### 3. Customer Management Module

**Purpose:** Comprehensive customer tracking with KYC and credit management

**Features:**
- **Customer Types**: B2B, B2C, Retail, Wholesale
- **KYC Verification**: Aadhaar, GSTIN, document upload
- **Credit Management**: Credit limits, current balance, payment terms
- **26+ Fields per Customer**:
  - Basic: name, email, business_name
  - Contact: contact_number, alternate_contact, whatsapp_number, address
  - KYC: photo_url, aadhaar_verified, kyc_document_type/number, gstin
  - Classification: customer_type, status (active/inactive/blocked)
  - Credit: credit_limit, current_balance, payment_terms_days, interest_rate
  - Analytics: last_transaction_date, total_purchases, loyalty_points
  - Custom: tags (array), notes

**CRUD Operations:**
- **Create**: Register new customer with optional KYC
- **Read**: View all customers with search/filter
- **Update**: Modify details, adjust credit limits
- **Delete**: Soft delete with attestation and reason

**API Endpoints:**
```
GET    /api/v1/customers           # List all customers
POST   /api/v1/customers           # Create new customer
GET    /api/v1/customers/:id       # Get specific customer
PUT    /api/v1/customers/:id       # Update customer
DELETE /api/v1/customers/:id       # Soft delete customer
```

**Demo Data:** 15 customers including:
- Raj Supermarket (B2B, ₹50,000 credit limit)
- Green Valley Grocers (Wholesale)
- Fresh Foods Ltd (B2B, large credit)
- Various retail customers

### 4. Transaction Ledger Module

**Purpose:** Digital patti book with multi-item sales tracking

**Features:**
- **Multi-Item Transactions**: Each transaction includes multiple sale items
- **Real Sale Items**: Links to actual inventory with quantities and prices
- **Payment Tracking**: Cash, credit, UPI, card with payment status
- **21+ Fields per Transaction**:
  - Core: customer_id, transaction_type (sale/purchase/return/adjustment)
  - Financial: total_amount, discount_amount, tax_amount, final_amount
  - Payment: payment_method, payment_status (pending/partial/paid/overdue)
  - Delivery: delivery_status, delivery_date, delivery_address
  - Documents: invoice_number, invoice_url
  - Analytics: profit_margin, notes

**Sale Items** (sub-table):
- inventory_lot_id (links to inventory)
- item_name, quantity, unit
- unit_price, total_price
- discount_applied, tax_applied

**CRUD Operations:**
- **Create**: New transaction with multiple items
- **Read**: View all transactions with filters (date, customer, status)
- **Update**: Modify details, update payment status
- **Delete**: Soft delete with attestation

**API Endpoints:**
```
GET    /api/v1/transactions        # List all transactions
POST   /api/v1/transactions        # Create new transaction
GET    /api/v1/transactions/:id    # Get specific transaction
PUT    /api/v1/transactions/:id    # Update transaction
DELETE /api/v1/transactions/:id    # Soft delete transaction
GET    /api/v1/sale-items          # Get sale items for transaction
```

**Demo Data:** 12 transactions including:
- Multi-item sales (Tomatoes + Onions + Potatoes)
- Mix of payment methods (Cash, Credit, UPI)
- Various payment statuses (Paid, Pending, Partial)

### 5. Crate Management Module

**Purpose:** Track returnable crates issued to customers

**Features:**
- **Crate Tracking**: Issue and return with automatic balance calculation
- **Per-Customer Balances**: Track how many crates each customer has
- **Transaction Types**: 
  - Issue: Customer takes crates
  - Return: Customer returns crates
  - Adjustment: Manual corrections
- **Fields**:
  - customer_id, transaction_type
  - quantity (positive for issue, negative for return)
  - balance_after (auto-calculated)
  - transaction_date, notes

**API Endpoints:**
```
GET    /api/v1/crates              # List all crate transactions
POST   /api/v1/crates              # Create crate transaction
GET    /api/v1/crates/:id          # Get specific transaction
PUT    /api/v1/crates/:id          # Update transaction
DELETE /api/v1/crates/:id          # Soft delete transaction
```

**Demo Data:** 5 crate ledger entries:
- Issues to various customers
- Returns tracked
- Running balance per customer

### 6. Wastage Tracking Module

**Purpose:** Log and track damaged, expired, or contaminated items

**Features:**
- **Wastage Categories**: 
  - Spoiled, Damaged, Expired
  - Contaminated, Over-stocked, Other
- **Cost Tracking**: Automatically calculates financial impact
- **Photo Documentation**: Upload proof of wastage
- **Fields**:
  - inventory_lot_id (links to inventory)
  - item_name, quantity, unit
  - wastage_reason, wastage_date
  - cost_impact (quantity × cost_price)
  - responsible_person, photo_url, notes

**API Endpoints:**
```
GET    /api/v1/wastage             # List all wastage logs
POST   /api/v1/wastage             # Create wastage log
GET    /api/v1/wastage/:id         # Get specific log
PUT    /api/v1/wastage/:id         # Update log
DELETE /api/v1/wastage/:id         # Soft delete log
```

**Demo Data:** 7 wastage logs:
- Spoiled tomatoes (5 kg, ₹200 loss)
- Damaged potatoes during transport
- Expired spinach
- Over-stocked bananas

### 7. Expiry Alerts Module

**Purpose:** Automatic alerts for items approaching expiry

**Features:**
- **Urgency Levels**:
  - Critical: < 3 days until expiry (red)
  - Urgent: 3-7 days until expiry (orange)
  - Moderate: > 7 days until expiry (yellow)
- **Auto-Generation**: Created automatically by backend
- **Alert Actions**: Acknowledge, dismiss, or take action
- **Fields**:
  - inventory_lot_id, item_name
  - expiry_date, days_until_expiry
  - alert_level (critical/urgent/moderate)
  - is_acknowledged, acknowledged_at

**API Endpoints:**
```
GET    /api/v1/expiry-alerts       # List all alerts
GET    /api/v1/expiry-alerts/:id   # Get specific alert
PUT    /api/v1/expiry-alerts/:id   # Acknowledge alert
DELETE /api/v1/expiry-alerts/:id   # Dismiss alert
```

**Demo Data:** 12 expiry alerts:
- 5 critical (< 3 days)
- 4 urgent (3-7 days)
- 3 moderate (> 7 days)

### 8. AI Forecasting & Reports Module

**Purpose:** Demand prediction and comprehensive reporting

**AI Forecasting Features:**
- **Gemini AI Integration**: Uses Google Gemini 1.5 Pro
- **Historical Analysis**: Analyzes past sales patterns
- **Seasonal Adjustments**: Considers seasonal demand variations
- **Demand Prediction**: Forecasts next 30 days demand
- **Reorder Recommendations**: Suggests optimal order quantities

**Reporting Features:**
- **Sales Reports**: Revenue, items sold, profit margins by date range
- **Inventory Reports**: Stock levels, expiry summary, wastage analysis
- **Customer Reports**: Outstanding balances, top customers, credit utilization
- **Print-Friendly**: Optimized layouts for printing

**API Endpoints:**
```
POST   /api/v1/forecast            # Generate AI forecast
POST   /api/v1/reports/generate    # Generate custom report
GET    /api/v1/reports/:type       # Get predefined reports
```

**Report Types:**
- Sales summary (daily/weekly/monthly)
- Inventory valuation
- Customer outstanding balances
- Wastage cost analysis
- Profit margin analysis

---

## API Reference

### Authentication

**Register New User**
```http
POST /api/v1/register
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "SecurePass123"
}

Response: 201 Created
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "user_id": "uuid",
    "email": "john@example.com"
  }
}
```

**Login**
```http
POST /api/v1/login
Content-Type: application/json

{
  "email": "demo@vms.com",
  "password": "demo123"
}

Response: 200 OK
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "uuid",
      "email": "demo@vms.com",
      "name": "Demo User"
    }
  }
}
```

### Inventory Endpoints

**List Inventory**
```http
GET /api/v1/inventory
Authorization: Bearer <token>

Query Parameters:
- category: filter by category
- status: filter by status (fresh/expiring/expired)
- search: search by name
- limit: items per page (default: 100)
- offset: pagination offset

Response: 200 OK
{
  "success": true,
  "count": 45,
  "data": [
    {
      "id": 1,
      "name": "Tomato",
      "category": "Vegetable",
      "variant": "Cherry",
      "quantity": 50,
      "unit": "kg",
      "expiry_date": "2025-11-25T00:00:00Z",
      "status": "fresh",
      ...
    }
  ]
}
```

**Create Inventory Item**
```http
POST /api/v1/inventory
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Tomato",
  "category": "Vegetable",
  "variant": "Cherry",
  "quantity": 50,
  "unit": "kg",
  "cost_price": 30.00,
  "selling_price": 45.00,
  "expiry_date": "2025-12-01",
  "supplier_name": "Fresh Farms",
  "lot_number": "LOT-TOM-001"
}

Response: 201 Created
{
  "success": true,
  "message": "Inventory item created successfully",
  "data": { ... }
}
```

**Update Inventory Item**
```http
PUT /api/v1/inventory/:id
Authorization: Bearer <token>
Content-Type: application/json

{
  "quantity": 45,
  "selling_price": 50.00
}

Response: 200 OK
```

**Delete Inventory Item**
```http
DELETE /api/v1/inventory/:id
Authorization: Bearer <token>
Content-Type: application/json

{
  "deleted_by": "demo@vms.com",
  "deletion_reason": "Expired and disposed"
}

Response: 200 OK
{
  "success": true,
  "message": "Inventory item soft deleted successfully"
}
```

### Customer Endpoints

**List Customers**
```http
GET /api/v1/customers
Authorization: Bearer <token>

Response: 200 OK
{
  "success": true,
  "count": 15,
  "data": [ ... ]
}
```

**Create Customer**
```http
POST /api/v1/customers
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Raj Supermarket",
  "email": "raj@example.com",
  "contact_number": "+919876543210",
  "customer_type": "b2b",
  "credit_limit": 50000.00,
  "gstin": "29ABCDE1234F1Z5"
}

Response: 201 Created
```

### Transaction Endpoints

**Create Transaction**
```http
POST /api/v1/transactions
Authorization: Bearer <token>
Content-Type: application/json

{
  "customer_id": 1,
  "transaction_type": "sale",
  "payment_method": "credit",
  "items": [
    {
      "inventory_lot_id": 1,
      "quantity": 10,
      "unit_price": 45.00
    },
    {
      "inventory_lot_id": 2,
      "quantity": 5,
      "unit_price": 60.00
    }
  ],
  "discount_amount": 50.00,
  "tax_amount": 85.50
}

Response: 201 Created
{
  "success": true,
  "message": "Transaction created successfully",
  "data": {
    "transaction_id": 123,
    "invoice_number": "INV-20251121-001",
    "final_amount": 785.50
  }
}
```

### Dashboard Endpoints

**Get Dashboard Data**
```http
GET /api/v1/dashboard
Authorization: Bearer <token>

Response: 200 OK
{
  "success": true,
  "data": {
    "stats": {
      "total_customers": 15,
      "expiring_items_count": 21,
      "todays_sales": 12500.00,
      "outstanding_balances": 45000.00,
      "total_inventory_value": 189480.50
    },
    "sales_trend": [
      { "date": "2025-11-15", "sales": 8500 },
      { "date": "2025-11-16", "sales": 9200 },
      ...
    ],
    "top_products": [
      { "product": "Tomato - Cherry", "quantity": 150 },
      { "product": "Onion - Red", "quantity": 120 },
      ...
    ]
  }
}
```

### Error Responses

All endpoints may return error responses:

```http
400 Bad Request
{
  "success": false,
  "error": "Validation error: quantity must be positive"
}

401 Unauthorized
{
  "success": false,
  "error": "Invalid or expired token"
}

404 Not Found
{
  "success": false,
  "error": "Resource not found"
}

500 Internal Server Error
{
  "success": false,
  "error": "Internal server error"
}
```

**Complete API documentation:** [docs/API.md](docs/API.md)

---

## Database Schema

### Core Tables

#### customers
```sql
CREATE TABLE customers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE,
    business_name VARCHAR(255),
    contact_number VARCHAR(20) NOT NULL,
    alternate_contact VARCHAR(20),
    whatsapp_number VARCHAR(20),
    address TEXT,
    photo_url TEXT,
    aadhaar_verified BOOLEAN DEFAULT false,
    kyc_document_type VARCHAR(50),
    kyc_document_number VARCHAR(100),
    gstin VARCHAR(15),
    customer_type VARCHAR(50),
    status VARCHAR(20) DEFAULT 'active',
    credit_limit NUMERIC(12, 2) DEFAULT 0.00,
    current_balance NUMERIC(12, 2) DEFAULT 0.00,
    payment_terms_days INTEGER DEFAULT 0,
    interest_rate NUMERIC(5, 2) DEFAULT 0.00,
    last_transaction_date TIMESTAMPTZ,
    total_purchases NUMERIC(15, 2) DEFAULT 0.00,
    loyalty_points INTEGER DEFAULT 0,
    tags TEXT[],
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    deleted_by TEXT,
    deletion_reason TEXT
);
```

#### inventory_items
```sql
CREATE TABLE inventory_items (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    category VARCHAR(100) NOT NULL,
    subcategory VARCHAR(100),
    variant VARCHAR(100),
    quantity NUMERIC(10, 2) NOT NULL,
    unit VARCHAR(20) NOT NULL,
    weight_per_unit NUMERIC(10, 3),
    min_stock_level NUMERIC(10, 2),
    max_stock_level NUMERIC(10, 2),
    reorder_point NUMERIC(10, 2),
    cost_price NUMERIC(10, 2) NOT NULL,
    selling_price NUMERIC(10, 2) NOT NULL,
    margin_percentage NUMERIC(5, 2),
    lot_number VARCHAR(100),
    batch_number VARCHAR(100),
    purchase_date DATE,
    expiry_date DATE,
    days_until_expiry INTEGER,
    supplier_id INTEGER,
    supplier_name VARCHAR(255),
    purchase_invoice VARCHAR(100),
    storage_location VARCHAR(100),
    warehouse_section VARCHAR(50),
    rack_number VARCHAR(50),
    status VARCHAR(50),
    quality_grade VARCHAR(20),
    origin_country VARCHAR(100),
    is_organic BOOLEAN DEFAULT false,
    perishability_index NUMERIC(3, 2),
    hsn_code VARCHAR(20),
    gst_rate NUMERIC(5, 2),
    barcode VARCHAR(100),
    sku VARCHAR(100),
    image_url TEXT,
    total_sold NUMERIC(10, 2) DEFAULT 0.00,
    total_wasted NUMERIC(10, 2) DEFAULT 0.00,
    turnover_rate NUMERIC(5, 2),
    tags TEXT[],
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    deleted_by TEXT,
    deletion_reason TEXT
);
```

#### transactions
```sql
CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    customer_id INTEGER REFERENCES customers(id),
    transaction_type VARCHAR(50) NOT NULL,
    transaction_date TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    total_amount NUMERIC(12, 2) NOT NULL,
    discount_amount NUMERIC(12, 2) DEFAULT 0.00,
    tax_amount NUMERIC(12, 2) DEFAULT 0.00,
    final_amount NUMERIC(12, 2) NOT NULL,
    payment_method VARCHAR(50),
    payment_status VARCHAR(50) DEFAULT 'pending',
    payment_due_date DATE,
    invoice_number VARCHAR(100),
    invoice_url TEXT,
    delivery_status VARCHAR(50),
    delivery_date DATE,
    delivery_address TEXT,
    notes TEXT,
    profit_margin NUMERIC(5, 2),
    created_by INTEGER,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    deleted_by TEXT,
    deletion_reason TEXT
);
```

#### sale_items
```sql
CREATE TABLE sale_items (
    id SERIAL PRIMARY KEY,
    transaction_id INTEGER REFERENCES transactions(id),
    inventory_lot_id INTEGER REFERENCES inventory_items(id),
    item_name VARCHAR(255) NOT NULL,
    quantity NUMERIC(10, 2) NOT NULL,
    unit VARCHAR(20) NOT NULL,
    unit_price NUMERIC(10, 2) NOT NULL,
    total_price NUMERIC(12, 2) NOT NULL,
    discount_applied NUMERIC(10, 2) DEFAULT 0.00,
    tax_applied NUMERIC(10, 2) DEFAULT 0.00,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
```

### Supporting Tables

**crate_ledger** - Returnable crate tracking
**wastage_log** - Damaged/expired item logging
**expiry_alerts** - Automatic expiry notifications
**users** - User authentication
**payment_schedules** - Installment tracking
**pricing_tiers** - Volume-based pricing
**price_history** - Price change audit trail

### Indexes

Performance-optimized indexes on:
- Foreign keys (customer_id, inventory_lot_id, etc.)
- Search fields (name, email, category)
- Date fields (expiry_date, transaction_date)
- Status fields (status, payment_status)

**Complete schema:** [infra/migrations/](infra/migrations/)

---

## Development Guide

### Project Structure

```
VMS/
├── backend/                      # Go backend
│   ├── main.go                  # Entry point
│   ├── internal/
│   │   ├── domain/              # Entities & business rules
│   │   ├── repository/          # Data access layer
│   │   │   └── postgres/        # PostgreSQL implementations
│   │   ├── service/             # Business logic
│   │   ├── handlers/            # HTTP handlers
│   │   ├── middleware/          # Auth, CORS, logging
│   │   ├── router/              # Route configuration
│   │   ├── config/              # Configuration
│   │   └── pkg/                 # Shared utilities
│   ├── docs/                    # Backend documentation
│   └── tests/                   # Backend tests
│
├── frontend/                     # React frontend
│   ├── src/
│   │   ├── App.js               # Main application
│   │   ├── pages/               # Page components
│   │   ├── components/          # Reusable UI components
│   │   ├── services/            # API client
│   │   ├── context/             # State management
│   │   ├── styles/              # CSS files
│   │   └── utils/               # Helper functions
│   └── public/                  # Static assets
│
├── infra/                        # Infrastructure
│   ├── migrations/              # Database migrations (001-007)
│   ├── local/                   # Demo/seed data
│   └── terraform/               # IaC (optional)
│
├── docs/                         # Documentation
│   ├── API.md                   # API reference
│   ├── DEPLOYMENT.md            # Deployment guide
│   ├── SRS.md                   # Requirements spec
│   └── ...
│
├── .artifacts/                   # Historical docs
├── docker-compose.yml           # Local orchestration
├── Makefile                     # Build commands
├── setup.sh                     # One-command setup
└── README.md                    # Project overview
```

### Setting Up Development Environment

**1. Clone Repository**
```bash
git clone https://github.com/vineethtatipatri19/VMS.git
cd VMS
```

**2. Install Dependencies**

Backend:
```bash
cd backend
go mod download
go mod tidy
```

Frontend:
```bash
cd frontend
npm install
```

**3. Configure Environment**

Create `.env` in project root:
```bash
# Database
DATABASE_URL=postgres://pgvms_user:pgvms_password@localhost:5432/pgvms?sslmode=disable

# Backend
JWT_SECRET=your-dev-secret-key-change-in-production
GEMINI_API_KEY=your-gemini-api-key-optional
MIGRATE_ON_START=true
PORT=8080

# Frontend
REACT_APP_API_URL=http://localhost:8080/api/v1
```

**4. Start Services**

Using Docker (recommended):
```bash
docker-compose up -d --build
```

Or manually:
```bash
# Terminal 1: Database
docker-compose up -d db

# Terminal 2: Backend
cd backend
go run main.go

# Terminal 3: Frontend
cd frontend
npm start
```

### Running Tests

**Backend Unit Tests:**
```bash
cd backend
go test ./... -v

# With coverage
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

**Backend Integration Tests:**
```bash
cd backend/tests/integration
go test -v
```

**Specific Tests:**
```bash
# Test specific package
go test -v ./internal/domain/

# Test specific function
go test -v -run TestCustomerValidation

# Run with race detection
go test -race ./...
```

**Frontend Tests (if configured):**
```bash
cd frontend
npm test
```

### Code Style Guidelines

**Go (Backend):**
- Follow Go standard formatting: `go fmt`
- Use `golangci-lint` for linting
- Write descriptive variable names
- Add comments for exported functions
- Keep functions focused and small

**React (Frontend):**
- Use functional components with hooks
- Follow ESLint + Prettier configuration
- PropTypes for component props
- Descriptive component and variable names
- Extract reusable logic into custom hooks

**Git Commit Messages:**
```
feat: Add AI forecasting module
fix: Resolve expiry date calculation bug
docs: Update API documentation
refactor: Simplify transaction service
test: Add unit tests for customer validation
```

### Debugging

**Backend:**
```bash
# Enable debug logging
export LOG_LEVEL=debug

# Use Delve debugger
go install github.com/go-delve/delve/cmd/dlv@latest
dlv debug main.go

# Check logs
docker logs pgvms-backend -f --tail 100
```

**Frontend:**
```bash
# React DevTools (browser extension)
# Redux DevTools (if using Redux)
# Check browser console for errors

# Frontend logs
docker logs pgvms-frontend -f
```

**Database:**
```bash
# Connect to database
docker exec -it pgvms-postgres psql -U pgvms_user -d pgvms

# View queries
PGVMS=# SET log_statement = 'all';

# Check table contents
PGVMS=# SELECT COUNT(*) FROM customers;
PGVMS=# SELECT * FROM inventory_items LIMIT 5;
```

### Database Migrations

**Create New Migration:**
```bash
# Create file: infra/migrations/008_your_migration.sql
-- 008_your_migration.sql
-- Description of what this migration does

ALTER TABLE customers 
ADD COLUMN IF NOT EXISTS new_field VARCHAR(100);

CREATE INDEX IF NOT EXISTS idx_customers_new_field 
ON customers(new_field) 
WHERE deleted_at IS NULL;
```

**Apply Migration:**
```bash
# Using setup script
bash setup.sh

# Manually
docker exec -i pgvms-postgres psql -U pgvms_user -d pgvms < infra/migrations/008_your_migration.sql
```

**Verify Migration:**
```bash
docker exec pgvms-postgres psql -U pgvms_user -d pgvms -c "\d customers"
```

### Common Development Tasks

**Reset Database:**
```bash
docker-compose down -v
docker-compose up -d
bash setup.sh
```

**Rebuild Backend:**
```bash
docker-compose up -d --build backend
```

**Rebuild Frontend:**
```bash
docker-compose up -d --build frontend
```

**View Logs:**
```bash
docker-compose logs -f [service_name]
```

**Database Backup:**
```bash
docker exec pgvms-postgres pg_dump -U pgvms_user pgvms > backup.sql
```

**Restore Backup:**
```bash
docker exec -i pgvms-postgres psql -U pgvms_user -d pgvms < backup.sql
```

---

## Deployment

### Prerequisites

- Google Cloud Platform account
- `gcloud` CLI installed and configured
- Docker installed locally
- Domain name (optional)

### Quick Deployment to Cloud Run

**1. Build and Push Images:**
```bash
# Set project ID
export PROJECT_ID=your-gcp-project-id
gcloud config set project $PROJECT_ID

# Build using Cloud Build
gcloud builds submit --config=infra/cloudbuild.yaml
```

**2. Deploy Backend:**
```bash
gcloud run deploy pgvms-backend \
  --image gcr.io/$PROJECT_ID/pgvms-backend:latest \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars "DATABASE_URL=postgres://user:pass@/db?host=/cloudsql/project:region:instance,JWT_SECRET=prod-secret,MIGRATE_ON_START=true" \
  --add-cloudsql-instances project:region:instance
```

**3. Deploy Frontend:**
```bash
# Build with production API URL
cd frontend
REACT_APP_API_URL=https://pgvms-backend-xxx.run.app/api/v1 npm run build

# Deploy to Cloud Storage + CDN
gsutil -m rsync -r build/ gs://your-bucket-name/
```

**4. Setup Cloud SQL (PostgreSQL):**
```bash
gcloud sql instances create pgvms-db \
  --database-version=POSTGRES_15 \
  --tier=db-f1-micro \
  --region=us-central1

gcloud sql databases create pgvms --instance=pgvms-db
gcloud sql users create pgvms_user --instance=pgvms-db --password=secure-password
```

**Complete deployment guide:** [docs/DEPLOYMENT.md](docs/DEPLOYMENT.md)

### Environment Variables (Production)

```bash
# Backend
DATABASE_URL=postgres://user:pass@/db?host=/cloudsql/instance
JWT_SECRET=strong-random-secret-min-32-chars
GEMINI_API_KEY=your-production-gemini-key
MIGRATE_ON_START=true
PORT=8080

# Frontend (build-time)
REACT_APP_API_URL=https://your-backend-url.run.app/api/v1
```

### Health Checks

Backend health endpoint:
```bash
curl https://your-backend-url.run.app/api/v1/health
```

Expected response:
```json
{
  "status": "ok",
  "timestamp": "2025-11-21T10:30:00Z",
  "database": "connected"
}
```

### Monitoring

- **Cloud Run Metrics**: CPU, memory, requests/second
- **Cloud SQL Monitoring**: Connections, queries, disk usage
- **Application Logs**: Structured logging to Cloud Logging
- **Error Tracking**: Configure error reporting service

---

## Troubleshooting

### Common Issues

#### 1. "docker: command not found"

**Windows:**
```powershell
# Install Docker Desktop from:
https://docs.docker.com/desktop/install/windows-install/

# Ensure Docker Desktop is running
# Check system tray for whale icon
```

**macOS:**
```bash
# Install via Homebrew
brew install --cask docker

# Or download from Docker website
# Start Docker Desktop from Applications
```

**Linux:**
```bash
# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Add user to docker group
sudo usermod -aG docker $USER
# Logout and login
```

#### 2. "Port already in use"

**Check what's using the port:**
```bash
# Linux/macOS
lsof -i :3000
lsof -i :8080
lsof -i :5432

# Windows (PowerShell as Admin)
netstat -ano | findstr :3000
netstat -ano | findstr :8080
```

**Kill the process:**
```bash
# Linux/macOS
kill -9 <PID>

# Windows
taskkill /PID <PID> /F
```

**Or change ports in docker-compose.yml:**
```yaml
services:
  frontend:
    ports:
      - "3001:80"  # Changed from 3000
```

#### 3. "Cannot connect to Docker daemon"

**Windows:**
- Ensure Docker Desktop is running
- Check Settings → Resources → WSL Integration
- Restart Docker Desktop

**Linux:**
```bash
# Check Docker service
sudo systemctl status docker

# Start Docker service
sudo systemctl start docker

# Enable on boot
sudo systemctl enable docker
```

#### 4. Database connection errors

**Check database is running:**
```bash
docker-compose ps
docker logs pgvms-postgres
```

**Test connection:**
```bash
docker exec pgvms-postgres pg_isready -U pgvms_user -d pgvms
```

**Reset database:**
```bash
docker-compose down -v
docker-compose up -d db
bash setup.sh
```

#### 5. Frontend shows 401 errors

**Issue:** JWT token invalid or expired

**Solution:**
```bash
# Login again to get new token
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email":"demo@vms.com","password":"demo123"}'

# Or use the frontend login page
```

#### 6. Migrations not applied

**Symptoms:**
- "column deleted_at does not exist"
- Empty tables
- Backend errors on startup

**Solution:**
```bash
# Run migrations manually
for file in infra/migrations/*.sql; do 
  docker exec -i pgvms-postgres psql -U pgvms_user -d pgvms < "$file"
done

# Verify
docker exec pgvms-postgres psql -U pgvms_user -d pgvms -c "\dt"
```

#### 7. Frontend blank page

**Check browser console for errors**

**Common causes:**
- API URL misconfigured
- CORS issues
- Backend not running

**Solution:**
```bash
# Rebuild frontend with correct API URL
docker-compose build --build-arg REACT_APP_API_URL=http://localhost:8080/api/v1 frontend
docker-compose up -d frontend

# Check frontend logs
docker logs pgvms-frontend
```

#### 8. WSL 2 issues (Windows)

**Check WSL status:**
```powershell
wsl --status
wsl --list --verbose
```

**Update WSL:**
```powershell
wsl --update
wsl --set-default-version 2
```

**Restart WSL:**
```powershell
wsl --shutdown
# Then restart Docker Desktop
```

### Getting Help

**Check documentation:**
1. [SETUP.md](SETUP.md) - Complete setup guide
2. [QUICK_REFERENCE.md](QUICK_REFERENCE.md) - Common commands
3. [backend/docs/](backend/docs/) - Backend technical docs

**Enable debug logging:**
```bash
# Backend
export LOG_LEVEL=debug
docker-compose up backend

# Check all logs
docker-compose logs -f
```

**Report issues:**
- GitHub Issues: https://github.com/vineethtatipatri19/VMS/issues
- Include: OS, Docker version, error logs, steps to reproduce

---

## Contributing

### How to Contribute

1. **Fork the repository**
2. **Create a feature branch**
   ```bash
   git checkout -b feature/amazing-feature
   ```

3. **Make your changes**
   - Write clean, documented code
   - Follow existing code style
   - Add tests for new features

4. **Test your changes**
   ```bash
   # Backend tests
   cd backend && go test ./... -v
   
   # Integration tests
   cd backend/tests/integration && go test -v
   
   # Manual testing
   docker-compose up -d --build
   ```

5. **Commit with clear messages**
   ```bash
   git commit -m "feat: Add customer export to CSV"
   git commit -m "fix: Resolve expiry date timezone issue"
   git commit -m "docs: Update API documentation"
   ```

6. **Push to your fork**
   ```bash
   git push origin feature/amazing-feature
   ```

7. **Open a Pull Request**
   - Describe what your PR does
   - Link related issues
   - Include screenshots if UI changes

### Code Review Process

1. Automated checks run (tests, linting)
2. Maintainer reviews code
3. Address feedback if needed
4. PR merged after approval

### Areas to Contribute

**High Priority:**
- Mobile app (React Native)
- Additional report types
- Performance optimizations
- Test coverage improvements

**Medium Priority:**
- Additional AI features
- Enhanced analytics
- Export/import functionality
- Notification system

**Documentation:**
- Tutorial videos
- API examples
- Deployment guides
- Translation (i18n)

### Code of Conduct

- Be respectful and inclusive
- Provide constructive feedback
- Help others learn
- Maintain professional communication

---

## License

MIT License - See [LICENSE](LICENSE) file for details

---

## Acknowledgments

- Built with **Go**, **React**, **PostgreSQL**, and **Docker**
- Charts powered by **Chart.js**
- AI forecasting via **Google Gemini API**
- Icons from **Lucide React**
- Community contributors

---

## Quick Links

- 📖 [README](README.md) - Project overview
- 🚀 [Setup Guide](SETUP.md) - Detailed installation
- ⚡ [Quick Reference](QUICK_REFERENCE.md) - Common commands
- 🔌 [API Documentation](docs/API.md) - Complete API reference
- 🏗️ [Architecture](backend/docs/ARCHITECTURE.md) - Technical architecture
- 🚢 [Deployment](docs/DEPLOYMENT.md) - Production deployment
- 📊 [Requirements](docs/SRS.md) - Software requirements spec

---

**Made with ❤️ for efficient perishable goods management**

*Get started in 5 minutes: `bash setup.sh`*
