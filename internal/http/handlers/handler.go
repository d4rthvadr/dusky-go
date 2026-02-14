package handlers

import "github.com/d4rthvadr/dusky-go/internal/store"

type Handler struct {
	store   store.Storage
	version string
}

func New(store store.Storage, version string) *Handler {
	return &Handler{store: store, version: version}
}
