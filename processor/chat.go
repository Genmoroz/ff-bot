package processor

import (
	"errors"
	"log"

	"ff-bot/bot"
	tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type ChatProcessor struct {
	chatID            int64
	tgBot             bot.Client
	updateChan        chan tgBotApi.Update
	stateProcessorMap map[string]StateProcessor
}

func New(chatID int64, tgBot bot.Client, stateProcessorMap map[string]StateProcessor) (*ChatProcessor, error) {
	if stateProcessorMap == nil {
		return nil, errors.New("stateProcessorMap cannot be nil")
	}
	if tgBot == nil {
		return nil, errors.New("tgBot cannot be nil")
	}
	return &ChatProcessor{
		chatID:            chatID,
		tgBot:             tgBot,
		updateChan:        make(chan tgBotApi.Update, 1),
		stateProcessorMap: stateProcessorMap,
	}, nil
}

func (p *ChatProcessor) PutUpdate(update tgBotApi.Update) error {
	if update.Message.Chat.ID != p.chatID {
		return errors.New("the message was not delivered, chatIDs do not match")
	}

	p.updateChan <- update

	return nil
}

func (p *ChatProcessor) Process() {
	for {
		update, ok := <-p.updateChan
		if !ok {
			log.Printf("updateChan is closed")
			return
		}

		text := update.Message.Text
		stateProcessor, exist := p.stateProcessorMap[text]
		if exist {
			if err := stateProcessor.Process(p.updateChan); err != nil {
				log.Printf("failed to handle the updateChan in processor[ID:%d]: %s", p.chatID, err.Error())
			}
		} else {
			msg := "Unknown command, you're in the main state. You may choose current state by command, to see all available commands enter /help"
			if err := p.tgBot.Send(msg, p.chatID); err != nil {
				log.Printf("failed to send the message to processor[ID:%d]: %s", p.chatID, err.Error())
			}
		}
	}
}

func (p ChatProcessor) GetChatID() int64 {
	return p.chatID
}

func (p *ChatProcessor) Close() error { // nolint:unparam
	p.stateProcessorMap = nil
	close(p.updateChan)

	return nil
}
