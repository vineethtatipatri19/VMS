# Backend Architecture Documentation

## Overview

This backend is built using **Clean Architecture** principles, ensuring:
- **Separation of concerns**: Each layer has a single responsibility
- **Testability**: Business logic is independent of frameworks
- **Maintainability**: Code is organized and easy to understand
- **Extensibility**: New features can be added without major refactoring

## Architecture Layers

```
┌─────────────────────────────────────────────────────────────┐
│                        main.go                               │
│                 (Dependency Injection)                       │
└────────────────────┬────────────────────────────────────────┘
                     │
        ┌────────────┴────────────┐
        │                         │
┌───────▼───────┐        ┌────────▼──────┐
│   HTTP Layer  │        │  Config Layer │
│  (Handlers)   │        │               │
└───────┬───────┘        └───────────────┘
        │
┌───────▼────────────────────────────────────┐
│         Service Layer (Business Logic)      │
│  - CustomerService                          │
│  - InventoryService                         │
│  - TransactionService                       │
│  - WastageService                           │
│  - ExpiryAlertService                       │
└───────┬─────────────────────────────────────┘
        │
┌───────▼────────────────────────────────────┐
│      Repository Layer (Data Access)         │
│  - PostgreSQL implementations               │
│  - Interface-based design                   │
└───────┬─────────────────────────────────────┘
        │
┌───────▼─────────────────────────────────────┐
│         Domain Layer (Entities)              │
│  - Customer, Inventory, Transaction          │
│  - Business rules and validations            │
│  - No external dependencies                  │
└──────────────────────────────────────────────┘
```

## Directory Structure

```
backend/
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point
│
├── internal/
│   ├── domain/                     # Domain layer (entities, business rules)
│   │   ├── customer.go             # Customer entity & validation
│   │   ├── inventory.go            # Inventory entity & validation
│   │   ├── transaction.go          # Transaction entity & validation
│   │   ├── entities.go             # Other entities (Crate, Wastage, etc.)
│   │   └── errors.go               # Domain errors
│   │
│   ├── repository/                 # Data access layer
│   │   ├── interfaces.go           # Repository interfaces
│   │   ├── postgres/               # PostgreSQL implementations
│   │   │   ├── customer.go
│   │   │   ├── inventory.go
│   │   │   ├── transaction.go
│   │   │   └── ... (other repos)
│   │   └── mock/                   # Mock implementations for testing
│   │
│   ├── service/                    # Business logic layer
│   │   ├── customer.go             # Customer service
│   │   ├── inventory.go            # Inventory service
│   │   ├── transaction.go          # Transaction service
│   │   ├── wastage.go              # Wastage service
│   │   ├── expiry_alert.go         # Expiry alert service
│   │   └── forecast.go             # Forecasting service
│   │
│   ├── handlers/                   # HTTP handlers
│   │   ├── customer.go             # Customer HTTP handlers
│   │   ├── inventory.go            # Inventory HTTP handlers
│   │   ├── transaction.go          # Transaction HTTP handlers
│   │   ├── auth.go                 # Authentication handlers
│   │   ├── dashboard.go            # Dashboard handlers
│   │   └── reports.go              # Reports handlers
│   │
│   ├── middleware/                 # HTTP middleware
│   │   ├── auth.go                 # JWT authentication
│   │   ├── cors.go                 # CORS handling
│   │   ├── logging.go              # Request logging
│   │   └── recovery.go             # Panic recovery
│   │
│   ├── router/                     # Route configuration
│   │   └── router.go               # Chi router setup
│   │
│   ├── config/                     # Configuration management
│   │   └── config.go               # Config loading & validation
│   │
│   └── httputil/                   # HTTP utilities
│       └── response.go             # Standard response helpers
│
├── pkg/                            # Public packages (reusable)
│   ├── jwt/                        # JWT utilities
│   ├── password/                   # Password hashing
│   └── validation/                 # Input validation
│
├── migrations/                     # Database migrations
│   ├── 001_init.sql
│   ├── 002_users.sql
│   ├── 003_add_indexes.sql
│   └── 004_enhance_entities.sql
│
├── go.mod
└── go.sum
```

## Layer Responsibilities

### 1. Domain Layer (`internal/domain/`)

**Purpose**: Contains business entities and core business rules

**Responsibilities**:
- Define entity structs (Customer, Inventory, Transaction, etc.)
- Implement business validation logic
- Define domain errors
- No dependencies on external packages

**Example**:
```go
// Customer entity with validation
func (c *Customer) Validate() error {
    if c.Name == "" {
        return ErrInvalidInput("name is required")
    }
    return nil
}

// Business rule
func (c *Customer) CanPurchase(amount float64) bool {
    return c.CurrentBalance + amount <= c.CreditLimit
}
```

**Rules**:
- ✅ Can depend on: Nothing (pure Go)
- ❌ Cannot depend on: Database, HTTP, external libraries

