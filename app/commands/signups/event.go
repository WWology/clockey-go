package signups

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"clockey/app"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/handler"
)

var Event = discord.SlashCommandCreate{
	Name:        "event",
	Description: "Create a new event for Gardeners to sign up for",
}

func EventCommandHandler(b *app.Bot) handler.SlashCommandHandler {
	return func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
		// Show modal to collect event details
		if err := e.Modal(discord.ModalCreate{
			CustomID: "event_modal",
			Title:    "Event Modal",
			Components: []discord.LayoutComponent{
				discord.LabelComponent{
					Label:       "Event Type",
					Description: "Select the type of event",
					Component: discord.StringSelectMenuComponent{
						CustomID: "event_type",
						Options:  OGGames,
						Required: true,
					},
				},
				discord.LabelComponent{
					Label:       "Event Name",
					Description: "Enter the name of the event",
					Component: discord.TextInputComponent{
						CustomID:    "event_name",
						Style:       discord.TextInputStyleShort,
						Placeholder: "OG vs <opp team name>",
						Required:    true,
					},
				},
				discord.LabelComponent{
					Label:       "Event Schedule",
					Description: "Enter the unix time for the start of this event",
					Component: discord.TextInputComponent{
						CustomID:    "event_time",
						Style:       discord.TextInputStyleShort,
						Required:    true,
						Placeholder: "Insert unix time from hammertime here",
					},
				},
				discord.LabelComponent{
					Label:       "Event duration",
					Description: "How many hours is this event",
					Component: discord.TextInputComponent{
						CustomID: "event_duration",
						Style:    discord.TextInputStyleShort,
						Required: true,
					},
				},
				discord.LabelComponent{
					Label:       "Event Banner",
					Description: "The banner for this event (if any, 800x320 px in size). ",
					Component: discord.FileUploadComponent{
						CustomID: "event_banner",
						Required: false,
					},
				},
			},
		}); err != nil {
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
						m.Client().Logger.Error("Failed to parse event_time", slog.Any("err", err))
						m.CreateMessage(discord.MessageCreate{
							Content: m.Data.Text("event_time") + " is not a valid unix time. Please try again.",
						})
						return
					}

					replyText := "Hey <@&" + GardenerRoleID.String() + ">\n\n" +
						"Event: " + m.Data.StringValues("event_type")[0] + " - " + m.Data.Text("event_name") + "\n" +
						"Time: <t:" + m.Data.Text("event_time") + ":F> (<t:" + m.Data.Text("event_time") + ":R>)\n" +
						"Hours: " + m.Data.Text("event_duration") + " hours\n" +
						"Please react with " + SignupEmojiString + " to sign up!."

					var banner *discord.Icon
					attachments, provided := m.Data.OptAttachments("event_banner")
					if provided && len(attachments) > 0 {
						m.Client().Logger.Info("Banner Image provided")
						banner = getBanner(attachments[0], m.Client().Logger)
					} else {
						banner = nil
					}

					switch m.Data.StringValues("event_type")[0] {
					case "Dota":
						if _, err := m.Client().Rest.CreateGuildScheduledEvent(*m.GuildID(), discord.GuildScheduledEventCreate{
							Name:               m.Data.StringValues("event_type")[0] + " - " + m.Data.Text("event_name"),
							EntityType:         discord.ScheduledEventEntityTypeVoice,
							PrivacyLevel:       discord.ScheduledEventPrivacyLevelGuildOnly,
							ChannelID:          Channels.VoiceChannels["Dota"],
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
							ChannelID:          Channels.VoiceChannels["CS"],
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
								Location: Channels.ExternalChannels["MLBB"],
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
								Location: Channels.ExternalChannels["HoK"],
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
							ChannelID:          Channels.StageChannel,
							ScheduledStartTime: time.Unix(unixValue, 0),
							Image:              banner,
						}); err != nil {
							m.Client().Logger.Error("Failed to create scheduled event", slog.Any("err", err))
						}
					}

					if err := m.CreateMessage(discord.MessageCreate{
						Content: replyText,
						AllowedMentions: &discord.AllowedMentions{
							Roles: discord.DefaultAllowedMentions.Roles,
							Users: discord.DefaultAllowedMentions.Users,
						},
					}); err != nil {
						m.Client().Logger.Error("Failed to send event message", slog.Any("err", err))
						return
					}

					msg, err := m.Client().Rest.GetInteractionResponse(m.ApplicationID(), m.Token())
					if err != nil {
						m.Client().Logger.Error("Failed to get interaction response", slog.Any("err", err))
						return
					}

					// Add reaction to the message
					if err := m.Client().Rest.AddReaction(msg.ChannelID, msg.ID, "OGPeepoYes"); err != nil {
						m.Client().Logger.Error("Failed to add reaction to event message", slog.Any("err", err))
						return
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

func getBanner(attachment discord.Attachment, Logger *slog.Logger) *discord.Icon {
	resp, err := http.Get(attachment.URL)
	if err != nil {
		Logger.Error("Failed to create banner icon", slog.Any("err", err))
		return nil
	}
	banner, err := discord.NewIcon(discord.IconTypeJPEG, resp.Body)
	if err != nil {
		Logger.Error("Failed to create banner icon", slog.Any("err", err))
		return nil
	}
	return banner
}
