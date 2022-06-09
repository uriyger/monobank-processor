package config

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

// Config holds required params for the service
type Config struct {
	HTTPPort       int    `validate:"required"`
	ChatID         string `validate:"required"`
	SendMessageURL string `validate:"required"`
}

// NewConfig creates config from env or .env file
func NewConfig() (*Config, error) {
	viper.AutomaticEnv()

	config := &Config{
		HTTPPort:       viper.GetInt("HTTP_PORT"),
		ChatID:         viper.GetString("TG_CHAT_ID"),
		SendMessageURL: viper.GetString("TG_SEND_MESSAGE_URL"),
	}

	if err := validator.New().Struct(config); err != nil {
		return nil, fmt.Errorf("config validation: %w", err)
	}

	return config, nil
}
