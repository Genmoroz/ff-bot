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

func (h *uploadHandler) Handle(updateChan tgBot.UpdatesChannel, chatID int64) error {
	if updateChan == nil {
		return errors.New("updateChan cannot be nil")
	}

	if err := h.tbBot.Send("You're in the upload state.", chatID); err != nil {
		log.Printf("failed to send the message to chat: %s", err.Error())
	}
	for {
		update := <-updateChan

		text := update.Message.Text
		if text == End {
			return h.tbBot.Send("End of the upload state.", chatID)
		}

		if err := h.tbBot.Send(text, chatID); err != nil {
			log.Printf("failed to send the message: %s", err.Error())
		}
	}
}
