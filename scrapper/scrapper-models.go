package scrapper

import (
	"f1-results-rod/data"
	"fmt"
	"strings"
)

type ScrResults struct {
	scrapper    *Scrapper
	resultItems []ScrResultItem
}

func NewResults(scrapper *Scrapper) *ScrResults {
	return &ScrResults{
		scrapper:    scrapper,
		resultItems: []ScrResultItem{},
	}
}

func (r *ScrResults) String() string {
	var sb strings.Builder

	for _, item := range r.resultItems {
		fmt.Fprintln(&sb, item.TrackName+":")

		if item.ResultPlaces == nil {
			fmt.Fprint(&sb, "\t(no data)\n")
			continue
		}

		for _, value := range item.ResultPlaces {
			fmt.Fprintf(&sb, "\t%+v\n", value)
		}
	}

	return sb.String()
}

type ScrResultItem struct {
	Position     int
	TrackName    string
	Link         string
	ResultPlaces []data.ResultPlace
}
