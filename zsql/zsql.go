package zsql

import (
	"context"
	"database/sql"
	"strings"

	"github.com/wyattis/z/zslice/zstrings"
)

// Execute db.Begin with a closure
func WithBegin(db Beginable, handler TxHandler) (err error) {
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

// Call db.BeginTx with a closure
func WithBeginTx(db BeginTxable, handler TxHandler, ctx context.Context, opts *sql.TxOptions) (err error) {
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

// WithBegin without actuallly starting a transaction
func WithBeginNOP(db Beginable, handler TxHandler) (err error) {
	return handler(nopExec{})
}

// Return a string of params the same length as the parameters
func ParamsFor(vals ...string) string {
	return strings.Join(zstrings.Fill(make([]string, len(vals)), "?"), ", ")
}
