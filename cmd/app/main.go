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
	localRun         = flag.Bool("localRun", false, "run program locally as opposed to running it in Docker")
	thisYear         = time.Now().Year()
	DSN              string
	ROD_MANAGER_ADDR string
)

type App struct {
	localRun bool
	scrpr    *scrapper.Scrapper
	db       data.DBIntf
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
		localRun: *localRun,
	}

	// init db
	db, err := data.NewPostgresDB(DSN)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	app.db = db

	// init scrapper
	scrpr := scrapper.New(app.localRun, ROD_MANAGER_ADDR)
	defer scrpr.Close()
	app.scrpr = scrpr

	// get already existing tracks
	existingTracks, err := app.db.GetTracksYear(thisYear)
	if err != nil {
		panic(err)
	}

	// get results
	res := scrapper.NewResults(app.scrpr)
	results, err := res.GetResults(thisYear, existingTracks)
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
		sameTrackCount, err := app.db.CountTrack(r.TrackName, thisYear)
		if err != nil {
			panic(err)
		}

		var trackID int64
		if sameTrackCount > 0 {
			trackID, err = app.db.GetTrackID(r.TrackName, thisYear)
			if err != nil {
				panic(err)
			}
		} else {
			trackID, err = app.db.InsertTrack(r.TrackName, r.Link, thisYear)
			if err != nil {
				panic(err)
			}
		}

		err = app.db.InsertResultPlaces(r.ResultPlaces, trackID)
		if err != nil {
			panic(err)
		}
	}
}
