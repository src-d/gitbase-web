package testing

import "database/sql"

// MockDB is a mock of *sql.DB
type MockDB struct{}

// Close closes the DB connection
func (db *MockDB) Close() error {
	return nil
}

// Ping pings the DB
func (db *MockDB) Ping() error {
	return nil
}

// Query sends a query to the DB
func (db *MockDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return nil, nil
}

// QueryRow executes a query that is expected to return at most one row
func (db *MockDB) QueryRow(query string, args ...interface{}) *sql.Row {
	return nil
}
