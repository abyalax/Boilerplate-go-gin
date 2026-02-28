package e2e

import (
	"context"
	"fmt"
	"testing"
)

// DatabaseUserFixture builder for creating test users
type DatabaseUserFixture struct {
	testDB *TestDB
	name   string
	email  string
	pwd    string
	t      *testing.T
}

// NewDatabaseUserFixture creates a new user fixture builder
func NewDatabaseUserFixture(testDB *TestDB) *DatabaseUserFixture {
	return &DatabaseUserFixture{
		testDB: testDB,
		name:   "Test User",
		email:  "test@example.com",
		pwd:    "hashedpassword",
	}
}

// WithName sets the name for the fixture
func (f *DatabaseUserFixture) WithName(name string) *DatabaseUserFixture {
	f.name = name
	return f
}

// WithEmail sets the email for the fixture
func (f *DatabaseUserFixture) WithEmail(email string) *DatabaseUserFixture {
	f.email = email
	return f
}

// WithPassword sets the password for the fixture
func (f *DatabaseUserFixture) WithPassword(pwd string) *DatabaseUserFixture {
	f.pwd = pwd
	return f
}

// Create inserts the user into the database and returns the user ID
func (f *DatabaseUserFixture) Create(ctx context.Context) (int64, error) {
	pool := f.testDB.GetConnection()

	query := `
		INSERT INTO users (name, email, password)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	var userID int64
	err := pool.QueryRow(ctx, query, f.name, f.email, f.pwd).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("failed to create user fixture: %w", err)
	}

	return userID, nil
}
