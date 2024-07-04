package main

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/maxmwang/jobet/api"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
	"gopkg.in/zeromq/goczmq.v4"
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
		msgBytes, err := subscriber.RecvMessage()
		if err != nil {
			log.Fatal().Err(err).Msg("failed to receive message")
		}

		msg := api.ScrapeBatch{}
		err = proto.Unmarshal(msgBytes[0], &msg)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to unmarshal message")
		}

		log.Info().
			Int64("priority", msg.Priority).
			Msg("received msg")
		slices.SortFunc(msg.Jobs, func(a, b *api.ScrapeBatch_Job) int {
			if a.UpdatedAt > b.UpdatedAt {
				return 1
			}
			return -1
		})
		for _, j := range msg.Jobs {
			if isTarget(j.Title) {
				fmt.Println(toString(j))
			}
		}
	}
}

func isTarget(title string) bool {
	if strings.Index(title, "Intern") == strings.Index(title, "Internal") && strings.Count(title, "Intern") == 1 {
		return false
	}
	if strings.Index(title, "Intern") == strings.Index(title, "International") && strings.Count(title, "Intern") == 1 {
		return false
	}
	if strings.Contains(title, "Software") && strings.Contains(title, "Intern") {
		return true
	}
	if strings.Contains(title, "Platform") && strings.Contains(title, "Intern") {
		return true
	}

	return false
}

func toString(j *api.ScrapeBatch_Job) string {
	if time.Unix(j.UpdatedAt, 0).IsZero() {
		return fmt.Sprintf("%46s:\t %v", j.Company, j.Title)
	} else {
		loc, _ := time.LoadLocation("America/Los_Angeles")
		t := time.Unix(j.UpdatedAt, 0)
		return fmt.Sprintf("%24s: %20s:\t %v", t.In(loc).Format(time.DateTime+" MST"), j.Company, j.Title)
	}
}

func setupLogger() {
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		NoColor:    false,
		TimeFormat: time.DateTime,
	})
}
