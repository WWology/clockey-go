package main

import (
	"context"
	"database/sql"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"clockey/app"
	"clockey/app/commands"
	"clockey/app/commands/predictions"
	"clockey/app/commands/signups"
	"clockey/database/sqlc"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/handler"
	_ "modernc.org/sqlite"
)

var (
	Version = "dev"
	Commit  = "unknown"
)

func main() {
	shouldSyncCommands := flag.Bool("sync-commands", true, "Whether to sync commands to discord")
	path := flag.String("config", "config.toml", "path to config")
	flag.Parse()

	cfg, err := app.LoadConfig(*path)
	if err != nil {
		slog.Error("Failed to read config", slog.Any("err", err))
		os.Exit(-1)
	}

	db, err := sql.Open("sqlite", "database/clockey.db")
	if err != nil {
		slog.Error("Failed to open database", slog.Any("err", err))
		os.Exit(-1)
	}
	defer db.Close()

	setupLogger(cfg.Log)
	slog.Info("Starting Clockey...", slog.String("version", Version), slog.String("commit", Commit))

	b := app.New(*cfg, Version, Commit, app.Database{
		Queries: sqlc.New(db),
		Conn:    db,
	})

	h := handler.New()
	// Signups
	h.MessageCommand("/Cancel Event", signups.CancelCommandHandler(b))
	h.SlashCommand("/edit", signups.EditCommandHandler(b))
	h.SlashCommand("/event", signups.EventCommandHandler())
	h.MessageCommand("/Roll Gardener", signups.GardenerCommandHandler(b))
	h.SlashCommand("/manual", signups.ManualCommandHandler(b))
	h.SlashCommand("/report", signups.ReportCommandHandler(b))
	// Predictions
	h.SlashCommand("/bo", predictions.BestOfCommandHandler())
	h.SlashCommand("/deletebo", predictions.DeleteBestOfCommandHandler())
	h.Autocomplete("/bo", predictions.BestOfAutocompleteHandler())
	h.Autocomplete("/deletebo", predictions.BestOfAutocompleteHandler())
	h.SlashCommand("/show", predictions.ShowCommandHandler(b))
	// Other
	h.SlashCommand("/next", commands.NextCommandHandler())

	if err = b.SetupBot(h, bot.NewListenerFunc(b.OnReady)); err != nil {
		slog.Error("Failed to setup bot", slog.Any("err", err))
		os.Exit(-1)
	}

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		b.Client.Close(ctx)
	}()

	if *shouldSyncCommands {
		slog.Info("Syncing commands", slog.Bool("sync", *shouldSyncCommands))
		slog.Info("Syncing commands", slog.Int("commands", len(commands.Commands)), slog.Any("guild_ids", cfg.Bot.DevGuilds))
		if err = handler.SyncCommands(b.Client, commands.Commands, cfg.Bot.DevGuilds); err != nil {
			slog.Error("Failed to sync commands", slog.Any("err", err))
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = b.Client.OpenGateway(ctx); err != nil {
		slog.Error("Failed to open gateway", slog.Any("err", err))
		os.Exit(-1)
	}

	slog.Info("Bot is running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
	<-s
	slog.Info("Shutting down bot...")
}

func setupLogger(cfg app.LogConfig) {
	opts := &slog.HandlerOptions{
		AddSource: cfg.AddSource,
		Level:     cfg.Level,
	}

	var sHandler slog.Handler
	switch cfg.Format {
	case "json":
		sHandler = slog.NewJSONHandler(os.Stdout, opts)
	case "text":
		sHandler = slog.NewTextHandler(os.Stdout, opts)
	default:
		slog.Error("Unknown log format", slog.String("format", cfg.Format))
		os.Exit(-1)
	}
	slog.SetDefault(slog.New(sHandler))
}
