# pollendata

Fetches and caches pollendata from https://pollenvarsel.naaf.no/charts/forecast.

## Endpoints

* /regions - lists available regions
* /pollen/{region} - pollen data for {region}
* /forecast/{region} - text forecast for {region}
* /combined/{region} - combined pollen data and text forecast for {region}

## About the levels

The levels describe the amount of pollen grains per m3 of air.

* 0 - No spread (0)
* 1 - Low spread (1 - 9)
* 2 - Moderate spread (10 - 99)
* 3 - Heavy spread (100 - 999)
* 4 - Extreme spread (1000+)