---

### 2. Repository Layer (`internal/repository/`)

**Purpose**: Handles all data persistence

**Responsibilities**:
- Define repository interfaces
- Implement database queries
- Handle data mapping (DB ↔ Domain)
- Manage transactions

**Example**:
```go
type CustomerRepository interface {
    Create(ctx context.Context, customer *domain.Customer) error
    GetByID(ctx context.Context, id string) (*domain.Customer, error)
    List(ctx context.Context, filters domain.CustomerFilters) ([]*domain.Customer, error)
    Update(ctx context.Context, customer *domain.Customer) error
    SoftDelete(ctx context.Context, id, deletedBy, reason string) error
}
```

**Rules**:
- ✅ Can depend on: Domain layer
- ❌ Cannot depend on: Service layer, HTTP handlers

---

### 3. Service Layer (`internal/service/`)

**Purpose**: Contains business logic and orchestrates operations

**Responsibilities**:
- Implement use cases
- Coordinate between repositories
- Apply business rules
- Handle transactions
- Validate business constraints

**Example**:
```go
type CustomerService struct {
    customerRepo repository.CustomerRepository
    txRepo       repository.TransactionRepository
}

func (s *CustomerService) CreateCustomer(ctx context.Context, customer *domain.Customer) error {
    // Validate entity
    if err := customer.Validate(); err != nil {
        return err
    }
    
    // Check business rules
    existing, _ := s.customerRepo.GetByEmail(ctx, customer.Email)
    if existing != nil {
        return domain.ErrAlreadyExists
    }
    
    // Persist
    return s.customerRepo.Create(ctx, customer)
}
```

**Rules**:
- ✅ Can depend on: Domain layer, Repository interfaces
- ❌ Cannot depend on: HTTP handlers, concrete repository implementations

---

### 4. Handler Layer (`internal/handlers/`)

**Purpose**: HTTP request/response handling

**Responsibilities**:
- Parse HTTP requests
- Call service methods
- Format HTTP responses
- Handle HTTP errors

**Example**:
```go
type CustomerHandler struct {
    service *service.CustomerService
}

func (h *CustomerHandler) Create(w http.ResponseWriter, r *http.Request) {
    var req CreateCustomerRequest
    if err := httputil.DecodeJSON(r, &req); err != nil {
        httputil.SendError(w, err)
        return
    }
    
    customer := req.ToDomain()
    if err := h.service.CreateCustomer(r.Context(), customer); err != nil {
        httputil.SendError(w, err)
        return
    }
    
    httputil.SendJSON(w, http.StatusCreated, customer)
}
```

**Rules**:
- ✅ Can depend on: Service layer, Domain layer
- ❌ Cannot depend on: Repository layer directly

---

## Design Principles Applied

### 1. **Dependency Inversion Principle (DIP)**
- High-level modules (services) don't depend on low-level modules (repositories)
- Both depend on abstractions (interfaces)
- Example: `CustomerService` depends on `CustomerRepository` interface, not concrete implementation

### 2. **Single Responsibility Principle (SRP)**
- Each struct/function has one reason to change
- Example: `CustomerHandler` only handles HTTP, `CustomerService` only has business logic

### 3. **Interface Segregation Principle (ISP)**
- Interfaces are small and focused
- Example: Separate repositories for each entity rather than one giant repository

### 4. **Open/Closed Principle (OCP)**
- Open for extension, closed for modification
- Example: New repositories can be added without changing existing code

### 5. **Dependency Injection**
- Dependencies are injected via constructors
- Makes testing easy (can inject mocks)
- Example:
```go
func NewCustomerService(repo repository.CustomerRepository) *CustomerService {
    return &CustomerService{customerRepo: repo}
}
```

---

## Error Handling Strategy

### 1. Domain Errors
Defined in `internal/domain/errors.go`:
- `ErrNotFound` - Resource not found
- `ErrAlreadyExists` - Resource already exists
- `ErrInvalidInput` - Validation error
- `ErrUnauthorized` - Authentication error
- `ErrInvalidAttestation` - Delete attestation failed

### 2. Error Propagation
```
Domain Layer → Service Layer → Handler Layer → HTTP Response
```

### 3. HTTP Error Mapping
`httputil.SendError()` automatically maps domain errors to HTTP status codes:
- `ErrNotFound` → 404
- `ErrAlreadyExists` → 409
- `ValidationError` → 400
- `ErrUnauthorized` → 401
- Generic errors → 500

---

## Testing Strategy

### Unit Tests
- **Domain layer**: Test business rules and validations
- **Service layer**: Test business logic with mocked repositories
- **Handler layer**: Test HTTP handling with mocked services

### Integration Tests
- Test full request flow with real database (Docker container)
- Test repository implementations against actual PostgreSQL

