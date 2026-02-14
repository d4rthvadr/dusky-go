package handlers

import (
	"github.com/d4rthvadr/dusky-go/internal/store"
	"github.com/d4rthvadr/dusky-go/internal/utils"
)

type Handler struct {
	store   store.Storage
	version string
	logger  utils.Logger
}

func New(store store.Storage, version string, logger utils.Logger) *Handler {
	return &Handler{store: store, version: version, logger: logger}
}
