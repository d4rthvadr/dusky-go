package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/d4rthvadr/dusky-go/internal/models"
)

type createUserPayload struct {
	Username        string `json:"username" validate:"required,min=3,max=50"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

const UserContextKey string = "user"

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

	user, ok := getUserFromContext(r.Context())
	if !ok {
		fmt.Println("user not found in context")
		app.internalServerError(w, r, nil)
		return
	}

	user.Password = "" // Clear the password before sending the response

	if err := writeResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {

	type FollowUser struct {
		UserID int64 `json:"userId" validate:"required"`
	}

	// Get the authenticated user from the context instead
	var payload FollowUser
	if err := readJSON(r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := validatorInstance.Struct(payload); err != nil {
		writeValidationError(w, err)
		return
	}

	fmt.Println("here", payload.UserID)

	user, ok := getUserFromContext(r.Context())
	if !ok {
		fmt.Println("user not found in context")
		app.internalServerError(w, r, nil)
		return
	}

	if err := app.store.Followers.Follow(r.Context(), user.ID, payload.UserID); err != nil {

		app.internalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	type FollowUser struct {
		UserID int64 `json:"userId" validate:"required"`
	}
	// Get the authenticated user from the context instead
	var payload FollowUser
	if err := readJSON(r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	user, ok := getUserFromContext(r.Context())
	if !ok {
		fmt.Println("user not found in context")
		app.internalServerError(w, r, nil)
		return
	}

	if err := app.store.Followers.Unfollow(r.Context(), user.ID, payload.UserID); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) userContextMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		userID, err := parseIDParam(r, UserIDKey)

		if err != nil {

			app.badRequestError(w, r, errors.New("invalid user ID"))
			return
		}

		// Retrieve the user from the database using the user ID
		// this works for small(less traffic) apps, but for larger apps, we need to lazy load user when its actually needed
		// or construct a callback rather to be called by the handler
		user, err := app.store.Users.GetByID(r.Context(), userID)
		if err != nil {
			app.badRequestError(w, r, err)
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserFromContext(ctx context.Context) (*models.User, bool) {

	user, ok := ctx.Value(UserContextKey).(*models.User)
	return user, ok
}
