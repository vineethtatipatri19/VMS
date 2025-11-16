# Backend Refactoring Implementation Status

## Progress Overview

### ✅ Phase 1: Foundation (COMPLETED)
- Domain models created with business logic
- Repository interfaces defined
- Configuration management setup
- HTTP utilities for responses
- Error handling framework
- Architecture documentation

### 🔄 Phase 2-8: In Progress

Due to the complexity of refactoring an entire monolithic application and the need to ensure the team can continue development, we're taking a **pragmatic, incremental approach**:

## Implementation Strategy

### Approach: Hybrid Gradual Migration

Instead of refactoring everything at once, we'll:

1. **Keep existing handlers working** while adding new architecture
2. **Refactor one feature at a time** (e.g., start with Customers)
3. **Test each refactored feature** before moving to next
4. **Remove old code gradually** after new code is proven

This allows:
- ✅ Continuous development by team
- ✅ Safer migration with less risk
- ✅ Ability to rollback if needed
- ✅ Learning and adapting the architecture

## Next Steps (Recommended Order)

### Step 1: Complete Repository Layer (1 week)
Create PostgreSQL repository implementations for:
- `internal/repository/postgres/customer.go`
- `internal/repository/postgres/inventory.go`
- `internal/repository/postgres/transaction.go`
- `internal/repository/postgres/crate.go`
- `internal/repository/postgres/wastage.go`
- `internal/repository/postgres/expiry_alert.go`
- `internal/repository/postgres/payment_schedule.go`

**Key Points:**
- Extract SQL queries from current handlers
- Handle nullable fields properly (`sql.NullString`, `sql.NullTime`)
- Return domain errors (`domain.ErrNotFound`, etc.)
- Add transaction support methods

**Files to reference:**
- Current handlers: `customers.go`, `inventory.go`, `dashboard.go`, etc.
- Domain models: `internal/domain/*.go`
- Repository interfaces: `internal/repository/interfaces.go`

### Step 2: Create Service Layer (1-2 weeks)
Implement business logic services:
- `internal/service/customer_service.go` - Customer CRUD + validation
- `internal/service/inventory_service.go` - FEFO logic, stock management
- `internal/service/transaction_service.go` - Sales, payments, multi-repo coordination
- `internal/service/crate_service.go` - Crate tracking
- `internal/service/wastage_service.go` - Wastage recording
- `internal/service/expiry_service.go` - Expiry alerts
- `internal/service/payment_service.go` - Payment schedules

**Key Points:**
- Inject repository dependencies
- Implement business rules (credit limits, FEFO, validation)
- Coordinate multiple repositories for complex operations
- Return domain errors

### Step 3: Refactor Handlers (1 week)
Create new thin handlers that use services:
- `internal/handlers/customer_handler.go`
- `internal/handlers/inventory_handler.go`
- `internal/handlers/transaction_handler.go`
- etc.

**Key Points:**
- Parse requests (JSON decode, URL params)
- Call service methods
- Use `httputil.SendJSON()` and `httputil.SendError()`
- Keep handlers thin (< 30 lines per handler)

### Step 4: Extract Middleware & Router (3 days)
- Move middleware to `internal/middleware/`
  - `auth.go` - JWT authentication
  - `cors.go` - CORS handling
  - `logging.go` - Request logging
- Create `internal/router/router.go` for clean route registration

### Step 5: Update main.go (2 days)
Wire everything with dependency injection:
```go
// Initialize database
db := initDB()

// Create repositories
customerRepo := postgres.NewCustomerRepository(db)
inventoryRepo := postgres.NewInventoryRepository(db)
txRepo := postgres.NewTransactionRepository(db)

// Create services
customerService := service.NewCustomerService(customerRepo, txRepo)
inventoryService := service.NewInventoryService(inventoryRepo)
txService := service.NewTransactionService(txRepo, customerRepo, inventoryRepo)

// Create handlers
customerHandler := handlers.NewCustomerHandler(customerService)
inventoryHandler := handlers.NewInventoryHandler(inventoryService)
txHandler := handlers.NewTransactionHandler(txService)

// Setup router
r := router.NewRouter(customerHandler, inventoryHandler, txHandler)

// Start server
http.ListenAndServe(":8080", r)
```

### Step 6: Testing (1 week)
- Unit tests for domain models
- Service tests with mock repositories
- Integration tests for repositories
- Handler tests with mock services

### Step 7: Documentation & Cleanup (2-3 days)
- Update README with new architecture
- Add godoc comments
- Remove old handler files
- Update API documentation

## Reference Implementation

### Example: Customer Feature End-to-End

See `REFACTORING_EXAMPLE.md` for a complete working example of:
1. Repository implementation
2. Service implementation
3. Handler implementation
4. Main.go wiring
5. Tests

This example shows how all layers work together.

## Migration Checklist

### For Each Feature:
- [ ] Extract SQL from old handler → New repository
- [ ] Create service with business logic
- [ ] Create new handler
- [ ] Add route in router (alongside old route)
- [ ] Test new endpoint
- [ ] Switch frontend to use new endpoint
- [ ] Remove old handler code
- [ ] Delete old route

### Features to Migrate (in order):
1. [ ] Customers (start here - simpler, good learning)
2. [ ] Inventory (complex - FEFO logic, status management)
3. [ ] Transactions (complex - multi-repo coordination)
4. [ ] Crates (simpler - good for practice)
5. [ ] Wastage (simpler)
6. [ ] Expiry Alerts (simpler)
7. [ ] Payment Schedules (medium)
8. [ ] Dashboard/Reports (last - aggregations)

