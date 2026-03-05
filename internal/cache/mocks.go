package cache

import (
	"context"
	"time"

	"github.com/d4rthvadr/dusky-go/internal/models"
)

func NewMockCache() CacheStorage {
	return CacheStorage{
		Users: &UserCacheMock{},
	}
}

type UserCacheMock struct {
}

func (m *UserCacheMock) Get(context.Context, int64) (*models.User, error) {
	return nil, nil
}
func (m *UserCacheMock) Set(context.Context, *models.User, time.Duration) error {
	return nil
}
func (m *UserCacheMock) Delete(context.Context, int64) error {
	return nil
}
