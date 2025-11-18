# Test Organization

This directory contains integration tests, end-to-end tests, and shared test utilities.

## Directory Structure

```
tests/
├── integration/          # Integration tests (database, external services)
│   ├── api/             # API integration tests
│   ├── repository/      # Database repository tests
│   └── fixtures/        # Test data and fixtures
├── e2e/                 # End-to-end API tests
└── testutil/            # Shared test utilities and helpers
```

## Test Types

### Unit Tests
**Location**: Alongside source code (e.g., `internal/service/*_test.go`)
- Test individual functions/methods
- Use mocks for dependencies
- Fast execution
- Run with: `go test ./internal/...`

### Integration Tests
**Location**: `tests/integration/`
- Test multiple components together
- May use real database/services
- Slower than unit tests
- Run with: `go test ./tests/integration/...`

### E2E Tests
**Location**: `tests/e2e/`
- Test complete workflows through API
- Use real server and database
- Slowest but highest confidence
- Run with: `go test ./tests/e2e/...`

## Running Tests

```bash
# Run all tests
go test ./...

# Run only unit tests
go test ./internal/...

# Run only integration tests
go test ./tests/integration/...

# Run only e2e tests
go test ./tests/e2e/...

# Run with coverage
go test -cover ./...

# Run with verbose output
go test -v ./...
```

## Test Utilities

The `testutil/` package provides shared helpers:
- Database setup/teardown
- Test server creation
- Common assertions
- Fixture loading
- Mock data generation

## Writing Tests

### Unit Test Example (in source directory)
```go
// internal/service/customer_service_test.go
package service

import (
    "testing"
    "github.com/example/pgvms/internal/repository/mocks"
)

func TestCustomerService_Create(t *testing.T) {
    // Use mocks for dependencies
    mockRepo := &mocks.MockCustomerRepository{...}
    service := NewCustomerService(mockRepo)
    // Test logic
}
```

### Integration Test Example (in tests/ directory)
```go
// tests/integration/repository/customer_repository_test.go
package repository_test

import (
    "testing"
    "github.com/example/pgvms/tests/testutil"
)

func TestCustomerRepository_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }
    
    db := testutil.SetupTestDB(t)
    defer testutil.TeardownTestDB(t, db)
    // Test with real database
}
```

## Best Practices

1. **Keep unit tests fast** - Use mocks, avoid I/O
2. **Make integration tests isolated** - Clean database between tests
3. **Use table-driven tests** - Test multiple scenarios efficiently
4. **Name tests clearly** - `TestServiceName_MethodName_Scenario`
5. **Use test fixtures** - Share common test data
6. **Skip slow tests in short mode** - `if testing.Short() { t.Skip() }`

## Coverage Goals

- Domain Layer: ≥ 80%
- Service Layer: ≥ 85%
- Handler Layer: ≥ 70%
- Middleware: ≥ 95%
- Repository: ≥ 70%
- Overall: ≥ 75%

## Current Status

✅ Domain: 63.3% (48 tests)
✅ Service: 53.9% (102 tests)
✅ Handlers: 33.9% (39 tests)
⏳ Middleware: Not started
⏳ Repository: Not started
⏳ Integration: Not started
⏳ E2E: Not started

**Total Unit Tests**: 189 (100% passing)
