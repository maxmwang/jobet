package scraper

import (
	"errors"
)

type defaultScraper struct {
	scrapers map[string]Scraper
}

func NewDefault() Scraper {
	return defaultScraper{
		scrapers: map[string]Scraper{
			SiteAshby:      newAshbyScraper(),
			SiteGreenhouse: newGreenhouseScraper(),
			SiteLever:      newLeverScraper(),
		},
	}
}

func (s defaultScraper) Scrape(companyName, site string) ([]Job, error) {
	siteScraper, ok := s.scrapers[site]
	if !ok {
		return nil, errors.New("scraper not registered")
	}

	jobs, err := siteScraper.Scrape(companyName, site)
	if err != nil {
		return nil, errors.New("failed to scrape company")
	}

	return jobs, nil
}

func (s defaultScraper) ScrapeAll(companyName string) ([]Job, error) {
	jobs := make([]Job, 0)

	for site, siteScraper := range s.scrapers {
		siteJobs, err := siteScraper.Scrape(companyName, site)
		if err != nil {
			return nil, errors.New("failed to scrape company")
		}
		jobs = append(jobs, siteJobs...)
	}

	return jobs, nil
}

func (s defaultScraper) Sites() []string {
	sites := make([]string, 0)
	for site, _ := range s.scrapers {
		sites = append(sites, site)
	}
	return sites
}
