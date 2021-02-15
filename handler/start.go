package handler

import (
	"ff-bot/bot"
	"fmt"
	tgBot "github.com/go-telegram-bot-api/telegram-bot-api"
)

type startHandler struct {
	tbBot bot.Client
}

func NewStartHandler(tbBot bot.Client) Handler {
	return &startHandler{
		tbBot: tbBot,
	}
}

func (h *startHandler) Handle(_ tgBot.UpdatesChannel, chatID int64) error {
	msg := "Hey there, this is First Frost Bot. Author genvmoroz@gmail.com. To list all available commands enter /help."
	if err := h.tbBot.Send(msg, chatID); err != nil {
		return fmt.Errorf("failed to send the message to chat: %w", err)
	}

	return nil
}
