package main

import (
	"f1-results-rod/data"
	"f1-results-rod/scrapper"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

var thisYear = time.Now().Year()
var DSN string

type App struct {
	scrapper *scrapper.Scrapper
	db       data.DBIntf
}

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	DSN = os.Getenv("DSN")
}

func main() {
	app := App{}

	// init db
	db, err := data.NewPostgresDB(DSN)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	app.db = db

	// get already existing tracks
	existingTracks, err := app.db.GetTracksYear(thisYear)
	if err != nil {
		panic(err)
	}

	// init scrapper
	app.scrapper = scrapper.GetInstance()
	defer app.scrapper.Close()

	// get results
	res := scrapper.NewResults(app.scrapper)
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
