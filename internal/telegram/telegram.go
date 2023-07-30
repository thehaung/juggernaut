package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/thehaung/juggernaut/config"
	"github.com/thehaung/juggernaut/internal/logger"
)

type Telegram struct {
	logger         logger.Interface
	conf           *config.Config
	telegramBotAPI *tgbotapi.BotAPI
}

func NewTelegram(conf *config.Config) *Telegram {
	t := &Telegram{
		logger: logger.GetLogger(),
		conf:   conf,
	}
	t.establish()
	return t
}

func (t *Telegram) establish() {
	bot, err := tgbotapi.NewBotAPI(t.conf.Telegram.Token)

	if err != nil {
		t.logger.Errorf("Telegram - tgbotapi.NewBotAPI(). Error: %s", err)
		panic(err)
	}
	t.logger.Infof("Authorized on account: %s", bot.Self.UserName)
	t.telegramBotAPI = bot
}

func (t *Telegram) GetUpdateChannel() tgbotapi.UpdatesChannel {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	return t.telegramBotAPI.GetUpdatesChan(u)
}

func (t *Telegram) GetTelegramBotAPI() *tgbotapi.BotAPI {
	return t.telegramBotAPI
}
