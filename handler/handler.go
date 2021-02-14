package handler

import tgBot "github.com/go-telegram-bot-api/telegram-bot-api"

const (
	// commands
	Upload = "/upload"
	End    = "/end"
)

type Handler interface {
	Handle(updateChan tgBot.UpdatesChannel) error
}
