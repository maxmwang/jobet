package helpers

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/maxmwang/jobet/api"
)

func JobIsTarget(title string) bool {
	if strings.Index(title, "Intern") == strings.Index(title, "Internal") && strings.Count(title, "Intern") == 1 {
		return false
	}
	if strings.Index(title, "Intern") == strings.Index(title, "International") && strings.Count(title, "Intern") == 1 {
		return false
	}
	if strings.Contains(title, "Software") && strings.Contains(title, "Intern") {
		return true
	}
	if strings.Contains(title, "Platform") && strings.Contains(title, "Intern") {
		return true
	}

	return false
}

func PrintBatchSorted(b *api.ScrapeBatch) {
	slices.SortFunc(b.Jobs, func(a, b *api.ScrapeBatch_Job) int {
		if a.UpdatedAt > b.UpdatedAt {
			return 1
		}
		return -1
	})
	for _, j := range b.Jobs {
		if JobIsTarget(j.Title) {
			fmt.Println(JobToString(j))
		}
	}
}

func JobToString(j *api.ScrapeBatch_Job) string {
	if time.Unix(j.UpdatedAt, 0).IsZero() {
		return fmt.Sprintf("%46s:\t %v", j.Company, j.Title)
	} else {
		loc, _ := time.LoadLocation("America/Los_Angeles")
		t := time.Unix(j.UpdatedAt, 0)
		return fmt.Sprintf("%24s: %20s:\t %v", t.In(loc).Format(time.DateTime+" MST"), j.Company, j.Title)
	}
}
