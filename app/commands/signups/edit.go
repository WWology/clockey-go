package signups

import (
	"regexp"
	"strings"

	"clockey/app"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/snowflake/v2"
)

var Edit = discord.SlashCommandCreate{
	Name:        "edit",
	Description: "Edit an existing event for Gardeners to sign up for",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionString{
			Name:        "message_id",
			Description: "The ID of the event message to edit",
			Required:    true,
		},
		discord.ApplicationCommandOptionString{
			Name:        "new_name",
			Description: "The new name for the event",
			Required:    false,
		},
		discord.ApplicationCommandOptionString{
			Name:        "new_time",
			Description: "The new unix time for the event",
			Required:    false,
		},
		discord.ApplicationCommandOptionString{
			Name:        "new_duration",
			Description: "The new duration (in hours) for the event",
			Required:    false,
		},
	},
}

func EditCommandHandler(b *app.Bot) handler.SlashCommandHandler {
	return func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
		msg, err := e.Client().Rest.GetMessage(e.Channel().ID(), snowflake.MustParse(data.String("message_id")))
		if err != nil {
			return err
		}
		replyText := "Updated event details: \n"

		if newName, provided := data.OptString("new_name"); provided {
			nameRegex := regexp.MustCompile(`Event: \w+ - (.+?)(?:\n|$)`)
			nameMatch := nameRegex.FindStringSubmatch(msg.Content)
			if len(nameMatch) > 1 {
				oldName := nameMatch[1]
				msg.Content = strings.Replace(msg.Content, oldName, newName, 1)
				replyText += oldName + " -> " + newName + "\n"
			}
		}

		if newTime, provided := data.OptString("new_time"); provided {
			timeRegex := regexp.MustCompile(`<t:([^:]+):F>`)
			timeMatch := timeRegex.FindStringSubmatch(msg.Content)
			if len(timeMatch) > 1 {
				oldTime := timeMatch[1]
				msg.Content = strings.Replace(msg.Content, oldTime, newTime, 2)
				replyText += "<t:" + oldTime + ":F> -> <t:" + newTime + ":F>\n"
			}
		}

		if newDuration, provided := data.OptString("new_duration"); provided {
			hoursRegex := regexp.MustCompile(`Hours: (\d+) hours`)
			hoursMatch := hoursRegex.FindStringSubmatch(msg.Content)
			if len(hoursMatch) > 1 {
				oldDuration := hoursMatch[1]
				msg.Content = hoursRegex.ReplaceAllString(msg.Content, "Hours: "+newDuration+" hours")
				replyText += oldDuration + " -> " + newDuration + " hours\n"
			}
		}

		e.Client().Rest.UpdateMessage(e.Channel().ID(), snowflake.MustParse(data.String("message_id")), discord.MessageUpdate{
			Content: &msg.Content,
		})

		e.CreateMessage(discord.MessageCreate{
			Content: replyText,
		})

		return nil
	}
}
