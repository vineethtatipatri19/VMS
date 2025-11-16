# Backend Refactoring Implementation Plan

## Overview

This document provides a step-by-step plan to refactor the VMS backend from its current monolithic structure to a Clean Architecture pattern.

## Current State Assessment

**Current Issues:**
1. ❌ All code in root package (main)
2. ❌ Handlers directly interact with database
3. ❌ Business logic mixed with HTTP handling
4. ❌ No clear separation of concerns
5. ❌ Hard to test
6. ❌ Difficult to maintain and extend

**Current Files:**
- `main.go` - Server setup + routing
- `auth.go` - Authentication handlers + JWT logic
- `customers.go` - Customer CRUD + database queries
- `inventory.go` - Inventory CRUD + database queries
- `transactions.go`, `crates.go`, etc. - Similar pattern
- `helpers.go` - Mixed utilities

## Target State

**Clean Architecture with:**
1. ✅ Domain layer (entities, business rules)
2. ✅ Repository layer (data access interfaces + implementations)
3. ✅ Service layer (business logic)
4. ✅ Handler layer (HTTP handling)
5. ✅ Proper dependency injection
6. ✅ Testable code
7. ✅ Clear separation of concerns

## Phase 1: Foundation Setup (Week 1)

### Step 1.1: Create Directory Structure

```bash
mkdir -p internal/domain
mkdir -p internal/repository/postgres
mkdir -p internal/repository/mock
mkdir -p internal/service
mkdir -p internal/handlers
mkdir -p internal/middleware
mkdir -p internal/router
mkdir -p internal/config
mkdir -p internal/httputil
mkdir -p pkg/jwt
mkdir -p pkg/password
mkdir -p pkg/validation
```

### Step 1.2: Extract Domain Models ✅ DONE

Files created:
- `internal/domain/customer.go` - Customer entity
- `internal/domain/inventory.go` - Inventory entity
- `internal/domain/transaction.go` - Transaction entity
- `internal/domain/entities.go` - Other entities
- `internal/domain/errors.go` - Domain errors

### Step 1.3: Create Configuration Package ✅ DONE

File created:
- `internal/config/config.go` - Configuration management

### Step 1.4: Create HTTP Utilities ✅ DONE

File created:
- `internal/httputil/response.go` - Standard response helpers

### Step 1.5: Define Repository Interfaces ✅ DONE

File created:
- `internal/repository/interfaces.go` - All repository interfaces

## Phase 2: Repository Layer (Week 1-2)

### Step 2.1: Implement Customer Repository

Create `internal/repository/postgres/customer.go`:

```go
package postgres

import (
    "context"
    "database/sql"
    "github.com/yourusername/vms/backend/internal/domain"
    "github.com/yourusername/vms/backend/internal/repository"
    "github.com/lib/pq"
)

type customerRepository struct {
    db *sql.DB
}

func NewCustomerRepository(db *sql.DB) repository.CustomerRepository {
    return &customerRepository{db: db}
}

func (r *customerRepository) Create(ctx context.Context, customer *domain.Customer) error {
    query := `
        INSERT INTO customers (
            id, name, email, address, contact_number, alternate_contact,
            whatsapp_number, photo_url, business_name, gstin, customer_type,
            aadhaar_verified, kyc_document_type, kyc_document_number,
            credit_limit, current_balance, payment_terms_days, interest_rate,
            status, notes, tags, total_purchases, loyalty_points
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14,
            $15, $16, $17, $18, $19, $20, $21, $22, $23
        )`
    
    _, err := r.db.ExecContext(ctx, query,
        customer.ID, customer.Name, customer.Email, customer.Address,
        customer.ContactNumber, customer.AlternateContact, customer.WhatsappNumber,
        customer.PhotoURL, customer.BusinessName, customer.GSTIN, customer.CustomerType,
        customer.AadhaarVerified, customer.KYCDocumentType, customer.KYCDocumentNumber,
        customer.CreditLimit, customer.CurrentBalance, customer.PaymentTermsDays,
        customer.InterestRate, customer.Status, customer.Notes, pq.Array(customer.Tags),
        customer.TotalPurchases, customer.LoyaltyPoints,
    )
    return err
}

func (r *customerRepository) GetByID(ctx context.Context, id string) (*domain.Customer, error) {
    query := `
        SELECT id, name, email, address, contact_number, alternate_contact,
               whatsapp_number, photo_url, business_name, gstin, customer_type,
               aadhaar_verified, kyc_document_type, kyc_document_number,
               credit_limit, current_balance, payment_terms_days, interest_rate,
               status, notes, tags, last_transaction_date, total_purchases,
               loyalty_points, created_at, updated_at
        FROM customers
        WHERE id = $1 AND deleted_at IS NULL`
    
    var c domain.Customer
    var email, address, contactNumber sql.NullString
    // ... map nullable fields
    
    err := r.db.QueryRowContext(ctx, query, id).Scan(/* all fields */)
    if err == sql.ErrNoRows {
        return nil, domain.ErrNotFound
    }
    return &c, err
}

// Implement List, Update, SoftDelete, etc.
```

