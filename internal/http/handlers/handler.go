package handlers

import (
	"net/http"

	"github.com/d4rthvadr/dusky-go/internal/auth"
	"github.com/d4rthvadr/dusky-go/internal/cache"
	"github.com/d4rthvadr/dusky-go/internal/config"
	"github.com/d4rthvadr/dusky-go/internal/mailer"
	ratelimiter "github.com/d4rthvadr/dusky-go/internal/ratelmiter"
	"github.com/d4rthvadr/dusky-go/internal/store"
	"github.com/d4rthvadr/dusky-go/internal/utils/logger"
)

type Handler struct {
	store            store.Storage
	cache            cache.CacheStorage
	version          string
	logger           logger.Logger
	mailConfig       config.MailConfig
	mailer           mailer.Client
	isProdEnv        bool
	jwtAuthenticator *auth.JWTAuthenticator
	rateLimiter      ratelimiter.Limiter
}

type HandlerOptions struct {
	Store            store.Storage
	Version          string
	Logger           logger.Logger
	MailConfig       config.MailConfig
	Mailer           mailer.Client
	JWTAuthenticator *auth.JWTAuthenticator
	Cache            cache.CacheStorage
	IsProdEnv        bool
	RateLimiter      ratelimiter.Limiter
}

func New(opts HandlerOptions) *Handler {
	return &Handler{
		store:            opts.Store,
		cache:            opts.Cache,
		version:          opts.Version,
		logger:           opts.Logger,
		mailConfig:       opts.MailConfig,
		mailer:           opts.Mailer,
		isProdEnv:        opts.IsProdEnv,
		jwtAuthenticator: opts.JWTAuthenticator,
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
