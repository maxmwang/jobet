package scrape

import (
	"errors"

	"github.com/maxmwang/jobet/internal/db"
	"github.com/maxmwang/jobet/internal/proto"
)

type defaultScraper struct {
	scrapers map[string]Scraper
}

func NewDefault() Scraper {
	return defaultScraper{
		scrapers: map[string]Scraper{
			string(db.SiteAshby):      newAshbyScraper(),
			string(db.SiteGreenhouse): newGreenhouseScraper(),
			string(db.SiteLever):      newLeverScraper(),
		},
	}
}

func (s defaultScraper) Scrape(companyName, site string) ([]*proto.ScrapeBatch_Job, error) {
	siteScraper, ok := s.scrapers[site]
	if !ok {
		return nil, errors.New("scrape not registered")
	}

	jobs, err := siteScraper.Scrape(companyName, site)
	if err != nil {
		return nil, errors.New("failed to scrape company")
	}

	return jobs, nil
}

func (s defaultScraper) ScrapeAll(companyName string) ([]*proto.ScrapeBatch_Job, error) {
	jobs := make([]*proto.ScrapeBatch_Job, 0)

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
	for site := range s.scrapers {
		sites = append(sites, site)
	}
	return sites
}
