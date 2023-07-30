package config

import (
	"fmt"
	"github.com/thehaung/juggernaut/internal/utils/typeutil"

	"github.com/spf13/viper"
)

type Telegram struct {
	Token string `mapstructure:"TELEGRAM_BOT_TOKEN" validate:"required"`
}

func loadFromEnv(config *Config) error {
	var telegram Telegram
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("can't find the file app.env: %w", err)
	}

	if err := loadForTarget(&telegram); err != nil {
		return fmt.Errorf("environment can't be loaded: %w", err)
	}

	config.Telegram = telegram
	return nil
}

func loadForTarget(targets ...interface{}) error {
	validator := typeutil.GetValidator()
	for _, target := range targets {
		err := viper.Unmarshal(&target)
		if err != nil {
			return fmt.Errorf("environment can't be loaded: %w", err)
		}

		if err = validator.Struct(target); err != nil {
			return fmt.Errorf("required environment can't be loaded: %w", err)
		}
	}

	return nil
}
