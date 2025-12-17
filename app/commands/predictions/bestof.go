package predictions

import (
	"fmt"
	"log/slog"
	"slices"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/omit"
)

var BestOf = discord.SlashCommandCreate{
	Name:        "bo",
	Description: "Create prediction roles for selected game",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionString{
			Name:         "game",
			Description:  "The game to create prediction roles for",
			Required:     true,
			Autocomplete: true,
		},
		discord.ApplicationCommandOptionInt{
			Name:        "series_length",
			Description: "The length of the series",
			Required:    true,
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
		if err := e.DeferCreateMessage(true); err != nil {
			slog.Error("DisGo error(failed to defer interaction response)", slog.Any("err", err))
			return err
		}

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

		if _, err := e.UpdateInteractionResponse(discord.MessageUpdate{
			Content: omit.Ptr("Prediction roles created for " + game + " best of " + fmt.Sprint(seriesLength)),
		}); err != nil {
			return err
		}
		return nil
	}
}

var DeleteBestOf = discord.SlashCommandCreate{
	Name:        "deletebo",
	Description: "Delete all prediction roles for selected game",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionString{
			Name:         "game",
			Description:  "The game to delete prediction roles for",
			Required:     true,
			Autocomplete: true,
		},
		discord.ApplicationCommandOptionInt{
			Name:        "series_length",
			Description: "The length of the series",
			Required:    true,
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

func DeleteBestOfCommandHandler() handler.SlashCommandHandler {
	return func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
		if err := e.DeferCreateMessage(true); err != nil {
			slog.Error("DisGo error(failed to defer interaction response)", slog.Any("err", err))
			return err
		}

		game := data.String("game")
		seriesLength := data.Int("series_length")
		switch seriesLength {
		case 1:
			rolesToBeDeleted := []string{
				game + "1-0",
				game + "0-1",
			}
			if roles, err := e.Client().Rest.GetRoles(*e.GuildID()); err == nil {
				for _, role := range roles {
					if slices.Contains(rolesToBeDeleted, role.Name) {
						e.Client().Rest.DeleteRole(*e.GuildID(), role.ID)
					}
				}
			} else {
				return err
			}
		case 2:
			rolesToBeDeleted := []string{
				game + "2-0",
				game + "1-1",
				game + "0-2",
			}
			if roles, err := e.Client().Rest.GetRoles(*e.GuildID()); err == nil {
				for _, role := range roles {
					if slices.Contains(rolesToBeDeleted, role.Name) {
						e.Client().Rest.DeleteRole(*e.GuildID(), role.ID)
					}
				}
			} else {
				return err
			}
		case 3:
			rolesToBeDeleted := []string{
				game + "2-0",
				game + "2-1",
				game + "1-2",
				game + "0-2",
			}
			if roles, err := e.Client().Rest.GetRoles(*e.GuildID()); err == nil {
				for _, role := range roles {
					if slices.Contains(rolesToBeDeleted, role.Name) {
						e.Client().Rest.DeleteRole(*e.GuildID(), role.ID)
					}
				}
			} else {
				return err
			}
		case 5:
			rolesToBeDeleted := []string{
				game + "3-0",
				game + "3-1",
				game + "3-2",
				game + "2-3",
				game + "1-3",
				game + "0-3",
			}
			if roles, err := e.Client().Rest.GetRoles(*e.GuildID()); err == nil {
				for _, role := range roles {
					if slices.Contains(rolesToBeDeleted, role.Name) {
						e.Client().Rest.DeleteRole(*e.GuildID(), role.ID)
					}
				}
			} else {
				return err
			}
		case 7:
			rolesToBeDeleted := []string{
				game + "4-0",
				game + "4-1",
				game + "4-2",
				game + "4-3",
				game + "3-4",
				game + "2-4",
				game + "1-4",
				game + "0-4",
			}
			if roles, err := e.Client().Rest.GetRoles(*e.GuildID()); err == nil {
				for _, role := range roles {
					if slices.Contains(rolesToBeDeleted, role.Name) {
						e.Client().Rest.DeleteRole(*e.GuildID(), role.ID)
					}
				}
			} else {
				return err
			}
		}
		if _, err := e.UpdateInteractionResponse(discord.MessageUpdate{
			Content: omit.Ptr("Prediction roles deleted for " + game + " best of " + fmt.Sprint(seriesLength)),
		}); err != nil {
			return err
		}
		return nil
	}
}

func BestOfAutocompleteHandler() handler.AutocompleteHandler {
	return func(e *handler.AutocompleteEvent) error {
		if len(e.Data.Focused().String()) > 0 {
			return e.AutocompleteResult([]discord.AutocompleteChoice{
				discord.AutocompleteChoiceString{
					Name:  "Dota",
					Value: "Dota",
				},
				discord.AutocompleteChoiceString{
					Name:  "CS",
					Value: "CS",
				},
				discord.AutocompleteChoiceString{
					Name:  "MLBB",
					Value: "MLBB",
				},
				discord.AutocompleteChoiceString{
					Name:  "HoK",
					Value: "HoK",
				},
				discord.AutocompleteChoiceString{
					Name:  "Extra",
					Value: "EX",
				},
				discord.AutocompleteChoiceString{
					Name:  e.Data.Focused().String(),
					Value: e.Data.Focused().String(),
				},
			})
		} else {
			return e.AutocompleteResult([]discord.AutocompleteChoice{})
		}
	}
}
