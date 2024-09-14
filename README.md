# F1 Results Rod Scraper

This project is a Formula 1 results scraper built with Go. It uses the Rod library for web scraping and stores the data in a PostgreSQL database.

## Features

- Scrapes F1 race results from the official Formula 1 website
- Stores race results in a PostgreSQL database
- Supports scraping results for a specific year
- Avoids duplicating data for tracks that have already been scraped

## Prerequisites

- Go
- PostgreSQL database
- Docker (optional)
- Rod-Manager (optional, docker-run link below)

## Installation

1. Clone the repository:
```
git clone https://github.com/kamildemocko/F1ResultsRod.git
cd F1ResultsRod
```

2. Install dependencies:
```
go mod tidy
```

3. Set up the `.env` file:
```
DSN=host=192.168.92.241 port=5432 user=postgres password=SECRET timezone=UTC connect_timeout=5 search_path=f1scrap
```

## Running

### Remote

We need a **network**:

```powershell
sudo docker network create f1-result-rod-net
```

**Rod-manager** must be running, easiest way to run it is with docker:

```Dockerfile
docker run -p 7317:7317 --network f1-result-rod-net ghcr.io/go-rod/rod
```

-  the IP of the server running the manager must be specified in _scrapper.go_'s **new** function

Dockerfile will build image, and leave you with image under 50MB

```powershell
docker build -t f1-results-rod .
```

assuming the Postgres database is on the same network, we can run our built image:

```powershell
docker run --network f1-result-rod-net f1-results-rod
```

### Local

As default is to run in Docker you must uncomment this lines (and comment coresponding lines for Docker) in _scrapper.go_'s **new** function
```
controlURL := launcher.New().Headless(true).Devtools(false).MustLaunch()
browser := rod.New().Timeout(120 * time.Second).ControlURL(controlURL).MustConnect()
```

```
go run ./cmd/app
```
