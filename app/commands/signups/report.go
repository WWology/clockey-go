package signups

import (
	"clockey/app"
	"clockey/database/sqlc"
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/omit"
)

var Report = discord.SlashCommandCreate{
	Name:        "report",
	Description: "Look at the report of gardener signups over a period of time",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionString{
			Name:        "start_date",
			Description: "Start date of the report, please use DD-MM-YYYY format",
			Required:    true,
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
		discord.ApplicationCommandOptionString{
			Name:        "end_date",
			Description: "End date of the report, please use DD-MM-YYYY format",
			Required:    false,
		},
	},
}

func ReportCommandHandler(b *app.Bot) handler.SlashCommandHandler {
	return func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
		e.DeferCreateMessage(true)
		startDateString := data.String("start_date")
		startDate, err := time.Parse("02-01-2006", startDateString)
		if err != nil {
			if _, err := e.UpdateInteractionResponse(discord.MessageUpdate{
				Content: omit.Ptr("Failed to parse " + startDateString + ", please try again"),
			}); err != nil {
				return err
			}
			return err
		}

		var endDate time.Time
		if endDateString, provided := data.OptString("end_date"); provided {
			endDate, err = time.Parse("02-01-2006", endDateString)
			if err != nil {
				if _, err := e.UpdateInteractionResponse(discord.MessageUpdate{
					Content: omit.Ptr("Failed to parse " + endDateString + ", please try again"),
				}); err != nil {
					return err
				}
				return err
			}
		} else {
			endDate = time.Now()
		}

		if data.String("report_option") == "gardener" {
			return GenerateGardenerReport(b, e, startDate, endDate)
		} else {
			return GenerateGameReport(b, e, startDate, endDate)
		}

	}
}

func GenerateGameReport(b *app.Bot, e *handler.CommandEvent, startDate time.Time, endDate time.Time) error {
	//
	//  TODO
	panic("unimplemented")
}

type Invoice struct {
	Gardener string
	Events   []sqlc.Event
}

