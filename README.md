
# PGVMS - Perishable Goods Vendor Management System

A complete, production-ready vendor management system for perishable goods with AI-powered demand forecasting.

## Features

### Backend (Go)
- ✅ Complete REST API with JWT authentication
- ✅ Module 1: Dashboard with KPIs (quick stats, sales metrics)
- ✅ Module 2: Inventory management with FEFO sorting and status badges
- ✅ Module 3: Customer management with KYC status
- ✅ Module 4: Transaction ledger (digital patti book) with filtering
- ✅ Module 5: Crate management for returnable assets
- ✅ Module 6: Wastage tracking and expiry alerts
- ✅ Module 7: AI-powered demand forecasting with Google Gemini
- ✅ Module 8: Comprehensive reporting (sales, inventory, customer)
- ✅ **Enterprise-grade audit trail** with soft deletes (nothing is truly deleted)
- ✅ **Attestation requirement** for all delete operations with deletion reason
- ✅ Complete edit/update functionality for all entities
- ✅ Transactional integrity for financial operations
- ✅ PostgreSQL database with migrations
- ✅ Production-ready error handling and validation

### Frontend (React Web)
- ✅ Modern React web application with responsive design
- ✅ All screens: Dashboard, Inventory, Customers, Transactions, Crates, Wastage, Expiry Alerts, Forecasting, Reports
- ✅ **DeleteConfirmationModal** with attestation requirement ("I CONFIRM DELETE")
- ✅ Complete CRUD operations with edit/delete buttons on all entity tables
- ✅ Deletion reason capture and audit trail warnings
- ✅ JWT authentication flow with protected routes
- ✅ Custom UI component library (Badge, Button, Card, Input, Select, Modal, Toast)
- ✅ API client with axios and error handling
- ✅ State management with Context API
- ✅ Form validation and comprehensive error handling
- ✅ Print functionality for reports
- ✅ Loading states and user feedback

### DevOps
- ✅ Docker Compose for local development
- ✅ Separate Dockerfiles for backend and frontend
- ✅ PostgreSQL database container
- ✅ Environment configuration files
- ✅ Database migrations and initialization
- ✅ Ready for Google Cloud Run deployment

## Quick Start

### Local Development with Docker Compose (Recommended)

1. **Clone the repository**
```bash
git clone <repository-url>
cd VMS
```

2. **Set up environment variables**
```bash
cp .env.example .env
# Edit .env and set your JWT_SECRET and optionally GEMINI_API_KEY
```

3. **Start all services**
```bash
docker-compose up --build
```

This will start:
- PostgreSQL database on port 5432
- Backend API on http://localhost:8080
- Frontend web app on http://localhost:3000

4. **Access the application**
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080/api/v1
- Health check: http://localhost:8080/api/v1/health

5. **Create your first user**
- Navigate to http://localhost:3000
- Click "Register here"
- Create an account
- Login and start using the system

### Manual Setup (Without Docker)

#### Backend

1. **Prerequisites**
   - Go 1.21 or higher
   - PostgreSQL 15 or higher

2. **Setup database**
```bash
createdb pgvms
```

3. **Set environment variables**
```bash
export DATABASE_URL="postgres://user:password@localhost:5432/pgvms?sslmode=disable"
export JWT_SECRET="your-secret-key"
export GEMINI_API_KEY="your-gemini-api-key"  # Optional
export MIGRATE_ON_START="true"
export PORT="8080"
```

4. **Run backend**
```bash
cd backend
go mod download
go run .
```

#### Frontend

1. **Prerequisites**
   - Node.js 18 or higher
   - npm or yarn

2. **Install dependencies**
```bash
cd frontend
npm install
```

3. **Set environment variables**
```bash
export REACT_APP_API_URL="http://localhost:8080/api/v1"
```

4. **Run frontend**
```bash
npm start
```

Frontend will be available at http://localhost:3000

## Environment Variables

