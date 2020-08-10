package entity

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
