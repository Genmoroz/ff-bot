package dispatcher

import (
	"errors"
	"fmt"
	"log"

	"ff-bot/bot"
	"ff-bot/chat"
	"ff-bot/handler"
	tbBot "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Dispatcher struct {
	tbBot      bot.Client
	chatMap    map[int64]*chat.Chat
	handlerMap map[string]handler.Handler
}

func New(tbBot bot.Client, handlerMap map[string]handler.Handler) (*Dispatcher, error) {
	if tbBot == nil {
		return nil, errors.New("tbBot cannot be nil")
	}
	if handlerMap == nil {
		return nil, errors.New("handlerMap cannot be nil")
	}

	return &Dispatcher{
		tbBot:      tbBot,
		chatMap:    make(map[int64]*chat.Chat),
		handlerMap: handlerMap,
	}, nil
}

func (d *Dispatcher) Dispatch(updateChan tbBot.UpdatesChannel) error {
	if d == nil {
		return errors.New("dispatcher cannot be nil")
	}
	if updateChan == nil {
		return errors.New("updateChan cannot be nil")
	}

	for {
		update, ok := <-updateChan
		if !ok {
			log.Fatal("updateChan is closed")
		}

		chatID := update.Message.Chat.ID
		existedChat, exist := d.chatMap[chatID]
		if exist {
			go d.putUpdateIntoChatAndLog(existedChat, update, chatID)
		} else {
			newChat, err := chat.New(chatID, d.tbBot, d.handlerMap)
			if err != nil {
				return fmt.Errorf("failed to create a new chat[ID:%d]: %w", chatID, err)
			}
			d.chatMap[chatID] = newChat

			go newChat.Start()
			go d.putUpdateIntoChatAndLog(newChat, update, chatID)
		}
	}
}

func (d *Dispatcher) putUpdateIntoChatAndLog(c *chat.Chat, update tbBot.Update, chatID int64) {
	if err := c.PutUpdate(update); err != nil {
		log.Printf("failed to put the update into the chat[ID:%d]: %s", chatID, err.Error())
	}
}
