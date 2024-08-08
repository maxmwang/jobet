package main

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/zeromq/goczmq.v4"

	"github.com/maxmwang/jobet/internal/helpers"
	api "github.com/maxmwang/jobet/internal/proto"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
)

func main() {
	setupLogger()

	subscriber, err := goczmq.NewSub("tcp://127.0.0.1:5555", "")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start subscriber")
	}
	defer subscriber.Destroy()
	log.Info().Msg("subscriber started")

	for {
		msg, err := subscriber.RecvMessage()
		if err != nil {
			log.Fatal().Err(err).Msg("failed to receive message")
		}

		batch := &api.ScrapeBatch{}
		err = proto.Unmarshal(msg[0], batch)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to unmarshal message")
		}

		log.Info().
			Int32("priority<", batch.Priority).
			Msg("received msg")
		fmt.Print(helpers.BatchToStringSorted(batch))
	}
}

func setupLogger() {
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		NoColor:    false,
		TimeFormat: time.DateTime,
	})
}
