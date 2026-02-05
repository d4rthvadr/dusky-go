package store

import (
	"context"
	"database/sql"

	errCustom "github.com/d4rthvadr/dusky-go/internal/errors"
	"github.com/d4rthvadr/dusky-go/internal/models"
	"github.com/lib/pq"
)

type PostStore struct {
	db *sql.DB
}

// Create inserts a new post into the database and updates the post model with the generated ID and timestamps
func (p *PostStore) Create(ctx context.Context, post *models.Post) error {
	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeoutDuration)
	defer cancel()

	query := `
	INSERT INTO posts (title, content, user_id, tags) 
	VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at
	`
	err := p.db.QueryRowContext(ctx, query, post.Title, post.Content, post.UserID, pq.Array(post.Tags)).
		Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)

	return err
}

// GetByID retrieves a post by its ID, including its comments
func (p *PostStore) GetByID(ctx context.Context, id int64) (*models.Post, error) {
	query := `
	SELECT id, title, content, version, user_id, tags, created_at, updated_at
	FROM posts
	WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeoutDuration)
	defer cancel()

	var post models.Post
	err := p.db.QueryRowContext(ctx, query, id).
		Scan(&post.ID, &post.Title, &post.Content, &post.Version, &post.UserID, pq.Array(&post.Tags), &post.CreatedAt, &post.UpdatedAt)

	err = errCustom.HandleStorageError(err)

	if err != nil {
		return nil, err
	}

	return &post, err
}

func (p *PostStore) Update(ctx context.Context, post *models.Post) error {
	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeoutDuration)
	defer cancel()

	query := `
	UPDATE posts
	SET title = $1, content = $2, tags = $3, version = version + 1
	WHERE id = $4 and version = $5
	RETURNING updated_at, version
	`
	err := p.db.QueryRowContext(ctx, query, post.Title, post.Content, pq.Array(post.Tags), post.ID, post.Version).
		Scan(&post.UpdatedAt, &post.Version)
	return errCustom.HandleStorageError(err)
}

// Delete removes a post from the database by its ID
func (p *PostStore) Delete(ctx context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeoutDuration)
	defer cancel()

	query := `
	DELETE FROM posts
	WHERE id = $1
	`

	result, err := p.db.ExecContext(ctx, query, id)
	if err != nil {
		return errCustom.HandleStorageError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errCustom.HandleStorageError(err)
	}

	if rowsAffected == 0 {
		return errCustom.ErrResourceNotFound
	}

	return nil
}
