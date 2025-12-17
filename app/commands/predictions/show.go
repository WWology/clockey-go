package predictions

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"time"

	"clockey/app"
	"clockey/database/sqlc"
	"database/sql"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/omit"
	"github.com/disgoorg/snowflake/v2"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
)

var Show = discord.SlashCommandCreate{
	Name:        "show",
	Description: "Show the current prediction leaderboard or individual user score",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionString{
			Name:        "game",
			Description: "The game to show the leaderboard for",
			Required:    true,
			Choices: []discord.ApplicationCommandOptionChoiceString{
				{
					Name:  "Global",
					Value: "Global",
				},
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
		discord.ApplicationCommandOptionUser{
			Name:        "user",
			Description: "The user to show the score for (leave empty to show the full leaderboard)",
			Required:    false,
		},
	},
}

func ShowCommandHandler(b *app.Bot) handler.SlashCommandHandler {
	return func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
		if err := e.DeferCreateMessage(true); err != nil {
			slog.Error("DisGo error(failed to defer interaction response)", slog.Any("err", err))
			return err
		}

		game := data.String("game")
		user, provided := data.OptUser("user")

		if provided && game == "Global" {
			// Get Global Score for user
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if res, err := b.DB.Queries.GetMemberGlobalScore(ctx, int64(user.ID)); err == nil {
				e.UpdateInteractionResponse(discord.MessageUpdate{
					Content: omit.Ptr(fmt.Sprintf("The global prediction score for %s is %d, ranked at %d", user.Mention(), res.Score, res.Position)),
				})
				return nil
			} else if err == sql.ErrNoRows {
				// If user's not found
				e.UpdateInteractionResponse(discord.MessageUpdate{
					Content: omit.Ptr(fmt.Sprintf("%s isn't found on the %s scoreboard", user.Mention(), game)),
				})
				return nil
			} else {
				e.UpdateInteractionResponse(discord.MessageUpdate{
					Content: omit.Ptr("Something wrong has happened, please try again"),
				})
				return err
			}
		} else if provided && game != "Global" {
			// Get Game Score for user
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if res, err := b.DB.Queries.GetMemberScoreForGame(ctx, sqlc.GetMemberScoreForGameParams{
				Game:   sqlc.ScoreboardGame(game),
				Member: int64(user.ID),
			}); err == nil {
				e.UpdateInteractionResponse(discord.MessageUpdate{
					Content: omit.Ptr(fmt.Sprintf("The %s prediction score for %s is %d, ranked at %d", game, user.Mention(), res.Score, res.Position)),
				})
				return nil
			} else if err == sql.ErrNoRows {
				// If user's not found
				e.UpdateInteractionResponse(discord.MessageUpdate{
					Content: omit.Ptr(fmt.Sprintf("%s isn't found on the %s scoreboard", user.Mention(), game)),
				})
				return nil
			} else {
				e.UpdateInteractionResponse(discord.MessageUpdate{
					Content: omit.Ptr("Something wrong has happened, please try again"),
				})
				return err
			}
		}

		if game == "Global" {
			return generateGlobalLeaderboard(b, e)
		} else {
			return generateGameLeaderboard(b, e, game)
		}
	}
}

