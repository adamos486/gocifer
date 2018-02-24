package database

import (
	"database/sql"
)

type DBClient interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

type Client struct {
	DbClient DBClient
}

func NewClient(db DBClient) Client {
	return Client{DbClient: db}
}
