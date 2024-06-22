package scrape

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	SiteAshby      string = "ashby"
	SiteGreenhouse string = "greenhouse"
	SiteLever      string = "lever"
)

type Job struct {
	Title     string
	UpdatedAt time.Time
}

type Scraper interface {
	Scrape(companyName string) ([]Job, error)
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
