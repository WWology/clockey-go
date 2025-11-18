package predictions

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var BestOf = discord.SlashCommandCreate{
	Name:        "bo",
	Description: "Create prediction roles for selected game",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionString{
			Name:        "game",
			Description: "The game to create prediction roles for",
			Choices: []discord.ApplicationCommandOptionChoiceString{
				{
					Name:  "Dota",
					Value: "D",
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
		discord.ApplicationCommandOptionInt{
			Name:        "series_length",
			Description: "The length of the series",
			Choices: []discord.ApplicationCommandOptionChoiceInt{
				{
					Name:  "Bo1",
					Value: 1,
				},
				{
					Name:  "Bo2",
					Value: 2,
				},
				{
					Name:  "Bo3",
					Value: 3,
				},
				{
					Name:  "Bo5",
					Value: 5,
				},
				{
					Name:  "Bo7",
					Value: 7,
				},
			},
		},
	},
}

func BestOfCommandHandler() handler.SlashCommandHandler {
	return func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
		game := data.String("game")
		seriesLength := data.Int("series_length")

		switch seriesLength {
		case 1:
			rolesToBeCreated := []string{
				game + "1-0",
				game + "0-1",
			}
			for _, roleName := range rolesToBeCreated {
				if _, err := e.Client().Rest.CreateRole(*e.GuildID(), discord.RoleCreate{Name: roleName}); err != nil {
					return err
				}
			}
		case 2:
			rolesToBeCreated := []string{
				game + "2-0",
				game + "1-1",
				game + "0-2",
			}
			for _, roleName := range rolesToBeCreated {
				if _, err := e.Client().Rest.CreateRole(*e.GuildID(), discord.RoleCreate{Name: roleName}); err != nil {
					return err
				}
			}
		case 3:
			rolesToBeCreated := []string{
				game + "2-0",
				game + "2-1",
				game + "1-2",
				game + "0-2",
			}
			for _, roleName := range rolesToBeCreated {
				if _, err := e.Client().Rest.CreateRole(*e.GuildID(), discord.RoleCreate{Name: roleName}); err != nil {
					return err
				}
			}
		case 5:
			rolesToBeCreated := []string{
				game + "3-0",
				game + "3-1",
				game + "3-2",
				game + "2-3",
				game + "1-3",
				game + "0-3",
			}
			for _, roleName := range rolesToBeCreated {
				if _, err := e.Client().Rest.CreateRole(*e.GuildID(), discord.RoleCreate{Name: roleName}); err != nil {
					return err
				}
			}
		case 7:
			rolesToBeCreated := []string{
				game + "4-0",
				game + "4-1",
				game + "4-2",
				game + "4-3",
				game + "3-4",
				game + "2-4",
				game + "1-4",
				game + "0-4",
			}
			for _, roleName := range rolesToBeCreated {
				if _, err := e.Client().Rest.CreateRole(*e.GuildID(), discord.RoleCreate{Name: roleName}); err != nil {
					return err
				}
			}
		}

		return e.CreateMessage(discord.MessageCreate{
			Content: "Prediction roles created for " + game + " best of " + fmt.Sprint(seriesLength),
		})
	}
}
