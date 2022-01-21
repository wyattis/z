package zsql

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Exec interface {
	Exec(query string, arguments ...interface{}) (sql.Result, error)
}

type Execx interface {
	NamedExec(query string, arguments interface{}) (sql.Result, error)
}
type DB interface {
	Begin() (*sql.Tx, error)
}
type DBx interface {
	Beginx() (*sqlx.Tx, error)
}
type TxHandler = func(tx Exec) error
type TxxHandler = func(tx Execx) error

type nopExec struct{}

func (e *nopExec) Exec(query string, arguments ...interface{}) (r sql.Result, err error) {
	return
}
