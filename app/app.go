package app

import (
	"context"
	"log/slog"
	"time"

	"clockey/database/sqlc"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
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
