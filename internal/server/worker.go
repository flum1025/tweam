package server

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/flum1025/tweam/internal/config"
	"github.com/flum1025/tweam/internal/entity"
	"github.com/go-chi/render"
)

const timeout = time.Duration(25 * time.Second)

func (s *server) worker(w http.ResponseWriter, r *http.Request) {
	var params entity.MessageBody

	body, err := parse(r.Body, &params)
	if err != nil {
		log.Println(fmt.Sprintf("[ERROR] failed to parse body: %v", err))

		return
	}

	if params.ForUserID == "" {
		log.Println("for_user_id is required")

		return
	}

	target := func() *config.Account {
		for _, account := range s.config.Accounts {
			if account.ID == params.ForUserID {
				return &account
			}
		}

		return nil
	}()

	if target == nil {
		log.Println(fmt.Sprintf("not target user: %v", params.ForUserID))

		return
	}

	for _, webhook := range target.Webhooks {
		go func(webhook string) {
			log.Println(fmt.Sprintf("requesting to %v, user_id=%v", webhook, params.ForUserID))

			client := http.Client{Timeout: timeout}
			_, err := client.Post(webhook, "application/json", bytes.NewBuffer(body))

			if err != nil {
				log.Println(fmt.Sprintf("[ERROR] request to %v: %v", webhook, err))
				return
			}

			log.Println(fmt.Sprintf("delivered to %v, user_id=%v", webhook, params.ForUserID))
		}(webhook)
	}

	render.PlainText(w, r, "ok")
}
