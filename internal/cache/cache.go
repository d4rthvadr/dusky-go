package cache

import (
	"context"
	"time"

	"github.com/d4rthvadr/dusky-go/internal/models"
)

type CacheStorage struct {
	Users interface {
		Get(context.Context, int64) (*models.User, error)
		Set(context.Context, *models.User, time.Duration) error
		Delete(context.Context, int64) error
	}
}

func NewCache(rdb *RedisClient) *CacheStorage {
	return &CacheStorage{
		Users: &UserCache{rdb: rdb},
	}
}
