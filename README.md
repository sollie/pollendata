# pollendata

Fetch and cache pollendata from https://pollenvarsel.naaf.no/charts/forecast.
Expose as rest api.

## Endpoints

* /regions - lists available regions
* /pollen/{region} - pollen data for {region}
* /forecast/{region} - text forecast for {region}
