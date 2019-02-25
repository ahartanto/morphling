package morphling

import (
	"context"
	"database/sql"
	"time"
)

var (
	MySQLDriver = "mysql"
)

// DB is logical database object with main as master physical database
// and replica as slave database with loadbalancer
type DB struct {
	// main is master physical database
	main *sql.DB

	// replica can be a slave physical database, but for more than 1 slave replica
	// you can put loadbalancer on top of your replica sets that handle the load
	// distribution, which can be round robin or others
	replica *sql.DB
}

// Open opens master and slave database connection
func Open(driverName, dataSourceMainStr, dataSourceReplicaStr string) (*DB, error) {
	Morphling := DB{}

	dbMain, err := sql.Open(driverName, dataSourceMainStr)
	if err != nil {
		return nil, err
	}
	Morphling.main = dbMain

	dbReplica, err := sql.Open(driverName, dataSourceReplicaStr)
	if err != nil {
		return nil, err
	}
	Morphling.replica = dbReplica

	return &Morphling, nil
}

// Close closes all database, releasing any open resources.
//
// It is rare to Close a DB, as the DB handle is meant to be
// long-lived and shared between many goroutines.
func (m *DB) Close() error {

	err := m.main.Close()
	if err != nil {
		return err
	}

	err = m.replica.Close()
	if err != nil {
		return err
	}

	return nil
}

// Ping verifies connection to all database is still alive,
// establishing a connection if necessary.
func (m *DB) Ping() error {
	err := m.main.Ping()
	if err != nil {
		return err
	}

	err = m.replica.Ping()
	if err != nil {
		return err
	}

	return nil
}

// QueryRow executes a query that is expected to return at most one row.
// QueryRow always returns a non-nil value. Errors are deferred until
// Row's Scan method is called.
func (m *DB) QueryRow(query string, args ...interface{}) *sql.Row {
	return m.replica.QueryRow(query, args...)
}

// Query executes a query that returns rows, typically a SELECT.
// The args are for any placeholder parameters in the query.
func (m *DB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return m.replica.Query(query, args...)
}

// QueryRowMaster executes a query from master DB that is expected to return at most one row.
// QueryRow always returns a non-nil value. Errors are deferred until
// Row's Scan method is called.
func (m *DB) QueryRowMaster(query string, args ...interface{}) *sql.Row {
	return m.main.QueryRow(query, args...)
}

// QueryMaster executes a query from master DB that returns rows, typically a SELECT.
// The args are for any placeholder parameters in the query.
func (m *DB) QueryMaster(query string, args ...interface{}) (*sql.Rows, error) {
	return m.main.Query(query, args...)
}

// Exec executes a query without returning any rows.
// The args are for any placeholder parameters in the query.
func (m *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return m.main.Exec(query, args...)
}

// Prepare creates a prepared statement for later queries or executions.
// Multiple queries or executions may be run concurrently from the
// returned statement.
// The caller must call the statement's Close method
// when the statement is no longer needed.
func (m *DB) Prepare(query string) (*sql.Stmt, error) {
	return m.main.Prepare(query)
}

// Begin starts a transaction. The default isolation level is dependent on
// the driver.
func (m *DB) Begin() (*sql.Tx, error) {
	return m.main.Begin()
}

// BeginTx starts a transaction.
// The provided context is used until the transaction is committed or rolled back. If the context is canceled, the sql package will roll back the transaction. Tx.Commit will return an error if the context provided to BeginTx is canceled.
// The provided TxOptions is optional and may be nil if defaults should be used. If a non-default isolation level is used that the driver doesn't support, an error will be returned.
func (m *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return m.main.BeginTx(ctx, opts)
}

// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
//
// Expired connections may be closed lazily before reuse.
//
// If d <= 0, connections are reused forever.
func (m *DB) SetConnMaxLifetime(d time.Duration) {
	m.main.SetConnMaxLifetime(d)
	m.replica.SetConnMaxLifetime(d)
	return
}

// SetMaxIdleConns sets the maximum number of connections in the idle
// connection pool.
//
// If MaxOpenConns is greater than 0 but less than the new MaxIdleConns
// then the new MaxIdleConns will be reduced to match the MaxOpenConns limit
//
// If n <= 0, no idle connections are retained.
func (m *DB) SetMaxIdleConns(n int) {
	m.main.SetMaxIdleConns(n)
	m.replica.SetMaxIdleConns(n)
	return
}

// SetMaxOpenConns sets the maximum number of open connections to the database.
//
// If MaxIdleConns is greater than 0 and the new MaxOpenConns is less than
// MaxIdleConns, then MaxIdleConns will be reduced to match the new
// MaxOpenConns limit
//
// If n <= 0, then there is no limit on the number of open connections.
// The default is 0 (unlimited).
func (m *DB) SetMaxOpenConns(n int) {
	m.main.SetMaxOpenConns(n)
	m.replica.SetMaxOpenConns(n)
	return
}
