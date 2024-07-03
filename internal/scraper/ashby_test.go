package scraper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAshbyScraper_Sites(t *testing.T) {
	s := newAshbyScraper()
	assert.Equal(t, s.Sites(), []string{SiteAshby})
}
