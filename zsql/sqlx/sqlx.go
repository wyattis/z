package sqlx

import (
	"context"
	"database/sql"
)

type TxHandler = func (tx Tx) (err error)

// Call db.Begin with a closure
func WithBeginx(db DBx, handler TxHandler) (err error) {
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

// Call db.BeginTxx with a closure
func WithBeginTxx(db DBx, handler TxHandler, ctx context.Context, opts *sql.TxOptions) (err error) {
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

// WithBeginx without actuallly starting a transaction
func WithBeginxNOP(db DBx, handler TxxHandler) (err error) {
	return handler(nopExec{})
}
