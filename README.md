# pollendata

Fetches and caches pollendata from https://pollenvarsel.naaf.no/charts/forecast.

## Running with Docker

```sh
docker build -t pollendata .
docker run -d -p 8080:8080 pollendata
```

Service listens on port 8080. On startup, the server waits for an initial data fetch from upstream before accepting requests. This typically takes a few seconds. Until data is ready, requests to data endpoints return `503 Service Unavailable`.

## Endpoints

* /regions - lists available regions
* /levels - pollen severity levels with plain-language descriptions and grain count ranges
* /pollen/{region} - pollen data for {region}
* /forecast/{region} - text forecast for {region}
* /combined/{region} - combined pollen data and text forecast for {region}

## About the levels

The levels describe the amount of pollen grains per m³ of air. See `/levels` for the machine-readable version.

| Level | Label | Grains/m³ | Description |
|-------|-------|-----------|-------------|
| 0 | No spread | 0 | No pollen in the air |
| 1 | Low spread | 1–9 | Unlikely to cause symptoms |
| 2 | Moderate spread | 10–99 | May cause symptoms in sensitive individuals |
| 3 | Heavy spread | 100–999 | Likely to cause symptoms in most allergy sufferers |
| 4 | Extreme spread | 1000+ | Severe symptoms expected |
