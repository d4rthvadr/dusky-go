package store

import (
	"context"
	"database/sql"

	"github.com/d4rthvadr/dusky-go/internal/models"
)

type Storage struct {
	Posts interface {
		Create(context.Context, *models.Post) error
	}
	Users interface {
		Create(context.Context, *models.User) error
	}
}

type PostgresStorage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts: &PostStore{db: db},
		Users: &UserStore{db: db},
	}
}
