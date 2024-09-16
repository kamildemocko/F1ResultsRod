package main

import (
	"f1-results-rod/data"
	"f1-results-rod/scrapper"
	"flag"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

var (
	LOCAL_RUN = flag.Bool("localRun", false, "run program locally as opposed to running it in Docker")
	THIS_YEAR = time.Now().Year()
	// THIS_YEAR        = time.Now().AddDate(-1, 0, 0).Year() //grab last year
	DSN              string
	ROD_MANAGER_ADDR string
)

type App struct {
	localRun   bool
	scrpr      *scrapper.Scrapper
	repository data.Repository
}

func init() {
	flag.Parse()
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	DSN = os.Getenv("DSN")
	ROD_MANAGER_ADDR = os.Getenv("ROD_MANAGER_ADDR")
}

func main() {
	// set up App
	app := App{
		localRun: *LOCAL_RUN,
	}

	// init db
	repo, err := initPostgresDB(DSN)
	if err != nil {
		panic(err)
	}
	defer repo.Close()

	app.repository = repo

	// init scrapper
	scrpr := scrapper.New(app.localRun, ROD_MANAGER_ADDR)
	defer scrpr.Close()
	app.scrpr = scrpr

	// get already existing tracks
	existingTracks, err := app.repository.GetTracksYear(THIS_YEAR)
	if err != nil {
		panic(err)
	}

	// get results
	res := scrapper.NewResults(app.scrpr)
	results, err := res.GetResults(THIS_YEAR, existingTracks)
	// results, err := res.GetResultsByTrackName(thisYear, "monaco", existingTracks)
	if err != nil {
		panic(err)
	}

	if len(results) == 0 {
		log.Println("no new results")
		return
	}

	// insert into db
	for _, r := range results {
		if len(r.ResultPlaces) == 0 {
			continue
		}

		sameTrackCount, err := app.repository.CountTrack(r.TrackName, THIS_YEAR)
		if err != nil {
			panic(err)
		}

		var trackID int64
		if sameTrackCount > 0 {
			trackID, err = app.repository.GetTrackID(r.TrackName, THIS_YEAR)
			if err != nil {
				panic(err)
			}
		} else {
			trackID, err = app.repository.InsertTrack(r.TrackName, r.Link, THIS_YEAR)
			if err != nil {
				panic(err)
			}
		}

		err = app.repository.InsertResultPlaces(r.ResultPlaces, trackID)
		if err != nil {
			panic(err)
		}
	}
}
