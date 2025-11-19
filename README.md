
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

Get up and running in under 5 minutes with demo data!

### 🚀 One-Command Setup (Local & Codespaces)

```bash
# Works everywhere - automatically detects your environment
bash setup.sh
```

This single command:
- ✅ Detects if you're on local or Codespaces
- ✅ Configures correct URLs automatically
- ✅ Builds and starts all services
- ✅ Runs database migrations (7 migrations)
- ✅ Loads demo data (15 customers, 45 items, 4 transactions)
- ✅ Creates demo user (demo@vms.com / demo123)
- ✅ Verifies everything is working

**That's it!** You now have:
- ✅ 15 customers (B2B, wholesale, retail)
- ✅ 45 inventory items (vegetables & fruits)
- ✅ 4 transactions with sale items
- ✅ 5 crate ledger entries
- ✅ 7 wastage logs
- ✅ 12 expiry alerts
- ✅ Working dashboard with real charts

### 📱 Access Your Application

**Local Development:**
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080/api/v1

**GitHub Codespaces:**
- Automatic URL configuration
- Click the 'Ports' tab → Globe icon next to port 3000
- Or visit the URL shown in setup output

**Login:**
- Email: `demo@vms.com`
- Password: `demo123`

**📖 Need Help?** See [SETUP.md](SETUP.md) for detailed instructions and troubleshooting

## What's Inside

This system includes 8 complete modules with real-time data:

### 1. Dashboard
- Real-time KPIs (customers, expiring items, sales, balances)
- **Sales Trend Chart** - Last 7 days of actual sales data
- **Top Products Chart** - Best sellers by quantity from real transactions
- Quick action buttons and recent activity feed

### 2. Inventory Management (45 items)
- FEFO (First Expired First Out) sorting
- Visual status badges (Fresh, Expiring Soon, Expired)
- Complete CRUD with 35+ fields per item
- Track lot numbers, suppliers, pricing, margins, storage locations

### 3. Customer Management (15 customers)
- B2B, wholesale, and retail customers
- KYC verification, GSTIN tracking
- Credit limits and current balances
- 26+ fields per customer with full CRUD

### 4. Transaction Ledger (12 transactions)
- Digital patti book with multi-item sales
- **Real sale items** - Each transaction includes actual products sold
- Payment tracking, discounts, taxes
- Invoice generation and delivery status
- 21+ fields per transaction

### 5. Crate Management
- Track returnable crates per customer
- Issue/return with balance tracking
- Complete transaction history

### 6. Wastage Tracking (5 logs)
- Log damaged, expired, contaminated items
- Cost impact tracking
- Photo documentation support
- Categorized by reason

### 7. Expiry Alerts (7 alerts)
- Automatic alerts for expiring items
- Urgency levels (Critical, Urgent, Moderate)
- Days until expiry calculation
- Acknowledge and dismiss functionality

### 8. AI Forecasting & Reports
- Demand prediction using Google Gemini AI
- Sales, inventory, and customer reports
- Date range filtering
- Print-friendly layouts

## Architecture

### Backend (Go)
- Clean architecture with repository pattern
- PostgreSQL with automatic migrations
- JWT authentication
- RESTful API with proper error handling
- **Soft deletes with audit trail** - Nothing truly deleted
- **Sale items integration** - Transactions include actual products

### Frontend (React)
- Modern React 18 with hooks
- Custom UI component library
- Chart.js for data visualization
- **Real-time charts** - Dashboard shows actual sales trends and top products
- Responsive design with mobile support
- Context API for state management

### DevOps
- Docker Compose for local development
- Ready for Cloud Run deployment
- Automated database migrations
- Health checks and monitoring

## Environment Variables

See [SETUP.md](SETUP.md) for complete environment configuration details.

**Backend:**
- `DATABASE_URL` - PostgreSQL connection string
- `JWT_SECRET` - Secret for JWT token signing
- `GEMINI_API_KEY` - Optional, for AI forecasting
- `MIGRATE_ON_START` - Auto-run migrations (default: true)

