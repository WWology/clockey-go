package app

import (
	"context"
	"log/slog"
	"slices"
	"time"

	"clockey/database/sqlc"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/snowflake/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func New(cfg Config, version string, commit string, db Database) *Bot {
	return &Bot{
		Cfg:     cfg,
		Version: version,
		Commit:  commit,
		DB:      db,
	}
}

type Database struct {
	Queries *sqlc.Queries
	Conn    *pgxpool.Pool
}

type Bot struct {
	Cfg     Config
	Client  *bot.Client
	Version string
	Commit  string
	DB      Database
}

func (b *Bot) SetupBot(listeners ...bot.EventListener) error {
	client, err := disgo.New(b.Cfg.Bot.Token,
		bot.WithGatewayConfigOpts(gateway.WithIntents(gateway.IntentGuilds, gateway.IntentGuildMessages, gateway.IntentMessageContent, gateway.IntentGuildMembers, gateway.IntentGuildPresences)),
		bot.WithCacheConfigOpts(cache.WithCaches(cache.FlagGuilds|cache.FlagMembers)),
		bot.WithEventListeners(listeners...),
	)
	if err != nil {
		return err
	}

	b.Client = client
	return nil
}

func (b *Bot) OnReady(e *events.Ready) {
	slog.Info("Logged in", slog.String("user", e.User.Username))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := b.Client.SetPresence(ctx, gateway.WithListeningActivity("you"), gateway.WithOnlineStatus(discord.OnlineStatusOnline)); err != nil {
		slog.Error("Failed to set presence", slog.Any("err", err))
	}
}

func (b *Bot) OnCommand(e *events.ApplicationCommandInteractionCreate) {
	slog.LogAttrs(
		context.Background(),
		slog.LevelInfo,
		"Command used",
		slog.String("command", e.Data.CommandName()),
		slog.String("user", e.Member().EffectiveName()),
		slog.Any("data", e.SlashCommandInteractionData().Options),
	)
}

func (b *Bot) OnModal(m *events.ModalSubmitInteractionCreate) {
	slog.LogAttrs(
		context.Background(),
		slog.LevelInfo,
		"Modal submitted",
		slog.String("modal", m.Data.CustomID),
		slog.GroupAttrs("modal",
			slog.String("event_type", m.Data.StringValues("event_type")[0]),
			slog.String("event_name", m.Data.Text("event_name")),
			slog.String("event_time", m.Data.Text("event_time")),
			slog.String("event_duration", m.Data.Text("event_duration")),
		),
	)
}

func (b *Bot) OnMessageCreate(e *events.MessageCreate) {
	// Honeypot: automatically ban any user who posts in the designated channel
	if e.ChannelID == snowflake.ID(1459863214434287798) {
		gardenerRoleID := snowflake.ID(720253636797530203)
		if slices.Contains(e.Message.Member.RoleIDs, gardenerRoleID) {
			return
		}
		err := e.Client().Rest.AddBan(*e.GuildID, e.Message.Author.ID, time.Duration(168*time.Hour))
		if err != nil {
			slog.Error("failed to ban user", "error", err)
		}

		slog.LogAttrs(
			context.Background(),
			slog.LevelInfo,
			"scam & bot detected",
			slog.String("user", e.Message.Author.Username),
			slog.String("content", e.Message.Content),
		)
	}

}
