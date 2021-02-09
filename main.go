package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("1579541221:AAFw0p6T2TjL0wuTBzzpjG7Sr45RfvhJvlA")
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, commandResolver(update.Message.Text))

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

		_, err = bot.Send(msg)
		if err != nil {
			log.Printf("failed: %w", err)
		}
	}
}

func commandResolver(command string) string{
	switch command {
	case "/start":
		return "Welcome! This is First Frost Bot"
	}

	return "unknown command"
}
