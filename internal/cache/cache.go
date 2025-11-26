package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Cache interface {
	Get(ctx context.Context, key string) string
	Put(ctx context.Context, key string, value any) error
}

type cache struct {
	rdb        *redis.Client
	sugar      *zap.SugaredLogger
	expiration time.Duration
}

func NewCacheRepository(rdb *redis.Client, logger *zap.Logger) *cache {
	return &cache{
		rdb:        rdb,
		sugar:      logger.Sugar(),
		expiration: 60 * time.Second,
	}
}

func (c *cache) Get(ctx context.Context, key string) (string, bool) {
	if key == "" {
		c.sugar.Warn("attempted to get cache with empty key")
		return "", false
	}

	res, err := c.rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			c.sugar.Warn("cache miss",
				"key", key,
			)
			return "", false
		}

		c.sugar.Error("failed to read from cache",
			"key", key,
			"error", err,
		)
		return "", false
	}

	return res, true
}

func (c *cache) Put(ctx context.Context, key string, value any) error {
	if key == "" {
		return errors.New("cache key cannot be empty")
	}

	err := c.rdb.Set(ctx, key, value, c.expiration).Err()
	if err != nil {
		c.sugar.Error("failed to write to cache",
			"key", key,
			"expiration", c.expiration,
			"error", err,
		)
		return fmt.Errorf("cache put failed: %w", err)
	}

	return nil
}
