.PHONY: help build run test clean migrate-up migrate-down docker-build docker-run

# OS detection
ifeq ($(OS),Windows_NT)
    SHELL := cmd.exe
    NULL := NUL
    SLEEP := timeout /t 5 > NUL
    RM := rmdir /s /q
    RMFILE := del /f /q
    CHECK_CMD = where
else
    NULL := /dev/null
    SLEEP := sleep 5
    RM := rm -rf
    RMFILE := rm -f
    CHECK_CMD = command -v
endif

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  db-up              Start PostgreSQL database with docker-compose"
	@echo "  db-down            Stop and remove PostgreSQL database"
	@echo "  migrate-up         Run database migrations"
	@echo "  migrate-down       Rollback database migrations"
	@echo "  migrate-status     Show migration status"
	@echo "  run-dev            Run the application in dev mode"
	@echo "  test               Run all tests (unit + E2E)"
	@echo "  test-unit          Run unit tests only"
	@echo "  test-e2e           Run E2E tests only"
	@echo "  test-coverage      Run tests with coverage report"
	@echo "  test-race          Run tests with race detector"
	@echo "  clean              Clean build artifacts"
	@echo "  sqlc               Generate Go code from SQL queries"
	@echo "  install-tools      Install required development tools"

# Database targets
db-up:
	docker compose -f docker-compose.yaml up -d
	@echo "Database is starting, waiting..."
	@$(SLEEP)

db-down:
	docker compose -f docker-compose.yaml down -v

db-logs:
	docker compose -f docker-compose.yaml logs -f postgres

# Migration targets
migrate-up: db-up
	@echo "Running migrations..."
	migrate -path db/migration -database "postgres://boilerplate_go_gin:boilerplate_go_gin@localhost:5432/db_boilerplate_go_gin?sslmode=disable" -verbose up

migrate-down:
	@echo "Rolling back migrations..."
	migrate -path db/migration -database "postgres://boilerplate_go_gin:boilerplate_go_gin@localhost:5432/db_boilerplate_go_gin?sslmode=disable" -verbose down

migrate-status:
	migrate -path db/migration -database "postgres://boilerplate_go_gin:boilerplate_go_gin@localhost:5432/db_boilerplate_go_gin?sslmode=disable" -verbose version

run-dev:
	@echo "Running application..."
	air

# Test targets
test: test-unit test-e2e
	@echo "All tests passed!"

test-unit:
	go test -v -cover ./test/unit/...

test-e2e: migrate-up
	go test -v -timeout 60s -cover ./test/e2e/...

test-coverage:
	go test -v -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html

test-race:
	go test -v -race ./test/unit/...
	go test -v -race -timeout 60s ./test/e2e/...

# Code generation
sqlc:
	@echo "Generating SQL code..."
	@$(CHECK_CMD) sqlc > $(NULL) 2>&1 || (echo "sqlc not installed. Installing..." && go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest)
	sqlc generate

# Clean targets
clean:
	@echo "Cleaning..."
	-$(RM) cmd/api/bin
	-$(RMFILE) coverage.out

# Install tools
install-tools:
	go install -tags postgres github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest