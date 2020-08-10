package entity

import (
	"github.com/dghubble/go-twitter/twitter"
)

type Tweet twitter.Tweet

func (e Tweet) NewMessage(userID string) Message {
	return Message{
		ForUserID: userID,
		Type:      MessageTypeTweet,
		Data:      e,
	}
}
