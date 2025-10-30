package signups

import (
	"clockey/app"
	"clockey/app/constants"
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/snowflake/v2"
)

var Event = discord.SlashCommandCreate{
	Name:        "event",
	Description: "Create a new event for Gardeners to sign up for",
}

func EventCommandHandler(b *app.Bot) handler.SlashCommandHandler {
	return func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
		if err := e.Modal(discord.ModalCreate{
			CustomID: "event_modal",
			Title:    "Event Modal",
			Components: []discord.LayoutComponent{
				discord.LabelComponent{
					Label:       "Event Type",
					Description: "Select the type of event",
					Component: discord.StringSelectMenuComponent{
						CustomID: "event_type",
						Options:  constants.OGGames,
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
					}

					replyText := "Hey <@&720253636797530203>\n\n" +
						"Event: " + m.Data.StringValues("event_type")[0] + " - " + m.Data.Text("event_name") + "\n" +
						"Time: <t:" + m.Data.Text("event_time") + ":F> (<t:" + m.Data.Text("event_time") + ":R>)\n" +
						"Hours: " + m.Data.Text("event_duration") + " hours\n" +
						"Please react with <:OGpeepoYes:730890894814740541> to sign up!."

					var banner *discord.Icon
					attachments, provided := m.Data.OptAttachments("event_banner")
					if provided && len(attachments) > 0 {
						m.Client().Logger.Info("Banner Image provided")
						resp, err := http.Get(attachments[0].URL)
						if err != nil {
							m.Client().Logger.Error("Failed to fetch banner image", slog.Any("err", err))
							m.CreateMessage(discord.MessageCreate{
								Content: "Failed to fetch banner image, please try again",
							})
						}
						banner, err = discord.NewIcon(discord.IconTypeJPEG, resp.Body)
						if err != nil {
							m.Client().Logger.Error("Failed to create icon from banner image", slog.Any("err", err))
							m.CreateMessage(discord.MessageCreate{
								Content: "Failed to process banner image, please try again",
							})
						}
					} else {
						banner = nil
					}

					switch m.Data.StringValues("event_type")[0] {
					case "Dota":
						if _, err := m.Client().Rest.CreateGuildScheduledEvent(*m.GuildID(), discord.GuildScheduledEventCreate{
							Name:               m.Data.StringValues("event_type")[0] + " - " + m.Data.Text("event_name"),
							EntityType:         discord.ScheduledEventEntityTypeVoice,
							PrivacyLevel:       discord.ScheduledEventPrivacyLevelGuildOnly,
							ChannelID:          constants.Channels["Dota"].(snowflake.ID),
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
							ChannelID:          constants.Channels["CS"].(snowflake.ID),
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
								Location: constants.Channels["MLBB"].(string),
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
								Location: constants.Channels["HoK"].(string),
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
							ChannelID:          constants.Channels["OGStage"].(snowflake.ID),
							ScheduledStartTime: time.Unix(unixValue, 0),
							Image:              banner,
						}); err != nil {
							m.Client().Logger.Error("Failed to create scheduled event", slog.Any("err", err))
						}
					}

					m.CreateMessage(discord.MessageCreate{
						Content: replyText,
						AllowedMentions: &discord.AllowedMentions{
							Roles: discord.DefaultAllowedMentions.Roles,
							Users: discord.DefaultAllowedMentions.Users,
						},
					})

					msg, err := m.Client().Rest.GetInteractionResponse(m.ApplicationID(), m.Token())
					if err != nil {
						m.Client().Logger.Error("Failed to get interaction response", slog.Any("err", err))
					}

					// Add reaction to the message
					m.Client().Rest.AddReaction(msg.ChannelID, msg.ID, "OGPeepoYes")
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
