package main

import (
	"context"
	"os"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/maxmwang/jobet/internal/config"
	"github.com/maxmwang/jobet/internal/db"
	"github.com/maxmwang/jobet/internal/prober"
	"github.com/maxmwang/jobet/internal/scrape"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	setupLogger()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.LoadEnv()

	scraper := scrape.NewDefault()

	conn, err := db.Connect(ctx, cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start db client")
	}
	dbClient := db.New(conn)

	proberServer := prober.NewServer(ctx, cfg, scraper, dbClient)

	eg := errgroup.Group{}
	eg.Go(func() error {
		log.Info().Msg("starting prober server")
		return proberServer.Start(ctx)
	})

	err = eg.Wait()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start prober server")
	}
	log.Info().Msg("shutting down with no errors")
}

func setupLogger() {
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		NoColor:    false,
		TimeFormat: time.DateTime,
	})
}
