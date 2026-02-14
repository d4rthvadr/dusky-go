package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func parseIDParam(r *http.Request, param string) (int64, error) {
	idStr := chi.URLParam(r, param)
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("invalid %s parameter", param)
	}

	return int64(idInt), nil
}
