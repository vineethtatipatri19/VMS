# Clean Architecture Implementation

## Overview
This document describes the complete clean architecture refactoring of the VMS (Vendor Management System) backend, implemented in Phases 2-6.

## Architecture Layers

### Layer 1: Domain (`internal/domain/`)
**Purpose**: Core business entities and rules
**Dependencies**: None (innermost layer)
**Files**:
- `customer.go` - Customer entity with business rules (CanPurchase, IsOverdue)
- `inventory.go` - Inventory item with expiry/stock validation
- `transaction.go` - Transaction entity with type validation
- `entities.go` - Sale items, crates, wastage, expiry alerts, payment schedules
- `errors.go` - Domain-specific errors (ErrNotFound, ErrInvalidInput, etc.)

**Principles**:
- No external dependencies
- Pure business logic
- Framework-agnostic
- Database-agnostic

### Layer 2: Repository Interfaces (`internal/repository/`)
**Purpose**: Data access abstractions
**Dependencies**: Domain layer only
**Files**:
- `interfaces.go` - 8 repository interfaces (51 total methods)
- `postgres/*.go` - PostgreSQL implementations (9 files, ~1,903 lines)

**Repository Pattern Benefits**:
- Testable service layer (can mock repositories)
- Database-agnostic service layer
- SOLID principles (Dependency Inversion)
- Easy to swap implementations

**Repositories**:
1. **CustomerRepository** - CRUD + balance management
2. **InventoryRepository** - CRUD + stock management + FEFO
3. **TransactionRepository** - CRUD + filtering by type/date
4. **SaleItemRepository** - Sale line items
5. **CrateRepository** - Crate issue/return tracking
6. **WastageRepository** - Wastage logging
7. **ExpiryAlertRepository** - Expiry notifications
8. **PaymentScheduleRepository** - Payment tracking

### Layer 3: Service (`internal/service/`)
**Purpose**: Business logic orchestration
**Dependencies**: Domain + Repository interfaces
**Files**: 8 service files (~1,258 lines)

**Services**:
1. **CustomerService** - Customer management + credit limits
2. **InventoryService** - Stock management + FEFO allocation
3. **TransactionService** - Multi-repo sales/payment coordination
4. **SaleItemService** - Sale item management
5. **CrateService** - Crate balance tracking
6. **WastageService** - Wastage with inventory deduction
7. **ExpiryService** - Alert generation
8. **PaymentService** - Payment schedule management

**Key Characteristics**:
- Accept repository interfaces (DI)
- Coordinate multiple repositories
- Implement business rules
- Return domain models/errors
- No HTTP, no SQL, no framework code

### Layer 4: Handlers (`internal/handlers/`)
**Purpose**: HTTP request/response handling
**Dependencies**: Service layer + httputil
**Files**: 8 handler files (~1,028 lines)

**Handlers**:
1. **CustomerHandler** - 7 HTTP endpoints
2. **InventoryHandler** - 8 HTTP endpoints  
3. **TransactionHandler** - 8 HTTP endpoints
4. **SaleItemHandler** - 5 HTTP endpoints
5. **CrateHandler** - 7 HTTP endpoints
6. **WastageHandler** - 6 HTTP endpoints
7. **ExpiryHandler** - 7 HTTP endpoints
8. **PaymentHandler** - 7 HTTP endpoints

**Responsibilities**:
- JSON decode/encode
- HTTP status codes
- Query parameter parsing
- Delegate to services
- Use httputil for responses

### Layer 5: Middleware (`internal/middleware/`)
**Purpose**: Cross-cutting concerns
**Files**: 4 middleware (~138 lines)

**Middleware**:
1. **CORS** - Frontend access headers
2. **Logging** - Request/response logging
3. **Recovery** - Panic recovery
4. **Auth** - JWT token validation

### Layer 6: Router (`internal/router/`)
**Purpose**: Route organization and registration
**Files**: 9 route files (~249 lines)

