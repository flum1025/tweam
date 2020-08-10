package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/flum1025/tweam/internal/entity"
	"github.com/go-chi/render"
)

type event struct {
	ForUserID           string            `json:"for_user_id"`
	TweetCreateEvents   []twitter.Tweet   `json:"tweet_create_events"`
	DirectMessageEvents json.RawMessage   `json:"direct_message_events"`
	FavoriteEvents      []json.RawMessage `json:"favorite_events"`
	FollowEvents        []json.RawMessage `json:"follow_events"`
	TweetDeleteEvents   []json.RawMessage `json:"tweet_delete_events"`
}

func (s *server) twistributer(w http.ResponseWriter, r *http.Request) {
	var params event

	_, err := parse(r.Body, &params)
	if err != nil {
		log.Println(fmt.Sprintf("[ERROR] failed to parse body: %v", err))
		return
	}

	messages := parseEvent(params)

	// for _, message := range messages {
	// 	body, err := json.Marshal(message.Body)
	// 	if err != nil {
	// 		log.Println(fmt.Sprintf("[ERROR] failed to marshal json: %v", err))
	// 		return
	// 	}
	//
	// 	spew.Dump(string(body))
	//
	// 	if err = s.sqsClient.SendMessage(
	// 		string(body),
	// 		message.GroupID,
	// 		message.DeduplicationID,
	// 	); err != nil {
	// 		log.Println(fmt.Sprintf("[ERROR] failed to send message: %v", err))
	// 		return
	// 	}
	// }

	render.PlainText(w, r, "ok")
}

func parseEvent(params event) []entity.Message {
	messages := make([]entity.Message, 0)

	for _, tweet := range params.TweetCreateEvents {
		messages = append(messages, entity.Tweet(tweet).NewMessage(params.ForUserID))
	}

	// TODO: direct_message_events

	for _, event := range params.FavoriteEvents {
		messages = append(
			messages,
			entity.Message{
				ForUserID: params.ForUserID,
				Type:      entity.MessageTypeFavorite,
				Data:      event,
			},
		)
	}

	for _, event := range params.FollowEvents {
		messages = append(
			messages,
			entity.Message{
				ForUserID: params.ForUserID,
				Type:      entity.MessageTypeFollow,
				Data:      event,
			},
		)
	}

	for _, event := range params.TweetDeleteEvents {
		messages = append(
			messages,
			entity.Message{
				ForUserID: params.ForUserID,
				Type:      entity.MessageTypeTweetDelete,
				Data:      event,
			},
		)
	}

	return messages
}
