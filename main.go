package main

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	zl "github.com/pyr33x/proxio/internal/adapter/zap"
	"github.com/pyr33x/proxio/internal/proxy"
	"github.com/pyr33x/proxio/pkg/config"
	"go.uber.org/zap"

	_ "github.com/joho/godotenv/autoload"
)

func gracefulShutdown(logger *zap.Logger, proxy *http.Server, done chan struct{}) {
	stopCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-stopCtx.Done()
	logger.Info("shutdown triggered...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := proxy.Shutdown(shutdownCtx); err != nil {
		logger.Error("failed to shutdown proxy server",
			zap.Error(err),
		)
		return
	}

	logger.Info("exiting proxy...")
	done <- struct{}{}
}

func main() {
	done := make(chan struct{})

	cfg := config.New()
	logger := zl.New(cfg.Zap.Environment).GetLogger()
	proxy := proxy.NewProxyServer(context.Background(), cfg, logger)

	go gracefulShutdown(logger, proxy, done)

	logger.Info("attached proxy",
		zap.String("bind", proxy.Addr),
	)

	if err := proxy.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("failed to listen and serve to the proxy server",
			zap.Error(err),
		)
		close(done)
		return
	}

	<-done
	logger.Info("shutdown complete")
}
