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
	alias := flag.String("a", "", "alias")
	priority := flag.Int64("p", 5, "priority")
	dry := flag.Bool("d", false, "dry run")
	flag.Parse()

	if *company == "" {
		panic("please provide a company name with the -c flag")
	}
	if *alias == "" {
		alias = company
	}

	conn, err := grpc.NewClient("localhost:5001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := proto.NewJobetClient(conn)
	res, err := client.Probe(context.Background(), &proto.ProbeRequest{
		Name:     *company,
		Dry:      *dry,
		Alias:    *alias,
		Priority: *priority,
	})
	if err != nil {
		panic(fmt.Errorf("could not probe server: %w", err))
	}
	for _, r := range res.Results {
		log.Info().
			Int32("count", r.Count).
			Int32("target", r.Target).
			Bool("exists", r.Exists).
			Bool("added", r.Added).
			Msg(fmt.Sprintf("[site=%s]", r.Site))
	}
}

func setupLogger() {
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		NoColor:    false,
		TimeFormat: time.DateTime,
	})
}
