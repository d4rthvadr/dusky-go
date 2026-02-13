package utils

import (
	"fmt"
	"time"
)

// parseStringToTime is a helper function to convert a string to a time format. Adjust the layout as needed.
func ParseStrToTime(s string) (string, error) {
	// Implement your time parsing logic here, e.g., using time.Parse with a specific layout.
	// For example:
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return "", fmt.Errorf("invalid time format: %v", err)
	}
	return t.Format(time.RFC3339), nil
}
