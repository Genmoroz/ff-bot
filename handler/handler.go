package handler

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

const (
	// commands
	Upload = "/upload"
	End    = "/end"
)

type Handler interface {
	Handle(updateChan tgbotapi.UpdatesChannel) error
}
