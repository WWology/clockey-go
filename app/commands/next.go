package commands

import (
	"github.com/disgoorg/disgo/discord"
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
		},
	},
}
