package scrape

import (
	"context"
	"fmt"
	"log"
	"slices"
	"sync"
	"time"

	"github.com/maxmwang/jobet/internal/db"
	"github.com/maxmwang/jobet/internal/output"
)

func Daemon(ctx context.Context, q *db.Queries) {
	t := time.NewTicker(10 * time.Minute)
	count := 144

	for {
		var priority int64
		if count%144 == 0 { // 1d
			priority = 5
			count = 0 // reset to prevent overflow
		} else if count%36 == 0 { // 6hr
			priority = 4
		} else if count%6 == 0 { // 1hr
			priority = 3
		} else if count%2 == 0 { // 20min
			priority = 2
		} else {
			priority = 1
		}

		fmt.Printf("[%s] scraping priority<=%d\n", time.Now().Format(time.TimeOnly), priority)

		companies, err := q.GetCompaniesByMaxPriority(ctx, priority)
		if err != nil {
			log.Printf("failed to get companies from database: %v", err)
			continue
		}

		wgScrape := new(sync.WaitGroup)
		jobsChan := make(chan []output.Job)
		for _, c := range companies {
			wgScrape.Add(1)
			go func() {
				defer wgScrape.Done()
				jobs, err := scrape(c)
				if err != nil {
					log.Printf("failed to scrape company: %s", err.Error())
				}
				jobsChan <- jobs
			}()
		}

		wgCollect := new(sync.WaitGroup)
		sortedJobs := make([]output.Job, 0)
		wgCollect.Add(1)
		go func() {
			defer wgCollect.Done()
			for jobs := range jobsChan {
				sortedJobs = append(sortedJobs, jobs...)
			}
		}()

		wgScrape.Wait()
		close(jobsChan)
		wgCollect.Wait()

		slices.SortFunc(sortedJobs, func(a, b output.Job) int {
			return a.UpdatedAt.Compare(b.UpdatedAt)
		})

		output.Log(sortedJobs)

		<-t.C
		count++
	}
}
