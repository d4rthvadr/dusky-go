package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// parseIDParam extracts and parses an int64 ID parameter from the URL.
func parseIDParam(r *http.Request, param string) (int64, error) {
	idStr := chi.URLParam(r, param)
	var id int64
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("invalid %s parameter", param)
	}
	id = int64(idInt)
	return id, nil
}
