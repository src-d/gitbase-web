package testing

import (
	"context"
	"database/sql"
)

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

// QueryContext executes a query
func (db *MockDB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return nil, nil
}

// QueryRow executes a query that is expected to return at most one row
func (db *MockDB) QueryRow(query string, args ...interface{}) *sql.Row {
	return nil
}

// Exec executes a query without returning any rows
func (db *MockDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return nil, nil
}

// As returned by gitbase v0.17.0-rc.4, SELECT UAST('console.log("test")', 'JavaScript') AS uast
const (
	UASTMarshaled     = "\x00\x00\x02\x16\n\x04File\x1a\xfc\x03\n\aProgram\x12\x17\n\finternalRole\x12\aprogram\x12\x14\n\nsourceType\x12\x06module\x1a\xb0\x03\n\x13ExpressionStatement\x12\x14\n\finternalRole\x12\x04body\x1a\xf1\x02\n\x0eCallExpression\x12\x1a\n\finternalRole\x12\nexpression\x1a\xdc\x01\n\x10MemberExpression\x12\x16\n\finternalRole\x12\x06callee\x12\x11\n\bcomputed\x12\x05false\x1aC\n\nIdentifier\x12\x16\n\finternalRole\x12\x06object\x12\x0f\n\x04Name\x12\aconsole*\x04\x10\x01\x18\x012\x06\b\a\x10\x01\x18\b\x1aC\n\nIdentifier\x12\v\n\x04Name\x12\x03log\x12\x18\n\finternalRole\x12\bproperty*\x06\b\b\x10\x01\x18\t2\x06\b\v\x10\x01\x18\f*\x04\x10\x01\x18\x012\x06\b\v\x10\x01\x18\f:\x05\x02\x12\x01TU\x1aR\n\x06String\x12\r\n\x05Value\x12\x04test\x12\n\n\x06Format\x12\x00\x12\x19\n\finternalRole\x12\targuments*\x06\b\f\x10\x01\x18\r2\x06\b\x12\x10\x01\x18\x13:\x02T1*\x04\x10\x01\x18\x012\x06\b\x13\x10\x01\x18\x14:\x02\x12T*\x04\x10\x01\x18\x012\x06\b\x13\x10\x01\x18\x14:\x01\x13*\x04\x10\x01\x18\x012\x06\b\x13\x10\x01\x18\x14:\x019*\x04\x10\x01\x18\x012\x06\b\x13\x10\x01\x18\x14:\x01\""
	UASTMarshaledJSON = "[{\"InternalType\":\"File\",\"StartPosition\":{\"Offset\":0,\"Line\":1,\"Col\":1},\"EndPosition\":{\"Offset\":19,\"Line\":1,\"Col\":20},\"Roles\":[\"Unannotated\",\"File\"],\"Children\":[{\"InternalType\":\"Program\",\"Properties\":{\"internalRole\":\"program\",\"sourceType\":\"module\"},\"StartPosition\":{\"Offset\":0,\"Line\":1,\"Col\":1},\"EndPosition\":{\"Offset\":19,\"Line\":1,\"Col\":20},\"Roles\":[\"Module\"],\"Children\":[{\"InternalType\":\"ExpressionStatement\",\"Properties\":{\"internalRole\":\"body\"},\"StartPosition\":{\"Offset\":0,\"Line\":1,\"Col\":1},\"EndPosition\":{\"Offset\":19,\"Line\":1,\"Col\":20},\"Roles\":[\"Statement\"],\"Children\":[{\"InternalType\":\"CallExpression\",\"Properties\":{\"internalRole\":\"expression\"},\"StartPosition\":{\"Offset\":0,\"Line\":1,\"Col\":1},\"EndPosition\":{\"Offset\":19,\"Line\":1,\"Col\":20},\"Roles\":[\"Expression\",\"Call\"],\"Children\":[{\"InternalType\":\"MemberExpression\",\"Properties\":{\"computed\":\"false\",\"internalRole\":\"callee\"},\"StartPosition\":{\"Offset\":0,\"Line\":1,\"Col\":1},\"EndPosition\":{\"Offset\":11,\"Line\":1,\"Col\":12},\"Roles\":[\"Qualified\",\"Expression\",\"Identifier\",\"Call\",\"Callee\"],\"Children\":[{\"InternalType\":\"Identifier\",\"Properties\":{\"Name\":\"console\",\"internalRole\":\"object\"},\"StartPosition\":{\"Offset\":0,\"Line\":1,\"Col\":1},\"EndPosition\":{\"Offset\":7,\"Line\":1,\"Col\":8},\"Roles\":[],\"Children\":[]},{\"InternalType\":\"Identifier\",\"Properties\":{\"Name\":\"log\",\"internalRole\":\"property\"},\"StartPosition\":{\"Offset\":8,\"Line\":1,\"Col\":9},\"EndPosition\":{\"Offset\":11,\"Line\":1,\"Col\":12},\"Roles\":[],\"Children\":[]}]},{\"InternalType\":\"String\",\"Properties\":{\"Format\":\"\",\"Value\":\"test\",\"internalRole\":\"arguments\"},\"StartPosition\":{\"Offset\":12,\"Line\":1,\"Col\":13},\"EndPosition\":{\"Offset\":18,\"Line\":1,\"Col\":19},\"Roles\":[\"Call\",\"Argument\"],\"Children\":[]}]}]}]}]}]"
)
