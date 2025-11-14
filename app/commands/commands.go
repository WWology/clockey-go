package commands

import (
	"clockey/app/commands/signups"

	"github.com/disgoorg/disgo/discord"
)

var Commands = []discord.ApplicationCommandCreate{
	signups.Event,
	signups.Gardener,
	signups.Manual,
	signups.Edit,
	signups.Report,
	test,
	version,
}