## Key Architectural Patterns

### 1. Dependency Injection
```go
// Bad (tight coupling)
func (h *Handler) CreateCustomer() {
    db := getDB()
    db.Exec("INSERT ...")
}

// Good (dependency injection)
type Handler struct {
    service CustomerService
}

func NewHandler(service CustomerService) *Handler {
    return &Handler{service: service}
}
```

### 2. Interface Segregation
```go
// Repository depends on domain
type CustomerRepository interface {
    Create(ctx context.Context, c *domain.Customer) error
    GetByID(ctx context.Context, id string) (*domain.Customer, error)
}

// Service depends on repository interface
type CustomerService struct {
    repo CustomerRepository
}
```

### 3. Error Handling
```go
// Domain errors
if customer == nil {
    return domain.ErrNotFound
}

// Validation errors
if customer.Name == "" {
    return domain.ErrInvalidInput("name is required")
}

// Business errors
if balance > creditLimit {
    return domain.NewBusinessError("CREDIT_LIMIT", "...")
}

// In handler
if err != nil {
    httputil.SendError(w, err) // Maps to correct HTTP status
    return
}
```

### 4. Transaction Management
```go
// In service layer
func (s *Service) ComplexOperation(ctx context.Context) error {
    tx, err := s.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()

    // Multiple operations
    if err := s.repo1.Create(ctx, tx, obj1); err != nil {
        return err
    }
    if err := s.repo2.Update(ctx, tx, obj2); err != nil {
        return err
    }

    return tx.Commit()
}
```

## Testing Strategy

### Unit Tests (Domain)
```go
func TestCustomer_Validate(t *testing.T) {
    c := &domain.Customer{Name: ""}
    err := c.Validate()
    assert.Error(t, err)
}
```

### Service Tests (with Mocks)
```go
func TestCustomerService_CreateCustomer(t *testing.T) {
    mockRepo := &MockCustomerRepository{}
    service := NewCustomerService(mockRepo, nil)
    
    customer := &domain.Customer{Name: "Test"}
    err := service.CreateCustomer(context.Background(), customer)
    
    assert.NoError(t, err)
    assert.True(t, mockRepo.CreateCalled)
}
```

### Integration Tests (Real DB)
```go
func TestCustomerRepository_Create(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()
    
    repo := postgres.NewCustomerRepository(db)
    customer := &domain.Customer{ID: "123", Name: "Test"}
    
    err := repo.Create(context.Background(), customer)
    assert.NoError(t, err)
    
    // Verify in database
    found, err := repo.GetByID(context.Background(), "123")
    assert.NoError(t, err)
    assert.Equal(t, "Test", found.Name)
}
```

## Common Pitfalls to Avoid

### ❌ Don't:
1. Mix business logic in handlers
2. Access database directly from handlers
3. Return database errors to HTTP layer
4. Skip validation in services
5. Create circular dependencies (domain → repo)
6. Forget to handle nullable database fields
7. Skip transaction boundaries for multi-step operations

### ✅ Do:
1. Keep handlers thin (parse, call service, respond)
2. Put all business logic in services
3. Return domain errors from services
4. Validate input in both domain and service layers
5. Use interfaces for dependencies
6. Handle null fields with sql.NullString, etc.
7. Wrap multi-step operations in database transactions
8. Write tests for each layer
9. Document complex business rules
10. Keep functions small and focused

## Getting Help

### Questions to Ask:
1. "Which layer should this code go in?" → Use the architecture guide
2. "How do I handle this SQL null?" → Use `sql.NullString` and helpers
3. "How do I test this service?" → Use mock repositories
4. "This operation needs 3 tables" → Use transactions in service
5. "Should I add this field to domain?" → Yes, domain = database
6. "Should I add this calculation to domain?" → Yes if business rule
7. "How do I return errors?" → Domain errors, mapped in httputil

### Reference Files:
- **Architecture**: `ARCHITECTURE.md` - Complete guide
- **Plan**: `REFACTORING_PLAN.md` - Detailed steps
- **Examples**: `REFACTORING_EXAMPLE.md` - Working code
- **This File**: Implementation status and next steps

## Timeline Estimate

| Phase | Time (1 dev) | Time (3 devs) |
|-------|--------------|---------------|
| Repositories | 1 week | 3 days |
| Services | 2 weeks | 1 week |
| Handlers | 1 week | 3 days |
| Middleware/Router | 3 days | 2 days |
| Main.go DI | 2 days | 1 day |
| Testing | 1 week | 4 days |
| Documentation | 3 days | 2 days |
| **Total** | **6 weeks** | **3 weeks** |

## Success Criteria

### When migration is complete:
- ✅ All features work with new architecture
- ✅ All tests passing
- ✅ No old handler files remain
- ✅ main.go uses dependency injection
- ✅ Clean separation of layers
- ✅ Documentation up to date
- ✅ Team understands architecture
- ✅ New features easy to add

## Current Status: Foundation Complete + Implementation Plan Ready

The architectural foundation is in place. The team can now:
1. Follow this implementation plan
2. Start with Customer feature (simplest)
3. Use REFACTORING_EXAMPLE.md as reference
4. Migrate one feature at a time
5. Keep existing system running during migration

**Next Action**: Assign developers to implement repositories for Customer feature following the example.
