package store

import (
	"context"
	"database/sql"

	"github.com/d4rthvadr/dusky-go/internal/models"
)

type UserStore struct {
	db *sql.DB
}

func (u *UserStore) Create(ctx context.Context, user *models.User) error {
	// Implementation for creating a user in the database
	return nil
}
