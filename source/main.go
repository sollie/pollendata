package main

import (
	"log"
	"net/http"
	"net/http/pprof"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	cache        Pollendata
	lock         sync.RWMutex
	writeCacheCh = make(chan Pollendata, 1)
	lastUpdated  = time.Time{}
)

type Server struct {
	Router *chi.Mux
	Routes []string
}

func main() {
	log.Println("Starting server...")
	go updateCache()

	time.Sleep(10 * time.Second)

	log.Println("Starting webserver...")
	s := NewServer()
	s.MountHandlers()

	log.Fatal(http.ListenAndServe(":8080", s.Router))
}

func NewServer() *Server {
	s := &Server{}
	s.Router = chi.NewRouter()
	return s
}

func (s *Server) MountHandlers() {
	s.Router.Use(middleware.RealIP)
	s.Router.Use(middleware.NoCache)
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.NewCompressor(5, "application/json").Handler)
	s.Router.Use(middleware.Recoverer)

	s.Router.Get("/regions", getRegions)
	s.Router.Get("/pollen/{region}", getPollen)
	s.Router.Get("/forecast/{region}", getForecast)
	s.Router.Get("/combined/{region}", getCombined)

	// Add pprof routes
	s.Router.Mount("/debug/pprof", http.HandlerFunc(pprof.Index))
	s.Router.Mount("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	s.Router.Mount("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	s.Router.Mount("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	s.Router.Mount("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))

	// Restrict access to pprof routes to localhost
	s.Router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.RemoteAddr != "127.0.0.1" && r.RemoteAddr != "::1" {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	})
}
