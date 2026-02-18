package handlers

import (
	"github.com/d4rthvadr/dusky-go/internal/config"
	"github.com/d4rthvadr/dusky-go/internal/mailer"
	"github.com/d4rthvadr/dusky-go/internal/store"
	"github.com/d4rthvadr/dusky-go/internal/utils"
)

type Handler struct {
	store      store.Storage
	version    string
	logger     utils.Logger
	mailConfig config.MailConfig
	mailer     mailer.Client
	isProdEnv  bool
}

func New(store store.Storage, version string, logger utils.Logger, mailConfig config.MailConfig, mailer mailer.Client, isProdEnv bool) *Handler {
	return &Handler{
		store:      store,
		version:    version,
		logger:     logger,
		mailConfig: mailConfig,
		mailer:     mailer,
		isProdEnv:  isProdEnv,
	}
}
