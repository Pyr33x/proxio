package redis

import (
	"context"
	"fmt"

	"github.com/pyr33x/proxy/pkg/config"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Adapter struct {
	rdb *redis.Client
}

func New(ctx context.Context, cfg *config.Redis, logger *zap.Logger) *Adapter {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Username: cfg.Username,
		Password: cfg.Password,
		DB:       cfg.Database,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		logger.Error("failed to connect to redis", zap.Error(err))
		return nil
	}

	logger.Info("attached redis", zap.String("state", "connected"))
	return &Adapter{rdb: rdb}
}

func (a *Adapter) GetClient() *redis.Client {
	return a.rdb
}
