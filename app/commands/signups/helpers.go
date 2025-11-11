package signups

import (
	"log/slog"
	"net/http"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

type ChannelConfig struct {
	VoiceChannels    map[string]snowflake.ID
	ExternalChannels map[string]string
	StageChannel     snowflake.ID
}

// Production channel and role IDs

// var Channels = ChannelConfig{
// 	VoiceChannels: map[string]snowflake.ID{
// 		"Dota": snowflake.ID(738009797932351519),
// 		"CS":   snowflake.ID(746618267434614804),
// 	},
// 	ExternalChannels: map[string]string{
// 		"MLBB": "https://discord.com/channels/689865753662455829/1350252799019188236",
// 		"HoK":  "https://discord.com/channels/689865753662455829/1344676860562509955",
// 	},
// 	StageChannel: snowflake.ID(1186593338300842025),
// }

// var GardenerRoleID = snowflake.ID(720253636797530203)
// var SignupEmoji = "OGpeepoYes:730890894814740541"
// var ProcessedEmoji = "OGwecoo:787697278190223370"

// Dev channel and role IDs
var channels = ChannelConfig{
	VoiceChannels: map[string]snowflake.ID{
		"Dota": snowflake.ID(1435509993947398206),
		"CS":   snowflake.ID(1435510029846446080),
	},
	ExternalChannels: map[string]string{
		"MLBB": "https://discord.com/channels/738607619660578876/738607620566286397",
		"HoK":  "https://discord.com/channels/738607619660578876/940936483819380757",
	},
	StageChannel: snowflake.ID(991620472544440454),
}

const (
	gardenerRoleID = snowflake.ID(1435510452795871232)
	signupEmoji    = "khezuBrain:1329032244580323349"
	processedEmoji = "ruggahPain:951843834554376262"
)

var oGGames = []discord.StringSelectMenuOption{
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

var gardeners = []discord.ApplicationCommandOptionChoiceString{
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

var gardenerIDsMap = map[snowflake.ID]string{
	293360731867316225: "N1k",
	204923365205475329: "Kit",
	754724309276164159: "WW",
	172360818715918337: "Bonteng",
	332438787588227072: "Sam",
}

var eventModal = discord.ModalCreate{
	CustomID: "event_modal",
	Title:    "Event Modal",
	Components: []discord.LayoutComponent{
		discord.LabelComponent{
			Label:       "Event Type",
			Description: "Select the type of event",
			Component: discord.StringSelectMenuComponent{
				CustomID: "event_type",
				Options:  oGGames,
				Required: true,
			},
		},
		discord.LabelComponent{
			Label:       "Event Name",
			Description: "Enter the name of the event",
			Component: discord.TextInputComponent{
				CustomID:    "event_name",
				Style:       discord.TextInputStyleShort,
				Placeholder: "OG vs <opp team name>",
				Required:    true,
			},
		},
		discord.LabelComponent{
			Label:       "Event Schedule",
			Description: "Enter the unix time for the start of this event",
			Component: discord.TextInputComponent{
				CustomID:    "event_time",
				Style:       discord.TextInputStyleShort,
				Required:    true,
				Placeholder: "Insert unix time from hammertime here",
			},
		},
		discord.LabelComponent{
			Label:       "Event duration",
			Description: "How many hours is this event",
			Component: discord.TextInputComponent{
				CustomID: "event_duration",
				Style:    discord.TextInputStyleShort,
				Required: true,
			},
		},
		discord.LabelComponent{
			Label:       "Event Banner",
			Description: "The banner for this event (if any, 800x320 px in size). ",
			Component: discord.FileUploadComponent{
				CustomID: "event_banner",
				Required: false,
			},
		},
	},
}

func getBanner(attachment discord.Attachment, Logger *slog.Logger) *discord.Icon {
	resp, err := http.Get(attachment.URL)
	if err != nil {
		Logger.Error("Failed to create banner icon", slog.Any("err", err))
		return nil
	}
	banner, err := discord.NewIcon(discord.IconTypeJPEG, resp.Body)
	if err != nil {
		Logger.Error("Failed to create banner icon", slog.Any("err", err))
		return nil
	}
	return banner
}
