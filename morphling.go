package morphling

import (
	"database/sql"
)
// morphling is logical database object with main as master physical database
// and replica as slave database with load balancer
type morphling struct {
	Main *sql.DB
	Replica *sql.DB
}

// Open opens master and slave database connection
func Open(driverName, dataSourceMainStr, dataSourceReplicaStr string) (*morphling ,error) {
	morphling := morphling{}

	dbMain, err := sql.Open(driverName, dataSourceMainStr)
	if err != nil {
		return nil, err
	}
	morphling.Main = dbMain

	dbReplica, err := sql.Open(driverName, dataSourceReplicaStr)
	if err != nil {
		return nil, err
	}
	morphling.Replica = dbReplica

	return &morphling, nil
}

// Close closes all database, releasing any open resources.
//
// It is rare to Close a DB, as the DB handle is meant to be
// long-lived and shared between many goroutines.
func (m *morphling) Close () error {

	err := m.Main.Close()
	if err != nil {
		return err
	}

	err = m.Replica.Close()
	if err != nil {
		return err
	}

	return nil
}

// Ping verifies connection to all database is still alive,
// establishing a connection if necessary.
func (m *morphling) Ping () error{
	err := m.Main.Ping()
	if err != nil {
		return err
	}

	err = m.Replica.Ping()
	if err != nil {
		return err
	}

	return nil
}

// QueryRow executes a query that is expected to return at most one row.
// QueryRow always returns a non-nil value. Errors are deferred until
// Row's Scan method is called.
func (m *morphling) QueryRow (query string, args ...interface{}) *sql.Row{
	return m.Replica.QueryRow(query, args)
}

// Query executes a query that returns rows, typically a SELECT.
// The args are for any placeholder parameters in the query.
func (m *morphling) Query (query string, args ...interface{}) (*sql.Rows, error){
	return m.Replica.Query(query, args)
}

// Exec executes a query without returning any rows.
// The args are for any placeholder parameters in the query.
func (m *morphling) Exec (query string, args ...interface{}) (sql.Result, error){
	return m.Main.Exec(query, args)
}

// Prepare creates a prepared statement for later queries or executions.
// Multiple queries or executions may be run concurrently from the
// returned statement.
// The caller must call the statement's Close method
// when the statement is no longer needed.
func (m *morphling) Prepare (query string) (*sql.Stmt, error){
	return m.Main.Prepare(query)
}

// Begin starts a transaction. The default isolation level is dependent on
// the driver.
func (m *morphling) Begin (query string) (*sql.Tx, error){
	return m.Main.Begin()
}