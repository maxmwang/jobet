package main

import (
	"context"
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

	conn, err := db.Connect(false)
	if err != nil {
		panic(err)
	}
	q := db.New(conn)

	d := daemon.NewDefaultDaemon(context.Background(), q)
	s := server.NewJobetServer(q)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		d.Start()
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
