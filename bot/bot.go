package bot

import (
	"fmt"

	tbBot "github.com/go-telegram-bot-api/telegram-bot-api"
)

type (
	Client interface {
		Send(int64, string) error
		GetUpdateChannel(offset, limit, timeout int) (tbBot.UpdatesChannel, error)
	}

	client struct {
		bot *tbBot.BotAPI
	}
)

func NewTBBotClient(token string) (Client, error) {
	bot, err := tbBot.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("faield to create the bot: %w", err)
	}

	return &client{
		bot: bot,
	}, nil
}

func (c *client) Send(chatID int64, msg string) error {
	msgConfig := tbBot.NewMessage(chatID, msg)
	_, err := c.bot.Send(msgConfig)
	if err != nil {
		return err
	}
	return nil
}

func (c *client) GetUpdateChannel(offset, limit, timeout int) (tbBot.UpdatesChannel, error) {
	return c.bot.GetUpdatesChan(c.newUpdateConfig(offset, limit, timeout))
}

func (c *client) newUpdateConfig(offset, limit, timeout int) tbBot.UpdateConfig {
	updateConf := tbBot.NewUpdate(offset)
	updateConf.Limit = limit
	updateConf.Timeout = timeout

	return updateConf
}
