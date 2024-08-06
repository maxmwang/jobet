package main

import (
	"context"
	"flag"
	"os"
	"sync"
	"time"
	
	"github.com/maxmwang/jobet/internal/config"
	"github.com/maxmwang/jobet/internal/db"
	"github.com/maxmwang/jobet/internal/discord"
	"github.com/maxmwang/jobet/internal/prober"
	"github.com/maxmwang/jobet/internal/scrape"
	"github.com/maxmwang/jobet/internal/worker"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	setupLogger()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	useLogger := flag.Bool("l", false, "use logger publisher")
	useZeromq := flag.Bool("z", false, "use zeromq publisher")
	useDiscord := flag.Bool("d", false, "use discord publisher")
	flag.Parse()
	cfg := config.LoadEnv()

	scraper := scrape.NewDefault()

	conn, err := db.Connect(ctx, cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start db client")
	}
	dbClient := db.New(conn)

	proberClient, err := prober.NewClient(ctx, cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start prober client")
	}

	publishers := make([]worker.Publisher, 0)
	if *useLogger {
		publishers = append(publishers, worker.NewLoggerPublisher(ctx))
		log.Info().Msg("adding logger publisher")
	}
	if *useZeromq {
		publishers = append(publishers, worker.NewZeroMQPublisher(ctx))
		log.Info().Msg("adding zeromq publisher")
	}
	if *useDiscord {
		discordBot, err := discord.NewBot(ctx, cfg, proberClient, dbClient)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to start discord bot")
		}
		publishers = append(publishers, discordBot)
		log.Info().Msg("adding discord publisher")
	}

	scrapeWorker := worker.NewWorker(ctx, dbClient, scraper, publishers...)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		scrapeWorker.Start(ctx)
	}()

	wg.Wait()
}

func setupLogger() {
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		NoColor:    false,
		TimeFormat: time.DateTime,
	})
}
