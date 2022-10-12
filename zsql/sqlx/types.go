package sqlx

import (
	"context"
	"database/sql"

	"github.com/wyattis/z/zsql"
)

type Getable interface {
	Get(dest interface{}, query string, args ...interface{}) error
}

type GetableContext interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

type MustExecable interface {
	MustExec(query string, args ...interface{}) sql.Result
}

type MustExecableContext interface {
	MustExecContext(ctx context.Context, query string, args ...interface{}) sql.Result
}

type NamedExecable interface {
	NamedExec(query string, arg interface{}) (sql.Result, error)
}

type NamedExecableContext interface {
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
}

type NamedPreprable interface {
	PrepareNamed(query string) (NamedStmt, error)
}

type NamedPreprableContext interface {
	PrepareNamedContext(ctx context.Context, query string) (NamedStmt, error)
}

type Selectable interface {
	Select(dest interface{}, query string, args ...interface{}) error
}

type SelectableContext interface {
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

type Stmt interface{}
type NamedStmt interface{}

type Tx interface {
	zsql.Tx
	Getable
	GetableContext
	MustExecable
	MustExecableContext
	NamedExecable
	NamedExecableContext
	NamedPreprable
	NamedExecableContext
	Selectable
	SelectableContext
}

type DBx interface {
	zsql.DB
	BeginTxx(ctx context.Context, opts *sql.TxOptions) (Tx, error)
}
