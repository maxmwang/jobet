package scrape

import (
	"fmt"
	"net/http"

	"github.com/maxmwang/jobet/internal/db"
	"github.com/maxmwang/jobet/internal/proto"
)

type leverScraper struct {
	name string
}

func newLeverScraper() Scraper {
	return leverScraper{
		name: string(db.SiteLever),
	}
}

func (s leverScraper) Scrape(companyName, site string) ([]*proto.ScrapeBatch_Job, error) {
	if site != s.name {
		return nil, nil
	}
	return s.ScrapeAll(companyName)
}

func (s leverScraper) ScrapeAll(companyName string) ([]*proto.ScrapeBatch_Job, error) {
	res, err := http.Get(fmt.Sprintf("https://api.lever.co/v0/postings/%s?limit=999", companyName))

	body, err := checkThenDecode[[]struct {
		Title     string `json:"text"`
		UpdatedAt int    `json:"createdAt"`
	}](companyName, s.name, res, err)
	if err != nil {
		return nil, err
	}

	jobs := make([]*proto.ScrapeBatch_Job, 0)
	for _, j := range body {
		jobs = append(jobs, &proto.ScrapeBatch_Job{
			Company:   companyName,
			Title:     j.Title,
			UpdatedAt: int64(j.UpdatedAt),
		})
	}

	return jobs, nil
}

func (s leverScraper) Sites() []string {
	return []string{
		s.name,
	}
}
