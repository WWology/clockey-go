package signups

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

type ChannelConfig struct {
	VoiceChannels    map[string]snowflake.ID
	ExternalChannels map[string]string
	StageChannel     snowflake.ID
}

var Channels = ChannelConfig{
	VoiceChannels: map[string]snowflake.ID{
		"Dota": snowflake.ID(738009797932351519),
		"CS":   snowflake.ID(746618267434614804),
	},
	ExternalChannels: map[string]string{
		"MLBB": "https://discord.com/channels/689865753662455829/1350252799019188236",
		"HoK":  "https://discord.com/channels/689865753662455829/1344676860562509955",
	},
	StageChannel: snowflake.ID(1186593338300842025),
}

var OGGames = []discord.StringSelectMenuOption{
	{
		Label: "Dota",
		Value: "Dota",
	},
	{
		Label: "CS",
		Value: "CS",
	},
	{
		Label: "MLBB",
		Value: "MLBB",
	},
	{
		Label: "HoK",
		Value: "HoK",
	},
	{
		Label: "Other",
		Value: "Other",
	},
}

var Gardeners = []discord.ApplicationCommandOptionChoiceString{
	{
		Name:  "N1k",
		Value: "N1k",
	},
	{
		Name:  "Kit",
		Value: "Kit",
	},
	{
		Name:  "WW",
		Value: "WW",
	},
	{
		Name:  "Bonteng",
		Value: "Bonteng",
	},
	{
		Name:  "Sam",
		Value: "Sam",
	},
}

var GardenerRoleID = snowflake.ID(720253636797530203)
var SignupEmojiString = "<:OGpeepoYes:730890894814740541>"
