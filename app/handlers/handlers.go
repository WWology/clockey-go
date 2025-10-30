package handlers

import (
	"clockey/app"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
)

func MessageHandler(b *app.Bot) bot.EventListener {
	return bot.NewListenerFunc(func(e *events.MessageCreate) {
		// TODO: handle message
	})
}
