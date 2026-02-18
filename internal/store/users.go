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

func (u *UserStore) ActivateUser(ctx context.Context, token string) error {

	return WithTx(ctx, u.db, func(tx *sql.Tx) error {

		var user models.User

		// find the user invitation by token, if not found return an error
		getUserFromInvitationErr := u.getUserFromInvitation(ctx, tx, token, &user)
		if getUserFromInvitationErr != nil {
			return getUserFromInvitationErr
		}

		// update the user to set is_active to true
		user.IsActive = true

		if err := u.Update(ctx, tx, &user); err != nil {
			return err
		}

		if err := u.deleteUserInvitation(ctx, tx, user.ID); err != nil {
			return err
		}

		return nil

	})

}

func (u *UserStore) Update(ctx context.Context, tx *sql.Tx, user *models.User) error {

	query := `
	UPDATE users 
	SET username = $1, email = $2,  activated = $3, updated_at = now()
	WHERE id = $4
	RETURNING updated_at
	`

	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeoutDuration)
	defer cancel()

	err := tx.QueryRowContext(ctx, query, user.Username, user.Email, user.IsActive, user.ID).
		Scan(&user.UpdatedAt)

	return errCustom.HandleStorageError(err)
}

func (u *UserStore) getUserFromInvitation(ctx context.Context, tx *sql.Tx, token string, user *models.User) error {

	query := `
	SELECT u.id, u.username, u.email, u.created_at, u.activated
	FROM users u
	JOIN user_invitations ui ON u.id = ui.user_id
	WHERE ui.token = $1 AND ui.expires_at > now()
	`

	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeoutDuration)
	defer cancel()

	err := u.db.QueryRowContext(ctx, query, token).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.IsActive)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return errCustom.ErrResourceNotFound
		default:
			return errCustom.HandleStorageError(err)
		}
	}

	return nil

}

func (u *UserStore) Create(ctx context.Context, tx *sql.Tx, user *models.User) error {

	query := `
	INSERT INTO users (username, email, password_hash) 
	VALUES ($1, $2, $3) RETURNING id, created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeoutDuration)
	defer cancel()

	err := tx.QueryRowContext(ctx, query, user.Username, user.Email, user.Password.Hash).
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

func (u *UserStore) deleteUserInvitation(ctx context.Context, tx *sql.Tx, userID int64) error {

	query := `
	DELETE FROM user_invitations 
	WHERE user_id = $1 
	`

	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, userID)

	return errCustom.HandleStorageError(err)
}

func (u *UserStore) CreateAndInvite(ctx context.Context, user *models.User, token string, invitationExpiry time.Duration) error {

	return WithTx(ctx, u.db, func(tx *sql.Tx) error {

		if err := u.Create(ctx, tx, user); err != nil {
			return err
		}

		if err := u.createUserInvitation(ctx, tx, user.ID, token, invitationExpiry); err != nil {
			return err
		}

		return nil
	})
}