func generateGameLeaderboard(b *app.Bot, e *handler.CommandEvent, game string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	scores, err := b.DB.Queries.ShowScoreboardForGame(ctx, sqlc.ScoreboardGame(game))
	if err != nil {
		return err
	}

	var layouts [][]discord.LayoutComponent
	totalPage := len(scores)/10 + 1

	for i := 1; i <= totalPage; i++ {
		buf := new(bytes.Buffer)
		table := tablewriter.NewTable(buf,
			tablewriter.WithRenderer(renderer.NewMarkdown()),
			tablewriter.WithRendition(tw.Rendition{
				Borders: tw.Border{
					Left:  tw.Off,
					Right: tw.Off,
				},
			}),
			tablewriter.WithAlignment(tw.Alignment{tw.AlignCenter}),
		)
		table.Header([]string{"Rank", "Name", "Score"})
		offset := (i - 1) * 10
		end := min(offset+10, len(scores))

		for _, score := range scores[offset:end] {
			// Check if member exists in cache
			if cachedMember, exists := e.Client().Caches.Member(*e.GuildID(), snowflake.ID(score.Member)); exists {
				name := truncate(cachedMember.EffectiveName())
				table.Append([]string{fmt.Sprint(score.Position), name, fmt.Sprint(score.Score)})
			} else {
				// Make API calls if not in cache
				if member, err := e.Client().Rest.GetMember(*e.GuildID(), snowflake.ID(score.Member)); err == nil {
					name := truncate(member.EffectiveName())
					table.Append([]string{fmt.Sprint(score.Position), name, fmt.Sprint(score.Score)})
				} else {
					// If not, fetch user info
					if user, err := e.Client().Rest.GetUser(snowflake.ID(score.Member)); err == nil {
						name := truncate(user.EffectiveName())
						table.Append([]string{fmt.Sprint(score.Position), name, fmt.Sprint(score.Score)})
					} else {
						// Fallback to unknown user
						name := "Unknown User"
						table.Append([]string{fmt.Sprint(score.Position), name, fmt.Sprint(score.Score)})
					}
				}
			}
		}
		table.Render()

		layout := []discord.LayoutComponent{
			discord.TextDisplayComponent{
				Content: fmt.Sprintf("%s Prediction Leaderboard", game),
			},
			discord.SeparatorComponent{},
			discord.ContainerComponent{
				Components: []discord.ContainerSubComponent{
					discord.TextDisplayComponent{
						Content: fmt.Sprint("```\n" + buf.String() + "\n```"),
					},
				},
			},
			discord.ActionRowComponent{
				Components: []discord.InteractiveComponent{
					discord.ButtonComponent{
						Style:    discord.ButtonStyleSecondary,
						Label:    "⏮️",
						CustomID: "prev_show",
					},
					discord.ButtonComponent{
						Style:    discord.ButtonStyleSecondary,
						Label:    "⏭️",
						CustomID: "next_show",
					},
				},
			},
		}
		layouts = append(layouts, layout)

	}

	_, err = e.UpdateInteractionResponse(discord.MessageUpdate{
		Components: omit.Ptr(layouts[0]),
		Flags:      omit.Ptr(discord.MessageFlagIsComponentsV2),
	})
	if err != nil {
		return err
	}

	go func() {
		ch, cls := bot.NewEventCollector(e.Client(),
			func(c *events.ComponentInteractionCreate) bool {
				return c.Data.CustomID() == "next_show" || c.Data.CustomID() == "prev_show"
			},
		)
		defer cls()
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		currentPage := 0
		for {
			select {
			case <-ctx.Done():
				return
			case c := <-ch:
				if c.Data.CustomID() == "next_show" {
					currentPage++
					if currentPage >= totalPage {
						currentPage = 0
					}
				} else if c.Data.CustomID() == "prev_show" {
					currentPage--
					if currentPage < 0 {
						currentPage = totalPage - 1
					}
				}
				c.UpdateMessage(discord.MessageUpdate{
					Components: omit.Ptr(layouts[currentPage]),
					Flags:      omit.Ptr(discord.MessageFlagIsComponentsV2),
				})
			default:
				continue
			}
		}
	}()

	return nil
}

