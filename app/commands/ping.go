package commands

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var Ping = discord.SlashCommandCreate{
	Name:        "ping",
	Description: "Replies with Pong!",
}

func PingCommandHandler() handler.SlashCommandHandler {
	return func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
		latency := e.Client().Gateway.Latency()
		return e.CreateMessage(discord.MessageCreate{
			Content: fmt.Sprintf("Pong! Latency: %dms", latency.Milliseconds()),
		})
	}
}
