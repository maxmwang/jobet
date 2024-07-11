package daemon

import (
	"context"
	"fmt"

	"github.com/maxmwang/jobet/api"
	"github.com/maxmwang/jobet/internal/helpers"
)

type LoggerPublisher struct{}

func NewLoggerPublisher(ctx context.Context) *LoggerPublisher {
	return &LoggerPublisher{}
}

func (p *LoggerPublisher) Publish(ctx context.Context, batch *api.ScrapeBatch) error {
	fmt.Print(helpers.BatchToStringSorted(batch))
	return nil
}

func (p *LoggerPublisher) Shutdown(ctx context.Context) error {
	return nil
}
