package zsql

// Execute a transaction using a closure
func WithTx(db DB, handler TxHandler) (err error) {
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

// Do something with a closure using the same API as WithTx, but no transaction
func WithoutTx(db DB, handler TxHandler) (err error) {
	return handler(&nopExec{})
}

// Execute a sqlx transaction using a closure
func WithTxx(db DBx, handler TxxHandler) (err error) {
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
