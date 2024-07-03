package daemon

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/maxmwang/jobet/internal/db"
	"github.com/maxmwang/jobet/internal/scraper"
)

type Batch struct {
	Priority int64
	Jobs     []Job
}

func (b *Batch) Sort() {
	slices.SortFunc(b.Jobs, func(a, b Job) int {
		return a.UpdatedAt.Compare(b.UpdatedAt)
	})
}

type Job struct {
	Company   db.Company
	Title     string
	UpdatedAt time.Time
}

func jobFromScrape(j scraper.Job, company db.Company) Job {
	return Job{
		Company:   company,
		Title:     j.Title,
		UpdatedAt: j.UpdatedAt,
	}
}

func (j Job) String() string {
	if j.UpdatedAt.IsZero() {
		return fmt.Sprintf("%46s:\t %v", j.Company.Name, j.Title)
	} else {
		loc, _ := time.LoadLocation("America/Los_Angeles")
		return fmt.Sprintf("%24s: %20s:\t %v", j.UpdatedAt.In(loc).Format(time.DateTime+" MST"), j.Company.Name, j.Title)
	}
}

func (j Job) IsTarget() bool {
	if strings.Index(j.Title, "Intern") == strings.Index(j.Title, "Internal") && strings.Count(j.Title, "Intern") == 1 {
		return false
	}
	if strings.Index(j.Title, "Intern") == strings.Index(j.Title, "International") && strings.Count(j.Title, "Intern") == 1 {
		return false
	}
	if strings.Contains(j.Title, "Software") && strings.Contains(j.Title, "Intern") {
		return true
	}
	if strings.Contains(j.Title, "Platform") && strings.Contains(j.Title, "Intern") {
		return true
	}

	return false
}
