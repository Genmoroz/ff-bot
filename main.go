package main

import (
	"context"
	"log"
	"sync"

	"ff-bot/bot"
	"ff-bot/config"
	"ff-bot/dispatcher"
	"ff-bot/handler"
	"ff-bot/router"
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

	handlerMap := make(map[string]handler.Handler)
	handlerMap[handler.Upload] = handler.NewUploadHandler(tbBot)
	handlerMap[handler.Start] = handler.NewStartHandler(tbBot)

	disptch, err := dispatcher.New(tbBot, handlerMap)
	if err != nil {
		log.Fatalf("failed to create the dispatcher: %s", err.Error())
	}

	updateChan, err := tbBot.GetUpdateChannel(0, 0, 60)
	if err != nil {
		log.Fatalf("failed to create the Telegram bot: %s", err.Error())
	}

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		if err = disptch.Dispatch(updateChan); err != nil {
			log.Fatalf("failed to dispatch the updateChan: %s", err.Error())
		}
		wg.Done()
	}()

	r := router.New(cfg.Router.Port)
	wg.Add(1)
	go func() {
		if err = r.ListenAndServeWithContext(context.TODO()); err != nil {
			log.Fatalf("failed to start the router: %s", err.Error())
		}
		wg.Done()
	}()

	wg.Wait()
}
