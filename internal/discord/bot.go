package discord

import (
	"context"
	"fmt"
	"strings"

	"github.com/maxmwang/jobet/internal/config"
	"github.com/maxmwang/jobet/internal/db"
	"github.com/maxmwang/jobet/internal/proto"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

type Bot struct {
	*discordgo.Session
	proberClient proto.ProberClient
	dbClient     *db.Queries
}

func NewBot(ctx context.Context, cfg config.Config, proberClient proto.ProberClient, dbClient *db.Queries) (*Bot, error) {
	s, err := discordgo.New("Bot " + cfg.DiscordBotToken)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize discord bot session")
	}
	s.Identify.Intents = discordgo.IntentsGuildMessages
	if err = s.Open(); err != nil {
		return nil, errors.Wrap(err, "failed to start discord bot connection")
	}

	b := &Bot{
		Session:      s,
		proberClient: proberClient,
		dbClient:     dbClient,
	}

	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "subscribe",
			Description: "Subscribes this channel to scrape results.",
		},
		{
			Name:        "unsubscribe",
			Description: "Unsubscribes this channel to scrape results.",
		},
		{
			Name:        "probe",
			Description: "Probes a company without adding to the scrape list.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "company",
					Description: "Company name to probe. Must match the career site URL.",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "site",
					Description: "Site to probe. Leave empty to probe all sites.",
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "https://www.ashbyhq.com/",
							Value: db.SiteAshby,
						},
						{
							Name:  "https://www.greenhouse.com/",
							Value: db.SiteGreenhouse,
						},
						{
							Name:  "https://www.lever.co/",
							Value: db.SiteLever,
						},
					},
				},
			},
		},
	}
	for _, command := range commands {
		_, err := b.ApplicationCommandCreate(b.State.User.ID, "", command)
		if err != nil {
			return nil, errors.Wrap(err, "failed to register command "+command.Name)
		}
	}

	handlers := map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"subscribe": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			content := ""

			err := b.dbClient.AddChannel(ctx, i.ChannelID)
			if err != nil {
				content = "Error: failed to subscribe channel."
			} else {
				content = "Successfully subscribed!"
			}

			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: content,
				},
			})
			if err != nil {
				log.Error().Err(err).Msg("failed to respond to subscribe command")
			}
		},
		"unsubscribe": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			content := ""

			err := b.dbClient.RemoveChannel(ctx, i.ChannelID)
			if err != nil {
				content = "Error: failed to unsubscribe channel."
			} else {
				content = "Successfully unsubscribed!"
			}

			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: content,
				},
			})
			if err != nil {
				log.Error().Err(err).Msg("failed to respond to unsubscribe command")
			}
		},
		"probe": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			req := &proto.ProbeRequest{}
			for _, opt := range i.ApplicationCommandData().Options {
				switch opt.Name {
				case "company":
					req.Company = opt.StringValue()
					break
				case "site":
					req.Site = opt.StringValue()
					break
				}
			}

			res, err := b.proberClient.Probe(ctx, req)
			content := strings.Builder{}
			if err != nil {
				content.WriteString("Error: failed to probe service.")
			} else {
				for _, result := range res.Results {
					if result.Exists {
						content.WriteString(fmt.Sprintf(
							"[%s]: exists=%t; priority=%d; count=%d\n",
							result.Site,
							result.Exists,
							result.Priority,
							result.Count,
						))
					} else {
						content.WriteString(fmt.Sprintf(
							"[%s]: exists=%t; count=%d\n",
							result.Site,
							result.Exists,
							result.Count,
						))
					}
				}
			}
			if content.Len() == 0 {
				content.WriteString("Error: could not scrape.")
			}

			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Probe results for `" + req.Company + "`:\n```" + content.String() + "```",
				},
			})
			if err != nil {
				log.Error().Err(err).Msg("failed to respond to probe command")
			}
		},
	}
	b.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if handler, ok := handlers[i.ApplicationCommandData().Name]; ok {
			handler(s, i)
		}
	})

	return b, nil
}

func (b *Bot) Publish(ctx context.Context, batch *proto.ScrapeBatch) error {
	channels, err := b.dbClient.GetChannels(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get discord channels")
	}

	s := newSquashedJobs(batch)

	eg := errgroup.Group{}
	content := b.paginate(ctx, s.String())
	if len(content) > 0 {
		eg.Go(func() error {
			for _, channelId := range channels {
				for _, v := range content {
					_, err = b.ChannelMessageSend(channelId, fmt.Sprintf("priority<=%d\n```%s```", batch.Priority, v))
					if err != nil {
						return errors.Wrap(err, "failed to send message")
					}
				}
			}
			return nil
		})
	}
	err = eg.Wait()
	if err != nil {
		return err
	}
	return nil
}

func (b *Bot) paginate(ctx context.Context, content string) []string {
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

func (b *Bot) Shutdown(ctx context.Context) error {
	log.Info().Msg("shutting down discord bot")
	err := b.Close()
	if err != nil {
		return err
	}
	return nil
}
