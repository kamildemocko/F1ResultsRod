package scrapper

import (
	"log"
	"sync"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

const urlBase = "https://www.formula1.com"

var (
	instance *Scrapper
	once     sync.Once
)

type Scrapper struct {
	Browser *rod.Browser
	BaseUrl string
	url     string
	page    *rod.Page
}

func new() *Scrapper {
	controlURL := launcher.New().Headless(true).Devtools(false).MustLaunch()
	browser := rod.New().Timeout(120 * time.Second).ControlURL(controlURL).MustConnect()
	page := browser.MustPage()

	return &Scrapper{
		Browser: browser,
		BaseUrl: urlBase,
		page:    page,
	}
}

func GetInstance() *Scrapper {
	once.Do(func() {
		instance = new()
	})

	return instance
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

	err := s.page.Timeout(30 * time.Second).Navigate(s.url)
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
