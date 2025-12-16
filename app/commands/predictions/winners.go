package predictions

import (
	"clockey/app"
	"clockey/database/sqlc"
	"context"
	"fmt"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/omit"
	"github.com/disgoorg/snowflake/v2"
)

var Winners = discord.SlashCommandCreate{
	Name:        "winners",
	Description: "Give prediction winners their roles and remove previous winners",
}

func WinnersCommandHandler(b *app.Bot) handler.SlashCommandHandler {
	return func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
		e.DeferCreateMessage(false)

		// Remove previous winners
		guildMembers := e.Client().Caches.Members(*e.GuildID())
		for member := range guildMembers {
			for _, roleID := range member.RoleIDs {
				if roleID == theOracleRoleID || roleID == dotaOracleRoleID || roleID == csOracleRoleID || roleID == mlbbOracleRoleID || roleID == hokOracleRoleID {
					// Remove previous winner roles
					e.Client().Rest.RemoveMemberRole(*e.GuildID(), member.User.ID, roleID)
				}
			}
		}

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		// Get winners from database
		globalWinners, err := b.DB.Queries.GetGlobalWinner(ctx)
		if err != nil {
			return err
		}

		dotaWinners, err := b.DB.Queries.GetWinnerForGame(ctx, sqlc.ScoreboardGame("Dota"))
		if err != nil {
			return err
		}

		csWinners, err := b.DB.Queries.GetWinnerForGame(ctx, sqlc.ScoreboardGame("CS"))
		if err != nil {
			return err
		}

		mlbbWinners, err := b.DB.Queries.GetWinnerForGame(ctx, sqlc.ScoreboardGame("MLBB"))
		if err != nil {
			return err
		}

		hokWinners, err := b.DB.Queries.GetWinnerForGame(ctx, sqlc.ScoreboardGame("HoK"))
		if err != nil {
			return err
		}

		replyText := "THE ORACLE: "
		// Assign roles to winners
		for i, winner := range globalWinners {
			if i == 0 {
				replyText += fmt.Sprintf("<@%d>", winner.Member)
			} else {
				replyText += fmt.Sprintf(", <@%d>", winner.Member)
			}
			e.Client().Rest.AddMemberRole(*e.GuildID(), snowflake.ID(winner.Member), theOracleRoleID)
		}

		for i, winner := range dotaWinners {
			if i == 0 {
				replyText += fmt.Sprintf("\n\nDota: <@%d>", winner.Member)
			} else {
				replyText += fmt.Sprintf(", <@%d>", winner.Member)
			}
			e.Client().Rest.AddMemberRole(*e.GuildID(), snowflake.ID(winner.Member), dotaOracleRoleID)
		}

		for i, winner := range csWinners {
			if i == 0 {
				replyText += fmt.Sprintf("\n\nCS: <@%d>", winner.Member)
			} else {
				replyText += fmt.Sprintf(", <@%d>", winner.Member)
			}
			e.Client().Rest.AddMemberRole(*e.GuildID(), snowflake.ID(winner.Member), csOracleRoleID)
		}

		for i, winner := range mlbbWinners {
			if i == 0 {
				replyText += fmt.Sprintf("\n\nMLBB: <@%d>", winner.Member)
			} else {
				replyText += fmt.Sprintf(", <@%d>", winner.Member)
			}
			e.Client().Rest.AddMemberRole(*e.GuildID(), snowflake.ID(winner.Member), mlbbOracleRoleID)
		}

		for i, winner := range hokWinners {
			if i == 0 {
				replyText += fmt.Sprintf("\n\nHoK: <@%d>", winner.Member)
			} else {
				replyText += fmt.Sprintf(", <@%d>", winner.Member)
			}
			e.Client().Rest.AddMemberRole(*e.GuildID(), snowflake.ID(winner.Member), hokOracleRoleID)
		}

		e.UpdateInteractionResponse(discord.MessageUpdate{
			Content: omit.Ptr(replyText),
		})

		return nil
	}
}
