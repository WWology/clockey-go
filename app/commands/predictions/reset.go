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
						c.CreateMessage(discord.MessageCreate{
							Content: "Failed to reset scoreboard, please try again",
						})
					}
					c.CreateMessage(discord.MessageCreate{
						Content: "Monthly prediction leaderboard reset successfully",
					})
				} else {
					return
				}
			}
		}()
		return nil
	}
}
