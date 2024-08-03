package prober

import (
	"context"

	"github.com/maxmwang/jobet/internal/config"
	"github.com/maxmwang/jobet/internal/proto"
	
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewClient(ctx context.Context, cfg config.Config) (proto.ProberClient, error) {
	conn, err := grpc.NewClient("localhost:"+cfg.PollerPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create grpc client")
	}
	client := proto.NewProberClient(conn)

	return client, nil
}
