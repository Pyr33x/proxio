package proxy

import (
	"context"
	"net/http"
	"time"

	"github.com/pyr33x/proxio/internal/adapter/redis"
	"github.com/pyr33x/proxio/internal/cache"
	"github.com/pyr33x/proxio/pkg/config"
	"go.uber.org/zap"
)

type Server struct {
	Proxy  ProxyServer
	Origin OriginServer
	Cache  *cache.Cache
	logger *zap.Logger
}

type ProxyServer struct {
	Port string
}

type OriginServer struct {
	URL string
}

func NewProxyServer(ctx context.Context, cfg *config.Config, logger *zap.Logger) *http.Server {
	rdb := redis.New(ctx, &cfg.Redis, logger).GetClient()
	rdbCache := cache.NewRedisStore(rdb)

	srv := &Server{
		Proxy: ProxyServer{
			Port: cfg.Server.Proxy.Port,
		},
		Origin: OriginServer{
			URL: cfg.Server.Origin.URL,
		},
		Cache:  cache.NewCacheRepository(rdbCache, logger, 60*time.Second),
		logger: logger,
	}

	return &http.Server{
		Addr:         ":" + srv.Proxy.Port,
		Handler:      srv.Serve(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}
