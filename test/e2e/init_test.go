package e2e

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/abyalax/Boilerplate-go-gin/src/bootstrap"
	"github.com/abyalax/Boilerplate-go-gin/src/config/env"
)

// TestSuite manages E2E test setup and teardown
type TestSuite struct {
	app        *bootstrap.App
	httpClient *HTTPClient
	testDB     *TestDB
	t          *testing.T
	cancel     context.CancelFunc
}

func NewTestSuite(t *testing.T) *TestSuite {
	// Set environment variables for testing
	os.Setenv("SERVER_PORT", "4000")
	os.Setenv("SERVER_HOST", "0.0.0.0")
	os.Setenv("JWT_SECRET", "test-secret-key-for-e2e-tests")
	os.Setenv("ENV", "development")

	// Set database environment variables for testing
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "boilerplate_go_gin")
	os.Setenv("DB_PASSWORD", "boilerplate_go_gin")
	os.Setenv("DB_NAME", "db_boilerplate_go_gin")
	os.Setenv("DB_SSLMODE", "disable")

	baseURL := "http://localhost:4000"
	httpClient := NewHTTPClient(t, baseURL)

	// Check if server already running
	resp, err := httpClient.Get("/api/v1/health")
	if err == nil && resp.StatusCode == http.StatusOK {
		t.Log("Using existing running server on port 4000")

		// Always clear database even when reusing server
		testDB := NewTestDB(t)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := testDB.Clear(ctx); err != nil {
			t.Fatalf("Failed to clear database: %v", err)
		}
		// Don't close testDB here - let TestSuite.Close() handle it

		return &TestSuite{
			app:        nil,
			httpClient: httpClient,
			testDB:     testDB, // Keep testDB for cleanup
			t:          t,
			cancel:     cancel,
		}
	}

	testDB := NewTestDB(t)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := testDB.Clear(ctx); err != nil {
		t.Fatalf("Failed to clear database: %v", err)
	}

	// Load environment configuration
	cfg, err := env.Load()
	if err != nil {
		t.Fatalf("Failed to load environment config: %v", err)
	}

	app, err := bootstrap.NewApp(cfg)
	if err != nil {
		t.Fatalf("Failed to initialize app: %v", err)
	}

	go func() {
		if err := app.Start(); err != nil && err != http.ErrServerClosed {
			t.Logf("Server error: %v", err)
		}
	}()

	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		time.Sleep(500 * time.Millisecond)

		resp, err := httpClient.Get("/api/v1/health")
		if err == nil && resp.StatusCode == http.StatusOK {
			break
		}

		t.Logf("Waiting for server to be ready... (attempt %d/%d)", i+1, maxRetries)

		if i == maxRetries-1 {
			t.Fatalf("Server failed to become ready")
		}
	}

	return &TestSuite{
		app:        app,
		httpClient: httpClient,
		testDB:     testDB,
		t:          t,
		cancel:     cancel,
	}
}

func (ts *TestSuite) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if ts.app != nil {
		ts.app.Stop(ctx)
	}

	if ts.testDB != nil {
		ts.testDB.Close()
	}

	if ts.cancel != nil {
		ts.cancel()
	}
}

func (ts *TestSuite) AfterEach() {
	ts.testDB.AfterEach(ts.t)
}

// TestHealthCheck tests the health endpoint
func TestHealthCheck(t *testing.T) {
	suite := NewTestSuite(t)
	defer suite.Close()

	resp, err := suite.httpClient.Get("/api/v1/health")
	if err != nil {
		t.Fatalf("Failed to check health: %v", err)
	}

	resp.AssertStatusCode(t, http.StatusOK)
}

// TestReadyCheck tests the ready endpoint
func TestReadyCheck(t *testing.T) {
	suite := NewTestSuite(t)
	defer suite.Close()

	resp, err := suite.httpClient.Get("/api/v1/ready")
	if err != nil {
		t.Fatalf("Failed to check ready: %v", err)
	}

	// Should be OK if database is reachable
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusServiceUnavailable {
		t.Errorf("Expected status 200 or 503, got %d", resp.StatusCode)
	}
}
