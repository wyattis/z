package zsql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"time"
)

type Preparable interface {
	Prepare(query string) (Stmt, error)
}
type PreparableContext interface {
	PrepareContext(ctx context.Context, query string) (Stmt, error)
}
type Queryable interface {
	Query(query string, args ...interface{}) (rows *sql.Rows, err error)
}
type QueryableContext interface {
	QueryContext(ctx context.Context, query string, args ...interface{}) (rows *sql.Rows, err error)
}

type Execable interface {
	Exec(query string, arguments ...interface{}) (r sql.Result, err error)
}
type ExecableContext interface {
	ExecContext(ctx context.Context, query string, arguments ...interface{}) (r sql.Result, err error)
}

type Pingable interface {
	Ping() error
}

type PingableContext interface {
	PingContext(ctx context.Context) error
}

type Beginable interface {
	Begin() (Tx, error)
}

type BeginTxable interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (Tx, error)
}

type QueryRowable interface {
	QueryRow(query string, args ...interface{}) (row *sql.Row)
}
type QueryRowableContext interface {
	QueryRowContext(ctx context.Context, query string, args ...interface{}) (row *sql.Row)
}

type Closable interface {
	Close() error
}

type Conn interface {
	BeginTxable
	Closable
}

type Baseable interface {
	Execable
	Pingable
	Preparable
	Queryable
	QueryRowable
}

type BaseableContext interface {
	ExecableContext
	PingableContext
	PreparableContext
	QueryableContext
	QueryRowableContext
}

type DB interface {
	Beginable
	BeginTxable
	Baseable
	BaseableContext
	Closable
	Conn(ctx context.Context) (Conn, error)
}

type FullDB interface {
	DB
	Driver() driver.Driver
	SetConnMaxIdleTime(d time.Duration)
	SetConnMaxLifetime(d time.Duration)
	SetMaxIdleConns(n int)
	SetMaxOpenConns(n int)
	Stats() sql.DBStats
}

type Tx interface {
	Stmt(stmt Stmt) (s Stmt)
	StmtContext(ctx context.Context, stmt Stmt) (s Stmt)
	Rollback() (err error)
	Commit() (err error)
	Execable
	ExecableContext
	Preparable
	PreparableContext
	Queryable
	QueryableContext
	QueryRowable
	QueryRowableContext
}

type Stmt interface {
	Closable
	Execable
	ExecableContext
	Queryable
	QueryableContext
	QueryRowable
	QueryRowableContext
}

type ColumnType interface {
	DatabaseTypeName() string
	DecimalSize() (precision, scale int64, ok bool)
	Length() (length int64, ok bool)
	Name() string
	Nullable() (nullable, ok bool)
}

type Row interface {
	Err() error
	Scan(dest ...any) error
}

type Rows interface {
	Closable
	ColumnTypes() ([]*ColumnType, error)
	Columns() ([]string, error)
	Err() error
	Next() bool
	NextResultSet() bool
	Scan(dest ...any) error
}

type TxHandler = func(tx Tx) error
