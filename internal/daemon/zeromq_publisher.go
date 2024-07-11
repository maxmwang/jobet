package daemon

import (
	"context"

	"github.com/maxmwang/jobet/api"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
	"gopkg.in/zeromq/goczmq.v4"
)

type ZeroMQPublisher struct {
	p *goczmq.Sock
}

func NewZeroMQPublisher(ctx context.Context) *ZeroMQPublisher {
	publisher, err := goczmq.NewPub("tcp://127.0.0.1:5555")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start publisher")
	}

	return &ZeroMQPublisher{
		p: publisher,
	}
}

func (p *ZeroMQPublisher) Publish(ctx context.Context, batch *api.ScrapeBatch) error {
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

func (p *ZeroMQPublisher) Shutdown(ctx context.Context) error {
	p.p.Destroy()
	log.Info().Msg("publisher shutdown")

	return nil
}
