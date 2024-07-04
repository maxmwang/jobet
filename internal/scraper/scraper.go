package scraper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/maxmwang/jobet/api"
)

const (
	SiteAshby      string = "ashby"
	SiteGreenhouse string = "greenhouse"
	SiteLever      string = "lever"
)

type Job struct {
	Company   string
	Title     string
	UpdatedAt time.Time
}

func (j Job) ToApi() *api.ScrapeBatch_Job {
	return &api.ScrapeBatch_Job{
		Company:   j.Company,
		Title:     j.Title,
		UpdatedAt: j.UpdatedAt.Unix(),
	}
}

type Scraper interface {
	Scrape(companyName, site string) ([]Job, error)
	ScrapeAll(companyName string) ([]Job, error)
	Sites() []string
}

func checkThenDecode[T any](name, site string, res *http.Response, reqError error) (body T, err error) {
	if reqError != nil {
		return body, fmt.Errorf("failed to request site=%s for company=%s: %w", site, name, reqError)
	}
	if res.StatusCode != http.StatusOK {
		return body, fmt.Errorf("invalid site=%s for company=%s: code=%d", site, name, res.StatusCode)
	}

	if err = json.NewDecoder(res.Body).Decode(&body); err != nil {
		return body, fmt.Errorf("failed to decode response from site=%s for company=%s: %w", site, name, err)
	}

	return body, nil
}

// TODO: https://www.polymer.co/
