package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/d4rthvadr/dusky-go/internal/models"
	"github.com/redis/go-redis/v9"
)

type UserCache struct {
	rdb *RedisClient
}

func getUserCacheKey(userID int64) string {
	return fmt.Sprintf("user-%v", userID)
}

func (c *UserCache) Get(ctx context.Context, userID int64) (*models.User, error) {
	cacheKey := getUserCacheKey(userID)

	data, err := c.rdb.Get(ctx, cacheKey)

	if err == redis.Nil {
		// cache miss, return nil without error
		return nil, redis.Nil
	} else if err != nil {
		// other error while fetching from cache
		return nil, err
	}

	dataAsStr, err := c.getDataString(ctx, data, userID)
	if err != nil {
		// delete the cache if the data is not valid
		_ = c.Delete(ctx, userID)
		return nil, err
	}

	var user models.User
	if err = json.Unmarshal([]byte(dataAsStr), &user); err != nil {
		// delete the cache if the data is not valid or corrupted
		_ = c.Delete(ctx, userID)
		return nil, redis.Nil
	}

	return &user, nil
}

func (c *UserCache) getDataString(ctx context.Context, data interface{}, userID int64) (string, error) {
	// if the data is not a string, delete the cache and return nil
	dataStr, ok := data.(string)
	if !ok {
		return "", redis.Nil
	}

	dataStr = strings.TrimSpace(dataStr)
	if dataStr == "" || dataStr == "null" {
		return "", redis.Nil
	}

	return dataStr, nil

}

func (c *UserCache) Set(ctx context.Context, user *models.User, exp time.Duration) error {

	if user == nil {
		return errors.New("user cannot be nil")
	}

	cacheKey := getUserCacheKey(user.ID)

	// TODO: we need to store only necessary fields in the cache, we don't want to store the password hash and other sensitive information. We can create a separate struct for the cache and map the user model to the cache model before storing it in the cache. For simplicity, we are storing the entire user model in the cache for now.
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return c.rdb.Set(ctx, cacheKey, data, exp)
}

func (c *UserCache) Delete(ctx context.Context, userID int64) error {
	cacheKey := getUserCacheKey(userID)
	return c.rdb.Del(ctx, cacheKey)
}
