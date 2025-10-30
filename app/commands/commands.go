package commands

import (
	"clockey/app/commands/signups"

	"github.com/disgoorg/disgo/discord"
)

var Commands = []discord.ApplicationCommandCreate{
	signups.Event,
	test,
	version,
}
