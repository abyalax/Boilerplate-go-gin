package e2e

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/abyalax/Boilerplate-go-gin/src/config/env"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TestDB wraps database connection for E2E tests
type TestDB struct {
	pool *pgxpool.Pool
	t    *testing.T
}

// NewTestDB creates a test database connection
func NewTestDB(t *testing.T) *TestDB {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Load environment configuration
	cfg, err := env.Load()
	if err != nil {
		t.Fatalf("Failed to load environment config: %v", err)
	}

	pool, err := pgxpool.New(ctx, cfg.GetDatabaseURL())
	if err != nil {
		t.Fatalf("Failed to create connection pool: %v", err)
	}

	// Test the connection
	if err := pool.Ping(ctx); err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}

	t.Logf("Connected to database: %s", cfg.GetDatabaseURL())

	return &TestDB{
		pool: pool,
		t:    t,
	}
}

// Close closes the database connection
func (tdb *TestDB) Close() {
	if tdb.pool != nil {
		tdb.pool.Close()
	}
}

// Clear truncates all tables for a clean test state
func (tdb *TestDB) Clear(ctx context.Context) error {
	tables := []string{
		"user_roles",
		"role_permissions",
		"permissions",
		"roles",
		"users",
	}

	for _, table := range tables {
		if _, err := tdb.pool.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)); err != nil {
			return fmt.Errorf("failed to truncate %s: %w", table, err)
		}
	}

	// Reset sequences
	if _, err := tdb.pool.Exec(ctx, "ALTER SEQUENCE users_id_seq RESTART WITH 1"); err != nil {
		return fmt.Errorf("failed to reset users sequence: %w", err)
	}

	return nil
}

// GetConnection returns the pool for direct queries
func (tdb *TestDB) GetConnection() *pgxpool.Pool {
	return tdb.pool
}

// BeforeEach runs before each test
func (tdb *TestDB) BeforeEach(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := tdb.Clear(ctx); err != nil {
		t.Fatalf("Failed to clear database: %v", err)
	}
}

// AfterEach runs after each test
func (tdb *TestDB) AfterEach(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := tdb.Clear(ctx); err != nil {
		t.Errorf("Failed to clear database: %v", err)
	}
}
