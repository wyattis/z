package zsql

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Tx interface {
	Exec(query string, arguments ...interface{}) (r sql.Result, err error)
	ExecContext(ctx context.Context, query string, arguments ...interface{}) (r sql.Result, err error)
	Query(query string, args ...interface{}) (rows *sql.Rows, err error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (rows *sql.Rows, err error)
	QueryRow(query string, args ...interface{}) (row *sql.Row)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) (row *sql.Row)
	Stmt(stmt *sql.Stmt) (s *sql.Stmt)
	StmtContext(ctx context.Context, stmt *sql.Stmt) (s *sql.Stmt)
}

type Txx interface {
	Tx
	MustExec(query string, arguments ...interface{}) (r sql.Result)
	MustExecContext(ctx context.Context, query string, arguments ...interface{}) (r sql.Result)
	NamedExec(query string, arguments interface{}) (r sql.Result, err error)
	NamedExecContext(ctx context.Context, query string, arguments interface{}) (r sql.Result, err error)
	NamedStmt(stmt *sqlx.NamedStmt) (ns *sqlx.NamedStmt)
	NamedStmtContext(ctx context.Context, stmt *sqlx.NamedStmt) (ns *sqlx.NamedStmt)
	NamedQuery(query string, arguments interface{}) (rows *sqlx.Rows, err error)
}
type DB interface {
	Begin() (*sql.Tx, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}
type DBx interface {
	DB
	Beginx() (*sqlx.Tx, error)
	BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)
}
type TxHandler = func(tx Tx) error
type TxxHandler = func(tx Txx) error

type nopExec struct{}

func (e nopExec) Exec(query string, arguments ...interface{}) (r sql.Result, err error) {
	return
}
func (e nopExec) ExecContext(ctx context.Context, query string, arguments ...interface{}) (r sql.Result, err error) {
	return
}
func (e nopExec) Query(query string, args ...interface{}) (rows *sql.Rows, err error) {
	return
}
func (e nopExec) QueryContext(ctx context.Context, query string, args ...interface{}) (rows *sql.Rows, err error) {
	return
}
func (e nopExec) QueryRow(query string, args ...interface{}) (row *sql.Row) {
	return
}
func (e nopExec) QueryRowContext(ctx context.Context, query string, args ...interface{}) (row *sql.Row) {
	return
}
func (e nopExec) Stmt(stmt *sql.Stmt) (s *sql.Stmt) {
	return
}
func (e nopExec) StmtContext(ctx context.Context, stmt *sql.Stmt) (s *sql.Stmt) {
	return
}
func (e nopExec) MustExec(query string, arguments ...interface{}) (r sql.Result) {
	return
}
func (e nopExec) MustExecContext(ctx context.Context, query string, arguments ...interface{}) (r sql.Result) {
	return
}
func (e nopExec) NamedExec(query string, arguments interface{}) (r sql.Result, err error) {
	return
}
func (e nopExec) NamedExecContext(ctx context.Context, query string, arguments interface{}) (r sql.Result, err error) {
	return
}
func (e nopExec) NamedStmt(stmt *sqlx.NamedStmt) (ns *sqlx.NamedStmt) {
	return
}
func (e nopExec) NamedStmtContext(ctx context.Context, stmt *sqlx.NamedStmt) (ns *sqlx.NamedStmt) {
	return
}
func (e nopExec) NamedQuery(query string, arguments interface{}) (rows *sqlx.Rows, err error) {
	return
}
