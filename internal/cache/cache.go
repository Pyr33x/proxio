package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pyr33x/proxio/pkg/err"
	"go.uber.org/zap"
)

type Store interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Clear(ctx context.Context) error
}

type Cache struct {
	store      Store
	logger     *zap.Logger
	expiration time.Duration
}

type CacheValue struct {
	Status int
	Header http.Header
	Body   []byte
}

func NewCacheRepository(store Store, logger *zap.Logger, ttl time.Duration) *Cache {
	return &Cache{
		store:      store,
		logger:     logger,
		expiration: ttl,
	}
}

func (c *Cache) Get(ctx context.Context, key string) (*CacheValue, bool) {
	if key == "" {
		c.logger.Warn("attempted to get cache with empty key")
		return nil, false
	}

	raw, err := c.store.Get(ctx, key)
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

	if err := c.store.Set(ctx, key, b, c.expiration); err != nil {
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
	return c.store.Clear(ctx)
}
