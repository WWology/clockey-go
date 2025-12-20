package signups

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/handler"
)

var Event = discord.SlashCommandCreate{
	Name:        "event",
	Description: "Create a new event for Gardeners to sign up for",
}

func EventCommandHandler() handler.SlashCommandHandler {
	return func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
		// Show modal to collect event details
		if err := e.Modal(eventModal); err != nil {
			slog.Error("DisGo error(failed to send modal)", slog.Any("err", err))
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
					unixValue, err := strconv.ParseInt(m.Data.Text("event_time"), 0, 64)
					if err != nil {
						slog.Error("failed to parse event_time", slog.String("event_time", m.Data.Text("event_time")), slog.Any("err", err))
						if err := m.CreateMessage(discord.MessageCreate{
							Content: m.Data.Text("event_time") + " is not a valid unix time. Please try again.",
						}); err != nil {
							slog.Error("DisGo error(failed to send message)", slog.Any("err", err))
						}
						return
					}

					replyText := "Hey <@&" + gardenerRoleID.String() + ">\n\n" +
						"Event: " + m.Data.StringValues("event_type")[0] + " - " + m.Data.Text("event_name") + "\n" +
						"Time: <t:" + m.Data.Text("event_time") + ":F> (<t:" + m.Data.Text("event_time") + ":R>)\n" +
						"Hours: " + m.Data.Text("event_duration") + " hours\n" +
						"Please react with <:" + signupEmoji + "> to sign up!."

					var banner *discord.Icon
					attachments, provided := m.Data.OptAttachments("event_banner")
					if provided && len(attachments) > 0 {
						banner = getBanner(attachments[0])
					} else {
						banner = nil
					}

					if err := m.CreateMessage(discord.MessageCreate{
						Content: replyText,
						AllowedMentions: &discord.AllowedMentions{
							Parse: []discord.AllowedMentionType{
								discord.AllowedMentionTypeRoles,
								discord.AllowedMentionTypeUsers,
							},
						},
					}); err != nil {
						slog.Error("DisGo error(failed to send message)", slog.Any("err", err))
						return
					}

					msg, err := m.Client().Rest.GetInteractionResponse(m.ApplicationID(), m.Token())
					if err != nil {
						slog.Error("DisGo error(failed to get interaction response)", slog.Any("err", err))
					}

					// Add reaction to the message
					if err := m.Client().Rest.AddReaction(msg.ChannelID, msg.ID, signupEmoji); err != nil {
						slog.Error("DisGo error(failed to add reaction to event message)", slog.Any("err", err))
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
							slog.Error("DisGo error(failed to create scheduled event)", slog.Any("err", err))
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
							slog.Error("DisGo error(failed to create scheduled event)", slog.Any("err", err))
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
							slog.Error("DisGo error(failed to create scheduled event)", slog.Any("err", err))
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
							slog.Error("DisGo error(failed to create scheduled event)", slog.Any("err", err))
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
							slog.Error("DisGo error(failed to create scheduled event)", slog.Any("err", err))
						}
					default:
						panic(fmt.Sprintf("invalid state: %s", m.Data.StringValues("event_type")[0]))
					}
				},
				func() {
					if err := e.CreateMessage(discord.MessageCreate{
						Content: "Modal timed out. Please try again.",
					}); err != nil {
						slog.Error("DisGo error(failed to send modal timeout message)", slog.Any("err", err))
					}
				},
			)
		}()
		return nil
	}
}
