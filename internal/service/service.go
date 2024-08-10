package service

import (
	"MailingBot/internal/entity"
	"MailingBot/internal/repo"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot interface {
	MainBotFunc()
	SendingMailingList()
}

type Services struct {
	Bot Bot
}

type Dependencies struct {
	Repo          *repo.Repository
	Bot           *tgbotapi.BotAPI
	AllowedAdmins map[int64]entity.Admin
}

func NewServices(dep Dependencies) *Services {
	return &Services{
		Bot: NewBotService(dep.Bot, dep.Repo, dep.AllowedAdmins),
	}
}
