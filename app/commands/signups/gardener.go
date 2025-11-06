package signups

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"strings"
	"time"

	"clockey/app"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/omit"
)

var Gardener = discord.MessageCommandCreate{
	Name: "Roll Gardener",
	Contexts: []discord.InteractionContextType{
		discord.InteractionContextTypeGuild,
	},
}

func GardenerCommandHandler(b *app.Bot) handler.MessageCommandHandler {
	return func(data discord.MessageCommandInteractionData, e *handler.CommandEvent) error {
		// Check if message has already been processed
		if processed(data.TargetMessage()) {
			e.CreateMessage(discord.MessageCreate{
				Content: "This message has been processed for signups",
				Flags:   discord.MessageFlagEphemeral,
			})
			return nil
		}

		// Show gardener selection menu
		gardenerSelectMenu, err := gardenerSelectMenuBuilder(e, data.TargetMessage())
		if err != nil {
			return err
		}

		if err := e.CreateMessage(discord.MessageCreate{
			Components: []discord.LayoutComponent{
				discord.ActionRowComponent{
					Components: []discord.InteractiveComponent{
						gardenerSelectMenu,
					},
				},
			},
			Flags: discord.MessageFlagEphemeral,
		}); err != nil {
			return err
		}

		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			defer cancel()
			bot.WaitForEvent(e.Client(), ctx,
				func(s *events.ComponentInteractionCreate) bool {
					return s.Data.CustomID() == "gardener_select_menu"
				},
				func(s *events.ComponentInteractionCreate) {
					selectedGardenerID := s.Data.(discord.StringSelectMenuInteractionData).Values[0]
					// gardenerID, _ := strconv.ParseInt(selectedGardenerID, 10, 64)

					eventType, name, eventTime, hours, err := parseMessage(data.TargetMessage().Content)
					print(eventType, name, eventTime, hours)
					if err != nil {
						s.Client().Logger.Error("Failed to parse message", slog.Any("err", err))
						return
					}

					// b.DB.Queries.CreateEvent(ctx, database.CreateEventParams{
					// 	Type:     eventType,
					// 	Name:     name,
					// 	Time:     eventTime,
					// 	Hours:    hours,
					// 	Gardener: gardenerID,
					// })

					if err := s.Client().Rest.AddReaction(data.TargetMessage().ChannelID, data.TargetMessage().ID, processedEmoji); err != nil {
						s.Client().Logger.Error("Failed to add reaction", slog.Any("err", err))
					}

					s.UpdateMessage(discord.MessageUpdate{
						Content:    omit.Ptr("Hours added to the database"),
						Components: &[]discord.LayoutComponent{},
					})

					if _, err := s.Client().Rest.CreateMessage(s.Message.ChannelID, discord.MessageCreate{
						MessageReference: &discord.MessageReference{
							Type:      discord.MessageReferenceTypeForward,
							MessageID: omit.Ptr(data.TargetMessage().ID),
							ChannelID: omit.Ptr(data.TargetMessage().ChannelID),
						},
					}); err != nil {
						s.Client().Logger.Error("Failed to send message reference", slog.Any("err", err))
					}

					if _, err := s.Client().Rest.CreateMessage(s.Message.ChannelID, discord.MessageCreate{
						Content: "<@" + selectedGardenerID + "> will be working " + name,
					}); err != nil {
						s.Client().Logger.Error("Failed to send message", slog.Any("err", err))
					}

				},
				func() {
					if err := e.CreateMessage(discord.MessageCreate{
						Content: "Gardener selection timed out.",
						Flags:   discord.MessageFlagEphemeral,
					}); err != nil {
						e.Client().Logger.Error("Failed to send timeout message", slog.Any("err", err))
					}
				},
			)
		}()

		return nil
	}
}

func processed(msg discord.Message) bool {
	for _, reaction := range msg.Reactions {
		if reaction.Emoji.Name == processedEmoji {
			return true
		}
	}
	return false
}

func gardenerSelectMenuBuilder(e *handler.CommandEvent, msg discord.Message) (discord.StringSelectMenuComponent, error) {
	gardenersReacted, err := e.Client().Rest.GetReactions(msg.ChannelID, msg.ID, signupEmoji, discord.MessageReactionTypeNormal, 0, 6)
	if err != nil {
		return discord.StringSelectMenuComponent{}, err
	}

	for idx, gardener := range gardenersReacted {
		if gardener.ID == e.Client().ApplicationID {
			gardenersReacted = append(gardenersReacted[:idx], gardenersReacted[idx+1:]...)
		}
	}
	gardenerSelectMenuOptions := []discord.StringSelectMenuOption{}

	for _, gardener := range gardenersReacted {
		if name, exists := gardenerIDsMap[gardener.ID]; exists {
			gardenerSelectMenuOptions = append(gardenerSelectMenuOptions, discord.StringSelectMenuOption{
				Label: name,
				Value: gardener.ID.String(),
			})
		} else {
			return discord.StringSelectMenuComponent{}, fmt.Errorf("unknown gardener ID: %d", gardener.ID)
		}
	}

	return discord.StringSelectMenuComponent{
		CustomID:    "gardener_select_menu",
		Placeholder: "Select the gardener working this event",
		Options:     gardenerSelectMenuOptions,
	}, nil
}

func parseMessage(msg string) (string, string, int64, int64, error) {
	var eventType string
	if strings.Contains(msg, "Dota") {
		eventType = "Dota"
	} else if strings.Contains(msg, "CS") {
		eventType = "CS"
	} else if strings.Contains(msg, "MLBB") {
		eventType = "MLBB"
	} else if strings.Contains(msg, "HoK") {
		eventType = "HoK"
	} else if strings.Contains(msg, "Other") {
		eventType = "Other"
	} else {
		return "", "", 0, 0, fmt.Errorf("failed to parse event type")
	}

	var name string
	nameRegex := regexp.MustCompile(`Event: \w+ - (.+?)(?:\n|$)`)
	nameMatch := nameRegex.FindStringSubmatch(msg)
	if len(nameMatch) > 1 {
		name = nameMatch[1]
	} else {
		return "", "", 0, 0, fmt.Errorf("failed to parse event name")
	}

	var eventTime int64
	timeRegex := regexp.MustCompile(`<t:([^:]+):F>`)
	timeMatch := timeRegex.FindStringSubmatch(msg)
	if len(timeMatch) > 1 {
		parsedTime, err := strconv.ParseInt(timeMatch[1], 10, 64)
		if err != nil {
			return "", "", 0, 0, fmt.Errorf("failed to parse event time")
		}
		eventTime = parsedTime
	} else {
		return "", "", 0, 0, fmt.Errorf("failed to parse event time")
	}

	var hours int64
	hoursRegex := regexp.MustCompile(`Hours: (\d+) hours`)
	hoursMatch := hoursRegex.FindStringSubmatch(msg)
	if len(hoursMatch) > 1 {
		parsedHours, err := strconv.ParseInt(hoursMatch[1], 10, 64)
		if err != nil {
			return "", "", 0, 0, fmt.Errorf("failed to parse event hours")
		}
		hours = parsedHours
	}

	return eventType, name, eventTime, hours, nil

}
