package main

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	b, err := json.Marshal(v)
	if err != nil {
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(b)
}

func isCacheEmpty() bool {
	lock.RLock()
	defer lock.RUnlock()
	return len(cache.Props.PageProps.Data.RegionsData) == 0
}

func regionExists(id string) bool {
	for _, r := range cache.Props.PageProps.Data.RegionsData {
		if r.ID == id {
			return true
		}
	}
	return false
}

func healthz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func readyz(w http.ResponseWriter, r *http.Request) {
	if isCacheEmpty() {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"status": "not ready"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
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
	writeJSON(w, http.StatusOK, levels)
}

func getRegions(w http.ResponseWriter, r *http.Request) {
	if isCacheEmpty() {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "data not yet available"})
		return
	}
	slog.Info("getting regions")
	regions := []string{}

	lock.RLock()
	for _, region := range cache.Props.PageProps.Data.RegionsData {
		regions = append(regions, region.ID)
	}
	lock.RUnlock()

	writeJSON(w, http.StatusOK, regions)
}

func getPollen(w http.ResponseWriter, r *http.Request) {
	if isCacheEmpty() {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "data not yet available"})
		return
	}
	region := chi.URLParam(r, "region")

	lock.RLock()
	exists := regionExists(region)
	if !exists {
		lock.RUnlock()
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "region not found"})
		return
	}
	pollen := make(map[string]Pollen)
	for _, p := range cache.Props.PageProps.Data.ForecastData {
		for _, reg := range p.Regions {
			if reg.ID == region {
				pollen[p.Date] = reg.Pollen
			}
		}
	}
	lock.RUnlock()

	slog.Info("getting pollen", "region", region)
	writeJSON(w, http.StatusOK, pollen)
}

func getForecast(w http.ResponseWriter, r *http.Request) {
	if isCacheEmpty() {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "data not yet available"})
		return
	}
	region := chi.URLParam(r, "region")

	lock.RLock()
	exists := regionExists(region)
	if !exists {
		lock.RUnlock()
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "region not found"})
		return
	}
	var textForecast string
	for _, reg := range cache.Props.PageProps.Data.RegionsData {
		if reg.ID == region {
			textForecast = reg.TextForecast
			break
		}
	}
	lock.RUnlock()

	slog.Info("getting forecast", "region", region)
	writeJSON(w, http.StatusOK, map[string]string{region: textForecast})
}

func getCombined(w http.ResponseWriter, r *http.Request) {
	if isCacheEmpty() {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "data not yet available"})
		return
	}
	region := chi.URLParam(r, "region")

	lock.RLock()
	exists := regionExists(region)
	if !exists {
		lock.RUnlock()
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "region not found"})
		return
	}
	pollen := make(map[string]Pollen)
	for _, p := range cache.Props.PageProps.Data.ForecastData {
		for _, reg := range p.Regions {
			if reg.ID == region {
				pollen[p.Date] = reg.Pollen
			}
		}
	}
	var textForecast string
	for _, reg := range cache.Props.PageProps.Data.RegionsData {
		if reg.ID == region {
			textForecast = reg.TextForecast
			break
		}
	}
	lock.RUnlock()

	slog.Info("getting combined data", "region", region)
	writeJSON(w, http.StatusOK, map[string]any{
		"pollen":   pollen,
		"forecast": map[string]string{region: textForecast},
	})
}
