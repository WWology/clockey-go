package commands

import (
	"clockey/app/commands/predictions"
	"clockey/app/commands/signups"

	"github.com/disgoorg/disgo/discord"
)

var Commands = []discord.ApplicationCommandCreate{
	// Signups
	signups.Cancel,
	signups.Edit,
	signups.Event,
	signups.Gardener,
	signups.Manual,
	signups.Report,

	// Predictions
	predictions.Add,
	predictions.BestOf,
	predictions.DeleteBestOf,
	predictions.Reset,
	predictions.Show,
	predictions.Winners,

	// Other
	Next,
	Ping,
}
