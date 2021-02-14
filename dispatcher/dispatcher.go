package dispatcher

import (
	"errors"
	"log"

	"ff-bot/bot"
	"ff-bot/handler"
	tbBot "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Dispatcher struct {
	tbBot         bot.Client
	uploadHandler handler.Handler
}

func New(tbBot bot.Client, uploadHandler handler.Handler) (*Dispatcher, error) {
	if tbBot == nil {
		return nil, errors.New("tbBot cannot be nil")
	}
	if uploadHandler == nil {
		return nil, errors.New("uploadHandler cannot be nil")
	}

	return &Dispatcher{
		tbBot:         tbBot,
		uploadHandler: uploadHandler,
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
		update := <-updateChan

		text := update.Message.Text
		if text == handler.Upload {
			if err := d.tbBot.Send(update.Message.Chat.ID, "You're in the upload state."); err != nil {
				log.Printf("failed to send the message to chat: %s", err.Error())
			}
			if err := d.uploadHandler.Handle(updateChan); err != nil {
				log.Printf("failed to handler upload: %s", err.Error())
			}
		}

	}
}
