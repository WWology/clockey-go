package signups

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"clockey/app"
	"clockey/database/sqlc"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/omit"
)

var Cancel = discord.MessageCommandCreate{
	Name: "Cancel Event",
	Contexts: []discord.InteractionContextType{
		discord.InteractionContextTypeGuild,
	},
}

func CancelCommandHandler(b *app.Bot) handler.MessageCommandHandler {
	return func(data discord.MessageCommandInteractionData, e *handler.CommandEvent) error {
		if !processed(data.TargetMessage()) {
			return e.CreateMessage(discord.MessageCreate{
				Content: "This message has not been processed for signups yet",
				Flags:   discord.MessageFlagEphemeral,
			})
		}

		buttons := discord.ActionRowComponent{
			Components: []discord.InteractiveComponent{
				discord.ButtonComponent{
					Label:    "Yes",
					Style:    discord.ButtonStyleDanger,
					CustomID: "cancel_event_yes",
				},
				discord.ButtonComponent{
					Label:    "No",
					Style:    discord.ButtonStyleSecondary,
					CustomID: "cancel_event_no",
				},
			},
		}

		if err := e.CreateMessage(discord.MessageCreate{
			Content: "Are you sure you want to cancel signups for this event?",
			Components: []discord.LayoutComponent{
				buttons,
			},
		}); err != nil {
			slog.Error("DisGo error(failed to create cancel confirmation message)", slog.Any("err", err))
			return err
		}

		go func() {
			ch, cls := bot.NewEventCollector(e.Client(),
				func(c *events.ComponentInteractionCreate) bool {
					return c.Data.CustomID() == "cancel_event_yes" || c.Data.CustomID() == "cancel_event_no"
				})
			defer cls()
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			defer cancel()
			select {
			case <-ctx.Done():
				return
			case c := <-ch:
				if c.Data.CustomID() == "cancel_event_yes" {
					eventType, name, eventTime, hours, err := parseMessage(data.TargetMessage().Content)
					if err != nil {
						slog.Error("failed to parse message", slog.Any("err", err))
						return
					}

					if err := b.DB.Queries.DeleteEvent(ctx, sqlc.DeleteEventParams{
						Type:  eventType,
						Name:  name,
						Time:  eventTime,
						Hours: hours,
					}); err != nil {
						if err := c.CreateMessage(discord.MessageCreate{
							Content: "Error cancelling event, please try again",
						}); err != nil {
							slog.Error("DisGo error(failed to create error cancelling event message)", slog.Any("err", err))
						}
						return
					}
					if _, err := e.UpdateInteractionResponse(discord.MessageUpdate{
						Content:    omit.Ptr(fmt.Sprintf("%s - %s cancelled", eventType, name)),
						Components: &[]discord.LayoutComponent{},
					}); err != nil {
						slog.Error("DisGo error(failed to update interaction response)", slog.Any("err", err))
					}

					if err := c.Client().Rest.RemoveOwnReaction(c.Channel().ID(), data.TargetID(), processedEmoji); err != nil {
						slog.Error("DisGo error(failed to remove own reaction)", slog.Any("err", err))
					}
				} else if c.Data.CustomID() == "cancel_event_no" {
					if err := c.UpdateMessage(discord.MessageUpdate{
						Content:    omit.Ptr("Event cancellation aborted"),
						Components: &[]discord.LayoutComponent{},
					}); err != nil {
						slog.Error("DisGo error(failed to update message)", slog.Any("err", err))
					}
				} else {
					return
				}
			}
		}()

		return nil
	}
}
