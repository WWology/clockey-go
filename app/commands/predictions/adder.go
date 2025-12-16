package predictions

import (
	"clockey/app"
	"clockey/database/sqlc"
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/omit"
)

var Add = discord.SlashCommandCreate{
	Name:        "add",
	Description: "Add a prediction score to the chosen scoreboard",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionString{
			Name:        "game",
			Description: "The game leaderboard to update",
			Required:    true,
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
		discord.ApplicationCommandOptionRole{
			Name:        "role",
			Description: "The role to add score to",
		},
	},
}

func AddCommandHandler(b *app.Bot) handler.SlashCommandHandler {
	return func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
		e.DeferCreateMessage(false)
		guildMembers := e.Client().Caches.Members(*e.GuildID())

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		tx, err := b.DB.Conn.Begin(ctx)
		if err != nil {
			return err
		}
		defer tx.Rollback(ctx)
		count := 0
		for member := range guildMembers {
			if slices.Contains(member.RoleIDs, data.Role("role").ID) {
				if err := b.DB.Queries.WithTx(tx).UpdateScoreboardForGame(ctx, sqlc.UpdateScoreboardForGameParams{
					Member: int64(member.User.ID),
					Game:   sqlc.ScoreboardGame(data.String("game")),
				}); err != nil {
					return err
				}
				count++
			}
		}
		e.UpdateInteractionResponse(discord.MessageUpdate{
			Content: omit.Ptr(fmt.Sprintf("Added score for %d members to the %s scoreboard", count, data.String("game"))),
		})
		return tx.Commit(ctx)
	}
}
