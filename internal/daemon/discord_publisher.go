package daemon

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/maxmwang/jobet/api"
	"github.com/maxmwang/jobet/internal/helpers"
	"github.com/rs/zerolog/log"
)

type DiscordPublisher struct {
	s        *discordgo.Session
	channels []string
}

func NewDiscordPublisher(ctx context.Context, botToken string, channels ...string) (*DiscordPublisher, error) {
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

	p := &DiscordPublisher{
		s:        dg,
		channels: channels,
	}
	dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}
		if m.Content == ">subscribe" {
			if slices.Contains(p.channels, m.ChannelID) {
				p.s.ChannelMessageSend(m.ChannelID, "already subscribed")
			} else {
				p.channels = append(p.channels, m.ChannelID)
				log.Info().Str("channelId", m.ChannelID).Msg("adding subscribed channel")
				p.s.ChannelMessageSend(m.ChannelID, "subscribing")
			}
		}
	})

	return p, nil
}

func (p *DiscordPublisher) Publish(ctx context.Context, batch *api.ScrapeBatch) error {
	rawContent := helpers.BatchToStringSorted(batch)
	content := p.paginate(ctx, rawContent)
	if len(content) > 0 {
		for _, channelId := range p.channels {
			for _, v := range content {
				_, err := p.s.ChannelMessageSend(channelId, fmt.Sprintf("priority<=%d\n```%s```", batch.Priority, v))
				if err != nil {
					log.Error().
						Str("publisher", "discord").
						Err(err).
						Msg("failed to send message")
					return err
				}
			}
		}
	}
	return nil
}

func (p *DiscordPublisher) paginate(ctx context.Context, content string) []string {
	pages := make([]string, 0)

	lines := strings.Split(content, "\n")
	sb := strings.Builder{}
	for _, v := range lines {
		if sb.Len()+len(v) > 1950 {
			pages = append(pages, sb.String())
			sb.Reset()
		}
		sb.WriteString(v)
		sb.WriteRune('\n')
	}
	if sb.Len() > 0 {
		pages = append(pages, sb.String())
	}
	return pages
}

func (p *DiscordPublisher) Shutdown(ctx context.Context) error {
	if err := p.s.Close(); err != nil {
		return err
	}
	return nil
}
