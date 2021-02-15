package handler

import tgBot "github.com/go-telegram-bot-api/telegram-bot-api"

const (
	Upload = "/upload"
	End    = "/end"
)

type Handler interface {
	Handle(tgBot.UpdatesChannel, int64) error
}
