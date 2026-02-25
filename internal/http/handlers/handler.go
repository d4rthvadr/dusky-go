package handlers

import (
	"net/http"

	"github.com/d4rthvadr/dusky-go/internal/auth"
	"github.com/d4rthvadr/dusky-go/internal/config"
	"github.com/d4rthvadr/dusky-go/internal/mailer"
	"github.com/d4rthvadr/dusky-go/internal/store"
	"github.com/d4rthvadr/dusky-go/internal/utils"
)

type Handler struct {
	store            store.Storage
	version          string
	logger           utils.Logger
	mailConfig       config.MailConfig
	mailer           mailer.Client
	isProdEnv        bool
	jwtAuthenticator *auth.JWTAuthenticator
}

func New(store store.Storage, version string, logger utils.Logger, mailConfig config.MailConfig, mailer mailer.Client, jwtAuthenticator *auth.JWTAuthenticator, isProdEnv bool) *Handler {
	return &Handler{
		store:            store,
		version:          version,
		logger:           logger,
		mailConfig:       mailConfig,
		mailer:           mailer,
		isProdEnv:        isProdEnv,
		jwtAuthenticator: jwtAuthenticator,
	}
}

func (h *Handler) ValidateAndParseRequestBody(r *http.Request, w http.ResponseWriter, dst interface{}) error {

	if err := readJSON(r, dst); err != nil {
		h.badRequestError(w, r, err)
		return err
	}

	if err := validatorInstance.Struct(dst); err != nil {
		writeValidationError(w, err)
		return err
	}
	return nil
}
