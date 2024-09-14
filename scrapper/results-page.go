package scrapper

import (
	"f1-results-rod/data"
	"f1-results-rod/utils"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// returns table for specified year
func (r *ScrResults) GetResults(year int, skipTracks []data.Track) ([]ScrResultItem, error) {
	var urlResults = fmt.Sprintf("en/results/%d/races", year)
	var filterResults = fmt.Sprintf("en/results/%d/races/", year)
	var url = strings.Join([]string{r.scrapper.BaseUrl, urlResults}, "/")

	links, err := r.GetResultRaceLinks(url, filterResults)
	if err != nil {
		return nil, err
	}

	links = utils.RemoveListMatches(links, r.convertTracksToList(skipTracks))

	err = r.GetAndSetAllResultTables(links)
	if err != nil {
		return nil, err
	}

	return r.resultItems, nil
}

// returns table for specific year matching the track name
func (r *ScrResults) GetResultsByTrackName(year int, trackName string, skipTracks []data.Track) ([]ScrResultItem, error) {
	var urlResults = fmt.Sprintf("en/results/%d/races", year)
	var filterResults = fmt.Sprintf("en/results/%d/races/", year)
	var url = strings.Join([]string{r.scrapper.BaseUrl, urlResults}, "/")

	links, err := r.GetResultRaceLinks(url, filterResults)
	if err != nil {
		return nil, err
	}

	a := r.convertTracksToList(skipTracks)
	links = utils.RemoveListMatches(links, a)

	var linksTrack []string
	for _, track := range links {
		if strings.Contains(strings.ToLower(track), strings.ToLower(trackName)) {
			linksTrack = append(linksTrack, track)
		}
	}

	err = r.GetAndSetAllResultTables(linksTrack)
	if err != nil {
		return nil, err
	}

	return r.resultItems, nil
}

func (r *ScrResults) convertTracksToList(d []data.Track) []string {
	var output []string

	for _, val := range d {
		output = append(output, val.Link)
	}

	return output
}

// gets all links and filters only result race links
func (r *ScrResults) GetResultRaceLinks(url string, filter string) ([]string, error) {
	r.scrapper.SetUrl(url).Visit()

	allLinks := r.scrapper.GetAllBlockLinks()

	racingLinks := utils.FilterResultRaces(allLinks, filter)
	racingLinks = utils.FixRelativeLinks(racingLinks, r.scrapper.BaseUrl)

	if len(racingLinks) == 0 {
		return []string{}, fmt.Errorf("no results for %s", url)
	}

	return racingLinks, nil
}

func (r *ScrResults) GetAndSetAllResultTables(allLinks []string) error {
	var results []ScrResultItem

	for i, link := range allLinks {
		r.scrapper.SetUrl(link).Visit()

		title, err := r.parseResultsTitle()
		if err != nil {
			return err
		}

		table, _ := r.parseResultsTable()

		resultItem := ScrResultItem{
			Position:     i,
			TrackName:    title,
			Link:         link,
			ResultPlaces: table,
		}

		results = append(results, resultItem)
	}

	r.resultItems = results

	return nil
}

func (r *ScrResults) parseResultsTitle() (string, error) {
	title, err := r.scrapper.page.Element("h1.f1-heading")
	if err != nil {
		return "", err
	}

	txt := title.MustText()
	txt = strings.Replace(txt, " - RACE RESULT", "", -1)

	return txt, nil
}

func (r *ScrResults) parseResultsTable() ([]data.ResultPlace, error) {
	table, err := r.scrapper.page.Timeout(12 * time.Second).Element("table.f1-table")
	if err != nil {
		return nil, fmt.Errorf("timeout parsing results table")
	}

	log.Println("..parsing table")

	var resultRows []data.ResultPlace
	rows := table.MustElements("tbody > tr")

	for _, row := range rows {
		columns := row.MustElements("td")

		position, _ := strconv.Atoi(columns[0].MustText())
		driverNo, _ := strconv.Atoi(columns[1].MustText())
		laps, _ := strconv.Atoi(columns[4].MustText())
		points, _ := strconv.Atoi(columns[6].MustText())

		resultRow := data.ResultPlace{
			Position:      position,
			DriverNo:      driverNo,
			Driver:        columns[2].MustText(),
			Car:           columns[3].MustText(),
			Laps:          laps,
			TimeOrRetired: columns[5].MustText(),
			Points:        points,
		}

		resultRows = append(resultRows, resultRow)
	}

	return resultRows, nil
}
