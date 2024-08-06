package main

import (
	"context"
	"os"
	"sync"
	"time"
	
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

	wg := sync.WaitGroup{}

	proberServer := prober.NewServer(ctx, cfg, scraper, dbClient)
	wg.Add(1)
	go func() {
		defer wg.Done()
		proberServer.Start(ctx)
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
