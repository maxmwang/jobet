package worker

import (
	"context"

	"github.com/maxmwang/jobet/internal/proto"
)

type Publisher interface {
	Publish(ctx context.Context, batch *proto.ScrapeBatch) error
	Shutdown(ctx context.Context) error
}
