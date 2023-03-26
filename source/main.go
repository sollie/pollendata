package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/html"
)

type Pollendata struct {
	Props        Props         `json:"props"`
	Page         string        `json:"-"`
	Query        int           `json:"-"`
	BuildID      string        `json:"-"`
	IsFallback   bool          `json:"-"`
	Gssp         bool          `json:"-"`
	ScriptLoader []interface{} `json:"-"`
}
type Pollen struct {
	Bjork  int `json:"bjork"`
	Burot  int `json:"burot"`
	Gress  int `json:"gress"`
	Hassel int `json:"hassel"`
	Or     int `json:"or"`
	Salix  int `json:"salix"`
}
type Regions struct {
	ID     string `json:"id"`
	Pollen Pollen `json:"pollen"`
}
type ForecastData struct {
	Date    string    `json:"date"`
	Regions []Regions `json:"regions"`
}
type RegionsData struct {
	ID           string `json:"id"`
	TextForecast string `json:"textForecast"`
}
type Data struct {
	ForecastData []ForecastData `json:"forecastData"`
	RegionsData  []RegionsData  `json:"regionsData"`
}
type PageProps struct {
	Data Data `json:"data"`
}
type Props struct {
	PageProps PageProps `json:"pageProps"`
	NSsp      bool      `json:"-"`
}

func main() {
	// Read https://pollenvarsel.naaf.no/charts/forecast
	resp, err := http.Get("https://pollenvarsel.naaf.no/charts/forecast")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Fatal(err)
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

	// Unmarshal pollendata to Pollendata struct
	var pd Pollendata
	err = json.Unmarshal([]byte(pollendata), &pd)
	if err != nil {
		log.Fatal(err)
	}

	// Indent json
	b, err := json.MarshalIndent(pd.Props.PageProps.Data, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	// Print json
	fmt.Println(string(b))
}
