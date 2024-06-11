package scrape

import (
	"context"
	"fmt"
	"log"
	"slices"
	"strings"
	"time"

	"github.com/maxmwang/jobet/internal/db"
	"github.com/maxmwang/jobet/internal/output"
)

func Daemon(ctx context.Context, q *db.Queries) {
	t := time.NewTicker(10 * time.Minute)
	count := 0

	for {
		count++
		priority := 1

		if count%144 == 0 { // 1d
			priority = 5
			count = 0 // reset to prevent overflow
		} else if count%36 == 0 { // 6hr
			priority = 4
		} else if count%6 == 0 { // 1hr
			priority = 3
		} else if count%2 == 0 { // 20min
			priority = 2
		}

		fmt.Printf("[%s] scraping priority<=%d\n", time.Now().Format(time.TimeOnly), priority)

		companies, err := q.GetCompaniesByMaxPriority(ctx, int64(priority))
		if err != nil {
			log.Printf("failed to get companies from database: %v", err)
			continue
		}

		sortedJobs := make([]output.Job, 0)
		for _, c := range companies {
			jobs, err := scrape(c)
			if err != nil {
				log.Printf("failed to scrape company: %s", err.Error())
				continue
			}
			sortedJobs = append(sortedJobs, jobs...)
		}
		slices.SortFunc(sortedJobs, func(a, b output.Job) int {
			return strings.Compare(a.UpdatedAt.String(), b.UpdatedAt.String())
		})

		output.Log(sortedJobs)

		<-t.C
	}
}