**Action Items:**
1. Extract all SQL queries from `customers.go`
2. Move to `internal/repository/postgres/customer.go`
3. Implement all repository interface methods
4. Handle nullable fields properly
5. Return domain errors (not SQL errors)

### Step 2.2: Implement Inventory Repository

Create `internal/repository/postgres/inventory.go`:

Similar pattern to customer repository.

### Step 2.3: Implement Transaction Repository

Create `internal/repository/postgres/transaction.go`:

Include invoice number generation logic.

### Step 2.4: Implement Remaining Repositories

- `crate.go`
- `wastage.go`
- `expiry_alert.go`
- `sale_item.go`
- `payment_schedule.go`

**Testing:**
Create `internal/repository/postgres/customer_test.go` for each repository.

## Phase 3: Service Layer (Week 2-3)

### Step 3.1: Implement Customer Service

Create `internal/service/customer.go`:

```go
package service

import (
    "context"
    "github.com/google/uuid"
    "github.com/yourusername/vms/backend/internal/domain"
    "github.com/yourusername/vms/backend/internal/repository"
    "time"
)

type CustomerService struct {
    repo repository.CustomerRepository
}

func NewCustomerService(repo repository.CustomerRepository) *CustomerService {
    return &CustomerService{repo: repo}
}

func (s *CustomerService) CreateCustomer(ctx context.Context, customer *domain.Customer) error {
    // Validate
    if err := customer.Validate(); err != nil {
        return err
    }
    
    // Generate ID
    customer.ID = uuid.New().String()
    customer.CreatedAt = time.Now()
    customer.UpdatedAt = time.Now()
    
    // Set defaults
    if customer.Status == "" {
        customer.Status = "active"
    }
    if customer.CustomerType == "" {
        customer.CustomerType = "b2c"
    }
    
    // Create
    return s.repo.Create(ctx, customer)
}

func (s *CustomerService) GetCustomer(ctx context.Context, id string) (*domain.Customer, error) {
    return s.repo.GetByID(ctx, id)
}

func (s *CustomerService) ListCustomers(ctx context.Context, filters domain.CustomerFilters) ([]*domain.Customer, error) {
    return s.repo.List(ctx, filters)
}

func (s *CustomerService) UpdateCustomer(ctx context.Context, customer *domain.Customer) error {
    // Validate
    if err := customer.Validate(); err != nil {
        return err
    }
    
    // Check exists
    existing, err := s.repo.GetByID(ctx, customer.ID)
    if err != nil {
        return err
    }
    if existing == nil {
        return domain.ErrNotFound
    }
    
    customer.UpdatedAt = time.Now()
    return s.repo.Update(ctx, customer)
}

func (s *CustomerService) DeleteCustomer(ctx context.Context, id, deletedBy string, req *domain.DeleteRequest) error {
    // Validate attestation
    if err := req.Validate(); err != nil {
        return err
    }
    
    // Check exists
    _, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return err
    }
    
    // Soft delete
    return s.repo.SoftDelete(ctx, id, deletedBy, req.Reason)
}

func (s *CustomerService) CheckCreditLimit(ctx context.Context, customerID string, amount float64) error {
    customer, err := s.repo.GetByID(ctx, customerID)
    if err != nil {
        return err
    }
    
    if !customer.CanPurchase(amount) {
        return domain.NewBusinessError("CREDIT_LIMIT_EXCEEDED", 
            "Customer has reached credit limit")
    }
    
    return nil
}
```

**Action Items:**
1. Extract business logic from handlers
2. Move validation to service layer
3. Add business rule checks
4. Handle UUID generation
5. Implement soft delete with attestation
6. Add credit limit checking

### Step 3.2: Implement Inventory Service

Create `internal/service/inventory.go`:

