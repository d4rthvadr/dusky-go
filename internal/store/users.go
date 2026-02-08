package store

import (
	"context"
	"database/sql"

	errCustom "github.com/d4rthvadr/dusky-go/internal/errors"
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

func (u *UserStore) GetByID(ctx context.Context, id int64) (*models.User, error) {

	query := `
	SELECT id, username, email, password_hash, created_at, updated_at
	FROM users
	WHERE id = $1
	`
	var user models.User
	err := u.db.QueryRowContext(ctx, query, id).
		Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, errCustom.HandleStorageError(err)
	}

	return &user, nil
}
