package skit

import (
	"context"
	"database/sql"
)

// ExecTx executes the given function within a database transaction.
func ExecTx(ctx context.Context, dbc *sql.DB, opts *sql.TxOptions, fn func(tx *sql.Tx) error) (err error) {
	tx, err := dbc.BeginTx(ctx, opts)
	if err != nil {
		return ErrWrap(err, "begin tx")
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	if err = fn(tx); err != nil {
		return ErrWrap(err, "fn logic")
	}

	if err = tx.Commit(); err != nil {
		return ErrWrap(err, "commit tx")
	}

	return nil
}
