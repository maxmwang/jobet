package daemon

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/maxmwang/jobet/api"
	"github.com/maxmwang/jobet/internal/helpers"
	"github.com/rs/zerolog/log"
)

type discordPublisher struct {
	s        *discordgo.Session
	channels []string
}

func NewDiscordPublisher(ctx context.Context, botToken string) (Publisher, error) {
	dg, err := discordgo.New("Bot " + botToken)
	if err != nil {
		log.Error().Err(err).Msg("failed to create discord publisher")
		return nil, err
	}
	dg.Identify.Intents = discordgo.IntentsGuildMessages
	if err = dg.Open(); err != nil {
		log.Error().Err(err).Msg("failed to start discord publisher")
		return nil, err
	}

	p := &discordPublisher{
		s:        dg,
		channels: make([]string, 0),
	}
	dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}
		if m.Content == ">subscribe" {
			p.channels = append(p.channels, m.ChannelID)
			log.Info().Str("channelId", m.ChannelID).Msg("adding subscribed channel")
			p.s.ChannelMessageSend(m.ChannelID, "subscribing")
		}
	})

	return p, nil
}

func (p *discordPublisher) Publish(ctx context.Context, batch *api.ScrapeBatch) error {
	for _, channelId := range p.channels {
		_, err := p.s.ChannelMessageSend(channelId, "```"+helpers.BatchToStringSorted(batch)+"```")
		if err != nil {
			log.Error().
				Str("publisher", "discord").
				Err(err).
				Msg("failed to send message")
			return err
		}
	}
	return nil
}

func (p *discordPublisher) Shutdown(ctx context.Context) error {
	if err := p.s.Close(); err != nil {
		return err
	}
	return nil
}
