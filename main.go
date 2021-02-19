package main

import (
	"context"
	"log"
	"sync"

	"ff-bot/bot"
	"ff-bot/config"
	"ff-bot/dispatcher"
	"ff-bot/processor"
	"ff-bot/router"
)

func main() {
	cfg, err := config.ReadEnv()
	if err != nil {
		log.Fatalf("failed to read the envs: %s", err.Error())
	}

	tgBot, err := bot.NewTGBotClient(cfg.TelegramToken)
	if err != nil {
		log.Fatalf("failed to create the Telegram bot: %s", err.Error())
	}

	disptchr, err := dispatcher.New(tgBot, createStateProcessorMapFunc(cfg.FileStorePath))
	if err != nil {
		log.Fatalf("failed to create the dispatcher: %s", err.Error())
	}

	updateChan, err := tgBot.GetUpdateChannel(0, 0, 60)
	if err != nil {
		log.Fatalf("failed to create the Telegram bot: %s", err.Error())
	}

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		if err = disptchr.Dispatch(updateChan); err != nil {
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

func createStateProcessorMapFunc(fileStorePath string) func(tgBot bot.Client, chatID int64) map[string]processor.StateProcessor {
	return func(tgBot bot.Client, chatID int64) map[string]processor.StateProcessor {
		return map[string]processor.StateProcessor{
			processor.Start: processor.NewStartStateProcessor(tgBot, chatID),
			processor.Store: processor.NewStoreStateProcessor(tgBot, chatID, fileStorePath),
		}
	}
}
