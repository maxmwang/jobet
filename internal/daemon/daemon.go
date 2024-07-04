package daemon

import (
	"context"
	"sync"
	"time"

	"github.com/maxmwang/jobet/api"
	"github.com/maxmwang/jobet/internal/db"
	"github.com/maxmwang/jobet/internal/scraper"
	"github.com/rs/zerolog/log"
)

type Daemon struct {
	scraper    scraper.Scraper
	q          *db.Queries
	publishers []Publisher
}

func NewDefaultDaemon(ctx context.Context, q *db.Queries, publishers ...Publisher) *Daemon {
	return &Daemon{
		scraper:    scraper.NewDefault(),
		q:          q,
		publishers: publishers,
	}
}

func (d *Daemon) Start(ctx context.Context) {
	t := time.NewTicker(10 * time.Minute)
	count := 0

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
		jobsChan := make(chan []*api.ScrapeBatch_Job)
		for _, company := range companies {
			wgScrape.Add(1)
			go func() {
				defer wgScrape.Done()

				// scrape for all jobs
				scrapeJobs, err := d.scraper.Scrape(company.Name, company.Site)
				if err != nil {
					log.Error().
						Str("name", company.Name).
						Err(err).
						Msg("failed to scrape company")
					return
				}

				// map and reduce jobs
				outputJobs := make([]*api.ScrapeBatch_Job, 0)
				for _, job := range scrapeJobs {
					outputJobs = append(outputJobs, job.ToApi())
				}

				jobsChan <- outputJobs
			}()
		}

		wgCollect := new(sync.WaitGroup)
		allJobs := make([]*api.ScrapeBatch_Job, 0)
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

		batch := &api.ScrapeBatch{
			Priority: priority,
			Jobs:     allJobs,
		}
		for _, publisher := range d.publishers {
			err = publisher.Publish(ctx, batch)
			if err != nil {
				log.Error().Err(err).Msg("failed to publish batch")
			}
		}

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
