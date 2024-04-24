package main

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
	ID           string `json:"id"`
	Pollen       Pollen `json:"pollen"`
	TextForecast string `json:"forecast,-"`
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
