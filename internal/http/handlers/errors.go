package handlers

import (
	"errors"
	"log"
	"net/http"
)

func (h *Handler) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		err = errors.New("internal server error")
	}

	log.Printf("internal server\terror: %s path: %s error: %s", err.Error(), r.URL.Path, r.RemoteAddr)
	writeJSONError(w, http.StatusInternalServerError, err.Error())
}

func (h *Handler) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		err = errors.New("bad request")
	}

	log.Printf("bad request error: %s path: %s error: %s", err.Error(), r.URL.Path, r.RemoteAddr)
	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (h *Handler) notFoundError(w http.ResponseWriter, r *http.Request, _ error) {
	log.Printf("not found error: %s path: %s error: %s", "the requested resource could not be found", r.URL.Path, r.RemoteAddr)
	writeJSONError(w, http.StatusNotFound, "the requested resource could not be found")
}
