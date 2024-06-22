package scrape

import (
	"fmt"
	"net/http"
	"time"
)

type LeverScraper struct{}

func NewLeverScraper() LeverScraper {
	return LeverScraper{}
}

func (a LeverScraper) Scrape(companyName string) ([]Job, error) {
	res, err := http.Get(fmt.Sprintf("https://api.lever.co/v0/postings/%s?limit=999", companyName))

	body, err := checkThenDecode[[]struct {
		Title     string `json:"text"`
		UpdatedAt int    `json:"createdAt"`
	}](companyName, SiteLever, res, err)
	if err != nil {
		return nil, err
	}

	jobs := make([]Job, 0)
	for _, j := range body {
		jobs = append(jobs, Job{
			Title:     j.Title,
			UpdatedAt: time.UnixMilli(int64(j.UpdatedAt)),
		})
	}

	return jobs, nil
}
