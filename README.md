# Boilerplate Go Gin - Production Grade Backend

Production-grade Go backend boilerplate implementing Clean Architecture, Domain-Driven Design (DDD), and CQRS pattern.

## 📋 Requirements

- Go 1.22+
- PostgreSQL 16+
- Docker & Docker Compose (for local development)
- golang-migrate (for database migrations)

## 🚀 Quick Start

### 1. Install Dependencies

Install required tools:

```bash
make install-tools
```

### 2. Start Database

Start PostgreSQL with Docker Compose:

```bash
make db-up
```

### 3. Run Migrations

```bash
make migrate-up
```

### 4. Run Application

Development mode:

```bash
make run-dev
```

Production mode:

```bash
make build
make run
```

## 🏗 Architecture

### Clean Architecture Layers

- **Domain**: Pure business logic, no framework dependencies
  - Entities: User aggregate root
  - Value Objects: Email, Persona
  - Repository interfaces: UserRepository
  - Domain errors: ErrUserNotFound, etc.

- **Application**: Use cases, commands, and queries
  - Commands: CreateUser
  - Queries: GetUserByID, ListUsers
  - DTOs: UserDTO for responses

- **Infrastructure**: Framework implementations, database, HTTP
  - Persistence: PostgreSQL repository implementation
  - HTTP: Gin handlers and middleware
  - Bootstrap: Dependency injection and app setup

### CQRS Pattern

- **Commands**: Mutate state (CreateUser)
- **Queries**: Read state only (GetUserByID, ListUsers)
- Separated at application layer
- No business logic in handlers

## 📦 Project Structure

```
.
├── cmd/
│   └── api/
│       ├── main.go              # Entry point
│       └── bin/                 # Build output
├── internal/
│   ├── application/
│   │   ├── command/             # Commands
│   │   │   └── create_user.go
│   │   └── query/               # Queries
│   │       └── user_queries.go
│   ├── bootstrap/               # DI & app setup
│   │   └── app.go
│   ├── domain/
│   │   └── user/                # Domain logic
│   │       ├── aggregate.go
│   │       ├── entity.go
│   │       ├── email.go
│   │       ├── persona.go
│   │       ├── repository.go
│   │       └── errors.go
│   └── infrastructure/
│       ├── http/
│       │   ├── handler/         # HTTP handlers
│       │   │   └── user_handler.go
│       │   └── middleware/      # Middleware
│       │       └── middleware.go
│       └── persistence/
│           └── postgres/        # PostgreSQL impl
│               └── user_repository.go
├── db/
│   ├── migration/               # Database migrations
│   │   ├── 000001_create_users_table.up.sql
│   │   └── 000001_create_users_table.down.sql
│   └── query/                   # SQL queries (for sqlc)
│       └── users.sql
├── Dockerfile
├── docker-compose.yaml
├── Makefile
├── go.mod
└── README.md
```


```bash
GET /api/v1/ready

Response (200):
{
  "status": "ready"
}
```


## 📝 Database

### Migrations

* **Up**:

```bash
make migrate-up
```

* **Down**:

```bash
make migrate-down
```

* **Status**:

```bash
make migrate-status
```

---

### Adminer

Untuk mengelola database via UI menggunakan Adminer:

1. **Dapatkan IP container Postgres**

   ```bash
   docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' boilerplate_go_gin_postgres
   ```

2. **Buka Adminer di browser**

   * URL default: `http://localhost:8080`

3. **Isi form koneksi**

   * **System**: PostgreSQL
   * **Server**: IP container yang didapat dari langkah 1 (misal `172.21.0.2`)
   * **Username**: `boilerplate_go_gin`
   * **Password**: `boilerplate_go_gin`
   * **Database**: `db_boilerplate_go_gin`
   * **Schema/Namespace**: `public`

4. **Contoh URL langsung (opsional)**

```
http://localhost:8080/?pgsql=172.21.0.2&username=boilerplate_go_gin&db=db_boilerplate_go_gin&ns=public
```

---


## 🧪 Testing

This project includes comprehensive unit and end-to-end tests for all CRUD operations.

### Quick Test Commands

```bash
# Run all tests (unit + E2E)
make test

# Run unit tests only (fast, no database needed)
make test-unit

# Run E2E tests only (requires database)
make test-e2e

# Run tests with coverage report
make test-coverage

# Run tests with race detector
make test-race
```

### Test Structure

- **Unit Tests** (`test/unit/domain`): Test domain objects, value objects, and business rules
- **E2E Tests** (`test/e2e`): Full integration tests covering the complete request/response cycle

### Testing Guide

For detailed test documentation, setup instructions, and examples, see [TESTING.md](TESTING.md)

### Running Tests in Development

```bash
# Start database
make db-up

# Apply migrations
make migrate-up

# Run all tests
make test
```

## 🐳 Docker

### Build Image

```bash
make docker-build
```

### Run Container

```bash
make docker-run
```

### Push to Registry

```bash
make docker-push
```

## 📊 Observability

- **Logging**: Structured logging with zap
- **Health Checks**: `/api/v1/health` and `/api/v1/ready` endpoints
- **Request Logging**: All HTTP requests logged with method, path, status, and duration
- **Error Handling**: Graceful error handling with proper HTTP status codes

## 🔐 Domain Rules

- **Email Validation**: RFC-compliant email validation
- **Persona Values**: admin, user, or guest
- **User ID**: UUID generated at creation time, immutable
- **Created At**: Set at user creation, immutable

## 📚 Concepts Used

- **Bounded Context**: User context
- **Aggregate Root**: User aggregate
- **Value Objects**: Email, Persona
- **Repository Pattern**: UserRepository interface
- **Dependency Injection**: Constructor-based DI
- **CQRS**: Command (CreateUser) and Query (GetUserByID, ListUsers) separation
- **Clean Architecture**: Clear dependency direction and layering

## 🔧 Configuration

Environment variables:

- `DATABASE_URL`: PostgreSQL connection string (default: `postgres://postgres:postgres@localhost:5432/boilerplate_db?sslmode=disable`)
- `PORT`: Server port (default: `8080`)

## 📖 Development

### Makefile Commands

- `make help` - Show all available commands
- `make install-tools` - Install required tools
- `make db-up` - Start PostgreSQL
- `make db-down` - Stop PostgreSQL
- `make migrate-up` - Run migrations
- `make migrate-down` - Rollback migrations
- `make build` - Build application
- `make run-dev` - Run in development mode
- `make test` - Run tests
- `make clean` - Clean build artifacts

## 🚀 Production Deployment

1. Build Docker image: `make docker-build`
2. Push to registry: `make docker-push`
3. Deploy with orchestration platform (Kubernetes, etc.)

## 📝 License

MIT
