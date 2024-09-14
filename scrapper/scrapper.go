package scrapper

import (
	"log"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

const urlBase = "https://www.formula1.com"

type Scrapper struct {
	Browser *rod.Browser
	BaseUrl string
	url     string
	page    *rod.Page
}

// returns new Scrapper, needs local flag and rod manager addr, that can be empty if local
func New(local bool, rodManagerAddr string) *Scrapper {
	var browser *rod.Browser

	log.Println("starting up browser, local mode:", local)

	if local {
		controlURL := launcher.New().Headless(true).Devtools(false).MustLaunch()
		browser = rod.New().Timeout(240 * time.Second).ControlURL(controlURL).MustConnect()
	} else {
		controlURL := launcher.MustNewManaged(rodManagerAddr).MustClient()
		browser = rod.New().Timeout(240 * time.Second).Client(controlURL).MustConnect()
	}
	page := browser.MustPage()

	return &Scrapper{
		Browser: browser,
		BaseUrl: urlBase,
		page:    page,
	}
}

func (s *Scrapper) Close() {
	log.Println("closing browser")

	s.Browser.Close()
}

func (s *Scrapper) SetUrl(url string) *Scrapper {
	s.url = url

	return s
}

func (s *Scrapper) Visit() *Scrapper {
	log.Println("visting url:", s.url)

	err := s.page.Timeout(60 * time.Second).Navigate(s.url)
	if err != nil {
		panic(err)
	}

	s.page = s.page.MustWaitLoad()

	return s
}

// Gets all links from current page.
func (s *Scrapper) GetAllBlockLinks() []string {
	log.Println("..getting links from url", s.url)

	elements := s.page.MustElements("a.block")
	links := make([]string, 0, len(elements))

	for _, element := range elements {
		href, err := element.Attribute("href")
		if err != nil || href == nil {
			continue
		}

		links = append(links, *href)
	}

	return links
}
