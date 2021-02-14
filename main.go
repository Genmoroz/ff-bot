package main

import (
	"log"

	"ff-bot/bot"
	"ff-bot/config"
	"ff-bot/dispatcher"
	"ff-bot/handler"
)

func main() {
	cfg, err := config.ReadEnv()
	if err != nil {
		log.Fatalf("failed to read the envs: %s", err.Error())
	}

	tbBot, err := bot.NewTBBotClient(cfg.TelegramToken)
	if err != nil {
		log.Fatalf("failed to create the Telegram bot: %s", err.Error())
	}

	disptch, err := dispatcher.New(tbBot, handler.NewUploadHandler(tbBot))
	if err != nil {
		log.Fatalf("failed to create the dispatcher: %s", err.Error())
	}

	updateChan, err := tbBot.GetUpdateChannel(0, 0, 60)
	if err != nil {
		log.Fatalf("failed to create the Telegram bot: %s", err.Error())
	}

	if err = disptch.Dispatch(updateChan); err != nil {
		log.Fatalf("failed to dispatch the updateChan: %s", err.Error())
	}
}
