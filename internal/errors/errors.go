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

// TODO: map to domain errors

type ErrorCode string

const (
	CodeNotFound ErrorCode = "NOT_FOUND"
	CodeConflict ErrorCode = "CONFLICT"
	CodeInternal ErrorCode = "INTERNAL"
)

type AppError struct {
	Code    ErrorCode
	Message string
	Cause   error
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Cause
}

type StorageError struct {
	Message string
	Err     error
}

// WrapError maps storage layer errors to application-level errors.
// And returns error if it is not mapped to application-level error.
// that way we can log the original error in the handler and return a generic error message to the client.
// we dont want to leak internal error details to the client, but we want to log them for debugging purposes.
func (e *StorageError) WrapError() error {
	if e == nil || e.Err == nil {
		return nil
	}

	switch {
	case errors.Is(e.Err, sql.ErrNoRows):
		return &AppError{Code: CodeNotFound, Message: "resource not found", Cause: e.Err}
	default:
		var pqErr *pq.Error
		if errors.As(e.Err, &pqErr) {
			if pqErr.Code == "23505" {
				return &AppError{Code: CodeConflict, Message: "resource already exists", Cause: e.Err}
			}
		}

		publicMessage := e.Message
		if publicMessage == "" {
			publicMessage = "internal server error"
		}

		return &AppError{Code: CodeInternal, Message: publicMessage, Cause: e.Err}
	}
}

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
