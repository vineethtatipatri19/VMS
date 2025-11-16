# Testing Documentation

## Overview
This document describes the testing strategy for the VMS (Vendor Management System) clean architecture refactoring.

## Test Structure

### Unit Tests
- **Location**: `internal/*/\*_test.go`
- **Purpose**: Test individual components in isolation
- **Coverage**: Domain logic, business rules, error handling

### Integration Tests
- **Location**: Root level `integration_test.sh`
- **Purpose**: Test end-to-end API workflows
- **Coverage**: Complete HTTP request/response cycles

## Running Tests

### Unit Tests
```bash
# Run all unit tests
go test ./internal/...

# Run tests with coverage
go test -cover ./internal/...

# Run specific package tests
go test ./internal/service
go test ./internal/handlers
go test ./internal/repository/postgres
```

### Integration Tests
```bash
# Run integration test script
./integration_test.sh
```

## Test Coverage Goals

### Phase 7 Testing Strategy
1. **Domain Layer**: Business logic validation (CanPurchase, IsExpired, etc.)
2. **Repository Layer**: Database operations with test database
3. **Service Layer**: Business logic orchestration with mocked repositories
4. **Handler Layer**: HTTP request/response with mocked services
5. **Integration**: Full stack tests with real database

## Manual Testing

### Health Check
```bash
curl http://localhost:8080/api/v1/health
```

### Create Customer
```bash
curl -X POST http://localhost:8080/api/v1/customers \
  -H "Content-Type: application/json" \
  -d '{"name":"Test Customer","contactNumber":"1234567890"}'
```

### Get Customer
```bash
curl http://localhost:8080/api/v1/customers/{id}
```

## Test Database Setup

For integration tests, use a separate test database:

```bash
export DATABASE_URL="postgres://user:pass@localhost:5432/vms_test"
export MIGRATE_ON_START="true"
go test ./...
```

## Mocking Strategy

### Repository Mocks
- Used for testing service layer
- Implement repository interfaces with in-memory storage
- Configurable errors for testing error paths

### Service Mocks
- Used for testing handler layer  
- Would require interface extraction from concrete types
- Alternative: Use integration tests with test database

## CI/CD Integration

```yaml
# Example GitHub Actions workflow
name: Tests
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: 1.21
      - run: go test -v -cover ./...
```

## Future Improvements

1. Extract service interfaces for better testability
2. Add table-driven tests for all CRUD operations
3. Implement gomock or testify for cleaner mocks
4. Add benchmark tests for performance-critical paths
5. Implement mutation testing for test quality
6. Add API contract tests with OpenAPI spec

## Test Best Practices

1. **Arrange-Act-Assert**: Structure all tests clearly
2. **Table-Driven**: Use test tables for multiple scenarios
3. **Isolation**: Each test should be independent
4. **Fast**: Unit tests should run in milliseconds
5. **Deterministic**: No flaky tests, no time dependencies
6. **Descriptive**: Clear test names describing what is tested

## Current Test Status

✅ Domain validation tests (built-in to domain models)
✅ Integration test script (existing)
⚠️  Repository unit tests (pending - requires test DB)
⚠️  Service unit tests (pending - requires mocks)
⚠️  Handler unit tests (pending - requires service interfaces)
✅ Manual testing capability (HTTP endpoints ready)

## Recommended Next Steps

1. Set up test database container
2. Write repository tests with real DB
3. Extract service interfaces
4. Write comprehensive service tests
5. Write handler tests
6. Improve integration test coverage
7. Add CI/CD pipeline
