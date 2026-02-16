package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/d4rthvadr/dusky-go/internal/models"
	"github.com/google/uuid"
)

type RegisterUserPayload struct {
	Username        string `json:"username" validate:"required,max=120"`
	Email           string `json:"email" validate:"required,email,max=120"`
	Password        string `json:"password" validate:"required,min=8,max=255"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

// RegisterUser godoc
//
//	@Summary		Register a new user
//	@Description	Register a new user with the provided username, email, and password.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			user	body		RegisterUserPayload	true	"User payload"
//	@Success		201		{object}	models.User
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Router			/auth/register [post]
func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {

	var payload RegisterUserPayload
	if err := readJSON(r, &payload); err != nil {
		h.badRequestError(w, r, err)
		return
	}

	if err := validatorInstance.Struct(payload); err != nil {
		writeValidationError(w, err)
		return
	}

	userModel := models.User{
		Username: payload.Username,
		Email:    payload.Email,
	}

	// hash the password before saving to the database
	if err := userModel.Password.Set(payload.Password); err != nil {
		h.internalServerError(w, r, err)
		return
	}

	plainToken := uuid.New().String()
	hashedToken := hashAndEncodeToken(plainToken)

	invitationExpiry := time.Hour * 24

	if err := h.store.Users.CreateAndInvite(r.Context(), &userModel, hashedToken, invitationExpiry); err != nil {

		h.internalServerError(w, r, err)
		return
	}

	if err := writeResponse(w, http.StatusCreated, userModel); err != nil {
		h.internalServerError(w, r, err)
		return
	}
}

// hashAndEncodeToken takes a plain token string and returns its SHA-256 hash as a hexadecimal string
func hashAndEncodeToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])

}
