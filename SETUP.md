# VMS - Complete Setup Guide

This guide will help you get the Vendor Management System running locally in under 5 minutes.

## Prerequisites

- Docker and Docker Compose installed
- Git installed
- 4GB RAM available
- Ports 3000, 5432, and 8080 available

## Quick Start (Recommended)

### 1. Clone and Setup

```bash
# Clone the repository
git clone https://github.com/vineethtatipatri19/VMS.git
cd VMS

# Create environment file (optional - defaults work for local dev)
cat > .env << EOF
JWT_SECRET=local-dev-secret-change-in-production
GEMINI_API_KEY=your-gemini-api-key-optional
MIGRATE_ON_START=true
EOF
```

### 2. Start All Services

```bash
# Start the entire stack
docker-compose up -d --build
```

This command will:
- Build the backend Go application
- Build the frontend React application
- Start PostgreSQL database
- Start all services in the background

**Wait 30-60 seconds** for all services to fully start.

### 3. Run Database Migrations

> **Important**: Automatic migrations are currently not working due to a golang-migrate path issue. Migrations must be run manually:

```bash
# Run all migrations (001 through 007)
for file in infra/migrations/*.sql; do 
  docker exec -i pgvms-postgres psql -U pgvms_user -d pgvms < "$file"
done
```

This will create all required tables including:
- Customers, Inventory, Transactions, Sale Items
- Crate Ledger, Wastage Log, Expiry Alerts
- Users, Payment Schedules, Pricing Tiers
- Soft delete columns (migration 007)

### 4. Load Demo Data

```bash
# Copy demo data file to container
docker cp infra/local/demo_simple.sql pgvms-postgres:/tmp/

# Load demo data (15 customers, 45 inventory items, transactions)
docker exec pgvms-postgres psql -U pgvms_user -d pgvms -f /tmp/demo_simple.sql
```

This creates:
- 15 customers (B2B, wholesale, retail)
- 45 inventory items (vegetables, fruits, dairy, grains)
- 12 transactions with actual sale items
- 2 crate entries
- 5 wastage logs
- 7 expiry alerts

### 5. Create Demo User

```bash
# Create a demo user account
bash setup-demo-user.sh
```

This will create a user with:
- **Email:** demo@vms.com
- **Password:** demo123

### 6. Access the Application

Open your browser and navigate to:

**Frontend:** http://localhost:3000

**Login with:**
- Email: `demo@vms.com`
- Password: `demo123`

**API Backend:** http://localhost:8080/api/v1
- Health check: http://localhost:8080/api/v1/health

### 6. Explore the System

You can now:
- ✅ View Dashboard with real KPIs and charts
- ✅ Browse 45 inventory items with expiry tracking
- ✅ Manage 15 customers
- ✅ View transaction history with sale items
- ✅ Track crate balances
- ✅ Monitor wastage logs
- ✅ Review expiry alerts
- ✅ Generate reports
- ✅ Test AI forecasting (if GEMINI_API_KEY set)

## Troubleshooting

### Services Not Starting

```bash
# Check service status
docker-compose ps

# View logs for a specific service
docker logs pgvms-backend
docker logs pgvms-frontend
docker logs pgvms-postgres

# Restart all services
docker-compose restart
```

### Database Connection Issues

```bash
# Check if PostgreSQL is ready
docker exec pgvms-postgres pg_isready -U pgvms_user -d pgvms

# Connect to database manually
docker exec -it pgvms-postgres psql -U pgvms_user -d pgvms
```

### Port Conflicts

If ports are already in use, edit `docker-compose.yml`:

```yaml
services:
  db:
    ports:
      - "5433:5432"  # Change host port
  backend:
    ports:
      - "8081:8080"  # Change host port
  frontend:
    ports:
      - "3001:80"    # Change host port
```

### Frontend Not Loading

```bash
# Rebuild frontend
docker-compose up -d --build frontend

# Check frontend logs
docker logs -f pgvms-frontend
```

### Backend API Errors

```bash
# Check backend logs
docker logs -f pgvms-backend

# Verify database migrations ran
docker exec pgvms-postgres psql -U pgvms_user -d pgvms -c "\dt"

# Should see these tables:
# customers, inventory_items, transactions, sale_items
# crate_ledger, wastage_log, expiry_alerts
# users, payment_schedules, pricing_tiers, price_history
```

### Migrations Not Running Automatically

**Known Issue:** The backend's automatic migration system (golang-migrate) currently has a path resolution issue.

**Symptoms:**
- Backend logs show: `Migration warning: first .: file does not exist`
- Tables are not created automatically
- API returns `column "deleted_at" does not exist` errors

**Solution:** Run migrations manually:

```bash
# Run all 7 migrations
for file in infra/migrations/*.sql; do 
  docker exec -i pgvms-postgres psql -U pgvms_user -d pgvms < "$file"
done

# Verify migrations succeeded
docker exec pgvms-postgres psql -U pgvms_user -d pgvms -c "\dt"
```

**Note:** Migration 007 (`007_add_soft_delete.sql`) is critical - it adds the `deleted_at`, `deleted_by`, and `deletion_reason` columns required by the repository layer.

### API Returns Empty Data

If the API returns `{"success": false, "count": 0}`:

1. **Check if demo data was loaded:**
   ```bash
   docker exec pgvms-postgres psql -U pgvms_user -d pgvms -c "SELECT COUNT(*) FROM customers;"
   ```

