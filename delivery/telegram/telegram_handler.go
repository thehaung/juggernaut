package telegram

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/thehaung/juggernaut/domain"
	"github.com/thehaung/juggernaut/internal/logger"
)

var (
	_supportCommand = map[string]bool{
		"ip":       true,
		"serverIP": true,
		"publicIP": true,
	}
)

type telegramHandler struct {
	logger            logger.Interface
	juggernautUseCase domain.JuggernautUseCase
	telegramBotAPI    *tgbotapi.BotAPI
	telegramChann     tgbotapi.UpdatesChannel
}

func NewTelegramHandler(juggernautUseCase domain.JuggernautUseCase, telegramBotAPI *tgbotapi.BotAPI, telegramChann tgbotapi.UpdatesChannel) {
	handler := &telegramHandler{
		logger:            logger.GetLogger(),
		juggernautUseCase: juggernautUseCase,
		telegramBotAPI:    telegramBotAPI,
		telegramChann:     telegramChann,
	}

	for update := range telegramChann {
		go func(update tgbotapi.Update) {
			if _supportCommand[update.Message.Text] {
				resMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
				ip, err := handler.getServerPublicIP(context.TODO())
				if err != nil {
					resMsg.Text = err.Error()
				} else {
					resMsg.Text = ip
				}

				if _, err = handler.telegramBotAPI.Request(resMsg); err != nil {
					handler.logger.Errorf("Exec Telegram Event failed. Error: %s", err)
				}
			}
		}(update)
	}
}

func (h *telegramHandler) getServerPublicIP(ctx context.Context) (string, error) {
	return h.juggernautUseCase.GetServerPublicIP(ctx)
}
