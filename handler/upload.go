package handler

import (
	"errors"
	"log"

	"ff-bot/bot"
	tgBot "github.com/go-telegram-bot-api/telegram-bot-api"
)

type uploadHandler struct {
	tbBot bot.Client
}

func NewUploadHandler(tbBot bot.Client) Handler {
	return &uploadHandler{
		tbBot: tbBot,
	}
}

func (h *uploadHandler) Handle(updateChan tgBot.UpdatesChannel) error {
	if updateChan == nil {
		return errors.New("updateChan cannot be nil")
	}

	for {
		update := <-updateChan

		text := update.Message.Text
		chatID := update.Message.Chat.ID
		if text == End {
			return h.tbBot.Send(chatID, "End of the upload state.")
		}

		if err := h.tbBot.Send(chatID, text); err != nil {
			log.Printf("failed to send the message: %s", err.Error())
		}
	}
}
