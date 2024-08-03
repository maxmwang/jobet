package scrape

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/maxmwang/jobet/internal/db"
)

func TestGreenhouseScraper_Sites(t *testing.T) {
	s := newGreenhouseScraper()
	assert.Equal(t, s.Sites(), []string{string(db.SiteGreenhouse)})
}
