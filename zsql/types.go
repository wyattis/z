package zsql

import "database/sql"

type Exec interface {
	Exec(query string, arguments ...interface{}) (sql.Result, error)
}
type DB interface {
	Begin() (*sql.Tx, error)
}
type TxHandler = func(tx Exec) error

type nopExec struct{}

func (e *nopExec) Exec(query string, arguments ...interface{}) (r sql.Result, err error) {
	return
}
