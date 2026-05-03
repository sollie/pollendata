package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"golang.org/x/net/html"
)

var cacheReadyOnce sync.Once

const (
	refreshInterval = 60 * time.Minute
	fetchTimeout    = 30 * time.Second
	retryBase       = 5 * time.Second
	retryMax        = 5 * time.Minute
)

func updateCache() {
	slog.Info("starting cache update loop")
	ticker := time.NewTicker(refreshInterval)
	defer ticker.Stop()

	fetch := func() {
		retryDelay := retryBase
		for {
			pollendata, err := getRawData()
			if err != nil {
				slog.Error("failed to fetch raw data", "err", err, "retry_in", retryDelay)
				time.Sleep(retryDelay)
				retryDelay = min(retryDelay*2, retryMax)
				continue
			}

			var pd Pollendata
			if err := json.Unmarshal([]byte(pollendata), &pd); err != nil {
				slog.Error("failed to parse pollen data", "err", err, "retry_in", retryDelay)
				time.Sleep(retryDelay)
				retryDelay = min(retryDelay*2, retryMax)
				continue
			}

			lock.Lock()
			cache = pd
			lock.Unlock()

			cacheReadyOnce.Do(func() { close(cacheReady) })
			slog.Info("cache updated")
			return
		}
	}

	fetch()
	for range ticker.C {
		fetch()
	}
}

func getRawData() (string, error) {
	slog.Info("fetching data from pollenvarsel.naaf.no")
	client := &http.Client{Timeout: fetchTimeout}
	req, err := http.NewRequest("GET", "https://pollenvarsel.naaf.no/charts/forecast", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return "", err
	}

	var pollendata string
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "script" {
			for _, a := range n.Attr {
				if a.Key == "id" && a.Val == "__NEXT_DATA__" {
					pollendata = n.FirstChild.Data
					return
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(doc)

	return pollendata, nil
}
