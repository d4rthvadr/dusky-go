package errors

import (
	"database/sql"
	"errors"
)

// Predefined application-level errors.
var (
	ErrResourceNotFound = errors.New("resource not found")
	ErrInvalidInput     = errors.New("invalid input")
	ErrInternal         = errors.New("internal server error")
)

// HandleStorageError maps storage layer errors to application-level errors.
func HandleStorageError(err error) error {

	switch err {
	case sql.ErrNoRows:
		return ErrResourceNotFound
	default:
		return err
	}
}
