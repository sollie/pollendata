package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func getRegions(w http.ResponseWriter, r *http.Request) {
	log.Println("Getting regions...")
	regions := []string{}

	lock.RLock()
	for _, region := range cache.Props.PageProps.Data.RegionsData {
		regions = append(regions, region.ID)
	}
	lock.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(regions)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}

func getPollen(w http.ResponseWriter, r *http.Request) {
	region := chi.URLParam(r, "region")
	log.Printf("Getting pollen for %s...\n", region)
	var pollen = make(map[string]Pollen)

	lock.RLock()
	for _, p := range cache.Props.PageProps.Data.ForecastData {
		for _, r := range p.Regions {
			if r.ID == region {
				pollen[p.Date] = r.Pollen
			}
		}
	}
	lock.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(pollen)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}

func getForecast(w http.ResponseWriter, r *http.Request) {
	region := chi.URLParam(r, "region")
	log.Printf("Getting forecast for %s...\n", region)
	var forecast = make(map[string]string)

	lock.RLock()
	for _, r := range cache.Props.PageProps.Data.RegionsData {
		if r.ID == region {
			forecast[r.ID] = r.TextForecast
		}
	}
	lock.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(forecast)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}
