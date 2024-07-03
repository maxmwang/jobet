package scraper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLeverScraper_Sites(t *testing.T) {
	s := newLeverScraper()
	assert.Equal(t, s.Sites(), []string{SiteLever})
}
