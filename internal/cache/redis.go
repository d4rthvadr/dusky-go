package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	rdb *redis.Client
}

type RedisOptions struct {
	Addr     string
	Password string
	DB       int
}

// NewRedisClient initializes a new Redis client with the provided options.
func NewRedisClient(options *RedisOptions) *RedisClient {

	rdb := redis.NewClient(&redis.Options{
		Addr:     options.Addr,
		Password: options.Password,
		DB:       options.DB,
	})
	return &RedisClient{rdb: rdb}
}

func (r *RedisClient) Ping(ctx context.Context) error {
	return r.rdb.Ping(ctx).Err()
}

func (r *RedisClient) Close() error {
	return r.rdb.Close()
}

func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, exp time.Duration) error {

	err := r.rdb.Set(ctx, key, value, exp).Err()

	return err

}

func (r *RedisClient) Get(ctx context.Context, key string) (interface{}, error) {

	val, err := r.rdb.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return val, nil
}

func (r *RedisClient) Del(ctx context.Context, key string) error {

	err := r.rdb.Del(ctx, key).Err()
	return err
}
