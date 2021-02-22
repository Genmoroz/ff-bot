package state

import (
	"fmt"

	bot "github.com/genvmoroz/bot-engine/api"
	"github.com/genvmoroz/bot-engine/processor"
	tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const Start = "/start"

type startStateProcessor struct {
	tgBot  bot.Client
	chatID int64
}

func NewStartStateProcessor(tbBot bot.Client, chatID int64) processor.StateProcessor {
	return &startStateProcessor{
		tgBot:  tbBot,
		chatID: chatID,
	}
}

func (p *startStateProcessor) Process(_ <-chan tgBotApi.Update) error {
	msg := "Hey there, this is First Frost Bot. Author genvmoroz@gmail.com. To list all available commands enter /help."
	if err := p.tgBot.Send(msg, p.chatID); err != nil {
		return fmt.Errorf("failed to send the message: %w", err)
	}

	return nil
}
