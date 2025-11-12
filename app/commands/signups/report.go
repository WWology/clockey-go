package signups

import (
	"clockey/app"
	"clockey/database/sqlc"
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var Report = discord.SlashCommandCreate{
	Name: "report",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionString{
			Name:        "start_date",
			Description: "Start date of the report, please use DD-MM-YYYY format",
			Required:    true,
		},
		discord.ApplicationCommandOptionString{
			Name:        "end_date",
			Description: "End date of the report, please use DD-MM-YYYY format",
			Required:    false,
		},
		discord.ApplicationCommandOptionString{
			Name:        "report_option",
			Description: "Should the report be separated by gardener / game",
			Required:    true,
			Choices: []discord.ApplicationCommandOptionChoiceString{
				{
					Name:  "Per Gardener",
					Value: "gardener",
				},
				{
					Name:  "Per Game",
					Value: "game",
				},
			},
		},
	},
}

func ReportCommandHandler(b *app.Bot) handler.SlashCommandHandler {
	return func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
		startDateString := data.String("start_date")
		startDate, err := time.Parse("02-01-2006", startDateString)
		if err != nil {
			return e.CreateMessage(discord.MessageCreate{
				Content: "Failed to parse starting date " + startDateString + ", please try again",
			})
		}

		var endDate time.Time
		if endDateString, provided := data.OptString("end_date"); provided {
			endDate, err = time.Parse("02-01-2006", endDateString)
			if err != nil {
				return e.CreateMessage(discord.MessageCreate{
					Content: "Failed to parse ending date " + endDateString + ", please try again",
				})
			}
		} else {
			endDate = time.Now()
		}

		if data.String("report_option") == "gardener" {
			return GenerateGardenerReport(b, startDate, endDate)
		} else {
			return GenerateGameReport(b, startDate, endDate)
		}

	}
}

func GenerateGameReport(b *app.Bot, startDate time.Time, endDate time.Time) error {
	// TODO
	panic("not implemented")
}

func GenerateGardenerReport(b *app.Bot, startDate time.Time, endDate time.Time) error {
	var wg sync.WaitGroup

	invoices := make(chan []sqlc.Event, 5)
	for id, name := range gardenerIDsMap {
		wg.Go(func() {
			if rows, err := b.DB.Queries.GetEventsForGardener(context.TODO(), sqlc.GetEventsForGardenerParams{
				Start:    startDate.Unix(),
				End:      endDate.Unix(),
				Gardener: int64(id),
			}); err == nil {
				invoices <- rows
			} else {
				b.Client.Logger.Error("Failed to get invoice ", slog.Any("name", name))
			}
		})
	}

	for _, events := range <-invoices {
		panic("todo")
	}

	wg.Wait()
	return nil
}
