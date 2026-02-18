package handlers

import (
	"net/http"
	"time"

	"github.com/d4rthvadr/dusky-go/internal/mailer"
	"github.com/d4rthvadr/dusky-go/internal/models"
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

type emailDataEnvelope struct {
	Username      string
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

	// Ideally we would want to send the email asynchronously or into a message queue.
	// but for simplicity we'll do it synchronously here.

	sendEmailErr := h.sendUserInvitationEmail(userModel.Username, userModel.Email, plainToken, !h.isProdEnv)
	if sendEmailErr != nil {
		h.logger.Errorf("failed to send user invitation email: %s", sendEmailErr.Error())
		// We won't return an error to the client if the email fails to send, but we will log it for debugging purposes.

		h.internalServerError(w, r, sendEmailErr)
		return
	}

	if err := writeResponse(w, http.StatusCreated, UserInvitationWithToken{
		User:  userModel,
		Token: plainToken,
	}); err != nil {
		h.internalServerError(w, r, err)
		return
	}
}

func (h *Handler) sendUserInvitationEmail(username, email, token string, isSandBox bool) error {

	activationUrl := h.mailConfig.ApiUrl + "/auth/confirm?token=" + token
	emailData := emailDataEnvelope{
		Username:      username,
		ActivationURL: activationUrl,
	}

	return h.mailer.Send(mailer.TemplateUserInvitation, username, email, emailData, isSandBox)

}
