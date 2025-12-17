package predictions

import (
	"context"
	"log/slog"
	"time"

	"clockey/app"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/omit"
)

var Reset = discord.SlashCommandCreate{
	Name:        "reset",
	Description: "Reset the monthly prediction leaderboard",
}

func ResetCommandHandler(b *app.Bot) handler.SlashCommandHandler {
	return func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
		buttons := discord.ActionRowComponent{
			Components: []discord.InteractiveComponent{
				discord.ButtonComponent{
					Label:    "Yes",
					Style:    discord.ButtonStyleDanger,
					CustomID: "reset_leaderboard_yes",
				},
				discord.ButtonComponent{
					Label:    "No",
					Style:    discord.ButtonStyleSecondary,
					CustomID: "reset_leaderboard_no",
				},
			},
		}

		if err := e.CreateMessage(discord.MessageCreate{
			Content: "Are you sure you want to reset the monthly prediction leaderboard?",
			Components: []discord.LayoutComponent{
				buttons,
			},
		}); err != nil {
			slog.Error("DisGo error(failed to create message)", slog.Any("err", err))
			return err
		}

		go func() {
			ch, cls := bot.NewEventCollector(e.Client(),
				func(c *events.ComponentInteractionCreate) bool {
					return c.Data.CustomID() == "reset_leaderboard_yes" || c.Data.CustomID() == "reset_leaderboard_no"
				})
			defer cls()
			ctx, cancel := context.WithTimeout(context.TODO(), 2*time.Minute)
			defer cancel()
			select {
			case <-ctx.Done():
				return
			case c := <-ch:
				if c.Data.CustomID() == "reset_leaderboard_yes" {
					ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
					defer cancel()
					if err := b.DB.Queries.ClearScoreboard(ctx); err != nil {
						if err := c.UpdateMessage(discord.MessageUpdate{
							Content: omit.Ptr("Failed to reset scoreboard, please try again"),
						}); err != nil {
							slog.Error("DisGo error(failed to update message)", slog.Any("err", err))
						}
					}
					if err := c.UpdateMessage(discord.MessageUpdate{
						Content: omit.Ptr("Monthly prediction leaderboard reset successfully"),
					}); err != nil {
						slog.Error("DisGo error(failed to update message)", slog.Any("err", err))
					}
				} else if c.Data.CustomID() == "reset_leaderboard_no" {
					if err := c.UpdateMessage(discord.MessageUpdate{
						Content: omit.Ptr("Monthly prediction leaderboard reset cancelled"),
					}); err != nil {
						slog.Error("DisGo error(failed to update message)", slog.Any("err", err))
					}
					return
				} else {
					return
				}
			}
		}()
		return nil
	}
}
