package main

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	cache      Pollendata
	lock       sync.RWMutex
	cacheReady = make(chan struct{})
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	slog.Info("starting server")
	go updateCache()

	slog.Info("waiting for initial data fetch")
	<-cacheReady

	slog.Info("starting webserver")
	s := NewServer()
	s.MountHandlers()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: s.Router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "err", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("graceful shutdown failed", "err", err)
		os.Exit(1)
	}
	slog.Info("server stopped")
}

func NewServer() *Server {
	return &Server{Router: chi.NewRouter()}
}

type Server struct {
	Router *chi.Mux
}

func (s *Server) MountHandlers() {
	s.Router.Use(middleware.RealIP)
	s.Router.Use(middleware.NoCache)
	s.Router.Use(requestLogger)
	s.Router.Use(middleware.NewCompressor(5, "application/json").Handler)
	s.Router.Use(middleware.Recoverer)

	s.Router.Get("/healthz", healthz)
	s.Router.Get("/readyz", readyz)
	s.Router.Get("/regions", getRegions)
	s.Router.Get("/levels", getLevels)
	s.Router.Get("/pollen/{region}", getPollen)
	s.Router.Get("/forecast/{region}", getForecast)
	s.Router.Get("/combined/{region}", getCombined)

	pprofRouter := chi.NewRouter()
	pprofRouter.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			host, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil || (host != "127.0.0.1" && host != "::1") {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	})

	pprofRouter.HandleFunc("/", pprof.Index)
	pprofRouter.HandleFunc("/cmdline", pprof.Cmdline)
	pprofRouter.HandleFunc("/profile", pprof.Profile)
	pprofRouter.HandleFunc("/symbol", pprof.Symbol)
	pprofRouter.HandleFunc("/trace", pprof.Trace)

	s.Router.Mount("/debug/pprof", pprofRouter)
}

func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		start := time.Now()
		next.ServeHTTP(ww, r)
		slog.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", ww.Status(),
			"duration_ms", time.Since(start).Milliseconds(),
			"remote_addr", r.RemoteAddr,
		)
	})
}
