package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/d4rthvadr/dusky-go/internal/auth"
	errCustom "github.com/d4rthvadr/dusky-go/internal/errors"
	"github.com/d4rthvadr/dusky-go/internal/mailer"
	"github.com/d4rthvadr/dusky-go/internal/models"
	"github.com/golang-jwt/jwt/v5"

	"github.com/google/uuid"
)

type RegisterUserPayload struct {
	Username        string `json:"username" validate:"required,max=120"`
	Email           string `json:"email" validate:"required,email,max=120"`
	Password        string `json:"password" validate:"required,min=8,max=255"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

type UserInvitationWithToken struct {
	models.User
	Token string `json:"token"`
}

type CreateUserTokenPayload struct {
	Email    string `json:"email" validate:"required,email,max=120"`
	Password string `json:"password" validate:"required,min=8,max=255"`
}

type emailDataEnvelope struct {
	UserName      string
	ActivationURL string
}

// RegisterUser godoc
//
//	@Summary		Register a new user
//	@Description	Register a new user with the provided username, email, and password.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			user	body		RegisterUserPayload	true	"User payload"
//	@Success		201		{object}	UserInvitationWithToken
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
		h.logger.Errorf("error hashing password for user: %s error: %s", payload.Email, err.Error())
		h.internalServerError(w, r, nil)
		return
	}

	plainToken := uuid.New().String()
	hashedToken := hashAndEncodeToken(plainToken)

	invitationExpiry := time.Hour * 24

	if err := h.store.Users.CreateAndInvite(r.Context(), &userModel, hashedToken, invitationExpiry); err != nil {

		h.logger.Errorf("error creating user and invitation: %s error: %s", payload.Email, err.Error())
		h.internalServerError(w, r, nil)
		return
	}

	// Ideally we would want to send the email asynchronously or into a message queue.
	// but for simplicity we'll do it synchronously here.

	sendEmailErr := h.sendUserInvitationEmail(userModel.Username, userModel.Email, plainToken, !h.isProdEnv)
	if sendEmailErr != nil {
		h.logger.Errorf("failed to send user invitation email: %s", sendEmailErr.Error())
		// We won't return an error to the client if the email fails to send, but we will log it for debugging purposes.

		h.internalServerError(w, r, nil)
		return
	}

	if err := writeResponse(w, http.StatusCreated, UserInvitationWithToken{
		User:  userModel,
		Token: plainToken,
	}); err != nil {
		h.logger.Errorf("error writing response for user registration: %s error: %s", payload.Email, err.Error())
		h.internalServerError(w, r, nil)
		return
	}
}

func (h *Handler) sendUserInvitationEmail(username, email, token string, isSandBox bool) error {

	activationUrl := h.mailConfig.ApiUrl + "/auth/confirm?token=" + token
	emailData := emailDataEnvelope{
		UserName:      username,
		ActivationURL: activationUrl,
	}

	return h.mailer.Send(mailer.TemplateUserInvitation, username, email, emailData, isSandBox)

}

// CreateUserToken godoc
//
//	@Summary		Create a new authentication token for a user
//	@Description	Generate a new authentication token for a user based on their email and password.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			credentials	body		CreateUserTokenPayload	true	"User credentials"
//	@Success		200			{object}	map[string]string
//	@Failure		400			{object}	error
//	@Failure		401			{object}	error
//	@Failure		500			{object}	error
//	@Router			/auth/token [post]
func (h *Handler) CreateUserToken(w http.ResponseWriter, r *http.Request) {

	var payload CreateUserTokenPayload
	err := h.ValidateAndParseRequestBody(r, w, &payload)
	if err != nil {
		return
	}

	var user models.User
	err = h.store.Users.GetByEmail(r.Context(), payload.Email, &user)

	// We don't want to reveal whether the email exists or not, so we'll return a generic error message for both cases.

	if err != nil {
		switch err {
		case errCustom.ErrResourceNotFound:
			h.logger.Errorf("user not found: %s", payload.Email)
			h.badRequestError(w, r, errors.New("invalid email or password"))
			return
		default:
			h.logger.Errorf("error fetching user by email: %s error: %s", payload.Email, err.Error())
			h.internalServerError(w, r, nil)
			return
		}

	}

	token, err := h.generateTokenForUser(user.ID, h.jwtAuthenticator)

	if err != nil {
		// TODO: map the error to a more user-friendly message if needed
		h.logger.Errorf("error generating JWT token for user: %s error: %s", payload.Email, err.Error())
		h.internalServerError(w, r, nil)
		return
	}

	if err := writeResponse(w, http.StatusOK, map[string]string{"token": token}); err != nil {
		h.logger.Errorf("error writing response for user token generation: %s error: %s", payload.Email, err.Error())
		h.internalServerError(w, r, nil)
		return
	}
}

// generateTokenForUser generates a JWT token for the given user ID using the provided JWTAuthenticator.
func (h *Handler) generateTokenForUser(userID int64, jwtAuthenticator *auth.JWTAuthenticator) (string, error) {

	claims := jwt.MapClaims{
		"sub": userID,
		"aud": jwtAuthenticator.Aud,
		"iss": jwtAuthenticator.Iss,
		"exp": jwtAuthenticator.Exp,
	}

	return jwtAuthenticator.GenerateToken(claims)
}
