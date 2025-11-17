package signups

import (
	"clockey/app"
	"context"
	"time"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/handler"
)

var Cancel = discord.MessageCommandCreate{
	Name: "Cancel Signup",
	Contexts: []discord.InteractionContextType{
		discord.InteractionContextTypeGuild,
	},
}

func CancelCommandHandler(b app.Bot) handler.MessageCommandHandler {
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
			return err
		}

		go func() {
			ch, cls := bot.NewEventCollector(e.Client(),
				func(c *events.ComponentInteractionCreate) bool {
					return c.Data.CustomID() == "cancel_event_yes" || c.Data.CustomID() == "cancel_event_no"
				})
			defer cls()
			ctx, cancel := context.WithTimeout(context.TODO(), 2*time.Minute)
			defer cancel()
			select {
			case <-ctx.Done():
				return
			case c := <-ch:
				if c.Data.CustomID() == "cancel_event_yes" {
					panic("todo: implement cancel event")
				} else {
					return
				}
			}
		}()

		return nil
	}
}
