package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/html"
)

func updateCache() {
	log.Println("Starting cache update loop...")
	for {
		if time.Since(lastUpdated) < 60*time.Minute {
			continue
		}

		pollendata, err := getRawData()
		if err != nil {
			log.Fatal(err)
		}

		var pd Pollendata
		err = json.Unmarshal([]byte(pollendata), &pd)
		if err != nil {
			log.Fatal(err)
		}

		lock.Lock()
		cache = pd
		lastUpdated = time.Now()
		lock.Unlock()

		// Todo?
		// writeCacheCh <- pd
		time.Sleep(10 * time.Second)
	}
}

func getRawData() (string, error) {
	log.Println("Getting raw data from pollenvarsel.naaf.no...")
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://pollenvarsel.naaf.no/charts/forecast", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("UserAgent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return "", err
	}

	var pollendata string
	var data func(*html.Node)
	data = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "script" {
			for _, a := range n.Attr {
				if a.Key == "id" && a.Val == "__NEXT_DATA__" {
					pollendata = n.FirstChild.Data
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			data(c)
		}
	}
	data(doc)

	return pollendata, nil
}
