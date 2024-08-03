package scrape

import (
	"testing"

	"github.com/maxmwang/jobet/internal/db"
	"github.com/maxmwang/jobet/internal/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestDefaultScraper_ScrapeAll(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockSite := "test-site"

	mockAshbyScraper := mocks.NewMockScraper(ctrl)
	mockGreenhouseScraper := mocks.NewMockScraper(ctrl)
	mockLeverScraper := mocks.NewMockScraper(ctrl)

	mockAshbyScraper.EXPECT().Scrape(mockSite, string(db.SiteAshby))
	mockGreenhouseScraper.EXPECT().Scrape(mockSite, string(db.SiteGreenhouse))
	mockLeverScraper.EXPECT().Scrape(mockSite, string(db.SiteLever))

	s := defaultScraper{
		scrapers: map[string]Scraper{
			string(db.SiteAshby):      mockAshbyScraper,
			string(db.SiteGreenhouse): mockGreenhouseScraper,
			string(db.SiteLever):      mockLeverScraper,
		},
	}
	_, err := s.ScrapeAll(mockSite)
	assert.Nil(t, err)
}

func TestDefaultScraper_Sites(t *testing.T) {
	s := NewDefault()
	assert.Subset(t, s.Sites(), []string{string(db.SiteAshby), string(db.SiteGreenhouse), string(db.SiteLever)})
}
