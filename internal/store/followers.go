package store

import (
	"context"
	"database/sql"

	errCustom "github.com/d4rthvadr/dusky-go/internal/errors"
)

type FollowerStore struct {
	db *sql.DB
}

// Follow allows a user to follow another user. It takes the follower's ID and the followee's ID as parameters and creates a new entry in the followers table.
func (f *FollowerStore) Follow(ctx context.Context, userID, followerID int64) error {

	query := `
	INSERT INTO user_followers (follower_id, user_id) 
	VALUES ($1, $2)
	`
	_, err := f.db.ExecContext(ctx, query, followerID, userID)
	return errCustom.HandleStorageError(err)
}

// Unfollow allows a user to unfollow another user. It takes the follower's ID and the followee's ID as parameters and deletes the corresponding entry from the followers table.
func (f *FollowerStore) Unfollow(ctx context.Context, userID, followerID int64) error {

	query := `
	DELETE FROM user_followers 
	WHERE follower_id = $1 AND user_id = $2
	`
	_, err := f.db.ExecContext(ctx, query, followerID, userID)
	return errCustom.HandleStorageError(err)
}
