package main

import (
	"context"
	"sync"

	"github.com/maxmwang/jobet/internal/db"
	"github.com/maxmwang/jobet/internal/scrape"
)

func main() {
	conn, err := db.Connect(false)
	if err != nil {
		panic(err)
	}
	q := db.New(conn)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		scrape.Daemon(context.Background(), q)
	}()
	wg.Wait()
}
