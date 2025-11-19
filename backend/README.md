# VMS Backend

Production-ready Go backend for the Perishable Goods Vendor Management System using Clean Architecture principles.

## 🏗️ Architecture

This backend follows **Clean Architecture** with clear separation of concerns:

```
backend/
├── internal/
│   ├── config/          # Configuration management
│   ├── domain/          # Business entities & rules
│   ├── handlers/        # HTTP request handlers
│   ├── middleware/      # HTTP middleware (auth, CORS, logging)
│   ├── repository/      # Data access layer
│   │   ├── interfaces.go
│   │   └── postgres/    # PostgreSQL implementation
│   ├── router/          # Route definitions
│   ├── service/         # Business logic
│   └── pkg/             # Shared utilities
│       ├── logger/      # Structured logging
│       ├── validator/   # Input validation
│       ├── errors/      # Error handling
│       └── httputil/    # HTTP response utilities
├── tests/
│   ├── e2e/            # End-to-end tests
│   ├── integration/    # Integration tests
│   └── testutil/       # Test utilities
├── docs/               # Documentation
├── main.go            # Application entry point
├── migrate.go         # Database migration utilities
└── startup.go         # Initialization logic
```

### Architectural Layers

1. **Domain Layer** (`internal/domain/`)
   - Core business entities
   - Business rules and validations
   - No dependencies on other layers

2. **Service Layer** (`internal/service/`)
   - Business logic implementation
   - Orchestrates domain entities
   - Depends only on domain and repository interfaces

3. **Repository Layer** (`internal/repository/`)
   - Data persistence abstraction
   - PostgreSQL implementation
   - Implements repository interfaces

4. **Handler Layer** (`internal/handlers/`)
   - HTTP request/response handling
   - Input validation
   - Calls service layer

## 🚀 Quick Start

### Prerequisites
- Go 1.21+
- PostgreSQL 15+
- Docker & Docker Compose (for containerized setup)

### Local Development

```bash
# Install dependencies
go mod download

# Run tests
go test ./...

# Build
go build -o pgvms .

# Run (requires DATABASE_URL)
export DATABASE_URL="postgres://user:pass@localhost:5432/pgvms?sslmode=disable"
export JWT_SECRET="your-secret-key"
./pgvms
```

### Docker Development

```bash
# From project root
docker-compose up --build backend
```

## 📦 Key Features

### Authentication & Authorization
- JWT-based authentication
- Password hashing with bcrypt
- Protected routes with middleware

### Entity Management
- **Customers**: Full CRUD with KYC tracking
- **Inventory**: FEFO sorting, expiry management
- **Transactions**: Multi-item sales with payment tracking
- **Crates**: Returnable asset management
- **Wastage**: Loss tracking with categorization
- **Expiry Alerts**: Automated notifications

### Data Management
- **Soft Deletes**: Audit trail with deletion reasons
- **Transactions**: ACID compliance for financial operations
- **Migrations**: Automated schema management

### Reporting & Analytics
- Sales reports
- Inventory reports
- Customer financial reports
- AI-powered demand forecasting (Google Gemini)

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific test suite
go test ./internal/service/...
go test ./tests/integration/...
go test ./tests/e2e/...

# Run integration tests (requires database)
./integration_test.sh
```

### Test Structure
- **Unit Tests**: `*_test.go` files alongside source
- **Integration Tests**: `tests/integration/` - database operations
- **E2E Tests**: `tests/e2e/` - full API workflows

See [docs/TESTING.md](docs/TESTING.md) for detailed testing guide.

## 🔧 Configuration

Environment variables:

```bash
# Database
DATABASE_URL=postgres://user:pass@host:5432/dbname?sslmode=disable
MIGRATE_ON_START=true  # Auto-run migrations on startup

# Server
PORT=8080

# Authentication
JWT_SECRET=your-secret-key-change-in-production

# AI Features (optional)
GEMINI_API_KEY=your-gemini-api-key

# Logging
DEBUG=true  # Enable debug logging
```

## 📚 API Documentation

### Base URL
```
http://localhost:8080/api/v1
```

### Authentication

**Register**
```http
POST /api/v1/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123",
  "name": "John Doe"
}
```

**Login**
```http
POST /api/v1/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}

