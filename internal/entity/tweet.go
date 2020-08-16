package entity

import (
	"github.com/dghubble/go-twitter/twitter"
)

type Tweet twitter.Tweet
type Tweets []Tweet

func (e Tweet) NewMessage(userID string) Message {
	return Message{
		ForUserID: userID,
		Type:      MessageTypeTweet,
		Data:      e,
	}
}

func (e Tweets) Messages(userID string) []Message {
	ms := make([]Message, 0, len(e))

	for _, m := range e {
		ms = append(ms, m.NewMessage(userID))
	}

	return ms
}
