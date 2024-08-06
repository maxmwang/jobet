package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/maxmwang/jobet/internal/proto"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	setupLogger()

	company := flag.String("c", "", "company name")
	site := flag.String("a", "", "site")
	flag.Parse()

	if *company == "" {
		panic("please provide a company name with the -c flag")
	}

	conn, err := grpc.NewClient("localhost:5001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := proto.NewProberClient(conn)
	res, err := client.Probe(context.Background(), &proto.ProbeRequest{
		Company: *company,
		Site:    *site,
	})
	if err != nil {
		panic(fmt.Errorf("could not probe server: %w", err))
	}
	for _, r := range res.Results {
		if r.Exists {
			log.Info().
				Bool("exists", r.Exists).
				Int32("priority", r.Priority).
				Int32("count", r.Count).
				Msg(fmt.Sprintf("[site=%s]", r.Site))
		} else {
			log.Info().
				Bool("exists", r.Exists).
				Int32("count", r.Count).
				Msg(fmt.Sprintf("[site=%s]", r.Site))
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
