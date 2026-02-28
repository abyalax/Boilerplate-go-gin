#!/usr/bin/env bash
# migrate.sh - Database migration helper script

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-boilerplate_go_gin}
DB_PASSWORD=${DB_PASSWORD:-boilerplate_go_gin}
DB_NAME=${DB_NAME:-boilerplate_go_gin}
MIGRATE_VERSION=${MIGRATE_VERSION:-latest}

# PostgreSQL connection string
DB_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"

echo -e "${YELLOW}▶ Database Migration Helper${NC}"
echo "================================"

# Function to check if migrate tool is installed
check_migrate() {
    if ! command -v migrate &> /dev/null; then
        echo -e "${RED}✗ migrate tool not found${NC}"
        echo "Install with: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"
        exit 1
    fi
    echo -e "${GREEN}✓ migrate tool found${NC}"
}

# Function to wait for database
wait_for_db() {
    echo -e "${YELLOW}⏳ Waiting for database at $DB_HOST:$DB_PORT...${NC}"
    
    max_attempts=30
    attempt=1
    
    while [ $attempt -le $max_attempts ]; do
        if pg_isready -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" 2>/dev/null; then
            echo -e "${GREEN}✓ Database is ready${NC}"
            return 0
        fi
        
        echo "Attempt $attempt/$max_attempts..."
        sleep 1
        attempt=$((attempt + 1))
    done
    
    echo -e "${RED}✗ Database not ready after ${max_attempts}s${NC}"
    return 1
}

# Function to show migration status
show_status() {
    echo -e "\n${YELLOW}📊 Migration Status:${NC}"
    migrate -path db/migration -database "$DB_URL" version
}

# Function to run migrations up
migrate_up() {
    echo -e "\n${YELLOW}⬆️  Running migrations up...${NC}"
    
    if migrate -path db/migration -database "$DB_URL" -verbose up; then
        echo -e "${GREEN}✓ Migrations applied successfully${NC}"
        show_status
    else
        echo -e "${RED}✗ Migration failed${NC}"
        exit 1
    fi
}

# Function to run migrations down
migrate_down() {
    echo -e "\n${YELLOW}⬇️  Running migrations down...${NC}"
    
    if migrate -path db/migration -database "$DB_URL" -verbose down; then
        echo -e "${GREEN}✓ Migrations rolled back successfully${NC}"
        show_status
    else
        echo -e "${RED}✗ Migration rollback failed${NC}"
        exit 1
    fi
}

# Function to run migrations to specific version
migrate_to() {
    local version=$1
    echo -e "\n${YELLOW}🎯 Migrating to version $version...${NC}"
    
    if migrate -path db/migration -database "$DB_URL" -verbose goto "$version"; then
        echo -e "${GREEN}✓ Migrated to version $version${NC}"
        show_status
    else
        echo -e "${RED}✗ Migration to version $version failed${NC}"
        exit 1
    fi
}

# Function to reset database
reset_db() {
    echo -e "\n${RED}⚠️  Resetting database (all data will be lost)...${NC}"
    read -p "Are you sure? (yes/no): " confirm
    
    if [ "$confirm" != "yes" ]; then
        echo "Cancelled."
        return
    fi
    
    if migrate -path db/migration -database "$DB_URL" -verbose down -all; then
        echo -e "${GREEN}✓ Database reset${NC}"
        show_status
    else
        echo -e "${RED}✗ Database reset failed${NC}"
        exit 1
    fi
}

# Main command parsing
case "${1:-up}" in
    up)
        check_migrate
        wait_for_db
        migrate_up
        ;;
    down)
        check_migrate
        migrate_down
        ;;
    status)
        check_migrate
        show_status
        ;;
    reset)
        check_migrate
        reset_db
        ;;
    goto)
        if [ -z "$2" ]; then
            echo "Usage: $0 goto <version>"
            exit 1
        fi
        check_migrate
        migrate_to "$2"
        ;;
    *)
        echo "Usage: $0 {up|down|status|reset|goto <version>}"
        echo ""
        echo "Commands:"
        echo "  up       - Apply all pending migrations (default)"
        echo "  down     - Rollback one migration"
        echo "  status   - Show migration status"
        echo "  reset    - Rollback all migrations"
        echo "  goto N   - Migrate to specific version"
        echo ""
        echo "Environment Variables:"
        echo "  DB_HOST      - Database host (default: localhost)"
        echo "  DB_PORT      - Database port (default: 5432)"
        echo "  DB_USER      - Database user (default: boilerplate_go_gin)"
        echo "  DB_PASSWORD  - Database password (default: boilerplate_go_gin)"
        echo "  DB_NAME      - Database name (default: boilerplate_go_gin)"
        exit 1
        ;;
esac

echo -e "\n${GREEN}✓ Done${NC}"