func GenerateGardenerReport(b *app.Bot, e *handler.CommandEvent, startDate time.Time, endDate time.Time) error {
	var wg sync.WaitGroup

	invoices := make(chan Invoice, 5)
	for id, name := range gardenerIDsMap {
		wg.Go(func() {
			if events, err := b.DB.Queries.GetEventsForGardener(context.TODO(), sqlc.GetEventsForGardenerParams{
				Start:    startDate.Unix(),
				End:      endDate.Unix(),
				Gardener: int64(id),
			}); err == nil {
				invoices <- Invoice{Gardener: name, Events: events}
			} else {
				b.Client.Logger.Error("Failed to get invoice ", slog.Any("name", name))
			}
			// events := test.GetTestEventsForGardener(int64(id))
			// invoices <- Invoice{Gardener: name, Events: events}
		})
	}
	wg.Wait()
	close(invoices)

	layouts := make(map[string][]discord.LayoutComponent, 5)
	for invoice := range invoices {
		var dotaEvents, csEvents, mlbbEvents, hokEvents, otherEvents string
		gardenerHours := 0
		for _, event := range invoice.Events {
			schedule := time.Unix(event.Time, 0).Format("02-01-2006")
			switch event.Type {
			case "Dota":
				dotaEvents += fmt.Sprintf("%s at %s - %d hours\n", event.Name, schedule, event.Hours)
			case "CS":
				csEvents += fmt.Sprintf("%s at %s - %d hours\n", event.Name, schedule, event.Hours)
			case "MLBB":
				mlbbEvents += fmt.Sprintf("%s at %s - %d hours\n", event.Name, schedule, event.Hours)
			case "HoK":
				hokEvents += fmt.Sprintf("%s at %s - %d hours\n", event.Name, schedule, event.Hours)
			case "Other":
				otherEvents += fmt.Sprintf("%s at %s - %d hours\n", event.Name, schedule, event.Hours)
			}
			gardenerHours += int(event.Hours)
		}

		layouts[invoice.Gardener] = []discord.LayoutComponent{
			discord.TextDisplayComponent{
				Content: fmt.Sprintf("# %s's Invoice\n**%s - %s**", invoice.Gardener, startDate.Month().String(), endDate.Month().String()),
			},
			discord.ContainerComponent{
				Components: []discord.ContainerSubComponent{
					discord.TextDisplayComponent{
						Content: "# Dota\n" + dotaEvents,
					},
					discord.SeparatorComponent{},
					discord.TextDisplayComponent{
						Content: "# CS\n" + csEvents,
					},
					discord.SeparatorComponent{},
					discord.TextDisplayComponent{
						Content: "# MLBB\n" + mlbbEvents,
					},
					discord.SeparatorComponent{},
					discord.TextDisplayComponent{
						Content: "# HoK\n" + hokEvents,
					},
					discord.SeparatorComponent{},
					discord.TextDisplayComponent{
						Content: "# Other\n" + otherEvents,
					},
					discord.SeparatorComponent{},
					discord.TextDisplayComponent{
						Content: fmt.Sprintf("**Total: %d**", gardenerHours),
					},
				},
			},
			discord.ActionRowComponent{
				Components: []discord.InteractiveComponent{
					discord.ButtonComponent{
						Label:    "N1k",
						Style:    discord.ButtonStyleSecondary,
						CustomID: "n1k_button",
						Disabled: strings.Contains(invoice.Gardener, "N1k"),
					},
					discord.ButtonComponent{
						Label:    "Kit",
						Style:    discord.ButtonStyleSecondary,
						CustomID: "kit_button",
						Disabled: strings.Contains(invoice.Gardener, "Kit"),
					},
					discord.ButtonComponent{
						Label:    "WW",
						Style:    discord.ButtonStyleSecondary,
						CustomID: "ww_button",
						Disabled: strings.Contains(invoice.Gardener, "WW"),
					},
					discord.ButtonComponent{
						Label:    "Bonteng",
						Style:    discord.ButtonStyleSecondary,
						CustomID: "bonteng_button",
						Disabled: strings.Contains(invoice.Gardener, "Bonteng"),
					},
					discord.ButtonComponent{
						Label:    "Sam",
						Style:    discord.ButtonStyleSecondary,
						CustomID: "sam_button",
						Disabled: strings.Contains(invoice.Gardener, "Sam"),
					},
				},
			},
		}

	}

	_, err := e.UpdateInteractionResponse(discord.MessageUpdate{
		Components: omit.Ptr(layouts["N1k"]),
		Flags:      omit.Ptr(discord.MessageFlagIsComponentsV2),
	})
	if err != nil {
		e.Client().Logger.Error("Failed to send invoice message", slog.Any("err", err))
		return err
	}

	go func() {
		ch, cls := bot.NewEventCollector(e.Client(),
			func(c *events.ComponentInteractionCreate) bool {
				return c.Data.CustomID() == "n1k_button" ||
					c.Data.CustomID() == "kit_button" ||
					c.Data.CustomID() == "ww_button" ||
					c.Data.CustomID() == "bonteng_button" ||
					c.Data.CustomID() == "sam_button"
			})
		defer cls()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()
		for {
			select {
			case <-ctx.Done():
				return
			case c := <-ch:
				fmt.Println(c.Data.CustomID())
				switch c.Data.CustomID() {
				case "n1k_button":
					if err := c.UpdateMessage(discord.MessageUpdate{
						Components: omit.Ptr(layouts["N1k"]),
						Flags:      omit.Ptr(discord.MessageFlagIsComponentsV2),
					}); err != nil {
						c.Client().Logger.Error("Failed to update message to N1k", slog.Any("err", err))
					}
				case "kit_button":
					if err := c.UpdateMessage(discord.MessageUpdate{
						Components: omit.Ptr(layouts["Kit"]),
						Flags:      omit.Ptr(discord.MessageFlagIsComponentsV2),
					}); err != nil {
						c.Client().Logger.Error("Failed to update message to Kit", slog.Any("err", err))
					}
				case "ww_button":
					if err := c.UpdateMessage(discord.MessageUpdate{
						Components: omit.Ptr(layouts["WW"]),
						Flags:      omit.Ptr(discord.MessageFlagIsComponentsV2),
					}); err != nil {
						c.Client().Logger.Error("Failed to update message to WW", slog.Any("err", err))
					}
				case "bonteng_button":
					if err := c.UpdateMessage(discord.MessageUpdate{
						Components: omit.Ptr(layouts["Bonteng"]),
						Flags:      omit.Ptr(discord.MessageFlagIsComponentsV2),
					}); err != nil {
						c.Client().Logger.Error("Failed to update message to Bonteng", slog.Any("err", err))
					}
				case "sam_button":
					if err := c.UpdateMessage(discord.MessageUpdate{
						Components: omit.Ptr(layouts["Sam"]),
						Flags:      omit.Ptr(discord.MessageFlagIsComponentsV2),
					}); err != nil {
						c.Client().Logger.Error("Failed to update message to Sam", slog.Any("err", err))
					}
				}

			}
		}
	}()

	return nil
}
