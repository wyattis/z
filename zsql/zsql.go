package zsql

import (
	"context"
	"database/sql"
	"strings"

	"github.com/wyattis/z/zslice/zstrings"
)

// Execute db.Begin with a closure
func WithBegin(db DB, handler TxHandler) (err error) {
	txn, err := db.Begin()
	if err != nil {
		return err
	}
	if err = handler(txn); err != nil {
		if err2 := txn.Rollback(); err2 != nil {
			return err2
		}
		return err
	}
	return txn.Commit()
}

// Call db.Begin with a closure
func WithBeginx(db DBx, handler TxxHandler) (err error) {
	txn, err := db.Beginx()
	if err != nil {
		return err
	}
	if err = handler(txn); err != nil {
		if err2 := txn.Rollback(); err2 != nil {
			return err2
		}
		return err
	}
	return txn.Commit()
}

// Call db.BeginTx with a closure
func WithBeginTx(db DB, handler TxHandler, ctx context.Context, opts *sql.TxOptions) (err error) {
	txn, err := db.BeginTx(ctx, opts)
	if err != nil {
		return err
	}
	if err = handler(txn); err != nil {
		if err2 := txn.Rollback(); err2 != nil {
			return err2
		}
		return err
	}
	return txn.Commit()
}

// Call db.BeginTxx with a closure
func WithBeginTxx(db DBx, handler TxxHandler, ctx context.Context, opts *sql.TxOptions) (err error) {
	txn, err := db.BeginTxx(ctx, opts)
	if err != nil {
		return err
	}
	if err = handler(txn); err != nil {
		if err2 := txn.Rollback(); err2 != nil {
			return err2
		}
		return err
	}
	return txn.Commit()
}

// WithBegin without actuallly starting a transaction
func WithBeginNOP(db DB, handler TxHandler) (err error) {
	return handler(nopExec{})
}

// WithBeginx without actuallly starting a transaction
func WithBeginxNOP(db DBx, handler TxxHandler) (err error) {
	return handler(nopExec{})
}

// Return a string of params the same length as the parameters
func ParamsFor(vals ...string) string {
	return strings.Join(zstrings.Fill(make([]string, len(vals)), "?"), ", ")
}
