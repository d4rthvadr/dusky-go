package store

import (
	"context"
	"database/sql"

	"github.com/d4rthvadr/dusky-go/internal/models"
)

type CommentStore struct {
	db *sql.DB
}

func (c *CommentStore) GetByPostID(ctx context.Context, postID int64) ([]models.Comment, error) {
	query := `
	SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, users.username
	FROM comments c left join users on c.user_id = users.id
	WHERE c.post_id = $1
	ORDER BY c.created_at ASC
	`

	rows, err := c.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment

	for rows.Next() {
		var comment models.Comment
		if err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreatedAt, &comment.User.Username); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func (c *CommentStore) Create(ctx context.Context, comment *models.Comment) error {
	query := `
	INSERT INTO comments (post_id, user_id, content) 
	VALUES ($1, $2, $3) RETURNING id, created_at
	`
	err := c.db.QueryRowContext(ctx, query, comment.PostID, comment.UserID, comment.Content).
		Scan(&comment.ID, &comment.CreatedAt)

	return err
}
