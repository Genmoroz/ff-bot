package processor

import (
	"ff-bot/bot"
	"log"
)

func sendAndPrint(msg string, chatID int64, tgBot bot.Client) {
	if err := tgBot.Send(msg, chatID); err != nil {
		log.Printf("failed to send the message[chatID:%d]: %s", chatID, err.Error())
	}
	log.Print(msg)
}
