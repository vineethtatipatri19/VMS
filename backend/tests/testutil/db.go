package testutil

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/lib/pq"
)

// TestDBConfig holds test database configuration
type TestDBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

// DefaultTestDBConfig returns default configuration for test database
func DefaultTestDBConfig() *TestDBConfig {
	return &TestDBConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "pgvms_user",
		Password: "pgvms_password",
		DBName:   "pgvms_test",
	}
}

// SetupTestDB creates a test database connection and runs migrations
func SetupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	config := DefaultTestDBConfig()
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DBName,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to ping test database: %v", err)
	}

	// Clean database before tests
	CleanTestDB(t, db)

	return db
}

// TeardownTestDB cleans up and closes the test database connection
func TeardownTestDB(t *testing.T, db *sql.DB) {
	t.Helper()

	CleanTestDB(t, db)

	if err := db.Close(); err != nil {
		t.Errorf("Failed to close test database: %v", err)
	}
}

// CleanTestDB removes all data from test tables
func CleanTestDB(t *testing.T, db *sql.DB) {
	t.Helper()

	ctx := context.Background()

	// Order matters due to foreign key constraints
	tables := []string{
		"payment_schedules",
		"expiry_alerts",
		"wastage_log",
		"crate_ledger",
		"sale_items",
		"transactions",
		"inventory_items",
		"customers",
		"users",
	}

	for _, table := range tables {
		query := fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)
		if _, err := db.ExecContext(ctx, query); err != nil {
			// Only log warning, don't fail - table might not exist yet
			t.Logf("Warning: Failed to truncate table %s: %v", table, err)
		}
	}
}

// ExecSQL executes SQL statements (for test data setup)
func ExecSQL(t *testing.T, db *sql.DB, query string, args ...interface{}) {
	t.Helper()

	_, err := db.Exec(query, args...)
	if err != nil {
		t.Fatalf("Failed to execute SQL: %v\nQuery: %s", err, query)
	}
}

// MustQuery executes a query and returns the result, failing the test on error
func MustQuery(t *testing.T, db *sql.DB, query string, args ...interface{}) *sql.Rows {
	t.Helper()

	rows, err := db.Query(query, args...)
	if err != nil {
		t.Fatalf("Failed to execute query: %v\nQuery: %s", err, query)
	}

	return rows
}

// MustExec executes a statement and returns the result, failing the test on error
func MustExec(t *testing.T, db *sql.DB, query string, args ...interface{}) sql.Result {
	t.Helper()

	result, err := db.Exec(query, args...)
	if err != nil {
		t.Fatalf("Failed to execute statement: %v\nQuery: %s", err, query)
	}

	return result
}
