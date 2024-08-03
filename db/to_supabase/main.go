package main

import (
	"context"
	"os"
	"time"

	"github.com/maxmwang/jobet/db/to_supabase/sqlite"
	"github.com/maxmwang/jobet/internal/config"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/maxmwang/jobet/internal/db"
)

func main() {
	setupLogger()

	ctx := context.Background()
	cfg := config.LoadEnv()
	sqliteConn, err := sqlite.Connect(false)
	if err != nil {
		panic(err)
	}
	postgresConn, err := db.Connect(ctx, cfg)
	if err != nil {
		panic(err)
	}

	sqliteQ := sqlite.New(sqliteConn)
	postgresQ := db.New(postgresConn)

	sqliteCompanies, err := sqliteQ.GetCompanies(ctx)

	for _, c := range sqliteCompanies {
		log.Info().Any("company", c).Msg("adding company")
		err := postgresQ.AddCompany(ctx, db.AddCompanyParams{
			Name:     c.Name,
			Alias:    c.Alias,
			Site:     db.Site(c.Site),
			Priority: int32(c.Priority),
		})
		if err != nil {
			log.Error().Err(err).Msg("failed to add company")
		}
	}
}

func setupLogger() {
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		NoColor:    false,
		TimeFormat: time.DateTime,
	})
}
