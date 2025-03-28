package main

import (
	"log"
	"net/http"
	"net/http/pprof"
	"strings"
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

	pprofRouter := chi.NewRouter()
	pprofRouter.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			host := strings.Split(r.RemoteAddr, ":")[0]
			if host != "127.0.0.1" && host != "::1" && host != "localhost" {
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
