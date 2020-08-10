package entity

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

type MessageType string

var (
	MessageTypeTweet         MessageType = "tweet"
	MessageTypeDirectMessage MessageType = "directmessage"
	MessageTypeFavorite      MessageType = "favorite"
	MessageTypeFollow        MessageType = "follow"
	MessageTypeTweetDelete   MessageType = "tweet_delete"
)

type Message struct {
	ForUserID string      `json:"for_user_id"`
	Type      MessageType `json:"type"`
	Data      interface{} `json:"data"`
}

func (m Message) Key() (string, error) {
	switch m.Type {
	case MessageTypeTweet:
		return fmt.Sprintf("%s_%s_%d", m.ForUserID, m.Type, m.Data.(Tweet).ID), nil
	case MessageTypeFavorite, MessageTypeFollow, MessageTypeTweetDelete:
		j, err := json.Marshal(m.Data)
		if err != nil {
			return "", fmt.Errorf("json marshal: %w", err)
		}

		hash := md5.Sum(j)
		return fmt.Sprintf("%s_%s_%s", m.ForUserID, m.Type, hex.EncodeToString(hash[:])), nil
	}

	return "", fmt.Errorf("unsupported type: %v", m.Type)
}
