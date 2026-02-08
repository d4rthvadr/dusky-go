package main

import (
	"errors"
	"net/http"

	"github.com/d4rthvadr/dusky-go/internal/models"
)

type createUserPayload struct {
	Username        string `json:"username" validate:"required,min=3,max=50"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

// UserIDKey is the key used to extract the user ID from the URL parameters.
const UserIDKey string = "userID"

// createUserHandler handles the creation of a new user. It reads the user data from the request body,
func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {

	var user createUserPayload
	if err := readJSON(r, &user); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validatorInstance.Struct(user); err != nil {
		writeValidationError(w, err)
		return
	}

	if user.Password != user.ConfirmPassword {
		writeJSONError(w, http.StatusBadRequest, "passwords do not match")
		return
	}

	userModel := models.User{
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
	}

	if err := app.store.Users.Create(r.Context(), &userModel); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := writeResponse(w, http.StatusCreated, userModel); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

// getUserHandler retrieves a user by their ID from the URL and returns it as JSON.
func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := parseIDParam(r, UserIDKey)

	if err != nil {

		app.badRequestError(w, r, errors.New("invalid user ID"))
		return
	}

	user, err := app.store.Users.GetByID(r.Context(), userID)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	user.Password = "" // Clear the password before sending the response

	if err := writeResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
