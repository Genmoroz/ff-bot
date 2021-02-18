package processor

import (
	"fmt"

	"ff-bot/bot"
	tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type startStateProcessor struct {
	baseStateProcessor
}

func NewStartStateProcessor(tbBot bot.Client, chatID int64) StateProcessor {
	return &startStateProcessor{
		baseStateProcessor: newBaseStateProcessor(tbBot, chatID),
	}
}

func (p *startStateProcessor) Process(_ tgBotApi.UpdatesChannel) error {
	msg := "Hey there, this is First Frost Bot. Author genvmoroz@gmail.com. To list all available commands enter /help."
	if err := p.tgBot.Send(msg, p.chatID); err != nil {
		errMsg := fmt.Sprintf("failed to send the message[chatID:%d]: %s", p.chatID, err.Error())
		sendAndPrint(errMsg, p.chatID, p.tgBot)
	}

	return nil
}