2. **Verify JWT token is valid:**
   ```bash
   TOKEN=$(cat .demo-token)
   curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/customers
   ```

3. **Check for soft delete column errors in logs:**
   ```bash
   docker logs pgvms-backend | grep "deleted_at"
   ```

If you see "column deleted_at does not exist", run migration 007 manually.

## Data Reset

To start fresh with clean data:

```bash
# Stop all services and remove volumes
docker-compose down -v

# Start services again
docker-compose up -d --build

# Wait 30 seconds for services to start
sleep 30

# Run migrations manually
for file in infra/migrations/*.sql; do 
  docker exec -i pgvms-postgres psql -U pgvms_user -d pgvms < "$file"
done

# Load demo data
docker cp infra/local/demo_simple.sql pgvms-postgres:/tmp/
docker exec pgvms-postgres psql -U pgvms_user -d pgvms -f /tmp/demo_simple.sql

# Create demo user
bash setup-demo-user.sh
```

## Manual Setup (Without Docker)

### Backend Setup

1. **Install Go 1.21+**
   ```bash
   # Verify installation
   go version
   ```

2. **Install PostgreSQL 15+**
   ```bash
   # Create database
   createdb pgvms
   ```

3. **Configure Environment**
   ```bash
   export DATABASE_URL="postgres://pgvms_user:pgvms_password@localhost:5432/pgvms?sslmode=disable"
   export JWT_SECRET="your-secret-key"
   export MIGRATE_ON_START="true"
   export PORT="8080"
   ```

4. **Run Backend**
   ```bash
   cd backend
   go mod download
   go run .
   ```

### Frontend Setup

1. **Install Node.js 18+**
   ```bash
   # Verify installation
   node --version
   npm --version
   ```

2. **Install Dependencies**
   ```bash
   cd frontend
   npm install
   ```

3. **Configure Environment**
   ```bash
   # Create .env file
   echo "REACT_APP_API_URL=http://localhost:8080/api/v1" > .env
   ```

4. **Run Frontend**
   ```bash
   npm start
   ```

### Load Demo Data

```bash
# Run SQL file
psql -U pgvms_user -d pgvms -f infra/local/demo_simple.sql

# Create demo user via API
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "demo@vms.com",
    "password": "demo123",
    "name": "Demo User"
  }'
```

## Development Workflow

### Running Tests

```bash
# Backend tests
cd backend
go test ./...

# Run specific test
go test -v -run TestCustomerCRUD

# Integration tests
cd backend/tests/integration
go test -v
```

### Building for Production

```bash
# Backend
cd backend
go build -o pgvms .

# Frontend
cd frontend
npm run build
```

### Database Migrations

Migrations run automatically when `MIGRATE_ON_START=true`. Manual control:

```bash
# Apply migrations
make migrate-up

# Rollback one migration
make migrate-down

# Check migration status
docker exec pgvms-postgres psql -U pgvms_user -d pgvms -c "SELECT * FROM schema_migrations;"
```

## Useful Commands

```bash
# View all containers
docker-compose ps

# Stop all services
docker-compose down

# Rebuild specific service
docker-compose up -d --build backend

# View real-time logs
docker-compose logs -f

# Access database
docker exec -it pgvms-postgres psql -U pgvms_user -d pgvms

# Check API health
curl http://localhost:8080/api/v1/health

# Test API endpoint (after login)
TOKEN="your-jwt-token"
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/customers
```

## Default Credentials

After running `setup-demo-user.sh`:

**Demo User:**
- Email: `demo@vms.com`
- Password: `demo123`

**Database:**
- Host: `localhost`
- Port: `5432`
- Database: `pgvms`
- User: `pgvms_user`
- Password: `pgvms_password`

## Next Steps

1. **Explore the Dashboard** - View KPIs, sales trends, and top products
2. **Browse Inventory** - See FEFO sorting and expiry alerts
3. **Create a Transaction** - Add a sale with multiple items
4. **Generate Reports** - Try sales, inventory, and customer reports
5. **Test AI Forecasting** - Set GEMINI_API_KEY and predict demand
6. **Customize** - Modify code and see live reload (frontend) or rebuild (backend)

## Production Deployment

See [DEPLOYMENT.md](docs/DEPLOYMENT.md) for:
- Google Cloud Run deployment
- Cloud SQL setup
- CI/CD with Cloud Build
- Environment configuration
- Security best practices

## Support

- **Issues:** https://github.com/vineethtatipatri19/VMS/issues
- **Documentation:** See `/docs` folder
- **API Reference:** [API.md](docs/API.md)

## System Requirements

**Minimum:**
- 2 CPU cores
- 4GB RAM
- 10GB disk space

**Recommended:**
- 4 CPU cores
- 8GB RAM
- 20GB disk space
- SSD storage

## Features Included

- ✅ Complete CRUD operations for all entities
- ✅ Real-time dashboard with charts
- ✅ JWT authentication
- ✅ Soft delete with audit trail
- ✅ Transaction management with sale items
- ✅ Inventory tracking with FEFO sorting
- ✅ Customer management with credit tracking
- ✅ Crate management for returnable assets
- ✅ Wastage logging with photos
- ✅ Expiry alerts with urgency levels
- ✅ AI-powered demand forecasting
- ✅ Comprehensive reporting
- ✅ Print functionality
- ✅ Responsive web design

Happy coding! 🚀
