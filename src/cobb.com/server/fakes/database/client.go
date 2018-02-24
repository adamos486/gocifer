package database

import (
	"database/sql"
	"database/sql/driver"
)

type FakeDBClient struct {
	ExecCall struct {
		Receives struct {
			Query string
			Args  []interface{}
		}
		Returns struct {
			Result sql.Result
			Error  error
		}
	}
	QueryRowCall struct {
		Receives struct {
			Query string
			Args  []interface{}
		}
		Returns struct {
			Row *sql.Row
		}
	}
	QueryCall struct {
		Receives struct {
			Query string
			Args  []interface{}
		}
		Returns struct {
			Rows  *sql.Rows
			Error error
		}
	}
}

func NewPositiveResult(lastInsertID int64, rowsAffected int64, err error) driver.Result {
	return &fakeResult{
		insertID:     lastInsertID,
		rowsAffected: rowsAffected,
		err:          nil,
	}
}

type fakeResult struct {
	insertID     int64
	rowsAffected int64
	err          error
}

func (r *fakeResult) LastInsertId() (int64, error) {
	return r.insertID, r.err
}
func (r *fakeResult) RowsAffected() (int64, error) {
	return r.rowsAffected, r.err
}

func (f *FakeDBClient) Exec(query string, args ...interface{}) (sql.Result, error) {
	f.ExecCall.Receives.Query = query
	f.ExecCall.Receives.Args = args

	return f.ExecCall.Returns.Result, f.ExecCall.Returns.Error
}

func (f *FakeDBClient) QueryRow(query string, args ...interface{}) *sql.Row {
	f.QueryRowCall.Receives.Query = query
	f.QueryRowCall.Receives.Args = args

	return f.QueryRowCall.Returns.Row
}

func (f *FakeDBClient) Query(query string, args ...interface{}) (*sql.Rows, error) {
	f.QueryCall.Receives.Query = query
	f.QueryCall.Receives.Args = args

	return f.QueryCall.Returns.Rows, f.QueryCall.Returns.Error
}
