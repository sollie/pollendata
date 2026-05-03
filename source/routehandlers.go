package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func isCacheEmpty() bool {
	lock.RLock()
	defer lock.RUnlock()
	return len(cache.Props.PageProps.Data.RegionsData) == 0
}

func cacheNotReady(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusServiceUnavailable)
	w.Write([]byte(`{"error":"data not yet available"}`))
}

func getLevels(w http.ResponseWriter, r *http.Request) {
	levels := []struct {
		Level       int    `json:"level"`
		Label       string `json:"label"`
		GrainsPerM3 string `json:"grains_per_m3"`
		Description string `json:"description"`
	}{
		{0, "No spread", "0", "No pollen in the air"},
		{1, "Low spread", "1-9", "Unlikely to cause symptoms"},
		{2, "Moderate spread", "10-99", "May cause symptoms in sensitive individuals"},
		{3, "Heavy spread", "100-999", "Likely to cause symptoms in most allergy sufferers"},
		{4, "Extreme spread", "1000+", "Severe symptoms expected"},
	}
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(levels)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}

func getRegions(w http.ResponseWriter, r *http.Request) {
	if isCacheEmpty() {
		cacheNotReady(w)
		return
	}
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
	if isCacheEmpty() {
		cacheNotReady(w)
		return
	}
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
	if isCacheEmpty() {
		cacheNotReady(w)
		return
	}
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

func getCombined(w http.ResponseWriter, r *http.Request) {
	if isCacheEmpty() {
		cacheNotReady(w)
		return
	}
	region := chi.URLParam(r, "region")
	log.Printf("Getting combined data for %s...\n", region)
	var pollen = make(map[string]Pollen)
	var forecast = make(map[string]string)

	lock.RLock()
	for _, p := range cache.Props.PageProps.Data.ForecastData {
		for _, r := range p.Regions {
			if r.ID == region {
				pollen[p.Date] = r.Pollen
			}
		}
	}
	for _, r := range cache.Props.PageProps.Data.RegionsData {
		if r.ID == region {
			forecast[r.ID] = r.TextForecast
		}
	}
	lock.RUnlock()

	combined := map[string]any{
		"pollen":   pollen,
		"forecast": forecast,
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(combined)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}