Include:
- FEFO sorting logic
- Status calculation
- Margin calculation
- Quantity management

### Step 3.3: Implement Transaction Service

Create `internal/service/transaction.go`:

Include:
- Invoice generation
- Balance calculation
- Multi-item transaction handling
- Inventory quantity updates

**Complex Logic:**
```go
func (s *TransactionService) CreateSale(ctx context.Context, tx *domain.Transaction, items []*domain.SaleItem, userID string) error {
    // Validate transaction
    if err := tx.Validate(); err != nil {
        return err
    }
    
    // Validate all items
    for _, item := range items {
        if err := item.Validate(); err != nil {
            return err
        }
    }
    
    // Check customer credit limit
    if err := s.customerService.CheckCreditLimit(ctx, tx.CustomerID, tx.TotalAmount); err != nil {
        return err
    }
    
    // Start database transaction
    dbTx, err := s.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer dbTx.Rollback()
    
    // Generate invoice number
    tx.InvoiceNumber, err = s.repo.GenerateInvoiceNumber(ctx)
    if err != nil {
        return err
    }
    
    // Create transaction
    if err := s.repo.CreateWithTx(ctx, dbTx, tx); err != nil {
        return err
    }
    
    // Create sale items and update inventory
    for _, item := range items {
        item.TransactionID = tx.ID
        item.CalculateProfit()
        item.CalculateTotal()
        
        if err := s.saleItemRepo.CreateWithTx(ctx, dbTx, item); err != nil {
            return err
        }
        
        // Decrement inventory
        if err := s.inventoryRepo.DecrementQuantityWithTx(ctx, dbTx, item.InventoryLotID, item.Quantity); err != nil {
            return err
        }
    }
    
    // Update customer balance
    newBalance := currentBalance + tx.TotalAmount
    if err := s.customerRepo.UpdateBalanceWithTx(ctx, dbTx, tx.CustomerID, newBalance); err != nil {
        return err
    }
    
    // Commit transaction
    return dbTx.Commit()
}
```

### Step 3.4: Implement Remaining Services

- `wastage.go`
- `expiry_alert.go`
- `crate.go`
- `forecast.go`
- `dashboard.go`
- `report.go`

## Phase 4: Handler Layer (Week 3-4)

### Step 4.1: Implement Customer Handlers

Create `internal/handlers/customer.go`:

```go
package handlers

import (
    "net/http"
    "github.com/go-chi/chi/v5"
    "github.com/yourusername/vms/backend/internal/domain"
    "github.com/yourusername/vms/backend/internal/httputil"
    "github.com/yourusername/vms/backend/internal/service"
)

type CustomerHandler struct {
    service *service.CustomerService
}

func NewCustomerHandler(service *service.CustomerService) *CustomerHandler {
    return &CustomerHandler{service: service}
}

// CreateCustomerRequest represents the HTTP request body
type CreateCustomerRequest struct {
    Name             string   `json:"name"`
    Email            string   `json:"email,omitempty"`
    ContactNumber    string   `json:"contactNumber,omitempty"`
    BusinessName     string   `json:"businessName,omitempty"`
    CustomerType     string   `json:"customerType"`
    CreditLimit      float64  `json:"creditLimit"`
    PaymentTermsDays int      `json:"paymentTermsDays"`
    Status           string   `json:"status"`
    // ... other fields
}

// ToDomain converts request to domain model
func (r *CreateCustomerRequest) ToDomain() *domain.Customer {
    return &domain.Customer{
        Name:             r.Name,
        Email:            r.Email,
        ContactNumber:    r.ContactNumber,
        BusinessName:     r.BusinessName,
        CustomerType:     r.CustomerType,
        CreditLimit:      r.CreditLimit,
        PaymentTermsDays: r.PaymentTermsDays,
        Status:           r.Status,
        // ... other fields
    }
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

func (h *CustomerHandler) GetByID(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    
    customer, err := h.service.GetCustomer(r.Context(), id)
    if err != nil {
        httputil.SendError(w, err)
        return
    }
    
    httputil.SendJSON(w, http.StatusOK, customer)
}

func (h *CustomerHandler) List(w http.ResponseWriter, r *http.Request) {
    filters := domain.CustomerFilters{
        Status:       r.URL.Query().Get("status"),
        CustomerType: r.URL.Query().Get("type"),
        SearchTerm:   r.URL.Query().Get("q"),
    }
    
    customers, err := h.service.ListCustomers(r.Context(), filters)
    if err != nil {
        httputil.SendError(w, err)
        return
    }
    
    httputil.SendJSON(w, http.StatusOK, customers)
}

func (h *CustomerHandler) Update(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    
    var req CreateCustomerRequest
    if err := httputil.DecodeJSON(r, &req); err != nil {
        httputil.SendError(w, err)
        return
    }
    
    customer := req.ToDomain()
    customer.ID = id
    
    if err := h.service.UpdateCustomer(r.Context(), customer); err != nil {
        httputil.SendError(w, err)
        return
    }
    
    httputil.SendJSON(w, http.StatusOK, customer)
}

func (h *CustomerHandler) Delete(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    
    // Get user ID from context (set by auth middleware)
    userID := r.Context().Value("userID").(string)
    
    var req domain.DeleteRequest
    if err := httputil.DecodeJSON(r, &req); err != nil {
        httputil.SendError(w, err)
        return
    }
    
    if err := h.service.DeleteCustomer(r.Context(), id, userID, &req); err != nil {
        httputil.SendError(w, err)
        return
    }
    
    httputil.SendJSON(w, http.StatusOK, map[string]string{
        "message": "Customer deleted successfully",
    })
}
```

