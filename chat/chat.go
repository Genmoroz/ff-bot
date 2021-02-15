package chat

import (
	"errors"
	"log"

	"ff-bot/bot"
	"ff-bot/handler"
	tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Chat struct {
	id         int64
	tgBot      bot.Client
	updateChan chan tgBotApi.Update
	handlerMap map[string]handler.Handler
}

func New(id int64, tgBot bot.Client, handlerMap map[string]handler.Handler) (Chat, error) {
	if handlerMap == nil {
		return Chat{}, errors.New("handlerMap cannot be nil")
	}
	if tgBot == nil {
		return Chat{}, errors.New("tgBot cannot be nil")
	}
	return Chat{
		id:         id,
		tgBot:      tgBot,
		updateChan: make(chan tgBotApi.Update, 1),
		handlerMap: handlerMap,
	}, nil
}

func (c Chat) PutUpdate(update tgBotApi.Update) error {
	if update.Message.Chat.ID != c.id {
		return errors.New("the message was not delivered, chatIDs do not match")
	}

	c.updateChan <- update

	return nil
}

func (c Chat) Start() {
	for {
		update := <-c.updateChan

		text := update.Message.Text
		h, exist := c.handlerMap[text]
		if exist {
			if err := h.Handle(c.updateChan, c.id); err != nil {
				log.Printf("failed to handle the updateChan in chat[ID:%d]: %s", c.id, err.Error())
			}
		} else {
			msg := "Unknown command, you're in the main state. You may choose current state by command, to see all available commands enter /help"
			if err := c.tgBot.Send(msg, c.id); err != nil {
				log.Printf("failed to send the message to chat[ID:%d]: %s", c.id, err.Error())
			}
		}
	}
}
