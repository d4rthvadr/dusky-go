package errors

import (
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

// Predefined application-level errors.
var (
	ErrResourceNotFound = errors.New("resource not found")
	ErrInvalidInput     = errors.New("invalid input")
	ErrInternal         = errors.New("internal server error")
	ErrConflict         = errors.New("resource already exists")
)

// HandleStorageError maps storage layer errors to application-level errors.
func HandleStorageError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return ErrResourceNotFound
	default:
		// Check for PostgreSQL specific errors
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			// 23505 is the error code for unique constraint violation
			if pqErr.Code == "23505" {
				return ErrConflict
			}
		}
		return err
	}
}
