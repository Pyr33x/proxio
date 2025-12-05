package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisStore struct {
	client *redis.Client
}

func NewRedisStore(client *redis.Client) *redisStore {
	return &redisStore{client: client}
}

func (r *redisStore) Get(ctx context.Context, key string) ([]byte, error) {
	return r.client.Get(ctx, key).Bytes()
}

func (r *redisStore) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *redisStore) Clear(ctx context.Context) error {
	return r.client.FlushDB(ctx).Err()
}
