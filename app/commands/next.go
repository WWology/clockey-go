package commands

import (
	"cmp"
	"fmt"
	"slices"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var Next = discord.SlashCommandCreate{
	Name:        "next",
	Description: "Next game for OG",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionString{
			Name:        "game",
			Description: "Which game do you want to know",
			Choices: []discord.ApplicationCommandOptionChoiceString{
				{
					Name:  "Dota",
					Value: "Dota",
				},
				{
					Name:  "CS",
					Value: "CS",
				},
				{
					Name:  "MLBB",
					Value: "MLBB",
				},
				{
					Name:  "HoK",
					Value: "HoK",
				},
			},
			Required: true,
		},
	},
}

func NextCommandHandler() handler.SlashCommandHandler {
	return func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
		game := data.String("game")
		eventList, err := e.Client().Rest.GetGuildScheduledEvents(*e.GuildID(), false)
		if err != nil {
			return err
		}

		sortedEvents := slices.SortedStableFunc(slices.Values(eventList), func(a, b discord.GuildScheduledEvent) int {
			return cmp.Compare(a.ScheduledStartTime.Unix(), b.ScheduledStartTime.Unix())
		})

		result := slices.IndexFunc(sortedEvents, func(event discord.GuildScheduledEvent) bool {
			return strings.Contains(event.Name, game)
		})

		if result != -1 {
			return e.CreateMessage(discord.MessageCreate{
				Content: fmt.Sprintf("https://discord.com/events/%s/%s", e.GuildID(), sortedEvents[result].ID),
			})
		} else {
			return e.CreateMessage(discord.MessageCreate{
				Content: "No upcoming game found for " + game,
			})
		}
	}
}
