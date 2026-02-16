package handlers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/d4rthvadr/dusky-go/internal/models"
)

type createUserPayload struct {
	Username        string `json:"username" validate:"required,min=3,max=50"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

type followUserPayload struct {
	UserID int64 `json:"userId" validate:"required"`
}

const userContextKey string = "user"
const UserIDKey string = "userID"

// CreateUser godoc
//
//	@Summary		Create a new user
//	@Description	Create a new user with the provided username, email, and password.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			user	body		createUserPayload	true	"User payload"
//	@Success		201		{object}	models.User
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Router			/users [post]
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var createUser createUserPayload
	if err := readJSON(r, &createUser); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validatorInstance.Struct(createUser); err != nil {
		writeValidationError(w, err)
		return
	}

	if createUser.Password != createUser.ConfirmPassword {
		writeJSONError(w, http.StatusBadRequest, "passwords do not match")
		return
	}

	userModel := models.User{
		Username: createUser.Username,
		Email:    createUser.Email,
	}

	// hash the password before saving to the database
	if err := userModel.Password.Set(createUser.Password); err != nil {
		h.internalServerError(w, r, err)
		return
	}

	if err := h.store.Users.CreateAndInvite(r.Context(), &userModel, "some-token", time.Hour*24); err != nil {
		h.internalServerError(w, r, err)
		return
	}

	if err := writeResponse(w, http.StatusCreated, userModel); err != nil {
		h.internalServerError(w, r, err)
		return
	}
}

// GetUser godoc
//
//	@Summary		Get a user by ID
//	@Description	Retrieve a user by their ID from the URL and returns it as JSON.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int64	true	"User ID"
//	@Success		200		{object}	models.User
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Router			/users/{userID} [get]
//	@Security		ApiKeyAuth
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	user, ok := getUserFromContext(r.Context())
	if !ok {
		h.internalServerError(w, r, errors.New("user not found in request context"))
		return
	}

	if err := writeResponse(w, http.StatusOK, user); err != nil {
		h.internalServerError(w, r, err)
		return
	}
}

// FollowUser godoc
//
//	@Summary		Follow a user by ID
//	@Description	Follow a user by their ID from the URL
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int64				true	"User ID"
//	@Param			payload	body		followUserPayload	true	"Follower payload"
//	@Success		204		{string}	string				""
//	@Failure		400		{object}	error				"Bad Request"
//	@failure		404		{object}	error				"User not found"
//	@Failure		500		{object}	error				"Internal Server Error"
//	@Router			/users/{userID}/follow [put]
//	@Security		ApiKeyAuth
func (h *Handler) FollowUser(w http.ResponseWriter, r *http.Request) {
	var payload followUserPayload
	if err := readJSON(r, &payload); err != nil {
		h.badRequestError(w, r, err)
		return
	}

	if err := validatorInstance.Struct(payload); err != nil {
		writeValidationError(w, err)
		return
	}

	user, ok := getUserFromContext(r.Context())
	if !ok {
		h.internalServerError(w, r, errors.New("user not found in request context"))
		return
	}

	if err := h.store.Followers.Follow(r.Context(), user.ID, payload.UserID); err != nil {
		h.internalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UnfollowUser godoc
//
//	@Summary		Unfollow a user by ID
//	@Description	Unfollow a user by their ID from the URL
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int64				true	"User ID"
//	@Param			payload	body		followUserPayload	true	"Follower payload"
//	@Success		204		{string}	string				""
//	@Failure		400		{object}	error				"Bad Request"
//	@failure		404		{object}	error				"User not found"
//	@Failure		500		{object}	error				"Internal Server Error"
//	@Router			/users/{userID}/unfollow [put]
//	@Security		ApiKeyAuth
func (h *Handler) UnfollowUser(w http.ResponseWriter, r *http.Request) {
	var payload followUserPayload
	if err := readJSON(r, &payload); err != nil {
		h.badRequestError(w, r, err)
		return
	}

	user, ok := getUserFromContext(r.Context())
	if !ok {
		h.internalServerError(w, r, errors.New("user not found in request context"))
		return
	}

	if err := h.store.Followers.Unfollow(r.Context(), user.ID, payload.UserID); err != nil {
		h.internalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) UserContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := parseIDParam(r, UserIDKey)
		if err != nil {
			h.badRequestError(w, r, errors.New("invalid user ID"))
			return
		}

		user, err := h.store.Users.GetByID(r.Context(), userID)
		if err != nil {
			h.badRequestError(w, r, err)
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserFromContext(ctx context.Context) (*models.User, bool) {
	user, ok := ctx.Value(userContextKey).(*models.User)
	return user, ok
}
