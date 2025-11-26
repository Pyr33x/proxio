package proxy

import (
	"net/http"
	"time"

	"github.com/pyr33x/proxy/internal/cache"
	"github.com/pyr33x/proxy/pkg/config"
	"github.com/redis/go-redis/v9"
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

func NewProxyServer(cfg *config.Server, rdb *redis.Client, logger *zap.Logger) *http.Server {
	srv := &Server{
		Proxy: ProxyServer{
			Port: cfg.Proxy.Port,
		},
		Origin: OriginServer{
			URL: cfg.Origin.URL,
		},
		Cache:  cache.NewCacheRepository(rdb, logger),
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
