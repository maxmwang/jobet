package scraper

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestDefaultScraper_ScrapeAll(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockSite := "test-site"

	mockAshbyScraper := NewMockScraper(ctrl)
	mockGreenhouseScraper := NewMockScraper(ctrl)
	mockLeverScraper := NewMockScraper(ctrl)

	mockAshbyScraper.EXPECT().Scrape(mockSite, SiteAshby)
	mockGreenhouseScraper.EXPECT().Scrape(mockSite, SiteGreenhouse)
	mockLeverScraper.EXPECT().Scrape(mockSite, SiteLever)

	s := defaultScraper{
		scrapers: map[string]Scraper{
			SiteAshby:      mockAshbyScraper,
			SiteGreenhouse: mockGreenhouseScraper,
			SiteLever:      mockLeverScraper,
		},
	}
	_, err := s.ScrapeAll(mockSite)
	assert.Nil(t, err)
}

func TestDefaultScraper_Sites(t *testing.T) {
	s := NewDefault()
	assert.Subset(t, s.Sites(), []string{SiteAshby, SiteGreenhouse, SiteLever})
}
