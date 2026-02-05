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

	query := `
	INSERT INTO users (username, email, password_hash) 
	VALUES ($1, $2, $3) RETURNING id, created_at, updated_at
	`
	err := u.db.QueryRowContext(ctx, query, user.Username, user.Email, user.Password).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	return err
}
