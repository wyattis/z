package zsql

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

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
func (e nopExec) Stmt(stmt Stmt) (s Stmt) {
	return
}
func (e nopExec) StmtContext(ctx context.Context, stmt Stmt) (s Stmt) {
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
func (e nopExec) Commit() (err error) {
	return
}
func (e nopExec) Rollback() (err error) {
	return
}
func (e nopExec) Prepare(query string) (stmt Stmt, err error) {
	return
}
func (e nopExec) PrepareContext(ctx context.Context, query string) (stmt Stmt, err error) {
	return
}
