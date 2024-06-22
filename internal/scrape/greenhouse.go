package scrape

import (
	"fmt"
	"net/http"
	"time"
)

type GreenhouseScraper struct{}

func NewGreenhouseScraper() GreenhouseScraper {
	return GreenhouseScraper{}
}

func (a GreenhouseScraper) Scrape(companyName string) ([]Job, error) {
	res, err := http.Get(fmt.Sprintf("https://api.greenhouse.io/v1/boards/%s/jobs", companyName))

	body, err := checkThenDecode[struct {
		Jobs []struct {
			Title     string `json:"title"`
			UpdatedAt string `json:"updated_at"`
		} `json:"jobs"`
	}](companyName, SiteGreenhouse, res, err)
	if err != nil {
		return nil, err
	}

	jobs := make([]Job, 0)
	for _, resJob := range body.Jobs {
		j := Job{
			Title: resJob.Title,
		}

		parsedTime, err := time.Parse(time.RFC3339, resJob.UpdatedAt)
		if err == nil {
			j.UpdatedAt = parsedTime
		}

		jobs = append(jobs, j)
	}

	return jobs, nil
}
