package proxy

import (
	"io"
	"maps"
	"net/http"

	"github.com/pyr33x/proxy/internal/cache"
	"go.uber.org/zap"
)

func (srv *Server) Serve() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", srv.ServeProxy)
	return mux
}

func (srv *Server) ServeProxy(w http.ResponseWriter, r *http.Request) {
	cacheKey := r.Method + ":" + r.URL.String()

	c, ok := srv.Cache.Get(r.Context(), cacheKey)
	if ok {
		srv.WriteHeaders(w, "HIT", c)
		srv.logger.Info("cache hit",
			zap.String("key", cacheKey),
			zap.String("state", "HIT"),
		)
		return
	}

	originURL := srv.Origin.URL + r.URL.String()
	resp, err := http.Get(originURL)
	if err != nil {
		srv.logger.Error("error forwarding request",
			zap.Error(err),
			zap.String("origin", originURL),
		)
		http.Error(w, "error forwarding request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		srv.logger.Error("error reading response body",
			zap.Error(err),
			zap.String("origin", originURL),
		)
		http.Error(w, "error reading response body", http.StatusInternalServerError)
		return
	}

	cached := cache.CacheValue{
		Status: resp.StatusCode,
		Header: resp.Header.Clone(),
		Body:   body,
	}

	if err := srv.Cache.Put(r.Context(), cacheKey, cached); err != nil {
		srv.logger.Error("failed to write to cache",
			zap.Error(err),
			zap.String("key", cacheKey),
		)
	}
	srv.WriteHeaders(w, "MISS", &cached)
}

func (srv *Server) WriteHeaders(w http.ResponseWriter, state string, cached *cache.CacheValue) {
	maps.Copy(w.Header(), cached.Header)
	w.Header().Set("X-Cache", state)
	w.WriteHeader(cached.Status)
	if _, err := w.Write(cached.Body); err != nil {
		srv.logger.Error("failed to write cached body",
			zap.Error(err),
		)
		return
	}
}
