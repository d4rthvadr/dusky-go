package handlers

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/d4rthvadr/dusky-go/internal/models"
	"github.com/go-chi/chi/v5"
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

// ActivateUserHandler godoc
//
//	@Summary		Activate a user account
//	@Description	Activate a user account using the provided activation token.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			token	path		string	true	"Activation token"
//	@Success		200		{string}	string	"Account activated successfully"
//	@Failure		400		{object}	error	"Bad Request"
//	@Failure		500		{object}	error	"Internal Server Error"
//	@Router			/users/activate/{token} [put]
func (h *Handler) ActivateUserHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	if token == "" {
		writeJSONError(w, http.StatusBadRequest, "activation token is required")
		return
	}

	hashedToken := hashAndEncodeToken(token)

	err := h.store.Users.ActivateUser(r.Context(), hashedToken)
	if err != nil {
		h.internalServerError(w, r, err)
		return

	}

	w.WriteHeader(http.StatusOK)

	if err := writeResponse(w, http.StatusOK, map[string]string{"message": "Account activated successfully"}); err != nil {
		h.internalServerError(w, r, err)
		return
	}
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

func (h *Handler) AuthTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			h.logger.Warnf("missing authorization header")
			h.unauthorizedError(w, r, errors.New("authorization header is required"))
			return
		}

		tokenStr := extractTokenFromHeader(authHeader)
		if tokenStr == "" {
			h.logger.Warnf("invalid authorization header format: %s", authHeader)
			h.unauthorizedError(w, r, errors.New("invalid authorization header format"))
			return
		}

		jwtToken, err := h.jwtAuthenticator.ValidateToken(tokenStr)
		if err != nil {
			h.logger.Warnf("invalid token: %s error: %s", tokenStr, err.Error())
			h.unauthorizedError(w, r, errors.New("invalid or expired token"))
			return
		}

		userID, err := h.jwtAuthenticator.GetUserIDFromClaims(jwtToken)
		if err != nil {
			h.logger.Warnf("failed to get user ID from token claims: %v", err.Error())
			h.unauthorizedError(w, r, err)
			return
		}

		// get user from database to ensure the user still exists and is active
		user, err := h.store.Users.GetByID(r.Context(), userID)
		if err != nil {
			h.logger.Warnf("failed to fetch user by ID from token: %d error: %s", userID, err.Error())
			h.unauthorizedError(w, r, errors.New("invalid or expired token"))
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// extractTokenFromHeader extracts a JWT token from the Authorization header.
// It accepts both "Bearer <token>" and raw "<token>" formats.
func extractTokenFromHeader(authHeader string) string {
	authHeader = strings.TrimSpace(authHeader)
	if authHeader == "" {
		return ""
	}

	const prefix = "Bearer "
	if len(authHeader) > len(prefix) && strings.EqualFold(authHeader[:len(prefix)], prefix) {
		token := strings.TrimSpace(authHeader[len(prefix):])
		if token == "" {
			return ""
		}
		return token
	}

	return authHeader
}

func getUserFromContext(ctx context.Context) (*models.User, bool) {
	user, ok := ctx.Value(userContextKey).(*models.User)
	return user, ok
}
