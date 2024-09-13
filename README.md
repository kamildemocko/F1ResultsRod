# F1 Results Scraper

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

## Local

### Installation

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

### Running

```
go run ./cmd/app
```
