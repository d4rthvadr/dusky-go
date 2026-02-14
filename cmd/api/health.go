package main

import (
	"log"
	"net/http"
)

// HealthCheck godoc
//
//	@Summary		Check the health of the API
//	@Description	Check if the API is running and healthy
//	@Tags			health
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string]string
//	@Failure		500	{object}	error
//	@Router			/health [get]
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {

	data := map[string]string{
		"status":  "ok",
		"version": version,
	}

	if err := writeJSON(w, http.StatusOK, data); err != nil {
		log.Print(err.Error())
		app.internalServerError(w, r, err)
		return
	}
}