**Action Items:**
1. Create request/response DTOs
2. Keep handlers thin (no business logic)
3. Delegate to service layer
4. Use httputil for responses
5. Extract user ID from context

### Step 4.2: Implement Remaining Handlers

- `inventory.go`
- `transaction.go`
- `crate.go`
- `wastage.go`
- `expiry_alert.go`
- `auth.go`
- `dashboard.go`
- `report.go`
- `forecast.go`

## Phase 5: Middleware & Router (Week 4)

### Step 5.1: Extract Middleware

Create `internal/middleware/auth.go`:
```go
package middleware

import (
    "context"
    "net/http"
    "strings"
    "github.com/yourusername/vms/backend/pkg/jwt"
)

func Auth(jwtSecret []byte) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            authHeader := r.Header.Get("Authorization")
            if authHeader == "" {
                http.Error(w, "missing authorization header", http.StatusUnauthorized)
                return
            }
            
            tokenString := strings.TrimPrefix(authHeader, "Bearer ")
            claims, err := jwt.ValidateToken(tokenString, jwtSecret)
            if err != nil {
                http.Error(w, "invalid token", http.StatusUnauthorized)
                return
            }
            
            // Add user ID to context
            ctx := context.WithValue(r.Context(), "userID", claims.UserID)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
```

Create `internal/middleware/cors.go`, `logging.go`, `recovery.go`

### Step 5.2: Create Router Package

Create `internal/router/router.go`:
```go
package router

import (
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "github.com/yourusername/vms/backend/internal/handlers"
    mw "github.com/yourusername/vms/backend/internal/middleware"
)

type Config struct {
    CustomerHandler     *handlers.CustomerHandler
    InventoryHandler    *handlers.InventoryHandler
    TransactionHandler  *handlers.TransactionHandler
    // ... other handlers
    JWTSecret          []byte
}

func Setup(cfg *Config) *chi.Mux {
    r := chi.NewRouter()
    
    // Global middleware
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)
    r.Use(mw.CORS())
    
    r.Route("/api/v1", func(r chi.Router) {
        // Public routes
        r.Post("/register", cfg.AuthHandler.Register)
        r.Post("/login", cfg.AuthHandler.Login)
        
        // Protected routes
        r.Group(func(r chi.Router) {
            r.Use(mw.Auth(cfg.JWTSecret))
            
            // Customers
            r.Get("/customers", cfg.CustomerHandler.List)
            r.Post("/customers", cfg.CustomerHandler.Create)
            r.Get("/customers/{id}", cfg.CustomerHandler.GetByID)
            r.Put("/customers/{id}", cfg.CustomerHandler.Update)
            r.Delete("/customers/{id}", cfg.CustomerHandler.Delete)
            
            // ... other routes
        })
    })
    
    return r
}
```

## Phase 6: Update main.go (Week 4)

### Step 6.1: Refactor main.go

