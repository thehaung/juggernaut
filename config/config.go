package config

import (
	"errors"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	HttpServer HttpServer
	Telegram   Telegram
}

type HttpServer struct {
	Port                 string
	Timeout              int
	ExcludeLogRouterPath map[string]bool
}

func Parse() (*Config, error) {
	var config Config
	v := viper.New()

	v.SetConfigName("./config/config")
	v.AddConfigPath(".")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	err := v.Unmarshal(&config)
	if err != nil {
		log.Printf("unable to decode into struct. error: %v", err)
		return nil, err
	}

	err = loadFromEnv(&config)
	if err != nil {
		log.Printf("unable to load env variable. error: %v", err)
		return nil, err
	}

	return &config, nil
}
