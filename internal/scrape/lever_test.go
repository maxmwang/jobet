package scrape

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/maxmwang/jobet/internal/db"
)

func TestLeverScraper_Sites(t *testing.T) {
	s := newLeverScraper()
	assert.Equal(t, s.Sites(), []string{string(db.SiteLever)})
}