Response: { "token": "jwt-token-here" }
```

### Protected Endpoints

All endpoints (except auth) require JWT token:
```http
Authorization: Bearer <token>
```

**Endpoints:**
- `GET /api/v1/customers` - List customers
- `POST /api/v1/customers` - Create customer
- `GET /api/v1/customers/:id` - Get customer
- `PUT /api/v1/customers/:id` - Update customer
- `DELETE /api/v1/customers/:id` - Soft delete customer

Similar patterns for:
- `/api/v1/inventory`
- `/api/v1/transactions`
- `/api/v1/crates`
- `/api/v1/wastage`
- `/api/v1/expiry-alerts`

See full API documentation: `/docs/API.md` (in project root)

## 🛠️ Development Guidelines

### Code Style
- Follow Go standard formatting: `go fmt`
- Run linter: `golangci-lint run`
- Write tests for new features
- Document exported functions

### Adding New Features

1. **Define Domain Entity** (`internal/domain/`)
   ```go
   type NewEntity struct {
       ID        string
       Name      string
       CreatedAt time.Time
   }
   ```

2. **Create Repository Interface** (`internal/repository/interfaces.go`)
   ```go
   type NewEntityRepository interface {
       Create(ctx context.Context, entity *domain.NewEntity) error
       GetByID(ctx context.Context, id string) (*domain.NewEntity, error)
   }
   ```

3. **Implement Repository** (`internal/repository/postgres/`)
   ```go
   func (r *NewEntityRepository) Create(ctx context.Context, entity *domain.NewEntity) error {
       // Implementation
   }
   ```

4. **Create Service** (`internal/service/`)
   ```go
   type NewEntityService struct {
       repo repository.NewEntityRepository
   }
   ```

5. **Create Handler** (`internal/handlers/`)
   ```go
   type NewEntityHandler struct {
       service *service.NewEntityService
   }
   ```

6. **Add Routes** (`internal/router/`)
   ```go
   r.Get("/new-entity", handler.List)
   r.Post("/new-entity", handler.Create)
   ```

7. **Write Tests** for each layer

### Database Migrations

Migrations are in `/infra/migrations/`:
```bash
# Applied automatically on startup if MIGRATE_ON_START=true
# Or manually:
for file in ../infra/migrations/*.sql; do
  psql $DATABASE_URL < "$file"
done
```

See [/infra/migrations/README.md](../infra/migrations/README.md) for details.

## 🔐 Security

- ✅ Password hashing with bcrypt
- ✅ JWT token authentication
- ✅ SQL injection protection (parameterized queries)
- ✅ CORS middleware
- ✅ Input validation
- ✅ Soft delete audit trail
- ⚠️ TODO: Rate limiting
- ⚠️ TODO: Request size limits
- ⚠️ TODO: HTTPS enforcement

## 📊 Monitoring & Logging

### Health Check
```http
GET /api/v1/health
Response: {"status": "ok"}
```

### Structured Logging
```go
import "github.com/example/pgvms/internal/pkg/logger"

log := logger.New("service")
log.Info("Operation completed", map[string]interface{}{
    "user_id": "123",
    "action": "create_customer",
})
```

### Error Handling
```go
import pkgerrors "github.com/example/pgvms/internal/pkg/errors"

err := pkgerrors.NotFound("Customer not found")
err := pkgerrors.Validation("Invalid email format")
err := pkgerrors.Database(dbErr)
```

## 🚢 Deployment

### Docker Production Build

```bash
# Build image
docker build -t pgvms-backend:latest .

# Run container
docker run -d \
  -p 8080:8080 \
  -e DATABASE_URL="postgres://..." \
  -e JWT_SECRET="..." \
  pgvms-backend:latest
```

### Environment Best Practices
- Use environment-specific configs
- Never commit secrets
- Use secrets management (Docker secrets, Kubernetes secrets)
- Enable HTTPS in production
- Set proper CORS origins
- Configure database connection pooling

## 📖 Additional Documentation

- [Architecture Details](docs/ARCHITECTURE.md)
- [Clean Architecture Principles](docs/CLEAN_ARCHITECTURE.md)
- [Testing Guide](docs/TESTING.md)
- [Test Structure](docs/TEST_STRUCTURE.md)

## 🤝 Contributing

1. Create feature branch
2. Write code following guidelines
3. Write/update tests
4. Run all tests
5. Submit PR with clear description

## 📝 License

[MIT License](../LICENSE)

## 🆘 Support

- Issues: https://github.com/vineethtatipatri19/VMS/issues
- Documentation: `/docs` folder

---

**Built with ❤️ using Go, PostgreSQL, and Clean Architecture**
