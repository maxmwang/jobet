package main

import (
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type env struct {
	botToken        string
	discordChannels []string
}

func loadEnv() env {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	return env{
		botToken:        os.Getenv("BOT_TOKEN"),
		discordChannels: strings.Split(os.Getenv("DISCORD_CHANNELS"), ","),
	}
}

func setupLogger() {
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		NoColor:    false,
		TimeFormat: time.DateTime,
	})
}
