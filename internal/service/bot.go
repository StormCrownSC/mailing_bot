package service

import (
	"MailingBot/internal/cconstant"
	"MailingBot/internal/entity"
	"MailingBot/internal/repo"
	"MailingBot/pkg/utils"
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
)

type botService struct {
	repo                 *repo.Repository
	bot                  *tgbotapi.BotAPI
	allowedUser          map[int64]entity.Admin
	registrationAdminMap map[int64]RegistrationModel
	registrationUserMap  map[int64]RegistrationModel
	mailingChan          chan uint32
	mailingResultChan    chan bool
}

func NewBotService(bot *tgbotapi.BotAPI, repo *repo.Repository, allowedUser map[int64]entity.Admin) *botService {
	return &botService{
		repo:                 repo,
		bot:                  bot,
		allowedUser:          allowedUser,
		registrationAdminMap: make(map[int64]RegistrationModel),
		registrationUserMap:  make(map[int64]RegistrationModel),
		mailingChan:          make(chan uint32),
		mailingResultChan:    make(chan bool),
	}
}

func (o *botService) MainBotFunc() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := o.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			go o.botProcessing(update)
		}
	}
}

func (o *botService) botProcessing(update tgbotapi.Update) {
	if user, ok := o.allowedUser[update.Message.From.ID]; ok && !user.IsDeleted {
		if update.Message.IsCommand() {
			o.CommandHandler(update)
			return
		} else {
			if info, ok := o.registrationUserMap[update.Message.From.ID]; ok {
				o.userRegistration(update, info)
				return
			} else if info, ok := o.registrationAdminMap[update.Message.From.ID]; ok {
				o.adminRegistration(update, info)
				return
			} else {
				o.send(update, cconstant.HelpText, false)
			}
		}
	}
}

func (o *botService) CommandHandler(update tgbotapi.Update) {
	command := update.Message.Command()
	commandText := update.Message.CommandArguments()
	switch command {
	case "sending":
		textID, err := strconv.Atoi(utils.TextFormatting(commandText))
		if err != nil {
			o.send(update, "Текст id должен быть введён числом", false)
			return
		}
		o.mailingChan <- uint32(textID)
		status := <-o.mailingResultChan
		if status == false {
			o.send(update, "Ошибка при отправке рассылки", false)
			return
		}
		o.send(update, "Рассылка успешно отправлена", false)
	case "help":
		o.send(update, cconstant.HelpText, false)
	case "adduser":
		o.registrationUserMap[update.Message.From.ID] = RegistrationModel{
			Iteration:  0,
			TelegramID: 0,
			Name:       "",
		}
		o.send(update, cconstant.RegisterStep[0], false)
		o.send(update, "Если хотите отменить добавление пользователя введите `stop`", false)
	case "removeuser":
		telegramID, err := strconv.Atoi(utils.TextFormatting(commandText))
		if err != nil {
			o.send(update, "Телеграмм id должен быть введён числом", false)
			return
		}
		err = o.repo.User.RemoveUser(context.Background(), entity.User{TelegramID: int64(telegramID)})
		if err != nil {
			o.send(update, "Не удалось выключить пользователя", false)
			return
		}
		o.send(update, "Пользователь успешно выключен", false)
	case "recoveruser":
		telegramID, err := strconv.Atoi(utils.TextFormatting(commandText))
		if err != nil {
			o.send(update, "Телеграмм id должен быть введён числом", false)
			return
		}
		err = o.repo.User.RecoverUser(context.Background(), entity.User{TelegramID: int64(telegramID)})
		if err != nil {
			o.send(update, "Не удалось включить пользователя", false)
			return
		}
		o.send(update, "Пользователь успешно включен", false)
	case "getusers":
		users, err := o.repo.User.GetUsers(context.Background())
		if err != nil {
			o.send(update, "Не удалось получить список пользователей", false)
			return
		}
		var answer = "Список пользователей:\n"
		for index, user := range users {
			answer += fmt.Sprintf("%d) Пользователь - %s с telegram id - %d", index+1, user.Name, user.TelegramID)
			if user.IsDeleted {
				answer += " - отключен"
			}
			answer += "\n"
		}
		o.send(update, answer, false)
	case "addadmin":
		o.registrationUserMap[update.Message.From.ID] = RegistrationModel{
			Iteration:  0,
			TelegramID: 0,
			Name:       "",
		}
		o.send(update, cconstant.RegisterStep[0], false)
		o.send(update, "Если хотите отменить регистрацию админа введите `stop`", false)
	case "removeadmin":
		telegramID, err := strconv.Atoi(utils.TextFormatting(commandText))
		if err != nil {
			o.send(update, "Телеграмм id должен быть введён числом", false)
			return
		}
		err = o.repo.Admin.RemoveAdmin(context.Background(), entity.Admin{TelegramID: int64(telegramID)})
		if err != nil {
			o.send(update, "Не удалось выключить администратора", false)
			return
		}
		admin := o.allowedUser[int64(telegramID)]
		admin.IsDeleted = true
		o.allowedUser[int64(telegramID)] = admin

		o.send(update, "Администратор успешно выключен", false)
	case "recoveradmin":
		telegramID, err := strconv.Atoi(utils.TextFormatting(commandText))
		if err != nil {
			o.send(update, "Телеграмм id должен быть введён числом", false)
			return
		}
		err = o.repo.Admin.RecoverAdmin(context.Background(), entity.Admin{TelegramID: int64(telegramID)})
		if err != nil {
			o.send(update, "Не удалось включить администратора", false)
			return
		}
		admin := o.allowedUser[int64(telegramID)]
		admin.IsDeleted = false
		o.allowedUser[int64(telegramID)] = admin
		o.send(update, "Админимтратор успешно включен", false)
	case "getadmins":
		admins, err := o.repo.Admin.GetAdmins(context.Background())
		if err != nil {
			o.send(update, "Не удалось получить список администраторов", false)
			return
		}
		var answer = "Список администраторов:\n"
		for index, admin := range admins {
			answer += fmt.Sprintf("%d) Администратор - %s с telegram id - %d", index+1, admin.Name, admin.TelegramID)
			if admin.IsDeleted {
				answer += " - отключен"
			}
			answer += "\n"
		}
		o.send(update, answer, false)
	case "addtext":
		if commandText == "" {
			o.send(update, "Необходимо ввести текст", false)
			return
		}
		err := o.repo.Text.AddText(context.Background(), entity.Text{Text: commandText})
		if err != nil {
			o.send(update, "Не удалось добавить текст", false)
			return
		}
		o.send(update, "Текст успешно добавлен", false)
	case "removetext":
		id, err := strconv.Atoi(utils.TextFormatting(commandText))
		if err != nil {
			o.send(update, "id должен быть введён числом", false)
			return
		}
		err = o.repo.Text.RemoveText(context.Background(), entity.Text{ID: uint32(id)})
		if err != nil {
			o.send(update, "Не удалось отключить текст рассылки", false)
		}
		o.send(update, "Текст рассылки успешно выключен", false)
	case "recovertext":
		id, err := strconv.Atoi(utils.TextFormatting(commandText))
		if err != nil {
			o.send(update, "id должен быть введён числом", false)
			return
		}
		err = o.repo.Text.RecoverText(context.Background(), entity.Text{ID: uint32(id)})
		if err != nil {
			o.send(update, "Не удалось включить текст рассылки", false)
			return
		}
		o.send(update, "Текстр рассылки успешно включен", false)
	case "gettexts":
		texts, err := o.repo.Text.GetTexts(context.Background())
		if err != nil {
			o.send(update, "Не удалось получить тексты рассылок", false)
			return
		}
		var answer = "Список текстов:\n"
		for _, text := range texts {
			answer += fmt.Sprintf("%d) Текст: ```%s```", text.ID, text.Text)
			if text.IsDeleted {
				answer += " - отключен"
			}
			answer += "\n"
		}
		o.send(update, answer, false)

	default:
		o.send(update, "Такой команды нет! Введите `/help`, чтобы получить список команд", false)
	}
	return
}

