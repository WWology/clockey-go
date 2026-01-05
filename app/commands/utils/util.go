package utils

import (
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var Util = discord.SlashCommandCreate{
	Name:        "util",
	Description: "Command that can be reused for multiple purposes, don't use",
}

func UtilCommandHandler() handler.SlashCommandHandler {
	return func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
		if err := e.DeferCreateMessage(true); err != nil {
			slog.Error("failed to defer interaction response", "error", err)
			return err
		}
		emojis, err := e.Client().Rest.GetEmojis(*e.GuildID())
		if err != nil {
			slog.Error("failed to get emojis", "error", err)
			return err
		}

		for _, emoji := range emojis {
			go func() {
				var opts discord.CDNOpt
				if emoji.Animated {
					opts = discord.WithFormat(discord.FileFormatAVIF)
				} else {
					opts = discord.WithFormat(discord.FileFormatPNG)
				}

				resp, err := http.Get(emoji.URL(opts))
				if err != nil {
					slog.Error("failed to get emoji URL", slog.Any("err", err))
					return
				}
				defer resp.Body.Close()

				var filename string
				if emoji.Animated {
					filename = filepath.Join("downloads", emoji.Name+".avif")
				} else {
					filename = filepath.Join("downloads", emoji.Name+".png")
				}
				out, err := os.Create(filename)
				if err != nil {
					slog.Error("failed to create file", slog.Any("err", err))
					return
				}
				defer out.Close()

				_, err = io.Copy(out, resp.Body)
				if err != nil {
					slog.Error("failed to save emoji to file", slog.Any("err", err))
					return
				}

			}()
		}

		return nil
	}
}
