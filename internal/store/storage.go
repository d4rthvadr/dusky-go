package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/d4rthvadr/dusky-go/internal/models"
)

const defaultQueryTimeoutDuration = 5 * time.Second

type Storage struct {
	Posts interface {
		Create(context.Context, *models.Post) error
		GetByID(context.Context, int64) (*models.Post, error)
		Update(context.Context, *models.Post) error
		Delete(context.Context, int64) error
		GetUserFeed(context.Context, int64, *PaginatedFeedQuery) ([]*PostWithMetadata, error)
	}
	Comments interface {
		GetByPostID(context.Context, int64) ([]models.Comment, error)
		Create(context.Context, *models.Comment) error
	}
	Users interface {
		Create(context.Context, *models.User) error
		GetByID(context.Context, int64) (*models.User, error)
	}
	Followers interface {
		Follow(context.Context, int64, int64) error
		Unfollow(context.Context, int64, int64) error
	}
}

type PostgresStorage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:     &PostStore{db: db},
		Comments:  &CommentStore{db: db},
		Users:     &UserStore{db: db},
		Followers: &FollowerStore{db: db},
	}
}
