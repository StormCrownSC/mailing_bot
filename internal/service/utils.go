package service

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (o *botService) sendMessageToChat(chatID int64, messageText string) {
	msg := tgbotapi.NewMessage(chatID, messageText)
	msg.ParseMode = "Markdown"
	_, err := o.bot.Send(msg)
	fmt.Println(err, chatID, messageText)
	if err != nil {
		msg.ParseMode = "HTML"
		o.bot.Send(msg)
	}
}

func (o *botService) send(update tgbotapi.Update, text string, replay bool) {
	var msg tgbotapi.MessageConfig

	msg = tgbotapi.NewMessage(update.Message.Chat.ID, text)
	msg.ParseMode = "Markdown"

	if replay {
		msg.ReplyToMessageID = update.Message.MessageID
	}

	_, err := o.bot.Send(msg)
	if err != nil {
		msg.ParseMode = "HTML"
		o.bot.Send(msg)
	}
}