func generateGlobalLeaderboard(b *app.Bot, e *handler.CommandEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	scores, err := b.DB.Queries.ShowGlobalScoreboard(ctx)
	if err != nil {
		return err
	}

	var layouts [][]discord.LayoutComponent
	totalPage := len(scores)/10 + 1

	for i := 1; i <= totalPage; i++ {
		buf := new(bytes.Buffer)
		table := tablewriter.NewTable(buf,
			tablewriter.WithRenderer(renderer.NewMarkdown()),
			tablewriter.WithRendition(tw.Rendition{
				Borders: tw.Border{
					Left:  tw.Off,
					Right: tw.Off,
				},
			}),
			tablewriter.WithAlignment(tw.Alignment{tw.AlignCenter}),
		)
		table.Header([]string{"Rank", "Name", "Score"})
		offset := (i - 1) * 10
		end := min(offset+10, len(scores))

		for _, score := range scores[offset:end] {
			// Check if member exists in cache
			if cachedMember, exists := e.Client().Caches.Member(*e.GuildID(), snowflake.ID(score.Member)); exists {
				name := truncate(cachedMember.EffectiveName())
				table.Append([]string{fmt.Sprint(score.Position), name, fmt.Sprint(score.Score)})
			} else {
				// Make API calls if not in cache
				if member, err := e.Client().Rest.GetMember(*e.GuildID(), snowflake.ID(score.Member)); err == nil {
					name := truncate(member.EffectiveName())
					table.Append([]string{fmt.Sprint(score.Position), name, fmt.Sprint(score.Score)})
				} else {
					// If not, fetch user info
					if user, err := e.Client().Rest.GetUser(snowflake.ID(score.Member)); err == nil {
						name := truncate(user.EffectiveName())
						table.Append([]string{fmt.Sprint(score.Position), name, fmt.Sprint(score.Score)})
					} else {
						// Fallback to unknown user
						name := "Unknown User"
						table.Append([]string{fmt.Sprint(score.Position), name, fmt.Sprint(score.Score)})
					}
				}
			}
		}
		table.Render()

		layout := []discord.LayoutComponent{
			discord.TextDisplayComponent{
				Content: "Global Prediction Leaderboard",
			},
			discord.SeparatorComponent{},
			discord.ContainerComponent{
				Components: []discord.ContainerSubComponent{
					discord.TextDisplayComponent{
						Content: fmt.Sprint("```\n" + buf.String() + "\n```"),
					},
				},
			},
			discord.ActionRowComponent{
				Components: []discord.InteractiveComponent{
					discord.ButtonComponent{
						Style:    discord.ButtonStyleSecondary,
						Label:    "⏮️",
						CustomID: "prev_show",
					},
					discord.ButtonComponent{
						Style:    discord.ButtonStyleSecondary,
						Label:    "⏭️",
						CustomID: "next_show",
					},
				},
			},
		}
		layouts = append(layouts, layout)

	}

	_, err = e.UpdateInteractionResponse(discord.MessageUpdate{
		Components: omit.Ptr(layouts[0]),
		Flags:      omit.Ptr(discord.MessageFlagIsComponentsV2),
	})
	if err != nil {
		return err
	}

	go func() {
		ch, cls := bot.NewEventCollector(e.Client(),
			func(c *events.ComponentInteractionCreate) bool {
				return c.Data.CustomID() == "next_show" || c.Data.CustomID() == "prev_show"
			},
		)
		defer cls()
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		currentPage := 0
		for {
			select {
			case <-ctx.Done():
				return
			case c := <-ch:
				if c.Data.CustomID() == "next_show" {
					currentPage++
					if currentPage >= totalPage {
						currentPage = 0
					}
				} else if c.Data.CustomID() == "prev_show" {
					currentPage--
					if currentPage < 0 {
						currentPage = totalPage - 1
					}
				}
				c.UpdateMessage(discord.MessageUpdate{
					Components: omit.Ptr(layouts[currentPage]),
					Flags:      omit.Ptr(discord.MessageFlagIsComponentsV2),
				})
			default:
				continue
			}
		}
	}()

	return fmt.Errorf("Todo")
}

func truncate(name string) string {
	if len(name) > 12 {
		return fmt.Sprintf("%s...", name[:9])
	}
	return name
}
