# Testing Guide

This document describes the test structure and how to run tests in the boilerplate project.

## Test Structure

```
test/
├── e2e/                          # End-to-end tests
│   ├── db.go                     # Test database setup and utilities
│   ├── fixtures.go               # Test data factories and builders
│   ├── http.go                   # HTTP client and response helpers
│   └── user_crud_test.go         # CRUD operation E2E tests
└── unit/
    └── domain/
        ├── value_objects_test.go # Email and Persona value object tests
        └── user_test.go          # User aggregate tests
```

## Unit Tests

Unit tests focus on domain logic validation and isolated behavior.

### Domain Layer Tests

The domain layer tests validate:
- **Email Value Object**: Format validation, immutability, equality
- **Persona Value Object**: Valid personas (admin, user, guest), validation
- **User Aggregate**: Creation, name/email/password validation, domain rules

#### Running Domain Tests

```bash
# Run all domain tests
go test ./test/unit/domain/... -v

# Run only value object tests
go test ./test/unit/domain/... -v -run TestEmail

# Run only user aggregate tests
go test ./test/unit/domain/... -v -run TestNewUser
```

#### Test Coverage

Value Objects:
- `TestNewEmail`: Valid/invalid email formats, length limits
- `TestEmailEquals`: Email comparison
- `TestNewPersona`: Valid persona creation, invalid persona rejection
- `TestPersonaIsValid`: Persona validation logic

User Aggregate:
- `TestNewUser`: User creation with validation, ID initialization
- `TestRestoreUser`: Aggregate reconstruction from persistence
- `TestSetID`: ID persistence callback
- `TestChangeEmail`: Email modification
- `TestChangeName`: Name modification
- `TestChangePassword`: Password modification

## End-to-End Tests

E2E tests validate the complete request/response cycle including:
- HTTP layer (handlers, DTOs)
- Application layer (commands, queries)
- Infrastructure layer (repository, database)
- Domain layer (validation, business rules)

### E2E Test Utilities

**TestDB** (`db.go`):
- Database connection pool management
- Transaction isolation via `Clear()` for test independence
- Automatic sequence reset for consistent IDs

**Fixtures** (`fixtures.go`):
- `UserFixture`: Builder pattern for domain User objects
- `DatabaseUserFixture`: Database operations (create, read, delete)
- Fluent API: `.WithName().WithEmail().WithPassword()`

**HTTPClient** (`http.go`):
- JSON request/response marshaling
- Method shortcuts: `.Post()`, `.Get()`, `.Put()`, `.Delete()`
- Response assertion helpers: `.AssertStatusCode()`
- Response parsing: `.UnmarshalJSON()`

### Running E2E Tests

#### Prerequisites

1. **Database Running**:
```bash
make db-up
```

2. **Migrations Applied**:
```bash
make migrate-up
```

#### Test Execution

```bash
# Run all E2E tests
go test ./test/e2e/... -v -timeout 30s

# Run specific test
go test ./test/e2e/... -v -run TestCreateUser

# Run with race detector
go test ./test/e2e/... -v -race

# Run with coverage
go test ./test/e2e/... -v -cover -coverprofile=coverage.out
go tool cover -html=coverage.out
```

#### E2E Test Suite

The test suite (`user_crud_test.go`) covers:

**Creation Tests**:
- `TestCreateUser`: Valid user creation with all fields
  - Valid user with name, email, password → 201 Created
  - Missing name → 400 Bad Request
  - Invalid email format → 400 Bad Request
  - Missing password → 400 Bad Request

- `TestCreateDuplicateEmail`: Email uniqueness constraint
  - Create user with email → 201 Created
  - Create another user with same email → 409 Conflict

**Read Tests**:
- `TestGetUser`: Retrieve existing user
  - User exists in database → 200 OK with user data
  - Verify ID, name, email in response

- `TestGetUserNotFound`: Handle non-existent user
  - Request user with invalid ID → 404 Not Found

**Update Tests**:
- `TestUpdateUser`: Partial updates
  - Update only name → 200 OK, name changed, email unchanged
  - Update only email → 200 OK, email changed, name unchanged
  - Update multiple fields → 200 OK, all changes applied

- `TestUpdateUserEmail`: Email uniqueness with self-exclusion
  - Create two users
  - Update user 2 email to user 1's email → 409 Conflict
  - Update user 2 email to new unique email → 200 OK

**Delete Tests**:
- `TestDeleteUser`: User deletion
  - Delete existing user → 204 No Content
  - GET deleted user → 404 Not Found
  - Verify user count decremented

**List Tests**:
- `TestListUsers`: Multiple user retrieval
  - Create 3 users
  - List all users → 200 OK with all 3 users
  - Verify list contains correct data

## Running All Tests

```bash
# Run all tests (unit + E2E)
make test

# Run with coverage report
make test-coverage

# Run tests in parallel (only unit tests, E2E requires DB state isolation)
go test ./test/unit/... -v -parallel 4
```

## Makefile Targets

```bash
make test              # Run all tests
make test-coverage     # Run tests with coverage report
make test-unit         # Run only unit tests
make test-e2e          # Run only E2E tests (requires DB running)
```

## Test Database

Each E2E test:
1. Calls `BeforeEach()` → Clears all tables, resets sequences
2. Executes test logic
3. Calls `AfterEach()` → Clears tables again

This ensures test isolation without requiring separate databases.

**Database URL**: Configured via `DATABASE_URL` environment variable
```bash
DATABASE_URL=postgres://user:password@localhost:5432/db go test ./test/e2e/...
```

Default: `postgres://boilerplate_go_gin:boilerplate_go_gin@localhost:5432/boilerplate_go_gin`

## Best Practices

### Unit Tests
- Test domain invariants and business rules
- Use table-driven tests for multiple scenarios
- Isolate from infrastructure

### E2E Tests
- Use fixtures for consistent test data
- Leverage `TestSuite` for setup/teardown
- Test happy path → error cases → edge cases
- Verify both response status and data

### Test Data
- Use `UserFixture.WithX()` builder pattern for clarity
- Create data directly via database for setup efficiency
- Let fixtures handle validation

## Debugging Tests

```bash
# Run tests with detailed output
go test ./test/e2e/... -v -count=1

# Run single test with detailed logs
go test ./test/e2e/... -v -run TestCreateUser -count=1

# Enable race detector (catches concurrent access bugs)
go test ./test/e2e/... -race -v

# Get execution time per test
go test ./test/e2e/... -v -count=1 | grep -E "RUN|PASS|FAIL|---"
```

## Continuous Integration

```bash
# CI-friendly test execution
go test ./... \
  -v \
  -timeout 60s \
  -coverprofile=coverage.out \
  -covermode=atomic \
  ./...

# Generate coverage report
go tool cover -html=coverage.out -o coverage.html
```

## Future Test Additions

- **Integration tests**: Application layer with mocked repositories
- **Load tests**: Concurrent CRUD operations
- **Contract tests**: API contract validation
- **Performance benchmarks**: Response time tracking
- **Security tests**: SQL injection, authentication validation
