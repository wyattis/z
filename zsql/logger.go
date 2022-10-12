package zsql

import (
	"context"
	"database/sql"
)

type Logger = func(query string, params ...interface{})

func BaseLogger(db DB, logger Logger) *baseLogger {
	return &baseLogger{
		log:  logger,
		base: db,
	}
}

func DBLogger(db DB, logger Logger) *dbLogger {
	return &dbLogger{
		baseLogger: BaseLogger(db, logger),
		DB:         db,
	}
}

type baseLogger struct {
	log  Logger
	base DB
}

func (l *baseLogger) Exec(query string, args ...any) (sql.Result, error) {
	l.log(query, args...)
	return l.base.Exec(query, args...)
}
func (l *baseLogger) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	l.log(query, args...)
	return l.base.ExecContext(ctx, query, args...)
}
func (l *baseLogger) Query(query string, args ...interface{}) (rows *sql.Rows, err error) {
	l.log(query, args...)
	return l.base.Query(query, args...)
}
func (l *baseLogger) QueryContext(ctx context.Context, query string, args ...interface{}) (rows *sql.Rows, err error) {
	l.log(query, args...)
	return l.base.QueryContext(ctx, query, args...)
}
func (l *baseLogger) QueryRow(query string, args ...interface{}) (row *sql.Row) {
	l.log(query, args...)
	return l.base.QueryRow(query, args...)
}
func (l *baseLogger) QueryRowContext(ctx context.Context, query string, args ...interface{}) (row *sql.Row) {
	l.log(query, args...)
	return l.base.QueryRowContext(ctx, query, args...)
}

type txLogger struct {
	Tx
	log Logger
}

type dbLogger struct {
	*baseLogger
	DB
}

func (l *dbLogger) Begin() (tx Tx, err error) {
	tx, err = l.DB.Begin()
	tx = txLogger{tx, l.log}
	return
}
func (l *dbLogger) BeginTx(ctx context.Context, opts *sql.TxOptions) (tx Tx, err error) {
	tx, err = l.DB.BeginTx(ctx, opts)
	tx = txLogger{tx, l.log}
	return
}