**Structure**:
- `router.go` - Main router setup with middleware
- 8 resource route files - Organized by domain

**Features**:
- Public routes (register, login)
- Protected routes (require JWT)
- Middleware chain application
- Clean route organization

### Layer 7: Main (`main.go`)
**Purpose**: Dependency injection and server startup
**Lines**: 128 lines

**Initialization Order**:
1. Load config from environment
2. Initialize database connection
3. Create 8 repository instances
4. Create 8 service instances
5. Create 8 handler instances
6. Setup router with middleware
7. Start HTTP server

## Dependency Flow

```
main.go
  ↓
router (+ middleware)
  ↓
handlers (HTTP layer)
  ↓
services (business logic)
  ↓
repositories (data access)
  ↓
domain (entities + rules)
  ↓
database
```

**Rules**:
- Outer layers depend on inner layers
- Inner layers never depend on outer layers
- Domain has zero dependencies
- Services depend on repository interfaces (not implementations)

## Statistics

### Code Metrics
```
Phase 1 - Domain Layer:     ~597 lines (5 files)
Phase 2 - Repository:       ~1,903 lines (10 files)
Phase 3 - Service:          1,258 lines (8 files)
Phase 4 - Handler:          1,028 lines (8 files)
Phase 5 - Middleware/Router: 387 lines (13 files)
Phase 6 - Main:              128 lines (1 file)

Total: ~5,301 lines across 47 files
Binary: 16MB (compiled)
```

### Endpoints
```
Public:   2 (register, login)
Protected: 55+ endpoints across 8 resources
Total: 57+ HTTP endpoints
```

## Design Patterns

1. **Repository Pattern** - Data access abstraction
2. **Dependency Injection** - Constructor injection via NewXXX functions
3. **Service Layer** - Business logic encapsulation
4. **Adapter Pattern** - Handlers adapt HTTP to service calls
5. **Strategy Pattern** - Repository implementations
6. **Middleware Pattern** - Request/response pipeline

## SOLID Principles

- **Single Responsibility**: Each layer has one responsibility
- **Open/Closed**: Services open for extension, closed for modification
- **Liskov Substitution**: Repository interfaces allow substitution
- **Interface Segregation**: Focused repository interfaces
- **Dependency Inversion**: Services depend on abstractions, not concretions

## Benefits

### Testability
- Mock repositories for service tests
- Mock services for handler tests (with interface extraction)
- Domain logic testable in isolation

### Maintainability
- Clear separation of concerns
- Easy to locate code
- Predictable structure

### Scalability
- Can swap database implementations
- Can add caching layer
- Can split into microservices

### Team Collaboration
- Clear boundaries between layers
- Can work on different layers independently
- Less merge conflicts

## Future Enhancements

1. **Service Interfaces** - Extract interfaces for better handler testing
2. **CQRS** - Separate read/write models for complex queries
3. **Event Sourcing** - Audit trail for all changes
4. **Caching** - Add Redis for frequently accessed data
5. **Message Queue** - Async processing for heavy operations
6. **API Versioning** - Support multiple API versions
7. **GraphQL** - Add GraphQL layer alongside REST
8. **Microservices** - Split by bounded contexts

## Migration Notes

### Old Code (deprecated)
- `main.go.old` - Legacy monolithic main
- `customers.go`, `inventory.go`, etc. - Old handler files
- Mixed concerns (SQL + HTTP + business logic)

### New Code (clean architecture)
- `internal/` - All application code
- Clear layer separation
- Testable components
- Dependency injection

## References

- Clean Architecture by Robert C. Martin
- Domain-Driven Design by Eric Evans
- Hexagonal Architecture (Ports and Adapters)
- The Clean Architecture Blog Post by Uncle Bob

## Contact

For questions about this architecture:
1. Read this document first
2. Check TESTING.md for test strategy
3. Review code in `internal/` packages
4. Follow the dependency flow diagram
