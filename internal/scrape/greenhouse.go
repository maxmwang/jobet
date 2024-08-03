package scrape

import (
	"fmt"
	"net/http"
	"time"

	"github.com/maxmwang/jobet/internal/db"
	"github.com/maxmwang/jobet/internal/proto"
)

type greenhouseScraper struct {
	name string
}

func newGreenhouseScraper() Scraper {
	return greenhouseScraper{
		name: string(db.SiteGreenhouse),
	}
}

func (s greenhouseScraper) Scrape(companyName, site string) ([]*proto.ScrapeBatch_Job, error) {
	if site != s.name {
		return nil, nil
	}
	return s.ScrapeAll(companyName)
}

func (s greenhouseScraper) ScrapeAll(companyName string) ([]*proto.ScrapeBatch_Job, error) {
	res, err := http.Get(fmt.Sprintf("https://api.greenhouse.io/v1/boards/%s/jobs", companyName))

	body, err := checkThenDecode[struct {
		Jobs []struct {
			Title     string `json:"title"`
			UpdatedAt string `json:"updated_at"`
		} `json:"jobs"`
	}](companyName, s.name, res, err)
	if err != nil {
		return nil, err
	}

	jobs := make([]*proto.ScrapeBatch_Job, 0)
	for _, resJob := range body.Jobs {
		j := &proto.ScrapeBatch_Job{
			Company: companyName,
			Title:   resJob.Title,
		}

		parsedTime, err := time.Parse(time.RFC3339, resJob.UpdatedAt)
		if err == nil {
			j.UpdatedAt = parsedTime.Unix()
		}

		jobs = append(jobs, j)
	}

	return jobs, nil
}

func (s greenhouseScraper) Sites() []string {
	return []string{
		s.name,
	}
}