**Frontend:**
- `REACT_APP_API_URL` - Backend API URL (default: http://localhost:8080/api/v1)

## API Documentation

Complete API reference: [docs/API.md](docs/API.md)

**Quick Reference:**
- **Authentication:** `POST /api/v1/register`, `POST /api/v1/login`
- **Dashboard:** `GET /api/v1/dashboard`, `GET /api/v1/dashboard/activity`
- **Inventory:** `GET /POST /PUT /DELETE /api/v1/inventory`
- **Customers:** `GET /POST /PUT /DELETE /api/v1/customers`
- **Transactions:** `GET /POST /PUT /DELETE /api/v1/transactions`
- **Crates:** `GET /POST /PUT /DELETE /api/v1/crates`
- **Wastage:** `GET /POST /PUT /DELETE /api/v1/wastage`
- **Expiry Alerts:** `GET /PUT /DELETE /api/v1/expiry-alerts`
- **Forecasting:** `POST /api/v1/forecast`
- **Reports:** `POST /api/v1/reports/generate`

All endpoints (except auth) require JWT: `Authorization: Bearer <token>`

## Documentation

- **[QUICK_REFERENCE.md](QUICK_REFERENCE.md)** - Essential commands and troubleshooting
- **[SETUP.md](SETUP.md)** - Complete setup guide with troubleshooting
- **[ARCHITECTURE.md](ARCHITECTURE.md)** - System architecture and design decisions
- **[docs/API.md](docs/API.md)** - Complete API reference
- **[docs/DEPLOYMENT.md](docs/DEPLOYMENT.md)** - Production deployment guide
- **[docs/ENHANCED_ENTITIES.md](docs/ENHANCED_ENTITIES.md)** - Entity field reference
- **[backend/TESTING.md](backend/TESTING.md)** - Testing guide

## Development

### Running Tests
```bash
# Backend tests
cd backend
go test ./... -v

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

### Code Style
- Backend: Follow Go standard formatting (`go fmt`)
- Frontend: ESLint + Prettier recommended
- Commit messages: Conventional commits preferred

## Deployment

For production deployment to Google Cloud Run, see [docs/DEPLOYMENT.md](docs/DEPLOYMENT.md)

**Quick deploy:**
```bash
# Build and push
gcloud builds submit --config=infra/cloudbuild.yaml

# Deploy backend
gcloud run deploy pgvms-backend \
  --image gcr.io/PROJECT_ID/pgvms-backend \
  --set-env-vars DATABASE_URL=...,JWT_SECRET=...

# Deploy frontend to Cloud Storage + CDN
gsutil -m rsync -r frontend/build gs://your-bucket
```

## Project Structure

```
VMS/
├── backend/                      # Go backend application
│   ├── main.go                  # Entry point with routing
│   ├── internal/                # Internal packages (clean architecture)
│   │   ├── domain/              # Business entities
│   │   ├── repository/          # Data access layer
│   │   │   └── postgres/        # PostgreSQL implementation
│   │   ├── service/             # Business logic
│   │   └── handlers/            # HTTP handlers
│   └── tests/                   # Backend tests
├── frontend/                     # React frontend application
│   ├── src/
│   │   ├── pages/               # Page components
│   │   ├── components/          # Reusable components
│   │   ├── services/api.js      # API client
│   │   └── context/             # React context
│   └── public/                  # Static files
├── infra/
│   ├── migrations/              # Database migrations
│   ├── local/                   # Demo/seed data
│   └── terraform/               # Infrastructure as Code
├── docs/                        # Documentation
├── SETUP.md                     # Setup guide
├── ARCHITECTURE.md              # Architecture docs
├── docker-compose.yml           # Local development
└── README.md                    # This file
```

For detailed architecture, see [ARCHITECTURE.md](ARCHITECTURE.md)

## Features Walkthrough

### Dashboard
- View key metrics: customers, expiring items, sales, outstanding balances
- **Real-time charts:** Sales trend (last 7 days) and top products (by quantity)
- Quick action buttons and recent activity feed

### Inventory Management (45 demo items)
- FEFO (First Expired First Out) sorting by expiry date
- Visual status badges (Fresh, Expiring Soon, Expired)
- Complete CRUD with 35+ fields including suppliers, pricing, storage

### Customer Management (15 demo customers)
- B2B, wholesale, and retail customer types
- KYC verification, GSTIN, credit limits
- 26+ fields with full CRUD operations

### Transaction Ledger (12 demo transactions)
- Digital patti book with multi-item sales
- **Real sale items:** Tracks actual products sold (Tomato, Onion, etc.)
- Payment tracking, discounts, taxes, delivery status

### Crate Management
- Track returnable crates per customer
- Issue/return with automatic balance calculation

### Wastage Tracking (5 demo logs)
- Log damaged, expired, contaminated items
- Cost impact tracking with photo support

### Expiry Alerts (7 demo alerts)
- Automatic alerts for items approaching expiry
- Urgency levels: Critical (< 3 days), Urgent (3-7 days), Moderate (> 7 days)

### AI Forecasting & Reports
- Demand prediction using Google Gemini AI
- Sales, inventory, and customer financial reports
- Print-friendly layouts

## Audit & Security Features

### Soft Delete with Audit Trail
- **Nothing is truly deleted** - All records preserved
- **Deletion reason mandatory** - Every delete requires explanation
- **Attestation required** - Users must type "I CONFIRM DELETE"
- **Track who, when, why** - Complete audit trail (deleted_at, deleted_by, deletion_reason)
- **Restorable** - Soft-deleted items can be restored by clearing deleted_at

### Security
- JWT authentication with password hashing (bcrypt)
- Protected routes requiring valid token
- Environment-based secrets (no hardcoded credentials)
- CORS enabled for frontend-backend communication
- SQL injection protection with parameterized queries
- Input validation and error handling on all endpoints

## Support

- **Issues:** https://github.com/vineethtatipatri19/VMS/issues
- **Quick Reference:** [QUICK_REFERENCE.md](QUICK_REFERENCE.md) - Essential commands
- **Setup Help:** [SETUP.md](SETUP.md) - Detailed troubleshooting

## Contributing

Contributions are welcome! Please:
1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

[MIT License](LICENSE) - See LICENSE file for details

## Acknowledgments

- Built with Go, React, PostgreSQL, and Docker
- Charts powered by Chart.js
- AI forecasting via Google Gemini API
- Icons from Lucide React

---

**Made with ❤️ for efficient perishable goods management**

Get started in 5 minutes: See [SETUP.md](SETUP.md) | Quick commands: [QUICK_REFERENCE.md](QUICK_REFERENCE.md)
