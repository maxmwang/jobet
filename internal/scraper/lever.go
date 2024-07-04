package scraper

import (
	"fmt"
	"net/http"
	"time"
)

type leverScraper struct {
	name string
}

func newLeverScraper() Scraper {
	return leverScraper{
		name: SiteLever,
	}
}

func (s leverScraper) Scrape(companyName, site string) ([]Job, error) {
	if site != s.name {
		return nil, nil
	}
	return s.ScrapeAll(companyName)
}

func (s leverScraper) ScrapeAll(companyName string) ([]Job, error) {
	res, err := http.Get(fmt.Sprintf("https://api.lever.co/v0/postings/%s?limit=999", companyName))

	body, err := checkThenDecode[[]struct {
		Title     string `json:"text"`
		UpdatedAt int    `json:"createdAt"`
	}](companyName, s.name, res, err)
	if err != nil {
		return nil, err
	}

	jobs := make([]Job, 0)
	for _, j := range body {
		jobs = append(jobs, Job{
			Company:   companyName,
			Title:     j.Title,
			UpdatedAt: time.UnixMilli(int64(j.UpdatedAt)),
		})
	}

	return jobs, nil
}

func (s leverScraper) Sites() []string {
	return []string{
		s.name,
	}
}
