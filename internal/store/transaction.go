package store

import (
	"context"
	"database/sql"
	"fmt"
)

// WitTx is a helper function that wraps the execution of a function within a database transaction.
func WithTx(ctx context.Context, db *sql.DB, fn func(*sql.Tx) error) error {

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		// handle panic and rollback transaction
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after Rollback
		}
	}()

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction rollback error: %v, original error: %w", rbErr, err)
		}
		return err
	}

	return tx.Commit()
}
