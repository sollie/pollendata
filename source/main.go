package main

import (
	"log"
	"net/http"
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

func main() {
	log.Println("Starting server...")
	// Start updateCache loop
	go updateCache()

	// Sleep for 10 seconds to let cache update
	time.Sleep(10 * time.Second)

	// Start webserver
	log.Println("Starting webserver...")
	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.NewCompressor(5, "br, gzip, deflate").Handler)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pollendata"))
	})
	r.Get("/regions", getRegions)
	r.Get("/pollen/{region}", getPollen)
	r.Get("/forecast/{region}", getForecast)

	log.Fatal(http.ListenAndServe(":8080", r))

}
