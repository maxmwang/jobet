package db

import (
	"context"

	"github.com/maxmwang/jobet/internal/config"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func Connect(ctx context.Context, cfg config.Config) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, cfg.PostgresURI)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to postgres")
		return nil, errors.Wrap(err, "failed to connect to postgres")
	}

	return conn, nil
}
