package daemon

import (
	"context"

	"github.com/maxmwang/jobet/api"
)

type Publisher interface {
	Publish(ctx context.Context, batch *api.ScrapeBatch) error
	Shutdown(ctx context.Context) error
}
