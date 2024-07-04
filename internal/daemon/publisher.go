package daemon

import (
	"context"

	"github.com/maxmwang/jobet/api"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
	"gopkg.in/zeromq/goczmq.v4"
)

type Publisher interface {
	Publish(ctx context.Context, batch *api.ScrapeBatch) error
	Shutdown(ctx context.Context) error
}

type zeroMQPublisher struct {
	p *goczmq.Sock
}

func NewZeroMQPublisher(ctx context.Context) Publisher {
	publisher, err := goczmq.NewPub("tcp://127.0.0.1:5555")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start publisher")
	}

	return zeroMQPublisher{
		p: publisher,
	}
}

func (p zeroMQPublisher) Publish(ctx context.Context, batch *api.ScrapeBatch) error {
	msg, err := proto.Marshal(batch)
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal message")
	}

	err = p.p.SendMessage([][]byte{msg})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to send message")
	}

	return nil
}

func (p zeroMQPublisher) Shutdown(ctx context.Context) error {
	p.p.Destroy()
	log.Info().Msg("publisher shutdown")

	return nil
}