### Test Structure
```
backend/
├── internal/
│   ├── domain/
│   │   ├── customer_test.go
│   │   └── inventory_test.go
│   ├── service/
│   │   ├── customer_test.go
│   │   └── inventory_test.go
│   └── handlers/
│       ├── customer_test.go
│       └── inventory_test.go
└── integration_test.go
```

---

## Database Transaction Management

### 1. Repository Level
- Simple CRUD operations don't need explicit transactions
- Complex operations use `tx *sql.Tx` parameter

### 2. Service Level
- Services coordinate multi-repository operations
- Use database transactions for consistency

Example:
```go
func (s *TransactionService) CreateSale(ctx context.Context, tx *domain.Transaction, items []*domain.SaleItem) error {
    // Start transaction
    dbTx, err := s.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer dbTx.Rollback()
    
    // Create transaction
    if err := s.txRepo.CreateWithTx(ctx, dbTx, tx); err != nil {
        return err
    }
    
    // Create sale items
    if err := s.saleItemRepo.CreateBatchWithTx(ctx, dbTx, items); err != nil {
        return err
    }
    
    // Update inventory quantities
    for _, item := range items {
        if err := s.inventoryRepo.DecrementQuantityWithTx(ctx, dbTx, item.InventoryLotID, item.Quantity); err != nil {
            return err
        }
    }
    
    // Commit
    return dbTx.Commit()
}
```

---

## Configuration Management

Configuration is loaded from environment variables in `internal/config/config.go`:

```go
cfg, err := config.Load()
```

### Environment Variables:
- `DATABASE_URL` - PostgreSQL connection string
- `PORT` - HTTP server port
- `JWT_SECRET` - JWT signing secret
- `GEMINI_API_KEY` - Google Gemini API key
- `MIGRATE_ON_START` - Run migrations on startup

---

## Adding New Features

### Example: Adding "Supplier" Feature

1. **Domain Layer** (`internal/domain/supplier.go`):
```go
type Supplier struct {
    ID   string
    Name string
    // ... fields
}

func (s *Supplier) Validate() error {
    // validation logic
}
```

2. **Repository Interface** (`internal/repository/interfaces.go`):
```go
type SupplierRepository interface {
    Create(ctx context.Context, supplier *domain.Supplier) error
    GetByID(ctx context.Context, id string) (*domain.Supplier, error)
    // ... other methods
}
```

3. **Repository Implementation** (`internal/repository/postgres/supplier.go`):
```go
type PostgresSupplierRepository struct {
    db *sql.DB
}

func (r *PostgresSupplierRepository) Create(ctx context.Context, supplier *domain.Supplier) error {
    // SQL implementation
}
```

4. **Service** (`internal/service/supplier.go`):
```go
type SupplierService struct {
    repo repository.SupplierRepository
}

func (s *SupplierService) CreateSupplier(ctx context.Context, supplier *domain.Supplier) error {
    // business logic
    if err := supplier.Validate(); err != nil {
        return err
    }
    return s.repo.Create(ctx, supplier)
}
```

5. **Handler** (`internal/handlers/supplier.go`):
```go
type SupplierHandler struct {
    service *service.SupplierService
}

func (h *SupplierHandler) Create(w http.ResponseWriter, r *http.Request) {
    // HTTP handling
}
```

6. **Wire in main.go**:
```go
supplierRepo := postgres.NewSupplierRepository(db)
supplierService := service.NewSupplierService(supplierRepo)
supplierHandler := handlers.NewSupplierHandler(supplierService)

// Register routes
r.Post("/api/v1/suppliers", supplierHandler.Create)
```

---

## Benefits of This Architecture

### 1. **Testability**
- Business logic can be tested without HTTP or database
- Mock repositories for unit testing
- Integration tests with real database

### 2. **Maintainability**
- Clear separation of concerns
- Easy to find and modify code
- Each file has a single purpose

### 3. **Extensibility**
- Add new features without modifying existing code
- Swap implementations (e.g., PostgreSQL → MongoDB)
- Easy to add new endpoints

### 4. **Team Collaboration**
- Different developers can work on different layers
- Clear interfaces between layers
- Less merge conflicts

### 5. **Production-Ready**
- Proper error handling
- Consistent API responses
- Transaction management
- Configuration management
- Middleware support

---

## Key Takeaways for Team

1. **Always validate in domain layer first**
2. **Use repository interfaces, not concrete implementations**
3. **Keep handlers thin - delegate to services**
4. **Use dependency injection**
5. **Return domain errors, map to HTTP at handler level**
6. **Write tests for each layer**
7. **Follow the existing patterns when adding features**

---

## Migration from Current Code

The refactoring will move:
- All structs → `internal/domain/`
- All database queries → `internal/repository/postgres/`
- Business logic → `internal/service/`
- HTTP handlers → `internal/handlers/`
- Middleware → `internal/middleware/`
- Configuration → `internal/config/`

This provides a clean, maintainable, and scalable codebase.
