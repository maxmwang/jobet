package scrape

import (
	"testing"

	"github.com/maxmwang/jobet/internal/db"

	"github.com/stretchr/testify/assert"
)

func TestAshbyScraper_Sites(t *testing.T) {
	s := newAshbyScraper()
	assert.Equal(t, s.Sites(), []string{string(db.SiteAshby)})
}
