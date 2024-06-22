package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/maxmwang/jobet/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	company := flag.String("c", "", "company name")
	alias := flag.String("a", "", "alias")
	priority := flag.Int64("p", 5, "priority")
	dry := flag.Bool("b", false, "dry run")
	flag.Parse()

	if *company == "" {
		panic("please provide a company name with the -c flag")
	}
	if *alias == "" {
		alias = company
	}

	conn, err := grpc.Dial("localhost:5001", grpc.WithTransportCredentials(insecure.NewCredentials()))
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
	fmt.Println(len(res.Results))
	for _, r := range res.Results {
		fmt.Printf("[site=%s]: count=%d, target=%d, exists=%t, added=%t\n", r.Site, r.Count, r.Target, r.Exists, r.Added)
	}
}