```go
package main

import (
    "context"
    "database/sql"
    "log"
    "net/http"
    "time"
    
    _ "github.com/jackc/pgx/v5/stdlib"
    "github.com/yourusername/vms/backend/internal/config"
    "github.com/yourusername/vms/backend/internal/handlers"
    "github.com/yourusername/vms/backend/internal/repository/postgres"
    "github.com/yourusername/vms/backend/internal/router"
    "github.com/yourusername/vms/backend/internal/service"
)

func main() {
    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("config: %v", err)
    }
    
    // Connect to database
    db, err := setupDatabase(cfg.Database)
    if err != nil {
        log.Fatalf("database: %v", err)
    }
    defer db.Close()
    
    // Initialize repositories
    customerRepo := postgres.NewCustomerRepository(db)
    inventoryRepo := postgres.NewInventoryRepository(db)
    transactionRepo := postgres.NewTransactionRepository(db)
    saleItemRepo := postgres.NewSaleItemRepository(db)
    crateRepo := postgres.NewCrateRepository(db)
    wastageRepo := postgres.NewWastageRepository(db)
    expiryAlertRepo := postgres.NewExpiryAlertRepository(db)
    
    // Initialize services
    customerService := service.NewCustomerService(customerRepo)
    inventoryService := service.NewInventoryService(inventoryRepo)
    transactionService := service.NewTransactionService(
        transactionRepo,
        saleItemRepo,
        inventoryRepo,
        customerRepo,
        db,
    )
    crateService := service.NewCrateService(crateRepo)
    wastageService := service.NewWastageService(wastageRepo, inventoryRepo)
    expiryAlertService := service.NewExpiryAlertService(expiryAlertRepo, inventoryRepo)
    dashboardService := service.NewDashboardService(customerRepo, inventoryRepo, transactionRepo)
    forecastService := service.NewForecastService(cfg.Gemini.APIKey)
    
    // Initialize handlers
    authHandler := handlers.NewAuthHandler(/* ... */)
    customerHandler := handlers.NewCustomerHandler(customerService)
    inventoryHandler := handlers.NewInventoryHandler(inventoryService)
    transactionHandler := handlers.NewTransactionHandler(transactionService)
    crateHandler := handlers.NewCrateHandler(crateService)
    wastageHandler := handlers.NewWastageHandler(wastageService)
    expiryAlertHandler := handlers.NewExpiryAlertHandler(expiryAlertService)
    dashboardHandler := handlers.NewDashboardHandler(dashboardService)
    forecastHandler := handlers.NewForecastHandler(forecastService)
    reportHandler := handlers.NewReportHandler(/* ... */)
    
    // Setup router
    r := router.Setup(&router.Config{
        AuthHandler:        authHandler,
        CustomerHandler:    customerHandler,
        InventoryHandler:   inventoryHandler,
        TransactionHandler: transactionHandler,
        CrateHandler:       crateHandler,
        WastageHandler:     wastageHandler,
        ExpiryAlertHandler: expiryAlertHandler,
        DashboardHandler:   dashboardHandler,
        ForecastHandler:    forecastHandler,
        ReportHandler:      reportHandler,
        JWTSecret:          cfg.JWT.Secret,
    })
    
    // Start server
    addr := cfg.Server.Host + ":" + cfg.Server.Port
    log.Printf("Server listening on %s", addr)
    if err := http.ListenAndServe(addr, r); err != nil {
        log.Fatalf("server: %v", err)
    }
}

func setupDatabase(cfg config.DatabaseConfig) (*sql.DB, error) {
    db, err := sql.Open("pgx", cfg.URL)
    if err != nil {
        return nil, err
    }
    
    db.SetMaxOpenConns(cfg.MaxOpenConns)
    db.SetMaxIdleConns(cfg.MaxIdleConns)
    db.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)
    
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    if err := db.PingContext(ctx); err != nil {
        return nil, err
    }
    
    // Run migrations if enabled
    if cfg.MigrateOnStart {
        log.Println("Running migrations...")
        // Migration logic here
    }
    
    return db, nil
}
```

## Phase 7: Testing (Week 5)

### Step 7.1: Unit Tests

Create tests for each layer:

