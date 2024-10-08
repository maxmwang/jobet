package helpers

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/maxmwang/jobet/internal/proto"
)

func JobIsTarget(title string) bool {
	if strings.Contains(title, "Intern") {
		index := strings.Index(title, "Intern")
		if index == strings.Index(title, "Internal") ||
			index == strings.Index(title, "International") {
			return false
		}

		if strings.Contains(title, "Software") ||
			strings.Contains(title, "Platform") ||
			strings.Contains(title, "Machine Learning") {
			return true
		}
	}

	if strings.Contains(title, "New Grad") {
		if strings.Contains(title, "Software") ||
			strings.Contains(title, "Platform") ||
			strings.Contains(title, "Machine Learning") {
			return true
		}
	}

	return false
}

func BatchToStringSorted(b *proto.ScrapeBatch) string {
	slices.SortFunc(b.Jobs, func(a, b *proto.ScrapeBatch_Job) int {
		if a.UpdatedAt > b.UpdatedAt {
			return 1
		}
		return -1
	})

	sb := strings.Builder{}
	for _, j := range b.Jobs {
		if JobIsTarget(j.Title) {
			sb.WriteString(JobToString(j))
			sb.WriteRune('\n')
		}
	}
	return sb.String()
}

func BatchToMarkdownTableSorted(b *proto.ScrapeBatch) string {
	slices.SortFunc(b.Jobs, func(a, b *proto.ScrapeBatch_Job) int {
		if a.UpdatedAt > b.UpdatedAt {
			return 1
		}
		return -1
	})

	sb := strings.Builder{}
	sb.WriteString("| Company | Title | UpdatedAt |\n")
	sb.WriteString("| ------- | ----- | --------- |\n")
	for _, job := range b.Jobs {
		updatedAt := ""
		if !time.Unix(job.UpdatedAt, 0).IsZero() {
			loc, _ := time.LoadLocation("America/Los_Angeles")
			t := time.Unix(job.UpdatedAt, 0)
			updatedAt = t.In(loc).Format(time.DateTime + " MST")
		}
		sb.WriteString(fmt.Sprintf("| %s | %s | %s |", job.Company, job.Title, updatedAt))
	}
	return sb.String()
}

func JobToString(j *proto.ScrapeBatch_Job) string {
	if j.UpdatedAt == 0 {
		return fmt.Sprintf("%46s:\t %v", j.Company, j.Title)
	} else {
		loc, _ := time.LoadLocation("America/Los_Angeles")
		t := time.Unix(j.UpdatedAt, 0)
		return fmt.Sprintf("%24s: %20s:\t %v", t.In(loc).Format(time.DateTime+" MST"), j.Company, j.Title)
	}
}
