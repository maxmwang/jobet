package daemon

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/maxmwang/jobet/internal/db"
	"github.com/maxmwang/jobet/internal/scraper"
	"github.com/rs/zerolog/log"
)

type Daemon struct {
	scraper scraper.Scraper
	q       *db.Queries
}

func NewDefaultDaemon(ctx context.Context, q *db.Queries) *Daemon {
	return &Daemon{
		scraper: scraper.NewDefault(),
		q:       q,
	}
}

func (d *Daemon) Start() {
	t := time.NewTicker(10 * time.Minute)
	count := 0
	ctx := context.Background()

	for {
		priority := getMaxPriority(count)

		companies, err := d.q.GetCompaniesByMaxPriority(ctx, priority)
		if err != nil {
			log.Error().
				Err(err).
				Msg("failed to get companies from database")
			continue
		}
		log.Info().
			Int64("priority<", priority).
			Int("count", len(companies)).
			Msg("scraping")

		wgScrape := new(sync.WaitGroup)
		jobsChan := make(chan []Job)
		for _, c := range companies {
			wgScrape.Add(1)
			go func() {
				defer wgScrape.Done()

				// scrape for all jobs
				scrapeJobs, err := d.scraper.Scrape(c.Name, c.Site)
				if err != nil {
					log.Error().
						Str("name", c.Name).
						Err(err).
						Msg("failed to scrape company")
					return
				}

				// map and reduce jobs
				outputJobs := make([]Job, 0)
				for _, j := range scrapeJobs {
					outputJob := jobFromScrape(j, c)
					if outputJob.IsTarget() {
						outputJobs = append(outputJobs, outputJob)
					}
				}

				jobsChan <- outputJobs
			}()
		}

		wgCollect := new(sync.WaitGroup)
		allJobs := make([]Job, 0)
		wgCollect.Add(1)
		go func() {
			defer wgCollect.Done()
			for j := range jobsChan {
				allJobs = append(allJobs, j...)
			}
		}()

		wgScrape.Wait()
		close(jobsChan)
		wgCollect.Wait()

		batch := Batch{
			Priority: priority,
			Jobs:     allJobs,
		}
		logTemp(batch)

		<-t.C
		count = (count + 1) % 144
	}
}

func getMaxPriority(count int) int64 {
	if count%144 == 0 { // 1d
		return 5
	} else if count%36 == 0 { // 6hr
		return 4
	} else if count%6 == 0 { // 1hr
		return 3
	} else if count%2 == 0 { // 20min
		return 2
	} else {
		return 1
	}
}

func logTemp(batch Batch) {
	batch.Sort()
	for _, j := range batch.Jobs {
		fmt.Println(j)
	}
}
