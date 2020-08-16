package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/flum1025/tweam/internal/config"
	"github.com/flum1025/tweam/internal/entity"
	"github.com/flum1025/tweam/internal/redis"
)

const timeout = time.Duration(25 * time.Second)

type App struct {
	config *config.Config
	redis  *redis.RedisClient
}

func NewApp(config *config.Config) (*App, error) {
	redis, err := redis.NewRedisClient(config)
	if err != nil {
		return nil, fmt.Errorf("get redis client: %w", err)
	}

	return &App{
		config: config,
		redis:  redis,
	}, nil
}

func (a *App) PublishMessages(messages []entity.Message) error {
	for _, message := range messages {
		key, err := message.Key()
		if err != nil {
			return fmt.Errorf("generate key: %w", err)
		}

		exists, err := a.redis.Exists(key)
		if err != nil {
			return fmt.Errorf("redis exists: %w", err)
		}

		if !exists {
			err := a.redis.Set(key, 1, time.Minute*30)
			if err != nil {
				return fmt.Errorf("redis set: %w", err)
			}

			if err := a.publishMessage(message); err != nil {
				log.Println(fmt.Sprintf("[ERROR] failed to publish message: %v", err))
			}
		}
	}

	return nil
}

func (a *App) publishMessage(message entity.Message) error {
	account := a.config.Accounts.Find(message.ForUserID)
	if account == nil {
		log.Println(fmt.Sprintf("not target user: %v", message.ForUserID))
		return nil
	}

	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("json marshal: %w", err)
	}

	var wg sync.WaitGroup

	for _, webhook := range account.Webhooks {
		wg.Add(1)

		go func(webhook string) {
			log.Println(fmt.Sprintf("requesting to %v, user_id=%v", webhook, message.ForUserID))

			client := http.Client{Timeout: timeout}
			_, err := client.Post(webhook, "application/json", bytes.NewBuffer(body))

			if err != nil {
				log.Println(fmt.Sprintf("[ERROR] request to %v: %v", webhook, err))
				wg.Done()
				return
			}

			wg.Done()
			log.Println(fmt.Sprintf("delivered to %v, user_id=%v", webhook, message.ForUserID))
		}(webhook)
	}

	wg.Wait()

	return nil
}
