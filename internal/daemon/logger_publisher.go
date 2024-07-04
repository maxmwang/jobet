package daemon

import (
	"context"
	"fmt"

	"github.com/maxmwang/jobet/api"
	"github.com/maxmwang/jobet/internal/helpers"
)

type loggerPublisher struct{}

func NewLoggerPublisher(ctx context.Context) Publisher {
	return loggerPublisher{}
}

func (p loggerPublisher) Publish(ctx context.Context, batch *api.ScrapeBatch) error {
	fmt.Print(helpers.BatchToStringSorted(batch))
	return nil
}

func (p loggerPublisher) Shutdown(ctx context.Context) error {
	return nil
}
