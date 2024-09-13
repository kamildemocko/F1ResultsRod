package utils

import (
	"net/url"
	"strings"
)

func FixRelativeLinks(links []string, baseURL string) []string {
	linksFixed := make([]string, 0, len(links))

	for _, link := range links {
		fl, _ := url.JoinPath(baseURL, link)
		linksFixed = append(linksFixed, fl)
	}

	return linksFixed
}

func FilterResultRaces(links []string, substrings ...string) []string {
	var foundSubstring bool
	linksFilter := make([]string, 0, len(links))

	for _, link := range links {
		foundSubstring = false

		for _, ss := range substrings {
			if strings.Contains(link, ss) {
				foundSubstring = true
				break
			}
		}

		if !foundSubstring {
			continue
		}

		linksFilter = append(linksFilter, link)
	}

	return linksFilter
}

func RemoveListMatches(listA, listB []string) []string {
	setB := make(map[string]struct{}) // this is effectively set
	for _, v := range listB {
		setB[v] = struct{}{}
	}

	var result []string
	for _, x := range listA {
		if _, found := setB[x]; !found {
			result = append(result, x)
		}
	}

	return result
}