**Domain Layer Tests:**
```go
// internal/domain/customer_test.go
func TestCustomer_Validate(t *testing.T) {
    tests := []struct {
        name    string
        customer domain.Customer
        wantErr bool
    }{
        {
            name: "valid customer",
            customer: domain.Customer{Name: "John Doe", CustomerType: "b2c"},
            wantErr: false,
        },
        {
            name: "missing name",
            customer: domain.Customer{CustomerType: "b2c"},
            wantErr: true,
        },
        // ... more test cases
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.customer.Validate()
            if (err != nil) != tt.wantErr {
                t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

**Service Layer Tests:**
```go
// internal/service/customer_test.go
func TestCustomerService_CreateCustomer(t *testing.T) {
    // Setup mock repository
    mockRepo := repository.NewMockCustomerRepository()
    service := service.NewCustomerService(mockRepo)
    
    // Test
    customer := &domain.Customer{Name: "John Doe"}
    err := service.CreateCustomer(context.Background(), customer)
    
    // Assert
    if err != nil {
        t.Errorf("CreateCustomer() error = %v", err)
    }
    if customer.ID == "" {
        t.Error("ID should be generated")
    }
}
```

### Step 7.2: Integration Tests

```go
// integration_test.go
func TestCustomerAPI_Integration(t *testing.T) {
    // Setup test database
    db := setupTestDB(t)
    defer db.Close()
    
    // Create service
    repo := postgres.NewCustomerRepository(db)
    service := service.NewCustomerService(repo)
    handler := handlers.NewCustomerHandler(service)
    
    // Create router
    r := chi.NewRouter()
    r.Post("/customers", handler.Create)
    
    // Test request
    body := `{"name":"John Doe","customerType":"b2c"}`
    req := httptest.NewRequest("POST", "/customers", strings.NewReader(body))
    w := httptest.NewRecorder()
    
    r.ServeHTTP(w, req)
    
    // Assert
    if w.Code != http.StatusCreated {
        t.Errorf("expected 201, got %d", w.Code)
    }
}
```

## Phase 8: Documentation & Cleanup (Week 5)

### Step 8.1: Update Documentation

1. Update `README.md` with new structure
2. Create `ARCHITECTURE.md` (already done ✅)
3. Add godoc comments to all public functions
4. Create migration guide

### Step 8.2: Remove Old Files

Once everything is working:
```bash
# Remove old monolithic files
rm backend/customers.go
rm backend/inventory.go
rm backend/transactions.go
# ... etc
```

## Migration Checklist

### Pre-Migration
- [ ] Review ARCHITECTURE.md
- [ ] Understand current code
- [ ] Set up development environment
- [ ] Create feature branch

### Domain Layer
- [x] Create domain models
- [x] Add validation logic
- [x] Define errors
- [ ] Write unit tests

### Repository Layer
- [x] Define interfaces
- [ ] Implement PostgreSQL repositories
- [ ] Handle nullable fields
- [ ] Add transaction support
- [ ] Write integration tests

### Service Layer
- [ ] Implement business logic
- [ ] Add validation
- [ ] Handle transactions
- [ ] Write unit tests with mocks

### Handler Layer
- [ ] Create HTTP handlers
- [ ] Add request/response DTOs
- [ ] Use httputil helpers
- [ ] Write handler tests

### Infrastructure
- [x] Config package
- [x] HTTP utilities
- [ ] Middleware extraction
- [ ] Router setup
- [ ] Update main.go

### Testing
- [ ] Domain tests
- [ ] Service tests
- [ ] Handler tests
- [ ] Integration tests
- [ ] End-to-end tests

### Deployment
- [ ] Update Docker build
- [ ] Update CI/CD
- [ ] Deploy to staging
- [ ] Load testing
- [ ] Deploy to production

## Timeline Summary

- **Week 1**: Foundation + Repository layer
- **Week 2-3**: Service layer
- **Week 3-4**: Handler layer + Middleware
- **Week 4**: Router + main.go refactor
- **Week 5**: Testing + Documentation
- **Total**: 5 weeks with 1 developer

With 2 developers: **3 weeks**
With 3 developers: **2 weeks**

## Success Criteria

1. ✅ All tests passing
2. ✅ Code coverage > 70%
3. ✅ No breaking changes to API
4. ✅ Performance maintained or improved
5. ✅ Documentation complete
6. ✅ Team trained on new structure

## Getting Help

- Review `ARCHITECTURE.md` for design decisions
- Check existing implementations for patterns
- Ask questions in team chat
- Pair program for complex parts

## Notes

- **Don't rush**: Quality over speed
- **Test as you go**: Don't wait until end
- **Review code**: Get feedback early
- **Document**: Write comments and docs
- **Commit often**: Small, focused commits

---

**Next Steps:**
1. Review this plan with the team
2. Set up development environment
3. Start with Phase 1
4. Work through phases systematically
5. Deploy when all phases complete

Good luck with the refactoring! 🚀
