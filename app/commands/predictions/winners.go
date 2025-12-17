package predictions

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"clockey/app"
	"clockey/database/sqlc"

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
		if err := e.DeferCreateMessage(true); err != nil {
			slog.Error("DisGo error(failed to defer interaction response)", slog.Any("err", err))
			return err
		}

		// Remove previous winners
		guildMembers := e.Client().Caches.Members(*e.GuildID())
		for member := range guildMembers {
			for _, roleID := range member.RoleIDs {
				if roleID == theOracleRoleID || roleID == dotaOracleRoleID || roleID == csOracleRoleID || roleID == mlbbOracleRoleID || roleID == hokOracleRoleID {
					// Remove previous winner roles
					if err := e.Client().Rest.RemoveMemberRole(*e.GuildID(), member.User.ID, roleID); err != nil {
						slog.Error("DisGo error(failed to remove member role)", slog.Any("err", err))
						return err
					}
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

		dotaWinners, err := b.DB.Queries.GetWinnerForGame(ctx, sqlc.ScoreboardGameDota)
		if err != nil {
			return err
		}

		csWinners, err := b.DB.Queries.GetWinnerForGame(ctx, sqlc.ScoreboardGameCS)
		if err != nil {
			return err
		}

		mlbbWinners, err := b.DB.Queries.GetWinnerForGame(ctx, sqlc.ScoreboardGameMLBB)
		if err != nil {
			return err
		}

		hokWinners, err := b.DB.Queries.GetWinnerForGame(ctx, sqlc.ScoreboardGameHoK)
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
			if err := e.Client().Rest.AddMemberRole(*e.GuildID(), snowflake.ID(winner.Member), theOracleRoleID); err != nil {
				slog.Error("DisGo error(failed to add member role)", slog.Any("err", err))
				return err
			}
		}

		for i, winner := range dotaWinners {
			if i == 0 {
				replyText += fmt.Sprintf("\n\nDota: <@%d>", winner.Member)
			} else {
				replyText += fmt.Sprintf(", <@%d>", winner.Member)
			}
			if err := e.Client().Rest.AddMemberRole(*e.GuildID(), snowflake.ID(winner.Member), dotaOracleRoleID); err != nil {
				slog.Error("DisGo error(failed to add member role)", slog.Any("err", err))
				return err
			}
		}

		for i, winner := range csWinners {
			if i == 0 {
				replyText += fmt.Sprintf("\n\nCS: <@%d>", winner.Member)
			} else {
				replyText += fmt.Sprintf(", <@%d>", winner.Member)
			}
			if err := e.Client().Rest.AddMemberRole(*e.GuildID(), snowflake.ID(winner.Member), csOracleRoleID); err != nil {
				slog.Error("DisGo error(failed to add member role)", slog.Any("err", err))
				return err
			}
		}

		for i, winner := range mlbbWinners {
			if i == 0 {
				replyText += fmt.Sprintf("\n\nMLBB: <@%d>", winner.Member)
			} else {
				replyText += fmt.Sprintf(", <@%d>", winner.Member)
			}
			if err := e.Client().Rest.AddMemberRole(*e.GuildID(), snowflake.ID(winner.Member), mlbbOracleRoleID); err != nil {
				slog.Error("DisGo error(failed to add member role)", slog.Any("err", err))
				return err
			}
		}

		for i, winner := range hokWinners {
			if i == 0 {
				replyText += fmt.Sprintf("\n\nHoK: <@%d>", winner.Member)
			} else {
				replyText += fmt.Sprintf(", <@%d>", winner.Member)
			}
			if err := e.Client().Rest.AddMemberRole(*e.GuildID(), snowflake.ID(winner.Member), hokOracleRoleID); err != nil {
				slog.Error("DisGo error(failed to add member role)", slog.Any("err", err))
				return err
			}
		}

		if _, err := e.UpdateInteractionResponse(discord.MessageUpdate{
			Content: omit.Ptr(replyText),
		}); err != nil {
			slog.Error("DisGo error(failed to update interaction response)", slog.Any("err", err))
			return err
		}

		return nil
	}
}
