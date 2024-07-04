package main

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type env struct {
	botToken string
}

func loadEnv() env {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	return env{
		botToken: os.Getenv("BOT_TOKEN"),
	}
}

func setupLogger() {
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		NoColor:    false,
		TimeFormat: time.DateTime,
	})
}
