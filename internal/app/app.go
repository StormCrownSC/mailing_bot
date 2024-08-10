package app

import (
	"MailingBot/config"
	"MailingBot/internal/entity"
	"MailingBot/internal/repo"
	"MailingBot/internal/service"
	"MailingBot/pkg/storage/postgres"
	"MailingBot/pkg/utils"
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os/signal"
	"sync"
	"syscall"
)

func Run(configFile string) {
	// initializing config
	cfg, err := config.NewConfig(configFile)

	if err != nil {
		log.Fatalf("error initializing config: %s", err.Error())
	}
	// initializing logger

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
	)
	defer cancel()

	var (
		botOnce sync.Once
		bot     *tgbotapi.BotAPI
	)

	botOnce.Do(func() {
		bot, err = tgbotapi.NewBotAPI(cfg.Telegram)
		if err != nil {
			log.Fatalf("Can't initialize bot %s", err)
		}
	})

	// initializing postgres storage
	var cPostgresData = postgres.Postgres{}
	utils.CopyData(cfg.Postgres, &cPostgresData)
	db := postgres.NewStorage(&cPostgresData)
	err = db.Connect(ctx)
	if err != nil {
		log.Fatalf("error creating postgres storage: %s", err.Error())
	}
	log.Println("using postgres storage")

	repository := repo.NewRepository(db.DB)

	admins, err := repository.Admin.GetAdmins(context.Background())
	if err != nil {
		log.Fatalf("error getting admins: %s", err.Error())
	}

	var allowedAdmins map[int64]entity.Admin = make(map[int64]entity.Admin)

	for _, admin := range admins {
		allowedAdmins[admin.TelegramID] = admin
	}

	// initializing services
	dep := service.Dependencies{
		Repo:          repository,
		Bot:           bot,
		AllowedAdmins: allowedAdmins,
	}

	services := service.NewServices(dep)

	go services.Bot.SendingMailingList()
	go services.Bot.MainBotFunc()

	select {
	case <-ctx.Done():
		log.Println("stopping mailing bot...")
		if err != nil {
			log.Println(fmt.Errorf("error stopping mailing bot: %w", err))
		}

		log.Println("bot is stopped")
	}
}
