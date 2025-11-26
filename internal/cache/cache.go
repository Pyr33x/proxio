package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pyr33x/proxy/pkg/err"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Caching interface {
	Get(ctx context.Context, key string) string
	Put(ctx context.Context, key string, value any) error
	Clear(ctx context.Context) error
}

type Cache struct {
	rdb        *redis.Client
	logger     *zap.Logger
	expiration time.Duration
}

type CacheValue struct {
	Status int
	Header http.Header
	Body   []byte
}

func NewCacheRepository(rdb *redis.Client, logger *zap.Logger) *Cache {
	return &Cache{
		rdb:        rdb,
		logger:     logger,
		expiration: 60 * time.Second,
	}
}

func (c *Cache) Get(ctx context.Context, key string) (*CacheValue, bool) {
	if key == "" {
		c.logger.Warn("attempted to get cache with empty key")
		return nil, false
	}

	raw, err := c.rdb.Get(ctx, key).Bytes()
	if err != nil {
		c.logger.Info("cache miss",
			zap.String("key", key),
			zap.String("state", "MISS"),
		)
		return nil, false
	}

	var val CacheValue
	if err := json.Unmarshal(raw, &val); err != nil {
		return nil, false
	}

	return &val, true
}

func (c *Cache) Put(ctx context.Context, key string, value CacheValue) error {
	if key == "" {
		return err.ErrEmptyCacheKey
	}

	b, err := json.Marshal(value)
	if err != nil {
		c.logger.Info("failed to marshal value",
			zap.String("key", key),
			zap.Any("value", value),
			zap.Error(err),
		)
		return err
	}

	err = c.rdb.Set(ctx, key, b, c.expiration).Err()
	if err != nil {
		c.logger.Error("failed to write to cache",
			zap.String("key", key),
			zap.Duration("expiration", c.expiration),
			zap.Error(err),
		)
		return fmt.Errorf("cache put failed: %w", err)
	}

	return nil
}

func (c *Cache) Clear(ctx context.Context) error {
	return c.rdb.FlushAll(ctx).Err()
}
