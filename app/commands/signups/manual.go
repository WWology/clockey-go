package signups

import (
	"clockey/app"
	"clockey/database/sqlc"
	"context"
	"log/slog"
	"strconv"
	"time"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/handler"
)

var Manual = discord.SlashCommandCreate{
	Name:        "manual",
	Description: "Manually assign gardeners to an event",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionString{
			Name:        "gardener",
			Description: "Gardener to work on the event",
			Required:    true,
			Choices: []discord.ApplicationCommandOptionChoiceString{
				{
					Name:  "N1k",
					Value: "N1k",
				},
				{
					Name:  "Kit",
					Value: "Kit",
				},
				{
					Name:  "WW",
					Value: "WW",
				},
				{
					Name:  "Bonteng",
					Value: "Bonteng",
				},
				{
					Name:  "Sam",
					Value: "Sam",
				},
			},
		},
	},
}

func ManualCommandHandler(b *app.Bot) handler.SlashCommandHandler {
	return func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
		// Show modal to collect event details
		if err := e.Modal(eventModal); err != nil {
			return err
		}

		// Handle modal submission
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			defer cancel()
			bot.WaitForEvent(e.Client(), ctx,
				func(m *events.ModalSubmitInteractionCreate) bool {
					return m.Data.CustomID == "event_modal"
				},
				func(m *events.ModalSubmitInteractionCreate) {
					// Handle the event details submission
					unixValue, err := strconv.ParseInt(m.Data.Text("event_time"), 0, 64)
					eventType := m.Data.StringValues("event_type")[0]
					name := m.Data.Text("event_name")
					hours, _ := strconv.ParseInt(m.Data.Text("event_duration"), 10, 64)
					gardener, _ := strconv.ParseInt(m.Data.Text("gardener"), 10, 64)

					if err != nil {
						m.Client().Logger.Error("Failed to parse event_time", slog.Any("err", err))
						m.CreateMessage(discord.MessageCreate{
							Content: m.Data.Text("event_time") + " is not a valid unix time. Please try again.",
						})
						return
					}

					if err := b.DB.Queries.CreateEvent(ctx, sqlc.CreateEventParams{
						Type:     eventType,
						Name:     name,
						Time:     unixValue,
						Hours:    hours,
						Gardener: gardener,
					}); err != nil {
						m.Client().Logger.Error("Failed to create event in database", slog.Any("err", err))
						return
					}

					replyText := "Event: " + eventType + " - " + name + "\n" +
						"Time: <t:" + m.Data.Text("event_time") + ":F> (<t:" + m.Data.Text("event_time") + ":R>)\n" +
						"Hours: " + m.Data.Text("event_duration") + " hours\n" +
						"Gardener: <@" + data.String("gardener") + ">"

					var banner *discord.Icon
					attachments, provided := m.Data.OptAttachments("event_banner")
					if provided && len(attachments) > 0 {
						banner = getBanner(attachments[0], m.Client().Logger)
					} else {
						banner = nil
					}

					if err := m.CreateMessage(discord.MessageCreate{
						Content: replyText,
					}); err != nil {
						m.Client().Logger.Error("Failed to send event message", slog.Any("err", err))
						return
					}

					switch m.Data.StringValues("event_type")[0] {
					case "Dota":
						if _, err := m.Client().Rest.CreateGuildScheduledEvent(*m.GuildID(), discord.GuildScheduledEventCreate{
							Name:               m.Data.StringValues("event_type")[0] + " - " + m.Data.Text("event_name"),
							EntityType:         discord.ScheduledEventEntityTypeVoice,
							PrivacyLevel:       discord.ScheduledEventPrivacyLevelGuildOnly,
							ChannelID:          channels.VoiceChannels["Dota"],
							ScheduledStartTime: time.Unix(unixValue, 0),
							Image:              banner,
						}); err != nil {
							m.Client().Logger.Error("Failed to create scheduled event", slog.Any("err", err))
						}
					case "CS":
						if _, err := m.Client().Rest.CreateGuildScheduledEvent(*m.GuildID(), discord.GuildScheduledEventCreate{
							Name:               m.Data.StringValues("event_type")[0] + " - " + m.Data.Text("event_name"),
							EntityType:         discord.ScheduledEventEntityTypeVoice,
							PrivacyLevel:       discord.ScheduledEventPrivacyLevelGuildOnly,
							ChannelID:          channels.VoiceChannels["CS"],
							ScheduledStartTime: time.Unix(unixValue, 0),
							Image:              banner,
						}); err != nil {
							m.Client().Logger.Error("Failed to create scheduled event", slog.Any("err", err))
						}
					case "MLBB":
						if _, err := m.Client().Rest.CreateGuildScheduledEvent(*m.GuildID(), discord.GuildScheduledEventCreate{
							Name:         m.Data.StringValues("event_type")[0] + " - " + m.Data.Text("event_name"),
							EntityType:   discord.ScheduledEventEntityTypeExternal,
							PrivacyLevel: discord.ScheduledEventPrivacyLevelGuildOnly,
							EntityMetaData: &discord.EntityMetaData{
								Location: channels.ExternalChannels["MLBB"],
							},
							ScheduledStartTime: time.Unix(unixValue, 0),
							Image:              banner,
						}); err != nil {
							m.Client().Logger.Error("Failed to create scheduled event", slog.Any("err", err))
						}
					case "HoK":
						if _, err := m.Client().Rest.CreateGuildScheduledEvent(*m.GuildID(), discord.GuildScheduledEventCreate{
							Name:         m.Data.StringValues("event_type")[0] + " - " + m.Data.Text("event_name"),
							EntityType:   discord.ScheduledEventEntityTypeExternal,
							PrivacyLevel: discord.ScheduledEventPrivacyLevelGuildOnly,
							EntityMetaData: &discord.EntityMetaData{
								Location: channels.ExternalChannels["HoK"],
							},
							ScheduledStartTime: time.Unix(unixValue, 0),
							Image:              banner,
						}); err != nil {
							m.Client().Logger.Error("Failed to create scheduled event", slog.Any("err", err))
						}
					case "Others":
						if _, err := m.Client().Rest.CreateGuildScheduledEvent(*m.GuildID(), discord.GuildScheduledEventCreate{
							Name:               m.Data.StringValues("event_type")[0] + " - " + m.Data.Text("event_name"),
							EntityType:         discord.ScheduledEventEntityTypeStageInstance,
							PrivacyLevel:       discord.ScheduledEventPrivacyLevelGuildOnly,
							ChannelID:          channels.StageChannel,
							ScheduledStartTime: time.Unix(unixValue, 0),
							Image:              banner,
						}); err != nil {
							m.Client().Logger.Error("Failed to create scheduled event", slog.Any("err", err))
						}
					default:
						panic("Unknown event type")
					}

				},
				func() {
					if err := e.CreateMessage(discord.MessageCreate{
						Content: "Modal timed out. Please try again.",
					}); err != nil {
						e.Client().Logger.Error("Failed to send modal timeout message", slog.Any("err", err))
					}
				},
			)
		}()
		return nil
	}
}
