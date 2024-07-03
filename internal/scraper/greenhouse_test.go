package scraper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGreenhouseScraper_Sites(t *testing.T) {
	s := newGreenhouseScraper()
	assert.Equal(t, s.Sites(), []string{SiteGreenhouse})
}
