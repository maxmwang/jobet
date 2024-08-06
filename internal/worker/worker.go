package worker

import (
	"context"
	"sync"
	"time"

	"github.com/maxmwang/jobet/internal/db"
	"github.com/maxmwang/jobet/internal/proto"
	"github.com/maxmwang/jobet/internal/scrape"

	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

type Worker struct {
	scraper    scrape.Scraper
	dbClient   *db.Queries
	publishers []Publisher
}

func NewWorker(ctx context.Context, dbClient *db.Queries, scraper scrape.Scraper, publishers ...Publisher) *Worker {
	return &Worker{
		scraper:    scraper,
		dbClient:   dbClient,
		publishers: publishers,
	}
}

func (d *Worker) Start(ctx context.Context) {
	t := time.NewTicker(30 * time.Minute)
	count := 0

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		priority := d.getMaxPriority(count)

		companies, err := d.dbClient.GetCompaniesByMaxPriority(ctx, priority)
		if err != nil {
			log.Error().
				Err(err).
				Msg("failed to get companies from database")
			continue
		}
		log.Info().
			Int32("priority<", priority).
			Int("count", len(companies)).
			Msg("scraping")

		scrapeWg := sync.WaitGroup{}
		jobsChan := make(chan []*proto.ScrapeBatch_Job)
		for _, company := range companies {
			scrapeWg.Add(1)
			go func() {
				defer scrapeWg.Done()

				// scrape for all jobs
				scrapeJobs, err := d.scraper.Scrape(company.Name, string(company.Site))
				if err != nil {
					log.Error().
						Str("name", company.Name).
						Err(err).
						Msg("failed to scrape company")
					return
				}

				jobsChan <- scrapeJobs
			}()
		}

		collectWg := sync.WaitGroup{}
		allJobs := make([]*proto.ScrapeBatch_Job, 0)
		collectWg.Add(1)
		go func() {
			defer collectWg.Done()
			for j := range jobsChan {
				allJobs = append(allJobs, j...)
			}
		}()

		scrapeWg.Wait()
		close(jobsChan)
		collectWg.Wait()

		batch := &proto.ScrapeBatch{
			Priority: priority,
			Jobs:     allJobs,
		}
		// swallow error
		_ = d.publish(ctx, batch)

		<-t.C
		count = (count + 1) % 144
	}
}

func (d *Worker) publish(ctx context.Context, batch *proto.ScrapeBatch) error {
	publishWg := errgroup.Group{}
	for _, publisher := range d.publishers {
		publishWg.Go(func() error {
			if err := publisher.Publish(ctx, batch); err != nil {
				log.Error().Err(err).Msg("failed to publish batch")
				return err
			}
			return nil
		})
	}
	err := publishWg.Wait()
	if err != nil {
		return err
	}
	return nil
}

func (d *Worker) getMaxPriority(count int) int32 {
	if count%48 == 0 { // 1d
		return 5
	} else if count%24 == 0 { // 12hr
		return 4
	} else if count%6 == 0 { // 3hr
		return 3
	} else if count%2 == 0 { // 1hr
		return 2
	} else {
		return 1
	}
}