### Backend
- `DATABASE_URL` - PostgreSQL connection string (required)
- `PORT` - Server port (default: 8080)
- `JWT_SECRET` - Secret key for JWT tokens (required for production)
- `GEMINI_API_KEY` - Google Gemini API key for AI forecasting (optional)
- `MIGRATE_ON_START` - Run migrations on startup (default: false)

### Frontend
- `REACT_APP_API_URL` - Backend API URL (default: http://localhost:8080/api/v1)

## API Documentation

### Authentication
- `POST /api/v1/register` - Register new user
- `POST /api/v1/login` - Login and get JWT token

### Dashboard
- `GET /api/v1/dashboard` - Get dashboard KPIs
- `GET /api/v1/dashboard/activity` - Get recent activity

### Inventory
- `GET /api/v1/inventory` - List inventory (supports ?status=expiring_soon|expired|fresh&sort=expiry)
- `POST /api/v1/inventory` - Create inventory item
- `GET /api/v1/inventory/{id}` - Get inventory item
- `PUT /api/v1/inventory/{id}` - Update inventory item
- `DELETE /api/v1/inventory/{id}` - Soft delete inventory item (requires reason & attestation)

### Customers
- `GET /api/v1/customers` - List customers
- `POST /api/v1/customers` - Create customer
- `GET /api/v1/customers/{id}` - Get customer
- `PUT /api/v1/customers/{id}` - Update customer
- `DELETE /api/v1/customers/{id}` - Soft delete customer (requires reason & attestation)

### Transactions
- `GET /api/v1/transactions` - List transactions (supports ?customerId=&type=sale|payment)
- `POST /api/v1/transactions` - Create transaction (sale or payment)
- `GET /api/v1/transactions/{id}` - Get transaction
- `PUT /api/v1/transactions/{id}` - Update transaction
- `DELETE /api/v1/transactions/{id}` - Soft delete transaction (requires reason & attestation)

### Crates
- `GET /api/v1/crates` - List crate ledger entries (supports ?customerId=)
- `POST /api/v1/crates` - Create crate entry
- `PUT /api/v1/crates/{id}` - Update crate entry
- `DELETE /api/v1/crates/{id}` - Soft delete crate entry (requires reason & attestation)
- `GET /api/v1/crates/balance/{customerId}` - Get crate balance for customer

### Wastage
- `GET /api/v1/wastage` - List wastage log entries (supports ?reason=expired|damaged|contaminated|other)
- `POST /api/v1/wastage` - Create wastage entry
- `PUT /api/v1/wastage/{id}` - Update wastage entry
- `DELETE /api/v1/wastage/{id}` - Soft delete wastage entry (requires reason & attestation)

### Expiry Alerts
- `GET /api/v1/expiry-alerts` - List expiry alerts (supports ?acknowledged=true|false)
- `PUT /api/v1/expiry-alerts/{id}/acknowledge` - Acknowledge an alert
- `DELETE /api/v1/expiry-alerts/{id}` - Soft delete expiry alert (requires reason & attestation)

### Forecasting (AI)
- `POST /api/v1/forecast` - Generate demand forecast using Google Gemini AI

### Reports
- `POST /api/v1/reports/generate` - Generate report (types: sales, inventory, customer)

## Database Migrations

Migrations are automatically run on startup when `MIGRATE_ON_START=true`.

Manual migration:
```bash
make migrate-up   # Apply migrations
make migrate-down # Rollback one migration
```

## Deployment

### Google Cloud Run

1. **Build and push to Container Registry**
```bash
gcloud builds submit --config=infra/cloudbuild.yaml
```

2. **Set up Cloud SQL PostgreSQL instance**
```bash
gcloud sql instances create pgvms-db --database-version=POSTGRES_15 --tier=db-f1-micro --region=us-central1
```

3. **Deploy**
```bash
gcloud run deploy pgvms-backend \
  --image gcr.io/[PROJECT_ID]/pgvms-backend \
  --platform managed \
  --region us-central1 \
  --set-env-vars DATABASE_URL=[CONNECTION_STRING],JWT_SECRET=[SECRET],GEMINI_API_KEY=[KEY]
```

## Development

### Running Tests
```bash
cd backend
go test ./...
```

### Building
```bash
# Backend
cd backend
go build -o pgvms .

# Frontend
cd frontend
npm run build
```

### Code Style
- Backend follows standard Go formatting (use `go fmt`)
- Frontend uses React best practices

## Project Structure

```
.
├── backend/                      # Go backend
│   ├── main.go                  # Main entry point with routing
│   ├── auth.go                  # Authentication handlers
│   ├── customers.go             # Customer CRUD handlers
│   ├── inventory.go             # Inventory CRUD handlers
│   ├── transaction_service.go   # Transaction list/create handlers
│   ├── transaction_update.go    # Transaction update handler
│   ├── crates.go                # Crate management handlers
│   ├── enhanced_entities.go     # Wastage and expiry alerts handlers
│   ├── delete_handlers.go       # Centralized soft delete handlers with attestation
│   ├── update_handlers.go       # Update handlers for customers, crates, wastage
│   ├── dashboard.go             # Dashboard KPI handlers
│   ├── forecasting.go           # AI forecasting with Gemini
│   ├── reports.go               # Report generation
│   ├── helpers.go               # Utility functions
│   ├── migrate.go               # Database migration logic
│   ├── startup.go               # Application initialization
│   └── Dockerfile               # Backend container
├── frontend/                     # React frontend
│   ├── src/
│   │   ├── pages/               # Page components
│   │   │   ├── Dashboard.js     # Dashboard with KPIs
│   │   │   ├── Inventory.js     # Inventory management with edit/delete
│   │   │   ├── Customers.js     # Customer management with edit/delete
│   │   │   ├── Transactions.js  # Transaction ledger with edit/delete
│   │   │   ├── Crates.js        # Crate management with edit/delete
│   │   │   ├── Wastage.js       # Wastage tracking with edit/delete
│   │   │   ├── ExpiryAlerts.js  # Expiry alerts with delete
│   │   │   ├── Forecasting.js   # AI forecasting
│   │   │   ├── Reports.js       # Report generation
│   │   │   ├── Login.js         # Login page
│   │   │   └── Register.js      # Registration page
│   │   ├── components/          # Reusable components
│   │   │   ├── Layout.js        # Main layout with sidebar
│   │   │   ├── DeleteConfirmationModal.js  # Delete confirmation with attestation
│   │   │   └── ui/              # UI component library
│   │   │       ├── Badge.js     # Badge component
│   │   │       ├── Button.js    # Button component
│   │   │       ├── Card.js      # Card component
│   │   │       ├── Input.js     # Input component
│   │   │       ├── Modal.js     # Modal component
│   │   │       ├── Select.js    # Select component
│   │   │       └── Toast.js     # Toast notification component
│   │   ├── services/            # API client
│   │   │   └── api.js           # Axios client with all API methods
│   │   ├── context/             # React context
│   │   │   └── AuthContext.js   # Authentication context
│   │   ├── styles/              # Global styles
│   │   │   └── variables.css    # CSS variables
│   │   └── App.js               # Main app with routing
│   ├── public/                  # Static files
│   ├── Dockerfile               # Frontend container
│   └── nginx.conf               # Nginx configuration
├── infra/
│   ├── migrations/              # Database migrations
│   │   ├── 001_init.sql        # Initial schema
│   │   ├── 002_users.sql       # Users table
│   │   ├── 003_add_indexes.sql # Performance indexes
│   │   └── 004_enhance_entities.sql  # Enhanced fields & audit columns
│   ├── local/                   # Local development data
│   │   ├── demo_simple.sql     # Demo data (45 inventory, 15 customers, etc.)
│   │   └── demo_seed.sql       # Alternative seed data
│   ├── terraform/               # Infrastructure as Code
│   └── cloudbuild.yaml          # Google Cloud Build config
├── docs/                        # Documentation
│   ├── API.md                   # API documentation
│   ├── DEPLOYMENT.md            # Deployment guide
│   ├── ENHANCED_ENTITIES.md     # Enhanced entity documentation
│   ├── ENTITY_FIELDS_SUMMARY.md # Field reference
│   └── ER.sql                   # Database schema reference
├── docker-compose.yml           # Local development setup
├── .env.example                 # Environment variables template
└── README.md                    # This file
```

## Features Walkthrough

### Dashboard
- View key metrics: customers, expiring items, sales, outstanding balances
- Quick action buttons for common tasks
- Recent activity feed
- Real-time inventory status overview

### Inventory Management
- FEFO (First Expired First Out) sorting
- Visual status badges (Fresh, Expiring Soon, Expired)
- Filter by status
- **Full CRUD operations**: Add, edit, update, soft delete with attestation
- Track 35+ fields including lot numbers, quantities, expiry dates, suppliers, pricing, storage location
- Comprehensive inventory details with cost price, selling price, margin calculation

### Customer Management
- Complete customer information with 26+ fields
- KYC verification status (Aadhaar, documents)
- Contact details (phone, WhatsApp, alternate)
- Business information (GSTIN, business name)
- Credit management (limit, current balance, payment terms)
- Tags and categorization (retail/wholesale)
- **Full CRUD operations** with edit and soft delete functionality

### Transaction Ledger (Digital Patti Book)
- Record sales and payments with 21+ fields
- Multi-item sales with automatic inventory deduction
- Payment tracking (method, reference, due dates)
- Discount and tax management
- Delivery status tracking
- Transaction filtering by type and customer
- **Full edit and soft delete** with audit trail
- Complete transaction history

### Crate Management
- Track returnable crates per customer
- Issue and return crates with balance tracking
- View balance for each customer
- **Edit and soft delete** crate entries
- Complete transaction history with notes

### Wastage Tracking
- Log damaged, expired, contaminated items
- Track cost impact of wastage
- Categorize by reason (expired, damaged, contaminated, spillage, other)
- Photo documentation support
- **Edit and soft delete** wastage entries
- Complete audit trail

### Expiry Alerts
- Automatic alerts for items approaching expiry
- Days until expiry calculation
- Urgency badges (Critical, Urgent, Moderate)
- Acknowledge alerts with timestamp
- **Soft delete** functionality with attestation
- Filter by acknowledged status

### AI Forecasting
- Demand prediction using Google Gemini AI
- Historical data analysis
- Configurable forecast period
- Confidence levels and summaries

### Reports
- Sales reports with top items
- Inventory status reports
- Customer financial statements
- Print-friendly layouts
- Date range filtering

### Audit & Compliance Features
- **Nothing is truly deleted** - All records preserved with soft delete markers
- **Complete audit trail** - Track who deleted, when, why, and attestation
- **Attestation requirement** - Users must type "I CONFIRM DELETE" exactly
- **Deletion reason mandatory** - Text explanation required for all deletions
- **Restorable records** - Soft-deleted items can be restored by clearing deleted_at
- **Updated tracking** - Track who updated records and when (updated_at, updated_by)

## Security

- JWT-based authentication with protected routes
- Password hashing with bcrypt
- Environment-based secrets (no hardcoded credentials)
- CORS enabled for frontend-backend communication
- Database transaction integrity
- Input validation and error handling
- **Audit trail compliance** - All deletions tracked with who, when, why
- **Soft delete protection** - No data loss, all records restorable
- **Attestation requirement** - Prevents accidental deletions
- SQL injection protection with parameterized queries

## Support

For issues and questions, please open an issue on GitHub.

## License

[Add your license here]
