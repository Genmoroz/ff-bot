package processor

import (
	"ff-bot/bot"
	tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	// Start defines a command for Start state
	Start = "/start"

	// Store defines a command for Upload state
	Store = "/store"

	// End defines a command that means end of any state
	End = "/end"
)

type (
	StateProcessor interface {
		Process(channel tgBotApi.UpdatesChannel) error
	}

	baseStateProcessor struct {
		tgBot  bot.Client
		chatID int64
	}
)

func newBaseStateProcessor(tgBot bot.Client, chatID int64) baseStateProcessor {
	return baseStateProcessor{
		tgBot:  tgBot,
		chatID: chatID,
	}
}
