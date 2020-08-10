package app

import (
	"github.com/flum1025/tweam/internal/config"
	"github.com/flum1025/tweam/internal/entity"
)

type App struct {
	config *config.Config
}

func NewApp(config *config.Config) *App {
	return &App{
		config: config,
	}
}

func (a *App) PublishMessages(messages []*entity.Message) error {
	return nil
}
