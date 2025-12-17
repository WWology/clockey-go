package signups

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"clockey/app"
	"clockey/database/sqlc"

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
			return e.CreateMessage(discord.MessageCreate{
				Content: "This message has been processed for signups",
				Flags:   discord.MessageFlagEphemeral,
			})
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
					gardenerID, _ := strconv.ParseInt(selectedGardenerID, 10, 64)

					eventType, name, eventTime, hours, err := parseMessage(data.TargetMessage().Content)
					if err != nil {
						s.Client().Logger.Error("Failed to parse message", slog.Any("err", err))
						return
					}

					if err := b.DB.Queries.CreateEvent(ctx, sqlc.CreateEventParams{
						Type:     eventType,
						Name:     name,
						Time:     eventTime,
						Hours:    hours,
						Gardener: gardenerID,
					}); err != nil {
						s.Client().Logger.Error("Failed to create event in database", slog.Any("err", err))
						return
					}

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
