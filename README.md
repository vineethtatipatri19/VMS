
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
- ✅ Module 6: AI-powered demand forecasting with Google Gemini
- ✅ Module 7: Comprehensive reporting (sales, inventory, customer)
- ✅ Transactional integrity for financial operations
- ✅ PostgreSQL database with migrations
- ✅ Production-ready error handling and validation

### Frontend (React Web)
- ✅ Modern React web application with responsive design
- ✅ All screens: Dashboard, Inventory, Customers, Transactions, Crates, Forecasting, Reports
- ✅ JWT authentication flow
- ✅ API client with axios
- ✅ State management with Context API
- ✅ Form validation and error handling
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
- `DELETE /api/v1/inventory/{id}` - Delete inventory item

### Customers
- `GET /api/v1/customers` - List customers
- `POST /api/v1/customers` - Create customer
- `GET /api/v1/customers/{id}` - Get customer
- `PUT /api/v1/customers/{id}` - Update customer
- `DELETE /api/v1/customers/{id}` - Delete customer

### Transactions
- `GET /api/v1/transactions` - List transactions (supports ?customerId=&type=sale|payment)
- `POST /api/v1/transactions` - Create transaction (sale or payment)
- `GET /api/v1/transactions/{id}` - Get transaction

### Crates
- `GET /api/v1/crates` - List crate ledger entries (supports ?customerId=)
- `POST /api/v1/crates` - Create crate entry
- `GET /api/v1/crates/balance/{customerId}` - Get crate balance for customer

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
├── backend/              # Go backend
│   ├── main.go          # Main entry point with routing
│   ├── auth.go          # Authentication handlers
│   ├── customers.go     # Customer CRUD handlers
│   ├── inventory.go     # Inventory CRUD handlers
│   ├── transaction_service.go  # Transaction handlers
│   ├── crates.go        # Crate management handlers
│   ├── dashboard.go     # Dashboard KPI handlers
│   ├── forecasting.go   # AI forecasting with Gemini
│   ├── reports.go       # Report generation
│   ├── helpers.go       # Utility functions
│   ├── migrate.go       # Database migration logic
│   └── Dockerfile       # Backend container
├── frontend/            # React frontend
│   ├── src/
│   │   ├── pages/       # Page components
│   │   ├── components/  # Reusable components
│   │   ├── services/    # API client
│   │   ├── context/     # React context (auth)
│   │   └── App.js       # Main app with routing
│   ├── public/          # Static files
│   ├── Dockerfile       # Frontend container
│   └── nginx.conf       # Nginx configuration
├── infra/
│   ├── migrations/      # Database migrations
│   └── cloudbuild.yaml  # Google Cloud Build config
├── docker-compose.yml   # Local development setup
├── .env.example         # Environment variables template
└── README.md            # This file
```

## Features Walkthrough

### Dashboard
- View key metrics: customers, expiring items, sales, outstanding balances
- Quick action buttons for common tasks
- Recent activity feed

### Inventory Management
- FEFO (First Expired First Out) sorting
- Visual status badges (Fresh, Expiring Soon, Expired)
- Filter by status
- Add, edit, delete inventory items
- Track lot numbers, quantities, and expiry dates

### Customer Management
- Complete customer information
- KYC verification status
- Contact details and address
- CRUD operations

### Transaction Ledger (Digital Patti Book)
- Record sales and payments
- Multi-item sales with automatic inventory deduction
- Transaction filtering by type and customer
- Complete transaction history

### Crate Management
- Track returnable crates per customer
- Issue and return crates
- View balance for each customer
- Transaction history

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

## Security

- JWT-based authentication
- Password hashing with bcrypt
- Environment-based secrets (no hardcoded credentials)
- CORS enabled for frontend-backend communication
- Database transaction integrity
- Input validation and error handling

## Support

For issues and questions, please open an issue on GitHub.

## License

[Add your license here]
