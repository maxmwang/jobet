package main

import (
	"context"
	"flag"
	"os"
	"sync"
	"time"

	"github.com/maxmwang/jobet/internal/daemon"
	"github.com/maxmwang/jobet/internal/db"
	"github.com/maxmwang/jobet/internal/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	setupLogger()

	useLogger := flag.Bool("l", false, "use logger publisher")
	useZeromq := flag.Bool("z", false, "use zeromq publisher")
	flag.Parse()

	ctx := context.Background()

	conn, err := db.Connect(false)
	if err != nil {
		panic(err)
	}
	q := db.New(conn)

	publishers := make([]daemon.Publisher, 0)
	if *useLogger {
		publishers = append(publishers, daemon.NewLoggerPublisher(ctx))
		log.Info().Msg("adding logger publisher")
	}
	if *useZeromq {
		publishers = append(publishers, daemon.NewZeroMQPublisher(ctx))
		log.Info().Msg("adding zeromq publisher")
	}

	d := daemon.NewDefaultDaemon(ctx, q, publishers...)
	s := server.NewJobetServer(q)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		d.Start(ctx)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.Start()
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
