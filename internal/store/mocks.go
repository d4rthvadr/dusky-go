package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/d4rthvadr/dusky-go/internal/models"
)

func NewMockStore() Storage {
	return Storage{
		Users: &UserStoreMock{},
	}
}

type UserStoreMock struct {
}

func (m *UserStoreMock) Create(context.Context, *sql.Tx, *models.User) error {
	return nil
}
func (m *UserStoreMock) GetByID(context.Context, int64) (*models.User, error) {
	return nil, nil
}
func (m *UserStoreMock) CreateAndInvite(context.Context, *models.User, string, time.Duration) error {
	return nil
}
func (m *UserStoreMock) ActivateUser(context.Context, string) error {
	return nil
}
func (m *UserStoreMock) GetByEmail(context.Context, string, *models.User) error {
	return nil
}
