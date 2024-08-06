package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	PostgresURI string

	PollerPort string

	DiscordBotToken string
}

func LoadEnv() Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	return Config{
		PostgresURI:     os.Getenv("POSTGRES_URI"),
		PollerPort:      os.Getenv("POLLER_PORT"),
		DiscordBotToken: os.Getenv("DISCORD_BOT_TOKEN"),
	}
}
