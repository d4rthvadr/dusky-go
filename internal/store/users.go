package store

import (
	"context"
	"database/sql"
	"time"

	errCustom "github.com/d4rthvadr/dusky-go/internal/errors"
	"github.com/d4rthvadr/dusky-go/internal/models"
)

type UserStore struct {
	db *sql.DB
}

func (u *UserStore) Create(ctx context.Context, tx *sql.Tx, user *models.User) error {

	query := `
	INSERT INTO users (username, email, password_hash) 
	VALUES ($1, $2, $3) RETURNING id, created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeoutDuration)
	defer cancel()

	err := tx.QueryRowContext(ctx, query, user.Username, user.Email, user.Password).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	//TODO: handle sql.Errors to user like duplicate email or username, not found, etc
	// before passing to the custom error handler

	return errCustom.HandleStorageError(err)
}

func (u *UserStore) GetByID(ctx context.Context, id int64) (*models.User, error) {

	query := `
	SELECT id, username, email, password_hash, created_at, updated_at
	FROM users
	WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeoutDuration)
	defer cancel()

	var user models.User
	err := u.db.QueryRowContext(ctx, query, id).
		Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, errCustom.HandleStorageError(err)
	}

	return &user, nil
}

func (u *UserStore) createUserInvitation(ctx context.Context, tx *sql.Tx, userId int64, token string, expiry time.Duration) error {

	query := `
	INSERT INTO user_invitations (user_id, token, expires_at) 
	VALUES ($1, $2, $3)
	`

	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeoutDuration)
	defer cancel()

	timeToExp := time.Now().Add(expiry)

	_, err := tx.ExecContext(ctx, query, userId, token, timeToExp)

	return errCustom.HandleStorageError(err)
}

func (u *UserStore) CreateAndInvite(ctx context.Context, user *models.User, token string, invitationExpiry time.Duration) error {

	return WitTx(ctx, u.db, func(tx *sql.Tx) error {

		if err := u.Create(ctx, tx, user); err != nil {
			return err
		}

		if err := u.createUserInvitation(ctx, tx, user.ID, token, invitationExpiry); err != nil {
			return err
		}

		return nil
	})
}
