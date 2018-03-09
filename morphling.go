package morphling

import (
	"database/sql"
)

type morphling struct {
	Main *sql.DB
	Replica *sql.DB
}

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


func (m *morphling) QueryRow (query string, args ...interface{}) *sql.Row{
	return m.Replica.QueryRow(query, args)
}

func (m *morphling) Query (query string, args ...interface{}) (*sql.Rows, error){
	return m.Replica.Query(query, args)
}

func (m *morphling) Exec (query string, args ...interface{}) (sql.Result, error){
	return m.Main.Exec(query, args)
}

func (m *morphling) Prepare (query string) (*sql.Stmt, error){
	return m.Main.Prepare(query)
}

func (m *morphling) Begin (query string) (*sql.Tx, error){
	return m.Main.Begin()
}