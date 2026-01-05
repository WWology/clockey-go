package main

import (
	"context"
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
	"clockey/app/commands/utils"
	"clockey/database/sqlc"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/handler"
	"github.com/jackc/pgx/v5/pgxpool"
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

	conn, err := pgxpool.New(context.Background(), cfg.Database.ConnectionString)
	if err != nil {
		slog.Error("Failed to connect to database", slog.Any("err", err))
		os.Exit(-1)
	}
	defer conn.Close()

	setupLogger(cfg.Log)
	slog.Info("Starting Clockey...", slog.String("version", Version), slog.String("commit", Commit))

	b := app.New(*cfg, Version, Commit, app.Database{
		Queries: sqlc.New(conn),
		Conn:    conn,
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
	h.SlashCommand("/add", predictions.AddCommandHandler(b))
	h.SlashCommand("/bo", predictions.BestOfCommandHandler())
	h.Autocomplete("/bo", predictions.BestOfAutocompleteHandler())
	h.SlashCommand("/deletebo", predictions.DeleteBestOfCommandHandler())
	h.Autocomplete("/deletebo", predictions.BestOfAutocompleteHandler())
	h.SlashCommand("/reset", predictions.ResetCommandHandler(b))
	h.SlashCommand("/show", predictions.ShowCommandHandler(b))
	h.SlashCommand("/winners", predictions.WinnersCommandHandler(b))
	// Utils
	h.SlashCommand("/util", utils.UtilCommandHandler())
	// Other
	h.SlashCommand("/ping", commands.PingCommandHandler())
	h.SlashCommand("/next", commands.NextCommandHandler())

	if err = b.SetupBot(h, bot.NewListenerFunc(b.OnReady), bot.NewListenerFunc(b.OnCommand), bot.NewListenerFunc(b.OnModal)); err != nil {
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

	// file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	// if err != nil {
	// 	slog.Error("Failed to open log file", slog.Any("err", err))
	// 	os.Exit(-1)
	// }

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
