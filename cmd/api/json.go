package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

func writeJSON(w http.ResponseWriter, status int, data any) error {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

var validatorInstance *validator.Validate

func init() {
	validatorInstance = validator.New()
}

func writeJSONError(w http.ResponseWriter, status int, message string) error {

	type envelope struct {
		Error string `json:"error"`
	}

	data := envelope{
		Error: message,
	}

	return writeJSON(w, status, data)

}

func readJSON(r *http.Request, dst any) error {

	maxBytes := 1_048_576 // 1 MB
	r.Body = http.MaxBytesReader(nil, r.Body, int64(maxBytes))
	decoder := json.NewDecoder(r.Body)

	// Disallow unknown fields
	decoder.DisallowUnknownFields()
	return decoder.Decode(dst)
}

// formatValidationError converts validator errors into human-readable messages
func formatValidationError(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			fieldName := fieldError.Field()

			var message string
			switch fieldError.Tag() {
			case "required":
				message = fmt.Sprintf("%s is required", fieldName)
			case "max":
				message = fmt.Sprintf("%s must be at most %s characters", fieldName, fieldError.Param())
			case "min":
				message = fmt.Sprintf("%s must be at least %s characters", fieldName, fieldError.Param())
			case "email":
				message = fmt.Sprintf("%s must be a valid email address", fieldName)
			case "url":
				message = fmt.Sprintf("%s must be a valid URL", fieldName)
			case "len":
				message = fmt.Sprintf("%s must be exactly %s characters", fieldName, fieldError.Param())
			default:
				message = fmt.Sprintf("%s is invalid", fieldName)
			}

			errors[strings.ToLower(fieldName)] = message
		}
	}

	return errors
}

// writeValidationError writes formatted validation errors as JSON
func writeValidationError(w http.ResponseWriter, err error) error {
	type envelope struct {
		Error  string            `json:"error"`
		Fields map[string]string `json:"fields,omitempty"`
	}

	data := envelope{
		Error:  "Validation failed",
		Fields: formatValidationError(err),
	}

	return writeJSON(w, http.StatusBadRequest, data)
}