func (o *botService) userRegistration(update tgbotapi.Update, info RegistrationModel) {
	messageText := update.Message.Text
	if messageText == "stop" {
		delete(o.registrationUserMap, update.Message.From.ID)
		o.send(update, "Процесс регистрации пользователя прерван", false)
		return
	}
	switch info.Iteration {
	case 0:
		telegramID, err := strconv.Atoi(utils.TextFormatting(messageText))
		if err != nil {
			o.send(update, "Телеграмм id должен быть введён числом, если хотите отменить регистрацию введите `stop`", false)
			return
		}
		info.TelegramID = int64(telegramID)
		info.Iteration += 1
		o.registrationUserMap[update.Message.From.ID] = info
		o.send(update, cconstant.RegisterStep[info.Iteration], false)
	case 1:
		name := messageText
		info.Name = name
		info.Iteration += 1
		err := o.repo.User.AddUser(context.Background(), entity.User{
			TelegramID: info.TelegramID,
			Name:       info.Name,
			IsDeleted:  false,
		})
		delete(o.registrationUserMap, update.Message.From.ID)
		if err != nil {
			o.send(update, "Не удалось зарегистрировать пользователя", false)
		} else {
			o.send(update, fmt.Sprintf("Пользователь %s успешно добавлен", info.Name), false)
		}
	default:
		delete(o.registrationUserMap, update.Message.From.ID)
		o.send(update, "Процесс добавления пользователя прошёл некорректно и отменён", false)
	}
}

func (o *botService) adminRegistration(update tgbotapi.Update, info RegistrationModel) {
	messageText := update.Message.Text
	if messageText == "stop" {
		delete(o.registrationUserMap, update.Message.From.ID)
		o.send(update, "Процесс регистрации администратора прерван", false)
		return
	}
	switch info.Iteration {
	case 0:
		telegramID, err := strconv.Atoi(utils.TextFormatting(messageText))
		if err != nil {
			o.send(update, "Телеграмм id должен быть введён числом, если хотите отменить регистрацию введите `stop`", false)
			return
		}
		info.TelegramID = int64(telegramID)
		info.Iteration += 1
		o.registrationUserMap[update.Message.From.ID] = info
		o.send(update, cconstant.RegisterStep[info.Iteration], false)
	case 1:
		name := messageText
		info.Name = name
		info.Iteration += 1
		admin := entity.Admin{TelegramID: info.TelegramID, Name: info.Name, IsDeleted: false}
		err := o.repo.Admin.AddAdmin(context.Background(), admin)
		delete(o.registrationUserMap, update.Message.From.ID)
		if err != nil {
			o.send(update, "Не удалось зарегистрировать администратора", false)
		} else {
			o.allowedUser[info.TelegramID] = admin
			o.send(update, fmt.Sprintf("Администратор %s успешно добавлен", info.Name), false)
		}
	default:
		delete(o.registrationUserMap, update.Message.From.ID)
		o.send(update, "Процесс добавления администратора прошёл некорректно и отменён", false)
	}
}
